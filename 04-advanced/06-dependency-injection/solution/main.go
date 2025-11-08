package main

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Domain Models
type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

// Interfaces for dependencies
type UserRepository interface {
	FindByID(ctx context.Context, id int) (*User, error)
	Save(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]*User, error)
}

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, err error)
	Debug(msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value interface{}
}

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	Clear()
}

type EmailSender interface {
	Send(to, subject, body string) error
}

// InMemoryUserRepository - concrete implementation
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[int]*User
	logger Logger
}

func NewInMemoryUserRepository(logger Logger) *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:  make(map[int]*User),
		logger: logger,
	}
}

func (r *InMemoryUserRepository) FindByID(ctx context.Context, id int) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if user, ok := r.users[id]; ok {
		r.logger.Debug("User found", Field{Key: "id", Value: id})
		return user, nil
	}

	return nil, fmt.Errorf("user %d not found", id)
}

func (r *InMemoryUserRepository) Save(ctx context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user.ID == 0 {
		// Auto-increment ID
		user.ID = len(r.users) + 1
	}

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	r.users[user.ID] = user
	r.logger.Info("User saved", Field{Key: "id", Value: user.ID}, Field{Key: "name", Value: user.Name})
	return nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return fmt.Errorf("user %d not found", id)
	}

	delete(r.users, id)
	r.logger.Info("User deleted", Field{Key: "id", Value: id})
	return nil
}

func (r *InMemoryUserRepository) List(ctx context.Context) ([]*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	r.logger.Debug("Users listed", Field{Key: "count", Value: len(users)})
	return users, nil
}

// SimpleLogger - concrete implementation
type SimpleLogger struct {
	mu     sync.Mutex
	prefix string
}

func NewSimpleLogger(prefix string) *SimpleLogger {
	return &SimpleLogger{prefix: prefix}
}

func (l *SimpleLogger) Info(msg string, fields ...Field) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Printf("[%s][INFO] %s %s\n", l.prefix, msg, formatFields(fields))
}

func (l *SimpleLogger) Error(msg string, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Printf("[%s][ERROR] %s: %v\n", l.prefix, msg, err)
}

func (l *SimpleLogger) Debug(msg string, fields ...Field) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Printf("[%s][DEBUG] %s %s\n", l.prefix, msg, formatFields(fields))
}

func formatFields(fields []Field) string {
	if len(fields) == 0 {
		return ""
	}

	parts := make([]string, len(fields))
	for i, f := range fields {
		parts[i] = fmt.Sprintf("%s=%v", f.Key, f.Value)
	}
	return "[" + strings.Join(parts, " ") + "]"
}

// NopLogger - logger that does nothing (useful for testing)
type NopLogger struct{}

func NewNopLogger() *NopLogger {
	return &NopLogger{}
}

func (l *NopLogger) Info(msg string, fields ...Field)  {}
func (l *NopLogger) Error(msg string, err error)       {}
func (l *NopLogger) Debug(msg string, fields ...Field) {}

// InMemoryCache - concrete implementation
type InMemoryCache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		items: make(map[string]cacheItem),
	}
}

func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := cacheItem{
		value: value,
	}

	if ttl > 0 {
		item.expiration = time.Now().Add(ttl)
	}

	c.items[key] = item
}

func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *InMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]cacheItem)
}

// MockEmailSender - concrete implementation for testing
type MockEmailSender struct {
	SentEmails []Email
	mu         sync.Mutex
}

type Email struct {
	To      string
	Subject string
	Body    string
}

func NewMockEmailSender() *MockEmailSender {
	return &MockEmailSender{
		SentEmails: make([]Email, 0),
	}
}

func (m *MockEmailSender) Send(to, subject, body string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.SentEmails = append(m.SentEmails, Email{
		To:      to,
		Subject: subject,
		Body:    body,
	})
	return nil
}

// UserService - service with constructor injection
type UserService struct {
	repo   UserRepository
	logger Logger
	cache  Cache
	email  EmailSender
	config *ServiceConfig
}

// ServiceConfig for functional options pattern
type ServiceConfig struct {
	CacheEnabled bool
	CacheTTL     time.Duration
	MaxRetries   int
	Timeout      time.Duration
}

// Option is a functional option for UserService
type Option func(*ServiceConfig)

func WithCache(enabled bool) Option {
	return func(c *ServiceConfig) {
		c.CacheEnabled = enabled
	}
}

func WithCacheTTL(ttl time.Duration) Option {
	return func(c *ServiceConfig) {
		c.CacheTTL = ttl
	}
}

func WithMaxRetries(retries int) Option {
	return func(c *ServiceConfig) {
		c.MaxRetries = retries
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *ServiceConfig) {
		c.Timeout = timeout
	}
}

// NewUserService creates a new UserService with constructor injection
func NewUserService(repo UserRepository, logger Logger, cache Cache, email EmailSender, opts ...Option) *UserService {
	config := &ServiceConfig{
		CacheEnabled: true,
		CacheTTL:     5 * time.Minute,
		MaxRetries:   3,
		Timeout:      30 * time.Second,
	}

	for _, opt := range opts {
		opt(config)
	}

	return &UserService{
		repo:   repo,
		logger: logger,
		cache:  cache,
		email:  email,
		config: config,
	}
}

func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
	// Try cache first if enabled
	if s.config.CacheEnabled {
		cacheKey := fmt.Sprintf("user:%d", id)
		if cached, ok := s.cache.Get(cacheKey); ok {
			s.logger.Debug("Cache hit", Field{Key: "id", Value: id})
			return cached.(*User), nil
		}
	}

	// Fetch from repository
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get user", err)
		return nil, err
	}

	// Cache the result
	if s.config.CacheEnabled {
		cacheKey := fmt.Sprintf("user:%d", id)
		s.cache.Set(cacheKey, user, s.config.CacheTTL)
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, name, email string) (*User, error) {
	user := &User{
		Name:  name,
		Email: email,
	}

	if err := s.repo.Save(ctx, user); err != nil {
		s.logger.Error("Failed to create user", err)
		return nil, err
	}

	// Send welcome email
	if err := s.email.Send(email, "Welcome!", fmt.Sprintf("Welcome, %s!", name)); err != nil {
		s.logger.Error("Failed to send welcome email", err)
		// Don't fail the operation if email fails
	}

	s.logger.Info("User created", Field{Key: "id", Value: user.ID}, Field{Key: "name", Value: name})
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *User) error {
	if err := s.repo.Save(ctx, user); err != nil {
		s.logger.Error("Failed to update user", err)
		return err
	}

	// Invalidate cache
	if s.config.CacheEnabled {
		cacheKey := fmt.Sprintf("user:%d", user.ID)
		s.cache.Delete(cacheKey)
	}

	s.logger.Info("User updated", Field{Key: "id", Value: user.ID})
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete user", err)
		return err
	}

	// Invalidate cache
	if s.config.CacheEnabled {
		cacheKey := fmt.Sprintf("user:%d", id)
		s.cache.Delete(cacheKey)
	}

	s.logger.Info("User deleted", Field{Key: "id", Value: id})
	return nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]*User, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		s.logger.Error("Failed to list users", err)
		return nil, err
	}

	s.logger.Info("Users listed", Field{Key: "count", Value: len(users)})
	return users, nil
}

// Container - simple DI container
type Container struct {
	logger      Logger
	cache       Cache
	emailSender EmailSender
	userRepo    UserRepository
	userService *UserService
}

// NewContainer creates a new DI container with all dependencies
func NewContainer() *Container {
	c := &Container{}

	// Initialize dependencies in order
	c.logger = NewSimpleLogger("app")
	c.cache = NewInMemoryCache()
	c.emailSender = NewMockEmailSender()
	c.userRepo = NewInMemoryUserRepository(c.logger)
	c.userService = NewUserService(c.userRepo, c.logger, c.cache, c.emailSender)

	return c
}

// NewTestContainer creates a container for testing
func NewTestContainer() *Container {
	c := &Container{}

	// Use nop logger for tests
	c.logger = NewNopLogger()
	c.cache = NewInMemoryCache()
	c.emailSender = NewMockEmailSender()
	c.userRepo = NewInMemoryUserRepository(c.logger)
	c.userService = NewUserService(
		c.userRepo,
		c.logger,
		c.cache,
		c.emailSender,
		WithCache(true),
		WithCacheTTL(1*time.Minute),
	)

	return c
}

// Getters for services
func (c *Container) UserService() *UserService {
	return c.userService
}

func (c *Container) Logger() Logger {
	return c.logger
}

func (c *Container) Cache() Cache {
	return c.cache
}

func (c *Container) EmailSender() EmailSender {
	return c.emailSender
}

// Provider pattern for dynamic dependency resolution
type Provider func() (interface{}, error)

type Registry struct {
	providers map[string]Provider
	instances map[string]interface{}
	mu        sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
		instances: make(map[string]interface{}),
	}
}

func (r *Registry) Register(name string, provider Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[name] = provider
}

func (r *Registry) Get(name string) (interface{}, error) {
	// Check if already instantiated (singleton pattern)
	r.mu.RLock()
	if instance, ok := r.instances[name]; ok {
		r.mu.RUnlock()
		return instance, nil
	}
	r.mu.RUnlock()

	// Get provider
	r.mu.RLock()
	provider, ok := r.providers[name]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("provider not found: %s", name)
	}

	// Double-check locking pattern
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check again in case another goroutine created it
	if instance, ok := r.instances[name]; ok {
		return instance, nil
	}

	// Create new instance
	instance, err := provider()
	if err != nil {
		return nil, err
	}

	r.instances[name] = instance
	return instance, nil
}

func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.instances = make(map[string]interface{})
}

// Lifecycle interface for start/stop operations
type Lifecycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// App demonstrates lifecycle management
type App struct {
	container  *Container
	components []Lifecycle
}

func NewApp(container *Container) *App {
	return &App{
		container:  container,
		components: make([]Lifecycle, 0),
	}
}

func (a *App) AddComponent(c Lifecycle) {
	a.components = append(a.components, c)
}

func (a *App) Start(ctx context.Context) error {
	for _, c := range a.components {
		if err := c.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	// Stop in reverse order
	for i := len(a.components) - 1; i >= 0; i-- {
		if err := a.components[i].Stop(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Wire-style provider pattern
func ProvideLogger() Logger {
	return NewSimpleLogger("wire")
}

func ProvideCache() Cache {
	return NewInMemoryCache()
}

func ProvideEmailSender() EmailSender {
	return NewMockEmailSender()
}

func ProvideUserRepository(logger Logger) UserRepository {
	return NewInMemoryUserRepository(logger)
}

func ProvideUserService(repo UserRepository, logger Logger, cache Cache, email EmailSender) *UserService {
	return NewUserService(repo, logger, cache, email)
}

// InitializeUserService demonstrates wire-style initialization
func InitializeUserService() *UserService {
	logger := ProvideLogger()
	cache := ProvideCache()
	email := ProvideEmailSender()
	repo := ProvideUserRepository(logger)
	service := ProvideUserService(repo, logger, cache, email)
	return service
}

// ErrorWithStack demonstrates error handling with dependency injection
type ErrorWithStack struct {
	message string
	cause   error
	stack   string
}

func NewErrorWithStack(message string, cause error) *ErrorWithStack {
	return &ErrorWithStack{
		message: message,
		cause:   cause,
		stack:   captureStack(),
	}
}

func (e *ErrorWithStack) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

func (e *ErrorWithStack) Unwrap() error {
	return e.cause
}

func (e *ErrorWithStack) StackTrace() string {
	return e.stack
}

func captureStack() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// Validation errors
var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("not found")
)

func main() {
	ctx := context.Background()

	// Create container
	container := NewContainer()

	// Use service
	user, err := container.UserService().CreateUser(ctx, "Alice", "alice@example.com")
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		return
	}

	fmt.Printf("Created user: %+v\n", user)

	// Retrieve user
	retrieved, err := container.UserService().GetUser(ctx, user.ID)
	if err != nil {
		fmt.Printf("Error retrieving user: %v\n", err)
		return
	}

	fmt.Printf("Retrieved user: %+v\n", retrieved)

	// List users
	users, err := container.UserService().ListUsers(ctx)
	if err != nil {
		fmt.Printf("Error listing users: %v\n", err)
		return
	}

	fmt.Printf("Total users: %d\n", len(users))
}
