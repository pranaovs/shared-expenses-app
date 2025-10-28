package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var nameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z .'\-]{0,63}[a-zA-Z]$`)

// ValidateName validates a user's name.
func ValidateName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("name is empty")
	}
	if !nameRegex.MatchString(name) {
		return "", errors.New("invalid name")
	}
	return name, nil
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail validates and normalizes an email.
// Returns a cleaned, lowercase email string or an error.
func ValidateEmail(email string) (string, error) {
	// Email validation regex
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	if email == "" {
		return "", errors.New("email is empty")
	}

	if !emailRegex.MatchString(email) {
		return "", errors.New("invalid email format")
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", errors.New("invalid email syntax")
	}

	return addr.Address, nil
}

// Passwords

// HashPassword hashes a plaintext password using bcrypt.
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("empty password provided")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword compares a plaintext password with its hashed version.
func CheckPassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

func randB64() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("failed to generate random bytes for JWT secret:", err)
	}

	return base64.StdEncoding.EncodeToString(b)
}

var jwtSecret = []byte(Getenv("JWT_SECRET", randB64()))

func GenerateJWT(userID string) (string, error) {
	expiryHours, _ := strconv.Atoi(Getenv("JWT_EXPIRY", "24"))
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(expiryHours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
