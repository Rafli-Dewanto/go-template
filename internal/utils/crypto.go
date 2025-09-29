package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

// CryptoConfig holds configuration for crypto operations
type CryptoConfig struct {
	JWTSecretKey string
	AESKey       string
	ScryptN      int
	ScryptR      int
	ScryptP      int
	ScryptKeyLen int
	BcryptCost   int
}

// DefaultCryptoConfig returns a default crypto configuration
func DefaultCryptoConfig() *CryptoConfig {
	return &CryptoConfig{
		ScryptN:      32768,
		ScryptR:      8,
		ScryptP:      1,
		ScryptKeyLen: 32,
		BcryptCost:   bcrypt.DefaultCost,
	}
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	if IsEmpty(password) {
		return "", errors.New("password cannot be empty")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// HashPasswordWithCost hashes a password using bcrypt with custom cost
func HashPasswordWithCost(password string, cost int) (string, error) {
	if IsEmpty(password) {
		return "", errors.New("password cannot be empty")
	}

	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return "", fmt.Errorf("invalid bcrypt cost: %d (must be between %d and %d)", cost, bcrypt.MinCost, bcrypt.MaxCost)
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// VerifyPassword verifies a password against its hash using bcrypt
func VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GenerateSalt generates a random salt for cryptographic operations
func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

// HashWithSalt creates a SHA256 hash with salt
func HashWithSalt(data, salt []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	hasher.Write(salt)
	return hex.EncodeToString(hasher.Sum(nil))
}

// HashSHA256 creates a SHA256 hash of the input string
func HashSHA256(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// DeriveKey derives a key from password using scrypt
func DeriveKey(password, salt []byte, config *CryptoConfig) ([]byte, error) {
	if config == nil {
		config = DefaultCryptoConfig()
	}

	key, err := scrypt.Key(password, salt, config.ScryptN, config.ScryptR, config.ScryptP, config.ScryptKeyLen)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	return key, nil
}

// EncryptAES encrypts data using AES-GCM
func EncryptAES(plaintext []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptAES decrypts data using AES-GCM
func DecryptAES(ciphertext []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptString encrypts a string and returns base64 encoded result
func EncryptString(plaintext string, key []byte) (string, error) {
	encrypted, err := EncryptAES([]byte(plaintext), key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DecryptString decrypts a base64 encoded string
func DecryptString(encryptedBase64 string, key []byte) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	decrypted, err := DecryptAES(encrypted, key)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role,omitempty"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT token with the given claims
func GenerateJWT(userID, email, role, secretKey string, duration time.Duration) (string, error) {
	if IsEmpty(secretKey) {
		return "", errors.New("secret key cannot be empty")
	}

	now := time.Now()
	claims := JWTClaims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		IssuedAt:  now,
		ExpiresAt: now.Add(duration),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString, secretKey string) (*JWTClaims, error) {
	if IsEmpty(tokenString) || IsEmpty(secretKey) {
		return nil, errors.New("token and secret key cannot be empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Additional validation for expiry
	if time.Now().After(claims.ExpiresAt) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

// GenerateRandomKey generates a random key of specified length
func GenerateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}
	return key, nil
}

// GenerateRandomString generates a random string of specified length using base64 encoding
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// SecureCompare performs constant-time comparison of two byte slices
func SecureCompare(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}

// SecureCompareString performs constant-time comparison of two strings
func SecureCompareString(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// GenerateAPIKey generates a secure API key
func GenerateAPIKey() (string, error) {
	key, err := GenerateRandomKey(32)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(key), nil
}

// HashAPIKey creates a hash of an API key for storage
func HashAPIKey(apiKey string) string {
	return HashSHA256(apiKey)
}

// VerifyAPIKey verifies an API key against its hash
func VerifyAPIKey(apiKey, hashedKey string) bool {
	return SecureCompareString(HashSHA256(apiKey), hashedKey)
}

// TokenPair represents an access and refresh token pair
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// GenerateTokenPair generates both access and refresh tokens
func GenerateTokenPair(userID, email, role, secretKey string, accessDuration, refreshDuration time.Duration) (*TokenPair, error) {
	accessToken, err := GenerateJWT(userID, email, role, secretKey, accessDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := GenerateJWT(userID, email, "refresh", secretKey, refreshDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessDuration.Seconds()),
	}, nil
}

// GenerateAPIID generates a unique API identifier for request tracking
func GenerateAPIID() string {
	return uuid.New().String()
}

// GenerateShortAPIID generates a shorter API identifier using timestamp and random bytes
func GenerateShortAPIID() string {
	now := time.Now().UnixNano()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	return fmt.Sprintf("api-%d-%x", now, randomBytes)
}

// GenerateCorrelationID generates a correlation ID for distributed tracing
func GenerateCorrelationID() string {
	return fmt.Sprintf("corr-%s", uuid.New().String())
}
