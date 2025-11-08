package main

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Use WAL mode for better concurrent access to in-memory database
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared&mode=rwc")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Enable WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		t.Fatalf("Failed to enable WAL mode: %v", err)
	}

	// Set busy timeout to handle locks
	if _, err := db.Exec("PRAGMA busy_timeout=5000"); err != nil {
		t.Fatalf("Failed to set busy timeout: %v", err)
	}

	migrator := NewMigrator(db)
	if err := migrator.Init(); err != nil {
		t.Fatalf("Failed to initialize migrator: %v", err)
	}
	if err := migrator.MigrateUp(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func TestQueryBuilder_Select(t *testing.T) {
	qb := NewQueryBuilder("users").Select("id", "name", "email")

	if len(qb.selectCols) != 3 {
		t.Errorf("Expected 3 select columns, got %d", len(qb.selectCols))
	}

	expected := []string{"id", "name", "email"}
	for i, col := range qb.selectCols {
		if col != expected[i] {
			t.Errorf("Expected column %s, got %s", expected[i], col)
		}
	}
}

func TestQueryBuilder_Where(t *testing.T) {
	qb := NewQueryBuilder("users").
		Where("name = ?", "John").
		Where("age > ?", 18)

	if len(qb.whereClauses) != 2 {
		t.Errorf("Expected 2 where clauses, got %d", len(qb.whereClauses))
	}

	if len(qb.whereArgs) != 2 {
		t.Errorf("Expected 2 where arguments, got %d", len(qb.whereArgs))
	}
}

func TestQueryBuilder_OrderBy(t *testing.T) {
	qb := NewQueryBuilder("users").OrderBy("created_at DESC")

	if qb.orderBy != "created_at DESC" {
		t.Errorf("Expected order by 'created_at DESC', got '%s'", qb.orderBy)
	}
}

func TestQueryBuilder_Limit(t *testing.T) {
	qb := NewQueryBuilder("users").Limit(10)

	if qb.limit != 10 {
		t.Errorf("Expected limit 10, got %d", qb.limit)
	}
}

func TestQueryBuilder_Offset(t *testing.T) {
	qb := NewQueryBuilder("users").Offset(20)

	if qb.offset != 20 {
		t.Errorf("Expected offset 20, got %d", qb.offset)
	}
}

func TestQueryBuilder_ToSQL_Simple(t *testing.T) {
	qb := NewQueryBuilder("users")
	query, args := qb.ToSQL()

	expected := "SELECT * FROM users"
	if query != expected {
		t.Errorf("Expected query '%s', got '%s'", expected, query)
	}

	if len(args) != 0 {
		t.Errorf("Expected no arguments, got %d", len(args))
	}
}

func TestQueryBuilder_ToSQL_WithSelect(t *testing.T) {
	qb := NewQueryBuilder("users").Select("id", "name")
	query, _ := qb.ToSQL()

	expected := "SELECT id, name FROM users"
	if query != expected {
		t.Errorf("Expected query '%s', got '%s'", expected, query)
	}
}

func TestQueryBuilder_ToSQL_WithWhere(t *testing.T) {
	qb := NewQueryBuilder("users").
		Select("*").
		Where("name = ?", "John")
	query, args := qb.ToSQL()

	expected := "SELECT * FROM users WHERE name = ?"
	if query != expected {
		t.Errorf("Expected query '%s', got '%s'", expected, query)
	}

	if len(args) != 1 || args[0] != "John" {
		t.Errorf("Expected args ['John'], got %v", args)
	}
}

func TestQueryBuilder_ToSQL_MultipleWhere(t *testing.T) {
	qb := NewQueryBuilder("users").
		Where("name = ?", "John").
		Where("age > ?", 18)
	query, args := qb.ToSQL()

	expected := "SELECT * FROM users WHERE name = ? AND age > ?"
	if query != expected {
		t.Errorf("Expected query '%s', got '%s'", expected, query)
	}

	if len(args) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(args))
	}
}

func TestQueryBuilder_ToSQL_Complete(t *testing.T) {
	qb := NewQueryBuilder("users").
		Select("id", "name", "email").
		Where("age > ?", 18).
		OrderBy("created_at DESC").
		Limit(10).
		Offset(5)

	query, args := qb.ToSQL()

	expected := "SELECT id, name, email FROM users WHERE age > ? ORDER BY created_at DESC LIMIT 10 OFFSET 5"
	if query != expected {
		t.Errorf("Expected query '%s', got '%s'", expected, query)
	}

	if len(args) != 1 || args[0] != 18 {
		t.Errorf("Expected args [18], got %v", args)
	}
}

func TestQueryBuilder_Chaining(t *testing.T) {
	qb := NewQueryBuilder("users")
	result := qb.Select("id").Where("name = ?", "test")

	if result != qb {
		t.Error("Expected method chaining to return the same instance")
	}
}

func TestScanStruct(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec(`INSERT INTO users (name, email) VALUES (?, ?)`, "Alice", "alice@example.com")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	rows, err := db.Query("SELECT id, name, email, created_at FROM users WHERE name = ?", "Alice")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Fatal("Expected at least one row")
	}

	var user User
	if err := ScanStruct(rows, &user); err != nil {
		t.Fatalf("ScanStruct failed: %v", err)
	}

	if user.Name != "Alice" {
		t.Errorf("Expected name 'Alice', got '%s'", user.Name)
	}

	if user.Email != "alice@example.com" {
		t.Errorf("Expected email 'alice@example.com', got '%s'", user.Email)
	}

	if user.ID == 0 {
		t.Error("Expected ID to be set")
	}
}

func TestGetStructFields(t *testing.T) {
	user := User{
		ID:    1,
		Name:  "John",
		Email: "john@example.com",
	}

	fields, values := GetStructFields(&user)

	if len(fields) != 4 {
		t.Errorf("Expected 4 fields, got %d", len(fields))
	}

	if len(values) != 4 {
		t.Errorf("Expected 4 values, got %d", len(values))
	}

	// Check that fields contain expected db tags
	expectedFields := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"created_at": true,
	}

	for _, field := range fields {
		if !expectedFields[field] {
			t.Errorf("Unexpected field: %s", field)
		}
	}
}

func TestMigrator_Init(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	migrator := NewMigrator(db)
	if err := migrator.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Verify schema_migrations table exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='schema_migrations'").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to check for schema_migrations table: %v", err)
	}

	if count != 1 {
		t.Error("schema_migrations table should exist")
	}
}

func TestMigrator_CurrentVersion_Empty(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	migrator := NewMigrator(db)
	if err := migrator.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	version, err := migrator.CurrentVersion()
	if err != nil {
		t.Fatalf("CurrentVersion failed: %v", err)
	}

	if version != 0 {
		t.Errorf("Expected version 0, got %d", version)
	}
}

func TestMigrator_MigrateUp(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Verify users table exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to check for users table: %v", err)
	}

	if count != 1 {
		t.Error("users table should exist after migration")
	}

	// Verify posts table exists
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='posts'").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to check for posts table: %v", err)
	}

	if count != 1 {
		t.Error("posts table should exist after migration")
	}
}

func TestMigrator_CurrentVersion_AfterMigration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	migrator := NewMigrator(db)
	version, err := migrator.CurrentVersion()
	if err != nil {
		t.Fatalf("CurrentVersion failed: %v", err)
	}

	if version != 2 {
		t.Errorf("Expected version 2, got %d", version)
	}
}

func TestMigrator_IdempotentMigrations(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	migrator := NewMigrator(db)

	// Run migrations again
	if err := migrator.MigrateUp(); err != nil {
		t.Fatalf("Second MigrateUp failed: %v", err)
	}

	version, err := migrator.CurrentVersion()
	if err != nil {
		t.Fatalf("CurrentVersion failed: %v", err)
	}

	if version != 2 {
		t.Errorf("Expected version to remain 2, got %d", version)
	}
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	user := &User{
		Name:  "Bob",
		Email: "bob@example.com",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected ID to be set after create")
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	user := &User{
		Name:  "Charlie",
		Email: "charlie@example.com",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.FindByID(user.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found.Name != user.Name {
		t.Errorf("Expected name '%s', got '%s'", user.Name, found.Name)
	}

	if found.Email != user.Email {
		t.Errorf("Expected email '%s', got '%s'", user.Email, found.Email)
	}
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	_, err := repo.FindByID(999)

	if err == nil {
		t.Error("Expected error for non-existent user")
	}

	if err != sql.ErrNoRows {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	user := &User{
		Name:  "David",
		Email: "david@example.com",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.FindByEmail(user.Email)
	if err != nil {
		t.Fatalf("FindByEmail failed: %v", err)
	}

	if found.Name != user.Name {
		t.Errorf("Expected name '%s', got '%s'", user.Name, found.Name)
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	user := &User{
		Name:  "Eve",
		Email: "eve@example.com",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	user.Name = "Eve Updated"
	user.Email = "eve.updated@example.com"

	if err := repo.Update(user); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	found, err := repo.FindByID(user.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found.Name != "Eve Updated" {
		t.Errorf("Expected name 'Eve Updated', got '%s'", found.Name)
	}

	if found.Email != "eve.updated@example.com" {
		t.Errorf("Expected email 'eve.updated@example.com', got '%s'", found.Email)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	user := &User{
		Name:  "Frank",
		Email: "frank@example.com",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if err := repo.Delete(user.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := repo.FindByID(user.ID)
	if err != sql.ErrNoRows {
		t.Error("Expected user to be deleted")
	}
}

func TestUserRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	users := []*User{
		{Name: "User1", Email: "user1@example.com"},
		{Name: "User2", Email: "user2@example.com"},
		{Name: "User3", Email: "user3@example.com"},
	}

	for _, user := range users {
		if err := repo.Create(user); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	all, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(all) != 3 {
		t.Errorf("Expected 3 users, got %d", len(all))
	}
}

func TestPostRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepository(db)
	user := &User{
		Name:  "George",
		Email: "george@example.com",
	}
	if err := userRepo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	postRepo := NewPostRepository(db)
	post := &Post{
		UserID:  user.ID,
		Title:   "Test Post",
		Content: "This is a test post",
	}

	if err := postRepo.Create(post); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if post.ID == 0 {
		t.Error("Expected ID to be set after create")
	}
}

func TestPostRepository_FindByUserID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepository(db)
	user := &User{
		Name:  "Hannah",
		Email: "hannah@example.com",
	}
	if err := userRepo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	postRepo := NewPostRepository(db)
	posts := []*Post{
		{UserID: user.ID, Title: "Post 1", Content: "Content 1"},
		{UserID: user.ID, Title: "Post 2", Content: "Content 2"},
	}

	for _, post := range posts {
		if err := postRepo.Create(post); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	found, err := postRepo.FindByUserID(user.ID)
	if err != nil {
		t.Fatalf("FindByUserID failed: %v", err)
	}

	if len(found) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(found))
	}
}

func TestUserModel_TableName(t *testing.T) {
	user := &User{}
	if user.TableName() != "users" {
		t.Errorf("Expected table name 'users', got '%s'", user.TableName())
	}
}

func TestPostModel_TableName(t *testing.T) {
	post := &Post{}
	if post.TableName() != "posts" {
		t.Errorf("Expected table name 'posts', got '%s'", post.TableName())
	}
}

func TestUserRepository_CreateTimestamp(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	user := &User{
		Name:  "Ian",
		Email: "ian@example.com",
	}

	before := time.Now()
	if err := repo.Create(user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	after := time.Now()

	found, err := repo.FindByID(user.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found.CreatedAt.Before(before.Add(-time.Second)) || found.CreatedAt.After(after.Add(time.Second)) {
		t.Error("CreatedAt timestamp should be set to current time")
	}
}

func TestConcurrentUserCreation(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Create users concurrently using errgroup for better error handling
	type result struct {
		err error
		id  int
	}
	results := make(chan result, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			user := &User{
				Name:  fmt.Sprintf("User%d", id),
				Email: fmt.Sprintf("user%d@example.com", id),
			}
			err := repo.Create(user)
			results <- result{err: err, id: id}
		}(i)
	}

	// Collect results
	var errors []error
	for i := 0; i < 10; i++ {
		res := <-results
		if res.err != nil {
			errors = append(errors, fmt.Errorf("User%d: %w", res.id, res.err))
		}
	}

	// Report any errors
	if len(errors) > 0 {
		for _, err := range errors {
			t.Errorf("Concurrent create failed: %v", err)
		}
	}

	users, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(users) != 10 {
		t.Errorf("Expected 10 users, got %d", len(users))
	}
}

func TestQueryBuilder_SQLInjectionPrevention(t *testing.T) {
	qb := NewQueryBuilder("users").
		Where("name = ?", "Robert'); DROP TABLE users; --")

	query, args := qb.ToSQL()

	// Query should use placeholders
	if !contains(query, "?") {
		t.Error("Query should use placeholders to prevent SQL injection")
	}

	// Dangerous input should be in args, not query
	if contains(query, "DROP TABLE") {
		t.Error("Dangerous SQL should not be in query string")
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(args))
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsRec(s, substr, 0))
}

func containsRec(s, substr string, start int) bool {
	if start+len(substr) > len(s) {
		return false
	}
	if s[start:start+len(substr)] == substr {
		return true
	}
	return containsRec(s, substr, start+1)
}

func TestUniqueEmailConstraint(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	user1 := &User{
		Name:  "Jack",
		Email: "duplicate@example.com",
	}
	if err := repo.Create(user1); err != nil {
		t.Fatalf("First create failed: %v", err)
	}

	user2 := &User{
		Name:  "Jill",
		Email: "duplicate@example.com",
	}
	err := repo.Create(user2)

	if err == nil {
		t.Error("Expected error when creating user with duplicate email")
	}
}
