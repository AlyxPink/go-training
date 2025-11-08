package main

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// User represents a user in the system
type User struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}

// TableName returns the table name for User
func (u *User) TableName() string {
	return "users"
}

// Post represents a blog post
type Post struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

// TableName returns the table name for Post
func (p *Post) TableName() string {
	return "posts"
}

// QueryBuilder provides a fluent API for building SQL queries
type QueryBuilder struct {
	table        string
	selectCols   []string
	whereClauses []string
	whereArgs    []interface{}
	orderBy      string
	limit        int
	offset       int
}

// NewQueryBuilder creates a new query builder for the given table
func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{
		table:      table,
		selectCols: []string{"*"},
	}
}

// Select specifies which columns to select
func (qb *QueryBuilder) Select(cols ...string) *QueryBuilder {
	qb.selectCols = cols
	return qb
}

// Where adds a WHERE clause
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	qb.whereClauses = append(qb.whereClauses, condition)
	qb.whereArgs = append(qb.whereArgs, args...)
	return qb
}

// OrderBy adds ORDER BY clause
func (qb *QueryBuilder) OrderBy(order string) *QueryBuilder {
	qb.orderBy = order
	return qb
}

// Limit adds LIMIT clause
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

// Offset adds OFFSET clause
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

// ToSQL generates the SQL query and arguments
func (qb *QueryBuilder) ToSQL() (string, []interface{}) {
	// Build SELECT clause
	query := fmt.Sprintf("SELECT %s FROM %s",
		strings.Join(qb.selectCols, ", "),
		qb.table,
	)

	// Add WHERE clause if present
	if len(qb.whereClauses) > 0 {
		query += " WHERE " + strings.Join(qb.whereClauses, " AND ")
	}

	// Add ORDER BY clause if present
	if qb.orderBy != "" {
		query += " ORDER BY " + qb.orderBy
	}

	// Add LIMIT clause if present
	if qb.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.limit)
	}

	// Add OFFSET clause if present
	if qb.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", qb.offset)
	}

	return query, qb.whereArgs
}

// ScanStruct scans a database row into a struct using reflection
func ScanStruct(rows *sql.Rows, dest interface{}) error {
	// Get the value and type of the destination
	destValue := reflect.ValueOf(dest).Elem()
	destType := destValue.Type()

	// Build a map of db tag names to struct field values
	fieldMap := make(map[string]reflect.Value)
	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			fieldMap[dbTag] = destValue.Field(i)
		}
	}

	// Get column names from the result set
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// Prepare scan destinations
	values := make([]interface{}, len(columns))
	for i, col := range columns {
		if field, ok := fieldMap[col]; ok {
			// Use the address of the field as scan destination
			values[i] = field.Addr().Interface()
		} else {
			// Column not in struct, use a dummy destination
			var dummy interface{}
			values[i] = &dummy
		}
	}

	// Scan the row into the prepared destinations
	return rows.Scan(values...)
}

// GetStructFields extracts field information from a struct
func GetStructFields(v interface{}) (fields []string, values []interface{}) {
	// Get the value and type
	val := reflect.ValueOf(v).Elem()
	typ := val.Type()

	// Iterate through struct fields
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		dbTag := field.Tag.Get("db")

		if dbTag != "" {
			fields = append(fields, dbTag)
			values = append(values, val.Field(i).Interface())
		}
	}

	return fields, values
}

// Migration represents a database migration
type Migration struct {
	Version int
	Up      string
	Down    string
}

// Migrator handles database migrations
type Migrator struct {
	db         *sql.DB
	migrations []Migration
}

// NewMigrator creates a new migrator
func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{
		db: db,
		migrations: []Migration{
			{
				Version: 1,
				Up: `
					CREATE TABLE users (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						name TEXT NOT NULL,
						email TEXT UNIQUE NOT NULL,
						created_at DATETIME DEFAULT CURRENT_TIMESTAMP
					)
				`,
				Down: `DROP TABLE users`,
			},
			{
				Version: 2,
				Up: `
					CREATE TABLE posts (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						user_id INTEGER NOT NULL,
						title TEXT NOT NULL,
						content TEXT,
						created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
						FOREIGN KEY (user_id) REFERENCES users(id)
					)
				`,
				Down: `DROP TABLE posts`,
			},
		},
	}
}

// Init initializes the migration tracking table
func (m *Migrator) Init() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

// CurrentVersion returns the current migration version
func (m *Migrator) CurrentVersion() (int, error) {
	var version int
	err := m.db.QueryRow(`
		SELECT COALESCE(MAX(version), 0) FROM schema_migrations
	`).Scan(&version)

	if err != nil {
		return 0, err
	}

	return version, nil
}

// MigrateUp runs all pending migrations
func (m *Migrator) MigrateUp() error {
	// Get current version
	current, err := m.CurrentVersion()
	if err != nil {
		return err
	}

	// Apply migrations that haven't been applied yet
	for _, migration := range m.migrations {
		if migration.Version > current {
			if err := m.applyMigration(migration); err != nil {
				return err
			}
		}
	}

	return nil
}

// applyMigration applies a single migration
func (m *Migrator) applyMigration(migration Migration) error {
	// Start transaction
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration
	if _, err := tx.Exec(migration.Up); err != nil {
		return fmt.Errorf("migration %d failed: %w", migration.Version, err)
	}

	// Record version in schema_migrations
	_, err = tx.Exec(
		`INSERT INTO schema_migrations (version) VALUES (?)`,
		migration.Version,
	)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

// UserRepository handles User persistence
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user
func (r *UserRepository) Create(user *User) error {
	err := r.db.QueryRow(
		`INSERT INTO users (name, email) VALUES (?, ?) RETURNING id`,
		user.Name, user.Email,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	// Fetch the created_at timestamp
	return r.db.QueryRow(
		`SELECT created_at FROM users WHERE id = ?`,
		user.ID,
	).Scan(&user.CreatedAt)
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id int) (*User, error) {
	var user User
	err := r.db.QueryRow(
		`SELECT id, name, email, created_at FROM users WHERE id = ?`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.QueryRow(
		`SELECT id, name, email, created_at FROM users WHERE email = ?`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *User) error {
	_, err := r.db.Exec(
		`UPDATE users SET name = ?, email = ? WHERE id = ?`,
		user.Name, user.Email, user.ID,
	)
	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}

// FindAll returns all users
func (r *UserRepository) FindAll() ([]*User, error) {
	rows, err := r.db.Query(`SELECT id, name, email, created_at FROM users`)
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
}

// PostRepository handles Post persistence
type PostRepository struct {
	db *sql.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create inserts a new post
func (r *PostRepository) Create(post *Post) error {
	err := r.db.QueryRow(
		`INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?) RETURNING id`,
		post.UserID, post.Title, post.Content,
	).Scan(&post.ID)

	if err != nil {
		return err
	}

	// Fetch the created_at timestamp
	return r.db.QueryRow(
		`SELECT created_at FROM posts WHERE id = ?`,
		post.ID,
	).Scan(&post.CreatedAt)
}

// FindByUserID finds all posts for a user
func (r *PostRepository) FindByUserID(userID int) ([]*Post, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, title, content, created_at FROM posts WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post Post
		if err := ScanStruct(rows, &post); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func main() {
	// Example usage
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Run migrations
	migrator := NewMigrator(db)
	if err := migrator.Init(); err != nil {
		panic(err)
	}
	if err := migrator.MigrateUp(); err != nil {
		panic(err)
	}

	// Create repositories
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)

	// Create a user
	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	if err := userRepo.Create(user); err != nil {
		panic(err)
	}
	fmt.Printf("Created user with ID: %d\n", user.ID)

	// Find user by ID
	foundUser, err := userRepo.FindByID(user.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found user: %s (%s)\n", foundUser.Name, foundUser.Email)

	// Create a post
	post := &Post{
		UserID:  user.ID,
		Title:   "My First Post",
		Content: "This is my first post!",
	}
	if err := postRepo.Create(post); err != nil {
		panic(err)
	}
	fmt.Printf("Created post with ID: %d\n", post.ID)

	// Query builder example
	qb := NewQueryBuilder("users").
		Select("id", "name", "email").
		Where("name = ?", "John Doe").
		OrderBy("created_at DESC").
		Limit(10)

	query, args := qb.ToSQL()
	fmt.Printf("Generated SQL: %s\n", query)
	fmt.Printf("Arguments: %v\n", args)
}
