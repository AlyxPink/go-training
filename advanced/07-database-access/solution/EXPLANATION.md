# Database Access Solution - Deep Dive

## Overview

This solution demonstrates production-grade database access patterns using Go's `database/sql` package, including connection pooling, prepared statements, transactions, proper error handling, and efficient scanning patterns.

## Architecture

### 1. Repository Pattern

```go
type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*User, error) {
    query := `SELECT id, name, email, created_at FROM users WHERE id = $1`

    var user User
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.CreatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, ErrUserNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("query user: %w", err)
    }

    return &user, nil
}
```

**Why repository pattern:**
- Abstracts database implementation
- Centralizes query logic
- Easy to test with interfaces
- Clean separation of concerns

### 2. Connection Pool Management

```go
func OpenDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }

    // Connection pool settings
    db.SetMaxOpenConns(25)                 // Maximum connections
    db.SetMaxIdleConns(5)                  // Idle connections to keep
    db.SetConnMaxLifetime(5 * time.Minute) // Max connection reuse time
    db.SetConnMaxIdleTime(10 * time.Minute) // Max idle time

    // Verify connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("ping database: %w", err)
    }

    return db, nil
}
```

**Pool tuning considerations:**
- `MaxOpenConns`: Should not exceed database max connections
- `MaxIdleConns`: Balance between connection reuse and resource usage
- `ConnMaxLifetime`: Prevents stale connections (especially with load balancers)
- `ConnMaxIdleTime`: Cleanup idle connections to free resources

### 3. Transaction Management

```go
func (r *UserRepository) CreateWithProfile(ctx context.Context, user *User, profile *Profile) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback() // Safe to call even after commit

    // Insert user
    err = tx.QueryRowContext(ctx,
        `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
        user.Name, user.Email,
    ).Scan(&user.ID)
    if err != nil {
        return fmt.Errorf("insert user: %w", err)
    }

    // Insert profile
    profile.UserID = user.ID
    _, err = tx.ExecContext(ctx,
        `INSERT INTO profiles (user_id, bio, avatar) VALUES ($1, $2, $3)`,
        profile.UserID, profile.Bio, profile.Avatar,
    )
    if err != nil {
        return fmt.Errorf("insert profile: %w", err)
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit transaction: %w", err)
    }

    return nil
}
```

**Transaction best practices:**
- Always defer `Rollback()` (idempotent after commit)
- Use `context.Context` for timeouts and cancellation
- Keep transactions short
- Handle errors at each step

## Key Patterns

### Pattern 1: Prepared Statements

```go
type UserRepository struct {
    db            *sql.DB
    findByIDStmt  *sql.Stmt
    findByEmailStmt *sql.Stmt
}

func NewUserRepository(db *sql.DB) (*UserRepository, error) {
    findByID, err := db.Prepare(`SELECT id, name, email FROM users WHERE id = $1`)
    if err != nil {
        return nil, err
    }

    findByEmail, err := db.Prepare(`SELECT id, name, email FROM users WHERE email = $1`)
    if err != nil {
        findByID.Close()
        return nil, err
    }

    return &UserRepository{
        db:              db,
        findByIDStmt:    findByID,
        findByEmailStmt: findByEmail,
    }, nil
}

func (r *UserRepository) Close() error {
    r.findByIDStmt.Close()
    r.findByEmailStmt.Close()
    return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*User, error) {
    var user User
    err := r.findByIDStmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.Name, &user.Email)
    return &user, err
}
```

**Benefits:**
- Query parsed once, executed many times
- Better performance for repeated queries
- Protection against SQL injection

**When to use:**
- Frequently executed queries
- Performance-critical paths
- Security-sensitive queries

**When not to use:**
- One-off queries
- Dynamic queries with varying structure
- Adds complexity for marginal benefit

### Pattern 2: Efficient Scanning

```go
// Scan single row
func (r *UserRepository) FindByID(ctx context.Context, id int) (*User, error) {
    var user User
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID, &user.Name, &user.Email,
    )
    return &user, err
}

// Scan multiple rows
func (r *UserRepository) FindAll(ctx context.Context) ([]*User, error) {
    rows, err := r.db.QueryContext(ctx, `SELECT id, name, email FROM users`)
    if err != nil {
        return nil, err
    }
    defer rows.Close() // CRITICAL: Always close rows

    var users []*User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
            return nil, err
        }
        users = append(users, &user)
    }

    // Check for errors during iteration
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return users, nil
}

// Scan with helper function
func scanUser(row interface{ Scan(...interface{}) error }) (*User, error) {
    var user User
    err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

**Common mistakes:**
- Forgetting to call `rows.Close()`
- Not checking `rows.Err()` after iteration
- Scanning into wrong types

### Pattern 3: NULL Handling

```go
type User struct {
    ID        int
    Name      string
    Email     string
    Bio       sql.NullString  // Handles NULL bio
    Age       sql.NullInt64   // Handles NULL age
    Avatar    sql.NullString
    UpdatedAt sql.NullTime
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*User, error) {
    var user User
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.Bio,      // sql.NullString handles NULL
        &user.Age,      // sql.NullInt64 handles NULL
        &user.Avatar,
        &user.UpdatedAt,
    )
    return &user, err
}

// Usage:
if user.Bio.Valid {
    fmt.Println(user.Bio.String)
} else {
    fmt.Println("No bio")
}

// Alternative: Custom types
type NullableString struct {
    String string
    Valid  bool
}

func (ns *NullableString) Scan(value interface{}) error {
    if value == nil {
        ns.String, ns.Valid = "", false
        return nil
    }
    ns.Valid = true
    return convertAssign(&ns.String, value)
}
```

### Pattern 4: Dynamic Queries

```go
func (r *UserRepository) Search(ctx context.Context, filters UserFilters) ([]*User, error) {
    query := `SELECT id, name, email FROM users WHERE 1=1`
    args := []interface{}{}
    argPos := 1

    if filters.Name != "" {
        query += fmt.Sprintf(` AND name LIKE $%d`, argPos)
        args = append(args, "%"+filters.Name+"%")
        argPos++
    }

    if filters.MinAge > 0 {
        query += fmt.Sprintf(` AND age >= $%d`, argPos)
        args = append(args, filters.MinAge)
        argPos++
    }

    if filters.Email != "" {
        query += fmt.Sprintf(` AND email = $%d`, argPos)
        args = append(args, filters.Email)
        argPos++
    }

    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    return scanUsers(rows)
}
```

**Better: Using query builder library (squirrel, goqu)**

```go
import sq "github.com/Masterminds/squirrel"

func (r *UserRepository) Search(ctx context.Context, filters UserFilters) ([]*User, error) {
    qb := sq.Select("id", "name", "email").
        From("users").
        PlaceholderFormat(sq.Dollar)

    if filters.Name != "" {
        qb = qb.Where(sq.Like{"name": "%" + filters.Name + "%"})
    }

    if filters.MinAge > 0 {
        qb = qb.Where(sq.GtOrEq{"age": filters.MinAge})
    }

    if filters.Email != "" {
        qb = qb.Where(sq.Eq{"email": filters.Email})
    }

    query, args, err := qb.ToSql()
    if err != nil {
        return nil, err
    }

    rows, err := r.db.QueryContext(ctx, query, args...)
    // ... scan as usual
}
```

## Error Handling

### Database-Specific Errors

```go
import (
    "github.com/lib/pq" // PostgreSQL driver
)

func (r *UserRepository) Create(ctx context.Context, user *User) error {
    _, err := r.db.ExecContext(ctx,
        `INSERT INTO users (email, name) VALUES ($1, $2)`,
        user.Email, user.Name,
    )

    if err != nil {
        // Check for unique constraint violation
        if pqErr, ok := err.(*pq.Error); ok {
            if pqErr.Code == "23505" { // unique_violation
                return ErrDuplicateEmail
            }
        }
        return fmt.Errorf("insert user: %w", err)
    }

    return nil
}
```

### Context Cancellation

```go
func (r *UserRepository) FindByID(ctx context.Context, id int) (*User, error) {
    var user User
    err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name)

    if err != nil {
        // Check if context was cancelled
        if ctx.Err() == context.Canceled {
            return nil, errors.New("operation cancelled")
        }
        if ctx.Err() == context.DeadlineExceeded {
            return nil, errors.New("operation timed out")
        }
        return nil, err
    }

    return &user, nil
}
```

## Performance Optimization

### 1. Batch Operations

```go
func (r *UserRepository) CreateBatch(ctx context.Context, users []*User) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    stmt, err := tx.PrepareContext(ctx,
        `INSERT INTO users (name, email) VALUES ($1, $2)`,
    )
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, user := range users {
        if _, err := stmt.ExecContext(ctx, user.Name, user.Email); err != nil {
            return err
        }
    }

    return tx.Commit()
}

// Better: Use COPY or batch insert for PostgreSQL
func (r *UserRepository) BulkInsert(ctx context.Context, users []*User) error {
    // Build multi-value INSERT
    valueStrings := make([]string, 0, len(users))
    valueArgs := make([]interface{}, 0, len(users)*2)

    for i, user := range users {
        valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
        valueArgs = append(valueArgs, user.Name, user.Email)
    }

    query := fmt.Sprintf(
        "INSERT INTO users (name, email) VALUES %s",
        strings.Join(valueStrings, ","),
    )

    _, err := r.db.ExecContext(ctx, query, valueArgs...)
    return err
}
```

### 2. Connection Pool Monitoring

```go
func MonitorDBStats(db *sql.DB, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for range ticker.C {
        stats := db.Stats()
        log.Printf("DB Stats: Open=%d InUse=%d Idle=%d WaitCount=%d WaitDuration=%s",
            stats.OpenConnections,
            stats.InUse,
            stats.Idle,
            stats.WaitCount,
            stats.WaitDuration,
        )

        if stats.WaitCount > 100 {
            log.Warn("High connection wait count - consider increasing MaxOpenConns")
        }
    }
}
```

### 3. Query Optimization

```go
// BAD: N+1 query problem
func (r *UserRepository) GetUsersWithProfiles(ctx context.Context) ([]*User, error) {
    users, err := r.FindAll(ctx)
    if err != nil {
        return nil, err
    }

    for _, user := range users {
        // N additional queries!
        profile, _ := r.GetProfile(ctx, user.ID)
        user.Profile = profile
    }

    return users, nil
}

// GOOD: Single JOIN query
func (r *UserRepository) GetUsersWithProfiles(ctx context.Context) ([]*User, error) {
    query := `
        SELECT
            u.id, u.name, u.email,
            p.bio, p.avatar
        FROM users u
        LEFT JOIN profiles p ON p.user_id = u.id
    `

    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    users := make([]*User, 0)
    for rows.Next() {
        user := &User{Profile: &Profile{}}
        err := rows.Scan(
            &user.ID, &user.Name, &user.Email,
            &user.Profile.Bio, &user.Profile.Avatar,
        )
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}
```

## Testing Strategies

### 1. In-Memory Database

```go
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()

    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatal(err)
    }

    // Run migrations
    _, err = db.Exec(`
        CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            email TEXT UNIQUE NOT NULL
        )
    `)
    if err != nil {
        t.Fatal(err)
    }

    t.Cleanup(func() { db.Close() })

    return db
}

func TestUserRepository(t *testing.T) {
    db := setupTestDB(t)
    repo := NewUserRepository(db)

    user := &User{Name: "John", Email: "john@test.com"}
    err := repo.Create(context.Background(), user)
    require.NoError(t, err)
    require.NotZero(t, user.ID)
}
```

### 2. Docker Test Containers

```go
func setupPostgresTest(t *testing.T) *sql.DB {
    // Use testcontainers-go
    ctx := context.Background()
    req := testcontainers.ContainerRequest{
        Image:        "postgres:15",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_PASSWORD": "test",
            "POSTGRES_DB":       "testdb",
        },
    }

    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    require.NoError(t, err)

    t.Cleanup(func() { container.Terminate(ctx) })

    // Get connection details and connect
    // ...
}
```

### 3. Transaction Rollback Testing

```go
func TestWithTransaction(t *testing.T) {
    db := setupTestDB(t)

    tx, _ := db.Begin()
    defer tx.Rollback() // Test data cleaned up automatically

    _, err := tx.Exec(`INSERT INTO users (name, email) VALUES (?, ?)`, "Test", "test@test.com")
    require.NoError(t, err)

    // Verify in transaction
    var count int
    tx.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
    assert.Equal(t, 1, count)

    // After rollback, data is gone (in real database)
}
```

## Common Pitfalls

### 1. Not Closing Rows

```go
// BAD: Resource leak
rows, _ := db.Query("SELECT * FROM users")
for rows.Next() {
    // ...
}
// rows never closed!

// GOOD
rows, _ := db.Query("SELECT * FROM users")
defer rows.Close() // ALWAYS defer Close()
for rows.Next() {
    // ...
}
```

### 2. SQL Injection

```go
// VULNERABLE
func (r *UserRepository) FindByEmail(email string) (*User, error) {
    query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
    // Attacker can pass: ' OR '1'='1
    rows, _ := r.db.Query(query)
    // ...
}

// SAFE: Use parameterized queries
func (r *UserRepository) FindByEmail(email string) (*User, error) {
    query := "SELECT * FROM users WHERE email = $1"
    rows, _ := r.db.Query(query, email) // Properly escaped
    // ...
}
```

### 3. Ignoring Context

```go
// BAD: Can't cancel long-running queries
func (r *UserRepository) FindAll() ([]*User, error) {
    rows, _ := r.db.Query("SELECT * FROM users")
    // ...
}

// GOOD: Respects cancellation and timeouts
func (r *UserRepository) FindAll(ctx context.Context) ([]*User, error) {
    rows, _ := r.db.QueryContext(ctx, "SELECT * FROM users")
    // ...
}
```

### 4. Transaction Leaks

```go
// BAD: Transaction never committed or rolled back
func (r *UserRepository) Update(user *User) error {
    tx, _ := r.db.Begin()
    tx.Exec("UPDATE users SET name = ? WHERE id = ?", user.Name, user.ID)
    // Oops, forgot to commit!
    return nil
}

// GOOD
func (r *UserRepository) Update(user *User) error {
    tx, _ := r.db.Begin()
    defer tx.Rollback() // Safe to call even after Commit

    _, err := tx.Exec("UPDATE users SET name = ? WHERE id = ?", user.Name, user.ID)
    if err != nil {
        return err
    }

    return tx.Commit()
}
```

## Production Checklist

- [ ] Connection pool properly configured
- [ ] All queries use `Context` for cancellation
- [ ] `rows.Close()` always called (use defer)
- [ ] `rows.Err()` checked after iteration
- [ ] Transactions use defer rollback pattern
- [ ] Prepared statements for frequently executed queries
- [ ] NULL values handled with sql.Null* types
- [ ] Database errors properly wrapped and logged
- [ ] No SQL injection vulnerabilities (use parameterized queries)
- [ ] Database driver imported with blank identifier
- [ ] Connection monitoring and metrics in place
- [ ] Migrations managed with tool (goose, migrate, atlas)

## Further Reading

- **database/sql:** https://pkg.go.dev/database/sql
- **Go database/sql tutorial:** http://go-database-sql.org/
- **sqlx:** https://github.com/jmoiron/sqlx (extensions to database/sql)
- **Squirrel:** https://github.com/Masterminds/squirrel (query builder)
- **PostgreSQL driver:** https://github.com/lib/pq
- **Migration tools:** https://github.com/pressly/goose, https://github.com/golang-migrate/migrate
