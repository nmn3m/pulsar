package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nmn3m/pulsar/backend/internal/config"
	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/middleware"
	"github.com/nmn3m/pulsar/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo                 repository.UserRepository
	orgRepo                  repository.OrganizationRepository
	config                   *config.Config
	emailVerificationService *EmailVerificationService
}

func NewAuthService(
	userRepo repository.UserRepository,
	orgRepo repository.OrganizationRepository,
	cfg *config.Config,
	emailVerificationService *EmailVerificationService,
) *AuthService {
	return &AuthService{
		userRepo:                 userRepo,
		orgRepo:                  orgRepo,
		config:                   cfg,
		emailVerificationService: emailVerificationService,
	}
}

type RegisterRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Username         string `json:"username" binding:"required,min=3,max=50"`
	Password         string `json:"password" binding:"required,min=8"`
	FullName         string `json:"full_name"`
	OrganizationName string `json:"organization_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User                      *domain.User         `json:"user"`
	Organization              *domain.Organization `json:"organization"`
	AccessToken               string               `json:"access_token"`
	RefreshToken              string               `json:"refresh_token"`
	RequiresEmailVerification bool                 `json:"requires_email_verification,omitempty"`
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	existingUser, _ = s.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, fmt.Errorf("user with this username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Check if email verification is enabled
	emailVerificationEnabled := s.emailVerificationService != nil && s.emailVerificationService.IsEmailServiceConfigured()

	// Create user
	user := &domain.User{
		ID:                      uuid.New(),
		Email:                   req.Email,
		Username:                req.Username,
		PasswordHash:            string(hashedPassword),
		FullName:                &req.FullName,
		Timezone:                "UTC",
		NotificationPreferences: make(map[string]interface{}),
		IsActive:                true,
		EmailVerified:           !emailVerificationEnabled, // Skip verification if email service not configured
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create organization with slug from name
	orgSlug := generateSlug(req.OrganizationName)
	org := &domain.Organization{
		ID:       uuid.New(),
		Name:     req.OrganizationName,
		Slug:     orgSlug,
		Plan:     string(domain.PlanFree),
		Settings: make(map[string]interface{}),
	}

	if err := s.orgRepo.Create(ctx, org); err != nil {
		// Rollback user creation if org creation fails
		// In production, use a transaction
		s.userRepo.Delete(ctx, user.ID)
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// Add user to organization as owner
	if err := s.orgRepo.AddUser(ctx, org.ID, user.ID, domain.RoleOwner); err != nil {
		// Rollback in production
		s.userRepo.Delete(ctx, user.ID)
		s.orgRepo.Delete(ctx, org.ID)
		return nil, fmt.Errorf("failed to add user to organization: %w", err)
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user.ID, user.Email, org.ID, string(domain.RoleOwner))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user.ID, user.Email, org.ID, string(domain.RoleOwner))
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Clear password hash before returning
	user.PasswordHash = ""

	// Send verification email if email verification is enabled
	if emailVerificationEnabled {
		if err := s.emailVerificationService.CreateAndSendOTP(ctx, user.ID, user.Email, user.Username); err != nil {
			// Log error but don't fail registration
			// In production, you might want to handle this differently
			fmt.Printf("Warning: failed to send verification email: %v\n", err)
		}
	}

	return &AuthResponse{
		User:                      user,
		Organization:              org,
		AccessToken:               accessToken,
		RefreshToken:              refreshToken,
		RequiresEmailVerification: emailVerificationEnabled,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is disabled")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Get user's organizations
	orgs, err := s.orgRepo.ListUserOrganizations(ctx, user.ID)
	if err != nil || len(orgs) == 0 {
		return nil, fmt.Errorf("user has no organizations")
	}

	// Use the first organization (in production, let user select)
	org := orgs[0]

	// Get user's role in the organization
	role, err := s.orgRepo.GetUserRole(ctx, org.ID, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user.ID, user.Email, org.ID, string(role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user.ID, user.Email, org.ID, string(role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Clear password hash before returning
	user.PasswordHash = ""

	// Check if email verification is required
	requiresVerification := !user.EmailVerified && s.emailVerificationService != nil && s.emailVerificationService.IsEmailServiceConfigured()

	return &AuthResponse{
		User:                      user,
		Organization:              org,
		AccessToken:               accessToken,
		RefreshToken:              refreshToken,
		RequiresEmailVerification: requiresVerification,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// Parse refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &middleware.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.RefreshSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(*middleware.Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	// Check if token is expired
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token expired")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Get organization
	org, err := s.orgRepo.GetByID(ctx, claims.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("organization not found")
	}

	// Generate new tokens
	accessToken, err := s.generateAccessToken(user.ID, user.Email, org.ID, claims.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.generateRefreshToken(user.ID, user.Email, org.ID, claims.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Clear password hash before returning
	user.PasswordHash = ""

	return &AuthResponse{
		User:         user,
		Organization: org,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) GetMe(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Clear password hash
	user.PasswordHash = ""

	return user, nil
}

func (s *AuthService) generateAccessToken(userID uuid.UUID, email string, orgID uuid.UUID, role string) (string, error) {
	claims := &middleware.Claims{
		UserID:         userID,
		Email:          email,
		OrganizationID: orgID,
		Role:           role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.config.JWT.AccessTTL) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *AuthService) generateRefreshToken(userID uuid.UUID, email string, orgID uuid.UUID, role string) (string, error) {
	claims := &middleware.Claims{
		UserID:         userID,
		Email:          email,
		OrganizationID: orgID,
		Role:           role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.config.JWT.RefreshTTL) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.RefreshSecret))
}

func generateSlug(name string) string {
	// Simple slug generation (in production, use a proper library)
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Add timestamp to ensure uniqueness
	slug = fmt.Sprintf("%s-%d", slug, time.Now().Unix())
	return slug
}
