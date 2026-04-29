package auth

import (
	"context"
	"controle_financas/internal/core/config"
	e "controle_financas/internal/core/domain/errors"
	"controle_financas/internal/core/security"
	"controle_financas/internal/core/validator"
	"controle_financas/internal/features/user"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AccessTokenExpiration  = 15 * time.Minute
	RefreshTokenExpiration = 7 * 24 * time.Hour
	TokenIssuer            = "controle-financas-api"
	TokenAudience          = "controle-financas-clients"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type TokenClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	IsAtivo  bool      `json:"is_ativo"`
	Type     TokenType `json:"type"`
	jwt.RegisteredClaims
}

type authService struct {
	userService user.UserService
	config      config.Config
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, userID uuid.UUID, err error)
	ExtractAuthenticatedUser(tokenString string) (security.UserDetails, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	ValidateToken(tokenString string, expectedType TokenType) (*TokenClaims, error)
	GetPublicKey() *rsa.PublicKey
}

func NewService(
	userService user.UserService,
	config config.Config,
) (AuthService, error) {
	service := &authService{
		userService: userService,
		config:      config,
	}

	if err := service.loadKeys(); err != nil {
		return nil, fmt.Errorf("failed to load RSA keys: %w", err)
	}

	return service, nil
}

func (s *authService) loadKeys() error {
	privateKey, err := loadRSAPrivateKey(s.config.Security.PrivateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load private key: %w", err)
	}
	s.privateKey = privateKey

	publicKey, err := loadRSAPublicKey(s.config.Security.PublicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load public key: %w", err)
	}
	s.publicKey = publicKey

	return nil
}

func loadRSAPrivateKey(path string) (*rsa.PrivateKey, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("private key file not found: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not RSA")
	}

	return rsaKey, nil
}

func loadRSAPublicKey(path string) (*rsa.PublicKey, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("public key file not found: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not RSA")
	}

	return rsaKey, nil
}

func (s *authService) Login(
	ctx context.Context,
	email, password string,
) (string, string, uuid.UUID, error) {
	v := validator.New()
	user.ValidatePasswordPlaintext(v, password)
	if !v.Valid() {
		return "", "", uuid.Nil, e.NewValidationError(v.Errors)
	}

	user, err := s.userService.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, e.ErrRecordNotFound) {
			return "", "", uuid.Nil, e.ErrInvalidCredentials
		}
		return "", "", uuid.Nil, err
	}

	if !user.IsAtivo {
		return "", "", uuid.Nil, e.ErrInactiveAccount
	}

	match, err := user.Senha.Matches(password)
	if err != nil {
		return "", "", uuid.Nil, err
	}

	if !match {
		return "", "", uuid.Nil, e.ErrInvalidCredentials
	}

	accessToken, err := s.createToken(user, TokenTypeAccess, AccessTokenExpiration)
	if err != nil {
		return "", "", uuid.Nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := s.createToken(user, TokenTypeRefresh, RefreshTokenExpiration)
	if err != nil {
		return "", "", uuid.Nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return accessToken, refreshToken, user.ID, nil
}

func (s *authService) ExtractAuthenticatedUser(
	tokenString string,
) (security.UserDetails, error) {
	claims, err := s.ValidateToken(tokenString, TokenTypeAccess)
	if err != nil {
		return nil, err
	}

	return security.NewAuthenticatedUser(
		claims.UserID,
		claims.Username,
		claims.IsAtivo,
	), nil
}

func (s *authService) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (string, error) {
	claims, err := s.ValidateToken(refreshToken, TokenTypeRefresh)
	if err != nil {
		return "", err
	}

	user, err := s.userService.FindByEmail(ctx, claims.Username)
	if err != nil {
		if errors.Is(err, e.ErrRecordNotFound) {
			return "", e.ErrInvalidCredentials
		}
		return "", err
	}

	if !user.IsAtivo {
		return "", e.ErrInactiveAccount
	}

	return s.createToken(user, TokenTypeAccess, AccessTokenExpiration)
}

func (s *authService) ValidateToken(
	tokenString string,
	expectedType TokenType,
) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&TokenClaims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return s.publicKey, nil
		},
		jwt.WithIssuer(TokenIssuer),
		jwt.WithAudience(TokenAudience),
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, e.ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenInvalidIssuer) || errors.Is(err, jwt.ErrTokenInvalidAudience) {
			return nil, e.ErrInvalidTokenClaims
		}

		if strings.Contains(err.Error(), "token has invalid claims") {
			return nil, e.ErrInvalidTokenClaims
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, e.ErrInvalidCredentials
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, e.ErrInvalidTokenClaims
	}

	if claims.Type != expectedType {
		return nil, e.ErrInvalidTokenType
	}

	if !claims.IsAtivo {
		return nil, e.ErrInactiveAccount
	}

	return claims, nil
}

func (s *authService) createToken(
	user security.UserDetails,
	tokenType TokenType,
	expiration time.Duration,
) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		UserID:   user.GetID(),
		Username: user.GetUsername(),
		IsAtivo:  user.GetIsAtivo(),
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    TokenIssuer,
			Audience:  jwt.ClaimStrings{TokenAudience},
			Subject:   user.GetID().String(),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *authService) GetPublicKey() *rsa.PublicKey {
	return s.publicKey
}
