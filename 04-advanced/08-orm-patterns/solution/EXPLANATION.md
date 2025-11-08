# ORM Patterns Solution - Deep Dive

## Overview

This solution demonstrates building a lightweight ORM (Object-Relational Mapper) in Go, including query builders, struct mapping, relationship handling, and a simple migration system. Understanding ORM internals helps you use existing ORMs effectively and know when to build custom solutions.

## Architecture

### 1. Query Builder

```go
type QueryBuilder struct {
    table      string
    selectCols []string
    whereClauses []string
    whereArgs  []interface{}
    orderBy    string
    limit      int
    offset     int
}

func (qb *QueryBuilder) Select(cols ...string) *QueryBuilder {
    qb.selectCols = cols
    return qb
}

func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
    qb.whereClauses = append(qb.whereClauses, condition)
    qb.whereArgs = append(qb.whereArgs, args...)
    return qb
}

func (qb *QueryBuilder) ToSQL() (string, []interface{}) {
    query := fmt.Sprintf("SELECT %s FROM %s",
        strings.Join(qb.selectCols, ", "),
        qb.table,
    )

    if len(qb.whereClauses) > 0 {
        query += " WHERE " + strings.Join(qb.whereClauses, " AND ")
    }

    if qb.orderBy != "" {
        query += " ORDER BY " + qb.orderBy
    }

    if qb.limit > 0 {
        query += fmt.Sprintf(" LIMIT %d", qb.limit)
    }

    return query, qb.whereArgs
}
```

**Benefits:**
- Fluent API for building queries
- Type-safe query construction
- Prevents SQL injection
- Easy to test and debug

### 2. Struct Mapping with Reflection

```go
type User struct {
    ID        int       `db:"id"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    CreatedAt time.Time `db:"created_at"`
}

func ScanStruct(rows *sql.Rows, dest interface{}) error {
    destValue := reflect.ValueOf(dest).Elem()
    destType := destValue.Type()

    // Build map of field names to struct fields
    fieldMap := make(map[string]reflect.Value)
    for i := 0; i < destType.NumField(); i++ {
        field := destType.Field(i)
        dbTag := field.Tag.Get("db")
        if dbTag != "" {
            fieldMap[dbTag] = destValue.Field(i)
        }
    }

    // Get column names from rows
    columns, _ := rows.Columns()

    // Prepare scan destinations
    values := make([]interface{}, len(columns))
    for i, col := range columns {
        if field, ok := fieldMap[col]; ok {
            values[i] = field.Addr().Interface()
        } else {
            var dummy interface{}
            values[i] = &dummy // Column not in struct
        }
    }

    return rows.Scan(values...)
}
```

**Trade-offs:**
- **Pro:** Generic, works with any struct
- **Pro:** Automatic mapping reduces boilerplate
- **Con:** Reflection overhead (slower than manual scanning)
- **Con:** Runtime errors vs compile-time safety

### 3. Active Record Pattern

```go
type Model interface {
    TableName() string
    Save(db *sql.DB) error
    Delete(db *sql.DB) error
}

type User struct {
    ID    int    `db:"id"`
    Name  string `db:"name"`
    Email string `db:"email"`
}

func (u *User) TableName() string {
    return "users"
}

func (u *User) Save(db *sql.DB) error {
    if u.ID == 0 {
        // INSERT
        return db.QueryRow(
            `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
            u.Name, u.Email,
        ).Scan(&u.ID)
    } else {
        // UPDATE
        _, err := db.Exec(
            `UPDATE users SET name = $1, email = $2 WHERE id = $3`,
            u.Name, u.Email, u.ID,
        )
        return err
    }
}

func (u *User) Delete(db *sql.DB) error {
    _, err := db.Exec(`DELETE FROM users WHERE id = $1`, u.ID)
    return err
}

// Usage:
user := &User{Name: "John", Email: "john@example.com"}
user.Save(db) // INSERT
user.Name = "Jane"
user.Save(db) // UPDATE
user.Delete(db)
```

**When to use:**
- Simple CRUD operations
- Rapid prototyping
- Small to medium applications

**Limitations:**
- Tight coupling to database
- Hard to test (models depend on DB)
- Difficult to optimize queries
- Doesn't scale to complex queries

## Key Patterns

### Pattern 1: Data Mapper

```go
type UserMapper struct {
    db *sql.DB
}

func (m *UserMapper) Find(id int) (*User, error) {
    var user User
    err := m.db.QueryRow(
        `SELECT id, name, email FROM users WHERE id = $1`,
        id,
    ).Scan(&user.ID, &user.Name, &user.Email)
    return &user, err
}

func (m *UserMapper) Insert(user *User) error {
    return m.db.QueryRow(
        `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
        user.Name, user.Email,
    ).Scan(&user.ID)
}

func (m *UserMapper) Update(user *User) error {
    _, err := m.db.Exec(
        `UPDATE users SET name = $1, email = $2 WHERE id = $3`,
        user.Name, user.Email, user.ID,
    )
    return err
}

func (m *UserMapper) Delete(id int) error {
    _, err := m.db.Exec(`DELETE FROM users WHERE id = $1`, id)
    return err
}
```

**Advantages over Active Record:**
- Separation of concerns (domain vs persistence)
- Easier to test (can mock mapper)
- More flexible query optimization
- Better for complex domains

### Pattern 2: Relationship Handling

```go
type User struct {
    ID       int
    Name     string
    Posts    []*Post  // One-to-many
    Profile  *Profile // One-to-one
}

type Post struct {
    ID      int
    UserID  int
    Title   string
    Content string
}

// Eager loading
func (m *UserMapper) FindWithPosts(id int) (*User, error) {
    user, err := m.Find(id)
    if err != nil {
        return nil, err
    }

    // Load posts in one query
    rows, err := m.db.Query(
        `SELECT id, title, content FROM posts WHERE user_id = $1`,
        user.ID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        post := &Post{UserID: user.ID}
        rows.Scan(&post.ID, &post.Title, &post.Content)
        user.Posts = append(user.Posts, post)
    }

    return user, nil
}

// Lazy loading
type LazyPosts struct {
    userID int
    mapper *PostMapper
    loaded bool
    posts  []*Post
}

func (lp *LazyPosts) Get() ([]*Post, error) {
    if !lp.loaded {
        posts, err := lp.mapper.FindByUserID(lp.userID)
        if err != nil {
            return nil, err
        }
        lp.posts = posts
        lp.loaded = true
    }
    return lp.posts, nil
}
```

### Pattern 3: Unit of Work

```go
type UnitOfWork struct {
    db      *sql.DB
    tx      *sql.Tx
    inserts []interface{}
    updates []interface{}
    deletes []interface{}
}

func NewUnitOfWork(db *sql.DB) *UnitOfWork {
    return &UnitOfWork{db: db}
}

func (uow *UnitOfWork) RegisterNew(entity interface{}) {
    uow.inserts = append(uow.inserts, entity)
}

func (uow *UnitOfWork) RegisterDirty(entity interface{}) {
    uow.updates = append(uow.updates, entity)
}

func (uow *UnitOfWork) RegisterDeleted(entity interface{}) {
    uow.deletes = append(uow.deletes, entity)
}

func (uow *UnitOfWork) Commit() error {
    tx, err := uow.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Execute all INSERTs
    for _, entity := range uow.inserts {
        if err := uow.insert(tx, entity); err != nil {
            return err
        }
    }

    // Execute all UPDATEs
    for _, entity := range uow.updates {
        if err := uow.update(tx, entity); err != nil {
            return err
        }
    }

    // Execute all DELETEs
    for _, entity := range uow.deletes {
        if err := uow.delete(tx, entity); err != nil {
            return err
        }
    }

    return tx.Commit()
}

// Usage:
uow := NewUnitOfWork(db)
uow.RegisterNew(user1)
uow.RegisterNew(user2)
uow.RegisterDirty(user3)
uow.RegisterDeleted(user4)
uow.Commit() // All in one transaction
```

### Pattern 4: Migration System

```go
type Migration struct {
    Version int
    Up      string
    Down    string
}

var migrations = []Migration{
    {
        Version: 1,
        Up: `
            CREATE TABLE users (
                id SERIAL PRIMARY KEY,
                name TEXT NOT NULL,
                email TEXT UNIQUE NOT NULL
            )
        `,
        Down: `DROP TABLE users`,
    },
    {
        Version: 2,
        Up: `
            CREATE TABLE posts (
                id SERIAL PRIMARY KEY,
                user_id INTEGER REFERENCES users(id),
                title TEXT NOT NULL,
                content TEXT
            )
        `,
        Down: `DROP TABLE posts`,
    },
}

type Migrator struct {
    db *sql.DB
}

func (m *Migrator) Init() error {
    _, err := m.db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version INTEGER PRIMARY KEY,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    return err
}

func (m *Migrator) CurrentVersion() (int, error) {
    var version int
    err := m.db.QueryRow(
        `SELECT COALESCE(MAX(version), 0) FROM schema_migrations`,
    ).Scan(&version)
    return version, err
}

func (m *Migrator) MigrateUp() error {
    current, err := m.CurrentVersion()
    if err != nil {
        return err
    }

    for _, migration := range migrations {
        if migration.Version > current {
            if err := m.applyMigration(migration); err != nil {
                return err
            }
        }
    }

    return nil
}

func (m *Migrator) applyMigration(migration Migration) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Execute migration
    if _, err := tx.Exec(migration.Up); err != nil {
        return fmt.Errorf("migration %d failed: %w", migration.Version, err)
    }

    // Record version
    _, err = tx.Exec(
        `INSERT INTO schema_migrations (version) VALUES ($1)`,
        migration.Version,
    )
    if err != nil {
        return err
    }

    return tx.Commit()
}
```

## Performance Considerations

### 1. N+1 Query Problem

```go
// BAD: N+1 queries
func GetUsersWithPosts() ([]*User, error) {
    users, _ := GetAllUsers() // 1 query

    for _, user := range users {
        posts, _ := GetPostsByUserID(user.ID) // N queries
        user.Posts = posts
    }

    return users, nil
}

// GOOD: JOIN or separate query
func GetUsersWithPosts() ([]*User, error) {
    query := `
        SELECT
            u.id, u.name, u.email,
            p.id, p.title, p.content
        FROM users u
        LEFT JOIN posts p ON p.user_id = u.id
        ORDER BY u.id
    `

    rows, _ := db.Query(query)
    defer rows.Close()

    userMap := make(map[int]*User)

    for rows.Next() {
        var userID, postID int
        var userName, userEmail, postTitle, postContent string

        rows.Scan(&userID, &userName, &userEmail, &postID, &postTitle, &postContent)

        user, ok := userMap[userID]
        if !ok {
            user = &User{ID: userID, Name: userName, Email: userEmail}
            userMap[userID] = user
        }

        if postID > 0 {
            user.Posts = append(user.Posts, &Post{
                ID:      postID,
                Title:   postTitle,
                Content: postContent,
            })
        }
    }

    users := make([]*User, 0, len(userMap))
    for _, user := range userMap {
        users = append(users, user)
    }

    return users, nil
}
```

### 2. Caching Strategies

```go
type CachedUserMapper struct {
    mapper *UserMapper
    cache  map[int]*User
    mu     sync.RWMutex
}

func (m *CachedUserMapper) Find(id int) (*User, error) {
    // Check cache
    m.mu.RLock()
    if user, ok := m.cache[id]; ok {
        m.mu.RUnlock()
        return user, nil
    }
    m.mu.RUnlock()

    // Load from database
    user, err := m.mapper.Find(id)
    if err != nil {
        return nil, err
    }

    // Store in cache
    m.mu.Lock()
    m.cache[id] = user
    m.mu.Unlock()

    return user, nil
}

func (m *CachedUserMapper) Update(user *User) error {
    if err := m.mapper.Update(user); err != nil {
        return err
    }

    // Invalidate cache
    m.mu.Lock()
    delete(m.cache, user.ID)
    m.mu.Unlock()

    return nil
}
```

### 3. Bulk Operations

```go
func (m *UserMapper) BulkInsert(users []*User) error {
    if len(users) == 0 {
        return nil
    }

    // Build multi-value INSERT
    placeholders := make([]string, len(users))
    values := make([]interface{}, len(users)*2)

    for i, user := range users {
        placeholders[i] = fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2)
        values[i*2] = user.Name
        values[i*2+1] = user.Email
    }

    query := fmt.Sprintf(
        "INSERT INTO users (name, email) VALUES %s",
        strings.Join(placeholders, ", "),
    )

    _, err := m.db.Exec(query, values...)
    return err
}
```

## Comparison with Popular ORMs

### GORM

```go
// GORM example
type User struct {
    gorm.Model
    Name  string
    Posts []Post `gorm:"foreignKey:UserID"`
}

// Usage
db.AutoMigrate(&User{})
db.Create(&User{Name: "John"})
db.Preload("Posts").First(&user, 1)
```

**Pros:**
- Feature-rich
- Active development
- Large community
- Comprehensive documentation

**Cons:**
- Heavy dependency
- Magic behavior can surprise
- Performance overhead
- Learning curve

### sqlx

```go
// sqlx is extension to database/sql, not full ORM
type User struct {
    ID    int    `db:"id"`
    Name  string `db:"name"`
    Email string `db:"email"`
}

db := sqlx.Connect("postgres", dsn)
users := []User{}
db.Select(&users, "SELECT * FROM users WHERE age > $1", 18)
```

**Pros:**
- Lightweight
- Similar to database/sql
- Named parameters
- Struct scanning

**Cons:**
- Not a full ORM
- Manual query writing
- No relationship handling

### ent

```go
// ent uses code generation
client, _ := ent.Open("postgres", dsn)
user, _ := client.User.
    Create().
    SetName("John").
    SetEmail("john@example.com").
    Save(ctx)
```

**Pros:**
- Type-safe
- Code generation
- Schema-first
- Good for complex schemas

**Cons:**
- Code generation step
- Opinionated structure
- Learning curve

## When to Build Custom ORM

### Build When:
- Simple CRUD, limited relationships
- Performance is critical (avoid reflection)
- Need full control over SQL
- Team prefers explicit over magic
- Small, focused domain

### Use Existing When:
- Complex relationships
- Rapid development needed
- Team familiar with ORM
- Standard patterns work
- Large application

## Common Pitfalls

### 1. Over-Abstraction

```go
// Don't hide too much
db.Model(&User{}).Where("age > ?", 18).Limit(10).Find(&users)

// Sometimes raw SQL is clearer
db.Query("SELECT * FROM users WHERE age > $1 LIMIT 10", 18)
```

### 2. Ignoring Indexes

```go
// Migration should include indexes
Up: `
    CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        email TEXT UNIQUE NOT NULL
    );
    CREATE INDEX idx_users_email ON users(email);
`
```

### 3. Lazy Loading in Loops

```go
// Creates N+1 problem
for _, user := range users {
    posts := user.GetPosts() // Lazy load
}

// Eager load instead
users := db.Preload("Posts").Find(&users)
```

## Production Checklist

- [ ] Migrations versioned and tested
- [ ] Indexes on frequently queried columns
- [ ] N+1 query problems identified and fixed
- [ ] Eager loading for associations
- [ ] Caching strategy for hot data
- [ ] Soft deletes if needed (deleted_at column)
- [ ] Timestamps (created_at, updated_at)
- [ ] Connection pooling configured
- [ ] Error handling for constraint violations
- [ ] Logging slow queries
- [ ] Benchmarks for critical paths

## Testing Concurrent Database Access

### SQLite WAL Mode for Better Concurrency

The test suite uses SQLite's WAL (Write-Ahead Logging) mode for better concurrent access:

```go
func setupTestDB(t *testing.T) *sql.DB {
    // Use WAL mode for better concurrent access to in-memory database
    db, err := sql.Open("sqlite3", "file::memory:?cache=shared&mode=rwc")

    // Enable WAL mode for better concurrency
    db.Exec("PRAGMA journal_mode=WAL")

    // Set busy timeout to handle locks
    db.Exec("PRAGMA busy_timeout=5000")

    return db
}
```

**Why WAL Mode?**
- Allows concurrent reads and writes
- Readers don't block writers
- Writers don't block readers
- Better performance under concurrent load

### Proper Goroutine Error Handling

The concurrent test uses buffered channels to properly collect errors from goroutines:

```go
func TestConcurrentUserCreation(t *testing.T) {
    type result struct {
        err error
        id  int
    }
    results := make(chan result, 10)

    for i := 0; i < 10; i++ {
        go func(id int) {
            err := repo.Create(user)
            results <- result{err: err, id: id}
        }(i)
    }

    // Collect and report errors properly
    var errors []error
    for i := 0; i < 10; i++ {
        res := <-results
        if res.err != nil {
            errors = append(errors, fmt.Errorf("User%d: %w", res.id, res.err))
        }
    }
}
```

**Best Practices:**
- Use buffered channels to prevent goroutine leaks
- Collect all results before assertions
- Report errors with context (which operation failed)
- Never use `t.Errorf` directly in goroutines (race condition)

## Further Reading

- **GORM:** https://gorm.io/
- **sqlx:** https://github.com/jmoiron/sqlx
- **ent:** https://entgo.io/
- **goose migrations:** https://github.com/pressly/goose
- **Design Patterns (Fowler):** Patterns of Enterprise Application Architecture
- **SQLite WAL Mode:** https://www.sqlite.org/wal.html
