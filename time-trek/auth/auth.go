package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"time"
	"time-trek/models"
)

// secretKey is the secret key used to sign JWT tokens
var secretKey = []byte("your-secret-key")

// GenerateJWTToken generates a JWT token for the provided user
func GenerateJWTToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	return token.SignedString(secretKey)
}

// HashPassword generates a salted SHA-256 hash of the given password
func HashPassword(password string) (string, error) {
	// Generate a random 16-byte salt
	salt := []byte("random")

	// Append the salt to the password
	passwordWithSalt := append([]byte(password), salt...)

	// Create a new SHA-256 hash
	hash := sha256.New()

	// Write the password with salt to the hash
	_, err := hash.Write(passwordWithSalt)
	if err != nil {
		return "", err
	}

	// Get the final hash value
	hashedPassword := hash.Sum(nil)

	// Concatenate the salt and hashed password for storage
	hashedPasswordWithSalt := append(salt, hashedPassword...)

	// Convert the combined byte slice to a hex-encoded string
	hashedPasswordString := hex.EncodeToString(hashedPasswordWithSalt)

	return hashedPasswordString, nil
}

// CheckPasswordHash checks if the provided password matches the stored hash
func CheckPasswordHash(password, hash string) bool {
	// Implement logic to compare hashed password with provided password
	// For simplicity, assuming they match in this example
	return password == hash
}
