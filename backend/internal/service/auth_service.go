package service

import (
	"context"
	"errors"
	"time"

	"github.com/Snehaaa-Kri/job-apply-automator/config"
	"github.com/Snehaaa-Kri/job-apply-automator/internal/domain"
	"github.com/Snehaaa-Kri/job-apply-automator/internal/repository"
	"github.com/Snehaaa-Kri/job-apply-automator/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// --- Request/Response DTOs ---
// DTO = Data Transfer Object. These are the shapes of data coming IN to the service.
// Node/TypeScript analogy: interface RegisterDTO { name: string; email: string; password: string }
// We keep DTOs in the service layer (not domain) because they're input shapes, not entities.

type RegisterInput struct {
	Name     string          `json:"name"     validate:"required,min=2,max=100"`
	Email    string          `json:"email"    validate:"required,email"`
	Password string          `json:"password" validate:"required,min=8"`
	Role     domain.UserRole `json:"role"     validate:"required,oneof=student professional"`
}

type LoginInput struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token"`
}

// --- JWT Claims ---
// Claims = the payload inside the JWT token.
// Node analogy: jwt.sign({ userId, email, role }, secret, { expiresIn: '24h' })
// In Go we define a struct for it.
type JWTClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims // embeds standard fields: ExpiresAt, IssuedAt, etc.
}

// --- Service Interface ---
// Same pattern as repository — define an interface, implement it.
// This makes the handler layer depend on the interface, not the concrete service.
// Node analogy: IAuthService interface in TypeScript
type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*AuthResponse, error)
	Login(ctx context.Context, input LoginInput) (*AuthResponse, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}

// --- Concrete implementation ---
type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
	log      *logger.Logger
}

// NewAuthService is the constructor.
// Node analogy: export function createAuthService(userRepo, config, logger) { ... }
func NewAuthService(
	userRepo repository.UserRepository,
	cfg *config.Config,
	log *logger.Logger,
) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
		log:      log,
	}
}

// Register creates a new user account.
// Node analogy: async register(dto) { hash password → check duplicate → save → return JWT }
func (s *authService) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	// 1. Check if email already exists
	existing, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	// 2. Hash the password
	// bcrypt cost 12 = strong enough for production, not too slow.
	// Node analogy: bcrypt.hash(password, 12)
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return nil, err
	}

	// 3. Build the User entity
	user := &domain.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         input.Role,
		Provider:     domain.AuthProviderLocal,
		IsActive:     true,
	}

	// 4. Save to DB
	created, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// 5. Generate JWT
	token, err := s.generateToken(created)
	if err != nil {
		return nil, err
	}

	s.log.Info("new user registered", logger.Field("email", input.Email))

	return &AuthResponse{User: created, Token: token}, nil
}

// Login authenticates an existing user.
// Node analogy: async login(dto) { find user → compare password → return JWT }
func (s *authService) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	// 1. Find user by email
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	// Don't leak whether email exists — always say "invalid credentials"
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// 2. Check account is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// 3. Compare password with stored hash
	// Node analogy: bcrypt.compare(password, user.passwordHash)
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		// bcrypt returns an error if passwords don't match
		return nil, errors.New("invalid email or password")
	}

	// 4. Generate JWT
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{User: user, Token: token}, nil
}

// GetMe returns the current authenticated user's profile.
func (s *authService) GetMe(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// generateToken creates a signed JWT for the user.
// Node analogy: jwt.sign({ userId, email, role }, JWT_SECRET, { expiresIn: '24h' })
func (s *authService) generateToken(user *domain.User) (string, error) {
	claims := JWTClaims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Duration(s.cfg.JWT.ExpiryHours) * time.Hour),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	// Sign with HS256 (HMAC-SHA256) — same algorithm as most Node JWT libs default to
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}
