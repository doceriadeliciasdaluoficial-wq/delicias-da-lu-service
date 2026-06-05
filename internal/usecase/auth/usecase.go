package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"delicias-da-lu-service.com/mod/internal/entity/user"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	userRepo "delicias-da-lu-service.com/mod/internal/repository/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type AuthUseCase interface {
	Login(ctx context.Context, username, password string) (*user.LoginResponse, error)
	RefreshToken(ctx context.Context, token string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, error)
}

type authUseCaseImpl struct {
	userRepository userRepo.UserRepository
	jwtSecret      string
}

func NewAuthUseCase(userRepository userRepo.UserRepository, jwtSecret string) AuthUseCase {
	return authUseCaseImpl{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
	}
}

func (a authUseCaseImpl) Login(ctx context.Context, username, password string) (*user.LoginResponse, error) {
	// Get user by username
	usr, err := a.userRepository.GetByUsername(ctx, username)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("user not found")
		return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "https://delicias-da-lu-service.com/docs/errors/invalid-credentials",
			Title:      "Invalid Credentials",
			Detail:     "The provided username or password is incorrect",
			HTTPStatus: http.StatusUnauthorized,
			Instance:   "localhost:8080/v1/auth/login",
			Severity:   problemdetails.Err,
		})
	}

	// Compare passwords (simple hash comparison - in production use bcrypt)
	passwordHash := hashPassword(password)
	if usr.Password != passwordHash {
		log.Warn().Str("username", username).Msg("invalid password")
		return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "https://delicias-da-lu-service.com/docs/errors/invalid-credentials",
			Title:      "Invalid Credentials",
			Detail:     "The provided username or password is incorrect",
			HTTPStatus: http.StatusUnauthorized,
			Instance:   "localhost:8080/v1/auth/login",
			Severity:   problemdetails.Err,
		})
	}

	// Update last login
	if err := a.userRepository.UpdateLastLogin(ctx, usr.ID); err != nil {
		log.Error().Err(err).Msg("failed to update last login")
	}

	// Generate JWT token
	token, err := a.generateToken(usr.ID)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate JWT token")
		return nil, err
	}

	// Remove password from response
	usr.Password = ""

	return &user.LoginResponse{
		Token: token,
		User:  *usr,
	}, nil
}

func (a authUseCaseImpl) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	// Validate and extract claims from existing token
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return "", problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "https://delicias-da-lu-service.com/docs/errors/invalid-token",
			Title:      "Invalid Token",
			Detail:     "The provided token is invalid or expired",
			HTTPStatus: http.StatusUnauthorized,
			Instance:   "localhost:8080/v1/auth/refresh",
			Severity:   problemdetails.Err,
		})
	}

	// Extract user ID and generate new token
	userID := (*claims)["sub"].(string)
	newToken, err := a.generateToken(userID)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

func (a authUseCaseImpl) ValidateToken(ctx context.Context, tokenString string) (string, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return "", problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "https://delicias-da-lu-service.com/docs/errors/invalid-token",
			Title:      "Invalid Token",
			Detail:     "The provided token is invalid or expired",
			HTTPStatus: http.StatusUnauthorized,
			Instance:   "localhost:8080/v1/auth",
			Severity:   problemdetails.Err,
		})
	}

	userID := (*claims)["sub"].(string)
	return userID, nil
}

func (a authUseCaseImpl) generateToken(userID string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.AddDate(0, 0, 7).Unix(), // 7 days expiration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
