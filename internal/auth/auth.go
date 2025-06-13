package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/jasonlovesdoggo/velo/internal/log"
	"github.com/jasonlovesdoggo/velo/internal/state"
	"google.golang.org/grpc"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenInvalid       = errors.New("invalid token")
)

// User represents a system user
type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"` // hashed
	Role     string    `json:"role"`     // admin, user
	Created  time.Time `json:"created"`
	Active   bool      `json:"active"`
}

// Token represents an authentication token
type Token struct {
	Value     string    `json:"value"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Created   time.Time `json:"created"`
}

// AuthService handles authentication and authorization
type AuthService struct {
	store state.StateStore
}

// NewAuthService creates a new authentication service
func NewAuthService(store state.StateStore) *AuthService {
	return &AuthService{
		store: store,
	}
}

// Initialize sets up the authentication service with default admin user
func (a *AuthService) Initialize() error {
	// Check if admin user exists
	users, err := a.GetAllUsers()
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	// Create default admin user if no users exist
	if len(users) == 0 {
		defaultPassword := "admin" // TODO: Generate random password
		hashedPassword, err := a.hashPassword(defaultPassword)
		if err != nil {
			return fmt.Errorf("failed to hash default password: %w", err)
		}

		admin := &User{
			ID:       generateID(),
			Username: "admin",
			Password: hashedPassword,
			Role:     "admin",
			Created:  time.Now(),
			Active:   true,
		}

		if err := a.CreateUser(admin); err != nil {
			return fmt.Errorf("failed to create default admin user: %w", err)
		}

		log.Info("Created default admin user", "username", "admin", "password", defaultPassword)
		log.Warn("Please change the default admin password immediately!")
	}

	return nil
}

// Authenticate validates credentials and returns a token
func (a *AuthService) Authenticate(username, password string) (*Token, error) {
	user, err := a.GetUserByUsername(username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.Active {
		return nil, ErrInvalidCredentials
	}

	if !a.verifyPassword(password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	tokenValue, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	token := &Token{
		Value:     tokenValue,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour expiry
		Created:   time.Now(),
	}

	// Store token
	if err := a.store.Set("token:"+token.Value, token); err != nil {
		return nil, fmt.Errorf("failed to store token: %w", err)
	}

	log.Info("User authenticated", "username", username, "userID", user.ID)
	return token, nil
}

// ValidateToken checks if a token is valid and returns the associated user
func (a *AuthService) ValidateToken(tokenValue string) (*User, error) {
	var token Token
	if err := a.store.Get("token:"+tokenValue, &token); err != nil {
		return nil, ErrTokenInvalid
	}

	if time.Now().After(token.ExpiresAt) {
		// Clean up expired token
		a.store.Delete("token:" + tokenValue)
		return nil, ErrTokenExpired
	}

	user, err := a.GetUserByID(token.UserID)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	if !user.Active {
		return nil, ErrTokenInvalid
	}

	return user, nil
}

// RevokeToken invalidates a token
func (a *AuthService) RevokeToken(tokenValue string) error {
	return a.store.Delete("token:" + tokenValue)
}

// CreateUser creates a new user
func (a *AuthService) CreateUser(user *User) error {
	// Check if username already exists
	if _, err := a.GetUserByUsername(user.Username); err == nil {
		return errors.New("username already exists")
	}

	return a.store.Set("user:"+user.ID, user)
}

// GetUserByID retrieves a user by ID
func (a *AuthService) GetUserByID(id string) (*User, error) {
	var user User
	if err := a.store.Get("user:"+id, &user); err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (a *AuthService) GetUserByUsername(username string) (*User, error) {
	users, err := a.GetAllUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

// GetAllUsers retrieves all users
func (a *AuthService) GetAllUsers() ([]User, error) {
	keys, err := a.store.List("user:")
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(keys))
	for _, key := range keys {
		var user User
		if err := a.store.Get(key, &user); err != nil {
			continue // Skip invalid entries
		}
		users = append(users, user)
	}

	return users, nil
}

// UpdateUser updates an existing user
func (a *AuthService) UpdateUser(user *User) error {
	// Verify user exists
	if _, err := a.GetUserByID(user.ID); err != nil {
		return errors.New("user not found")
	}

	return a.store.Set("user:"+user.ID, user)
}

// DeleteUser deactivates a user
func (a *AuthService) DeleteUser(id string) error {
	user, err := a.GetUserByID(id)
	if err != nil {
		return err
	}

	user.Active = false
	return a.UpdateUser(user)
}

// ChangePassword changes a user's password
func (a *AuthService) ChangePassword(userID, newPassword string) error {
	user, err := a.GetUserByID(userID)
	if err != nil {
		return err
	}

	hashedPassword, err := a.hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	return a.UpdateUser(user)
}

// Helper functions

func (a *AuthService) hashPassword(password string) (string, error) {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:]), nil
}

func (a *AuthService) verifyPassword(password, hash string) bool {
	passwordHash, err := a.hashPassword(password)
	if err != nil {
		return false
	}
	return passwordHash == hash
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Middleware for gRPC authentication
func (a *AuthService) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// For now, skip auth for basic functionality
	// TODO: Implement proper gRPC metadata token extraction
	return handler(ctx, req)
}
