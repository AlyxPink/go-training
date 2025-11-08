# ORM Patterns - Implementation Hints

## Architecture Overview

This exercise builds a lightweight ORM with several key components:

1. **Query Builder** - Fluent API for constructing SQL queries
2. **Struct Mapping** - Reflection-based mapping between structs and database rows
3. **Migrations** - Database schema versioning system
4. **Repositories** - Data access layer for entities

## Implementation Strategy

### 1. Query Builder Pattern

The query builder uses method chaining to construct SQL queries programmatically.

**Key Points:**
- Each method returns `*QueryBuilder` to enable chaining
- Store query components in struct fields
- Build final SQL in `ToSQL()` method
- Use placeholders (`?`) to prevent SQL injection

**Example Flow:**
```go
qb := NewQueryBuilder("users")
qb.Select("id", "name")         // Store ["id", "name"]
qb.Where("age > ?", 18)         // Store condition and arg
qb.OrderBy("created_at DESC")   // Store order clause
query, args := qb.ToSQL()       // Build final SQL
```

**ToSQL() Construction:**
```
SELECT <columns> FROM <table>
[WHERE <conditions>]
[ORDER BY <order>]
[LIMIT <limit>]
[OFFSET <offset>]
```

### 2. Reflection for Struct Mapping

Use Go's `reflect` package to dynamically map between structs and database rows.

**For ScanStruct:**
1. Get reflect.Value and reflect.Type of destination struct
2. Build a map of db tag names to reflect.Value fields
3. Get column names from sql.Rows
4. Create []interface{} slice matching column order
5. For each column, use field.Addr().Interface() as scan destination
6. Call rows.Scan() with the prepared destinations

**For GetStructFields:**
1. Use reflect.ValueOf and reflect.TypeOf
2. Iterate through struct fields with NumField()
3. Extract db tag with field.Tag.Get("db")
4. Collect field names and values using field.Interface()

**Key Reflection Methods:**
- `reflect.ValueOf(v)` - Get value wrapper
- `v.Elem()` - Dereference pointer
- `v.Type()` - Get type information
- `t.NumField()` - Number of struct fields
- `t.Field(i)` - Get field by index
- `field.Tag.Get("db")` - Extract struct tag
- `field.Addr().Interface()` - Get pointer to field for scanning

### 3. Migration System

Migrations manage database schema evolution across versions.

**Schema Migrations Table:**
```sql
CREATE TABLE schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
```

**Migration Flow:**
1. Init() - Create schema_migrations table
2. CurrentVersion() - Query MAX(version) from schema_migrations
3. MigrateUp() - Apply migrations with version > current
4. applyMigration() - Execute migration in transaction

**Transaction Safety:**
```go
tx, err := db.Begin()
defer tx.Rollback()  // Rollback if commit not reached

// Execute migration SQL
tx.Exec(migration.Up)

// Record version
tx.Exec("INSERT INTO schema_migrations ...")

tx.Commit()  // Commit on success
```

### 4. Repository Pattern

Repositories encapsulate data access logic for each entity type.

**CRUD Operations:**

**Create:**
- INSERT INTO users (...) VALUES (...)
- SQLite: Use AUTOINCREMENT, get ID with RETURNING clause or last_insert_rowid()
- Scan returned ID into struct

**Read:**
- FindByID: SELECT ... WHERE id = ?
- FindAll: SELECT ... (iterate rows)
- Use ScanStruct() for row mapping

**Update:**
- UPDATE users SET ... WHERE id = ?
- Use Exec() since no rows returned

**Delete:**
- DELETE FROM users WHERE id = ?
- Use Exec()

## Common Patterns

### Method Chaining
```go
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
    qb.whereClauses = append(qb.whereClauses, condition)
    qb.whereArgs = append(qb.whereArgs, args...)
    return qb  // Return self for chaining
}
```

### Handling Optional Clauses
```go
if len(qb.whereClauses) > 0 {
    query += " WHERE " + strings.Join(qb.whereClauses, " AND ")
}

if qb.orderBy != "" {
    query += " ORDER BY " + qb.orderBy
}

if qb.limit > 0 {
    query += fmt.Sprintf(" LIMIT %d", qb.limit)
}
```

### Variadic Arguments
```go
func (qb *QueryBuilder) Select(cols ...string) *QueryBuilder {
    qb.selectCols = cols  // cols is []string
    return qb
}
```

### Error Handling in Queries
```go
func (r *UserRepository) FindByID(id int) (*User, error) {
    var user User
    err := r.db.QueryRow(
        "SELECT id, name, email, created_at FROM users WHERE id = ?",
        id,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

    if err != nil {
        return nil, err  // Returns sql.ErrNoRows if not found
    }

    return &user, nil
}
```

### Iterating Rows
```go
rows, err := db.Query("SELECT ...")
if err != nil {
    return nil, err
}
defer rows.Close()

var users []*User
for rows.Next() {
    var user User
    if err := ScanStruct(rows, &user); err != nil {
        return nil, err
    }
    users = append(users, &user)
}

if err := rows.Err(); err != nil {
    return nil, err
}

return users, nil
```

## SQLite Specifics

### Auto-increment Primary Keys
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ...
)
```

### Getting Last Insert ID (Method 1 - RETURNING)
```go
err := db.QueryRow(
    "INSERT INTO users (name, email) VALUES (?, ?) RETURNING id",
    user.Name, user.Email,
).Scan(&user.ID)
```

### Getting Last Insert ID (Method 2 - LastInsertId)
```go
result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
id, err := result.LastInsertId()
user.ID = int(id)
```

### Foreign Keys
```sql
CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
)
```

### Default Timestamps
```sql
created_at DATETIME DEFAULT CURRENT_TIMESTAMP
```

## Testing Considerations

### In-Memory Database
```go
db, err := sql.Open("sqlite3", ":memory:")
```

### Setup Helper
```go
func setupTestDB(t *testing.T) *sql.DB {
    db := // open database
    migrator := NewMigrator(db)
    migrator.Init()
    migrator.MigrateUp()
    return db
}
```

### Test Data Isolation
Each test should use a fresh database or clean up after itself.

## Performance Tips

1. **Use Prepared Statements** - For repeated queries
2. **Batch Operations** - INSERT multiple rows at once
3. **Connection Pooling** - Configure db.SetMaxOpenConns()
4. **Eager Loading** - Avoid N+1 queries with JOINs
5. **Indexes** - Add indexes to frequently queried columns

## Common Pitfalls

1. **Forgetting defer rows.Close()** - Leaks connections
2. **Not checking rows.Err()** - May miss iteration errors
3. **SQL Injection** - Always use placeholders, never string concatenation
4. **Transaction Handling** - Always defer Rollback() before Commit()
5. **Pointer Receivers** - Use &struct when scanning into structs

## Extension Ideas

Once basic implementation works, consider:

1. **Caching Layer** - Cache frequently accessed entities
2. **Soft Deletes** - Add deleted_at column instead of hard delete
3. **Hooks** - Before/After Create/Update/Delete callbacks
4. **Validation** - Validate struct fields before insert/update
5. **Optimistic Locking** - Version field to prevent concurrent updates
6. **Query Logging** - Log all SQL queries for debugging
7. **Relationship Loading** - Eager/lazy loading for has-many relationships

## Debugging Tips

1. **Print Generated SQL:**
   ```go
   query, args := qb.ToSQL()
   fmt.Printf("SQL: %s\nArgs: %v\n", query, args)
   ```

2. **Check Row Scan Errors:**
   ```go
   if err := rows.Scan(...); err != nil {
       fmt.Printf("Scan error: %v\n", err)
   }
   ```

3. **Verify Table Schema:**
   ```go
   rows, _ := db.Query("SELECT sql FROM sqlite_master WHERE type='table' AND name='users'")
   ```

4. **Enable SQLite Logging:**
   ```go
   sql.Register("sqlite3_with_trace", &sqlite3.SQLiteDriver{
       ConnectHook: func(conn *sqlite3.SQLiteConn) error {
           conn.RegisterFunc("trace", func(sql string) {
               fmt.Println("SQL:", sql)
           }, true)
           return nil
       },
   })
   ```

## Solution Structure

Your implementation should:
1. ✅ Build valid SQL with proper syntax
2. ✅ Use parameterized queries (prevent SQL injection)
3. ✅ Handle errors appropriately
4. ✅ Support method chaining
5. ✅ Use reflection correctly for struct mapping
6. ✅ Apply migrations in transactions
7. ✅ Track migration versions
8. ✅ Implement full CRUD operations
9. ✅ Handle timestamps correctly
10. ✅ Pass all tests

Good luck!
