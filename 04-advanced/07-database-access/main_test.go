package main

import (
	"database/sql"
	"os"
	"testing"
)

func setupTestDB(t *testing.T) (*sql.DB, *UserRepository) {
	t.Helper()

	// Create temporary test database
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	repo := NewUserRepository(db)
	return db, repo
}

func TestInitDB(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - add proper TODOs")

	tests := []struct {
		name    string
		dbPath  string
		wantErr bool
	}{
		{
			name:    "Create in-memory database",
			dbPath:  ":memory:",
			wantErr: false,
		},
		{
			name:    "Create file-based database",
			dbPath:  "test_init.db",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := InitDB(tt.dbPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if db != nil {
				defer db.Close()
				// Clean up file-based database
				if tt.dbPath != ":memory:" {
					defer os.Remove(tt.dbPath)
				}

				// Verify table was created
				var tableName string
				err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users'").Scan(&tableName)
				if err != nil {
					t.Errorf("Table 'users' was not created: %v", err)
				}
			}
		})
	}
}

func TestUserRepositoryCreate(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - add proper TODOs")

	db, repo := setupTestDB(t)
	defer db.Close()

	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name: "Create valid user",
			user: &User{
				Name:  "Alice",
				Email: "alice@example.com",
				Age:   30,
			},
			wantErr: false,
		},
		{
			name: "Create user with different email",
			user: &User{
				Name:  "Bob",
				Email: "bob@example.com",
				Age:   25,
			},
			wantErr: false,
		},
		{
			name: "Create user with zero age",
			user: &User{
				Name:  "Charlie",
				Email: "charlie@example.com",
				Age:   0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.user.ID == 0 {
					t.Error("Create() did not set user ID")
				}
			}
		})
	}
}

func TestUserRepositoryCreateDuplicateEmail(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - add proper TODOs")

	db, repo := setupTestDB(t)
	defer db.Close()

	user1 := &User{Name: "Alice", Email: "duplicate@example.com", Age: 30}
	if err := repo.Create(user1); err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	user2 := &User{Name: "Bob", Email: "duplicate@example.com", Age: 25}
	err := repo.Create(user2)
	if err == nil {
		t.Error("Create() should fail with duplicate email")
	}
}

func TestUserRepositoryGet(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - add proper TODOs")

	db, repo := setupTestDB(t)
	defer db.Close()

	// Create a test user
	original := &User{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   30,
	}
	if err := repo.Create(original); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "Get existing user",
			id:      original.ID,
			wantErr: false,
		},
		{
			name:    "Get non-existent user",
			id:      99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.Get(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if user == nil {
					t.Error("Get() returned nil user")
					return
				}
				if user.ID != original.ID {
					t.Errorf("Get() ID = %d, expected %d", user.ID, original.ID)
				}
				if user.Name != original.Name {
					t.Errorf("Get() Name = %q, expected %q", user.Name, original.Name)
				}
				if user.Email != original.Email {
					t.Errorf("Get() Email = %q, expected %q", user.Email, original.Email)
				}
				if user.Age != original.Age {
					t.Errorf("Get() Age = %d, expected %d", user.Age, original.Age)
				}
			}
		})
	}
}

func TestUserRepositoryUpdate(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - add proper TODOs")

	db, repo := setupTestDB(t)
	defer db.Close()

	// Create a test user
	user := &User{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   30,
	}
	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name     string
		modified *User
		wantErr  bool
	}{
		{
			name: "Update user name",
			modified: &User{
				ID:    user.ID,
				Name:  "Alice Updated",
				Email: user.Email,
				Age:   user.Age,
			},
			wantErr: false,
		},
		{
			name: "Update user age",
			modified: &User{
				ID:    user.ID,
				Name:  "Alice Updated",
				Email: user.Email,
				Age:   31,
			},
			wantErr: false,
		},
		{
			name: "Update non-existent user",
			modified: &User{
				ID:    99999,
				Name:  "Ghost",
				Email: "ghost@example.com",
				Age:   0,
			},
			wantErr: false, // SQLite doesn't error on updating non-existent rows
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(tt.modified)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.modified.ID == user.ID {
				// Verify the update
				retrieved, err := repo.Get(tt.modified.ID)
				if err != nil {
					t.Errorf("Failed to retrieve updated user: %v", err)
					return
				}
				if retrieved.Name != tt.modified.Name {
					t.Errorf("Update() Name = %q, expected %q", retrieved.Name, tt.modified.Name)
				}
				if retrieved.Age != tt.modified.Age {
					t.Errorf("Update() Age = %d, expected %d", retrieved.Age, tt.modified.Age)
				}
			}
		})
	}
}

func TestUserRepositoryDelete(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - add proper TODOs")

	db, repo := setupTestDB(t)
	defer db.Close()

	// Create test users
	user1 := &User{Name: "Alice", Email: "alice@example.com", Age: 30}
	user2 := &User{Name: "Bob", Email: "bob@example.com", Age: 25}
	if err := repo.Create(user1); err != nil {
		t.Fatalf("Failed to create test user 1: %v", err)
	}
	if err := repo.Create(user2); err != nil {
		t.Fatalf("Failed to create test user 2: %v", err)
	}

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "Delete existing user",
			id:      user1.ID,
			wantErr: false,
		},
		{
			name:    "Delete non-existent user",
			id:      99999,
			wantErr: false, // SQLite doesn't error on deleting non-existent rows
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.id == user1.ID {
				// Verify deletion
				_, err := repo.Get(tt.id)
				if err == nil {
					t.Error("Delete() did not remove user from database")
				}
			}
		})
	}
}

func TestUserRepositoryList(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - add proper TODOs")

	db, repo := setupTestDB(t)
	defer db.Close()

	tests := []struct {
		name        string
		setupUsers  []*User
		wantCount   int
		wantErr     bool
	}{
		{
			name:       "List empty database",
			setupUsers: []*User{},
			wantCount:  0,
			wantErr:    false,
		},
		{
			name: "List single user",
			setupUsers: []*User{
				{Name: "Alice", Email: "alice@example.com", Age: 30},
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "List multiple users",
			setupUsers: []*User{
				{Name: "Alice", Email: "alice@example.com", Age: 30},
				{Name: "Bob", Email: "bob@example.com", Age: 25},
				{Name: "Charlie", Email: "charlie@example.com", Age: 35},
			},
			wantCount: 3,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before each test
			db.Exec("DELETE FROM users")

			// Setup test data
			for _, user := range tt.setupUsers {
				if err := repo.Create(user); err != nil {
					t.Fatalf("Failed to create setup user: %v", err)
				}
			}

			users, err := repo.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(users) != tt.wantCount {
					t.Errorf("List() returned %d users, expected %d", len(users), tt.wantCount)
				}

				// Verify all users have required fields
				for i, user := range users {
					if user.ID == 0 {
						t.Errorf("User %d has zero ID", i)
					}
					if user.Name == "" {
						t.Errorf("User %d has empty name", i)
					}
					if user.Email == "" {
						t.Errorf("User %d has empty email", i)
					}
				}
			}
		})
	}
}

func TestUserRepositoryCreateMultiple(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - add proper TODOs")

	db, repo := setupTestDB(t)
	defer db.Close()

	tests := []struct {
		name    string
		users   []*User
		wantErr bool
	}{
		{
			name: "Create multiple valid users",
			users: []*User{
				{Name: "Alice", Email: "alice@example.com", Age: 30},
				{Name: "Bob", Email: "bob@example.com", Age: 25},
				{Name: "Charlie", Email: "charlie@example.com", Age: 35},
			},
			wantErr: false,
		},
		{
			name:    "Create empty slice",
			users:   []*User{},
			wantErr: false,
		},
		{
			name: "Create with duplicate email should rollback all",
			users: []*User{
				{Name: "Dave", Email: "dave@example.com", Age: 40},
				{Name: "Eve", Email: "alice@example.com", Age: 28}, // Duplicate from previous test
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialCount, _ := repo.List()
			initialLen := len(initialCount)

			err := repo.CreateMultiple(tt.users)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMultiple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify all users were created
				for _, user := range tt.users {
					if user.ID == 0 {
						t.Error("CreateMultiple() did not set user IDs")
					}
				}

				// Verify count in database
				allUsers, err := repo.List()
				if err != nil {
					t.Errorf("Failed to list users: %v", err)
					return
				}
				expectedCount := initialLen + len(tt.users)
				if len(allUsers) != expectedCount {
					t.Errorf("CreateMultiple() created %d users, expected %d", len(allUsers)-initialLen, len(tt.users))
				}
			} else {
				// Verify transaction rollback - count should be unchanged
				allUsers, err := repo.List()
				if err != nil {
					t.Errorf("Failed to list users: %v", err)
					return
				}
				if len(allUsers) != initialLen {
					t.Errorf("CreateMultiple() failed but added users (transaction not rolled back)")
				}
			}
		})
	}
}

func TestUserRepositoryTransaction(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - add proper TODOs")

	db, repo := setupTestDB(t)
	defer db.Close()

	t.Run("Transaction commits on success", func(t *testing.T) {
		err := repo.Transaction(func(tx *sql.Tx) error {
			_, err := tx.Exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)", "Test", "test@example.com", 20)
			return err
		})

		if err != nil {
			t.Errorf("Transaction() error = %v", err)
		}

		users, _ := repo.List()
		if len(users) == 0 {
			t.Error("Transaction did not commit")
		}
	})

	t.Run("Transaction rolls back on error", func(t *testing.T) {
		initialUsers, _ := repo.List()
		initialCount := len(initialUsers)

		err := repo.Transaction(func(tx *sql.Tx) error {
			_, err := tx.Exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)", "Rollback", "rollback@example.com", 25)
			if err != nil {
				return err
			}
			// Force an error to trigger rollback
			return sql.ErrTxDone
		})

		if err == nil {
			t.Error("Transaction() should have returned error")
		}

		users, _ := repo.List()
		if len(users) != initialCount {
			t.Error("Transaction did not rollback on error")
		}
	})
}
