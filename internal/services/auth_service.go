package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"demo/internal/models"
	"demo/internal/repository"
)

var (
	ErrEmailAlreadyTaken  = errors.New("cet email est déjà utilisé")
	ErrInvalidCredentials = errors.New("email ou mot de passe invalide")
)

const tokenTTL = 24 * time.Hour

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: []byte(jwtSecret)}
}

func (s *AuthService) Register(email, password string) (*models.User, error) {
	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return nil, ErrEmailAlreadyTaken
	} else if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash mot de passe: %w", err)
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hash),
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return s.generateToken(user.ID)
}

func (s *AuthService) generateToken(userID uint) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprint(userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("génération token: %w", err)
	}

	return signed, nil
}

// ValidateToken vérifie la signature et l'expiration du token, et retourne l'ID utilisateur qu'il porte.
func (s *AuthService) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode de signature inattendue: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return 0, fmt.Errorf("token invalide: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return 0, errors.New("token invalide")
	}

	var userID uint
	if _, err := fmt.Sscan(claims.Subject, &userID); err != nil {
		return 0, fmt.Errorf("sujet du token invalide: %w", err)
	}

	return userID, nil
}
