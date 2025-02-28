package utils

import (
    "errors"
    "fmt"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/joho/godotenv"
)

var secretKey []byte
var refreshSecretKey []byte

// Load environment variables
func init() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        fmt.Println("Warning: No .env file found, using system environment variables.")
    }

    // Retrieve and validate JWT secrets
    secret := os.Getenv("JWT_SECRET")
    refreshSecret := os.Getenv("JWT_REFRESH_SECRET")

    if secret == "" || refreshSecret == "" {
        panic("ERROR: JWT_SECRET and JWT_REFRESH_SECRET must be set")
    }

    secretKey = []byte(secret)
    refreshSecretKey = []byte(refreshSecret)
}

// GenerateJWT generates a short-lived access token
func GenerateJWT(userID uint, role string) (string, int64, error) {
    expiration := time.Now().Add(time.Hour * 24).Unix() // Access token expires in 24 hours
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":   userID,
        "role":      role,
        "exp":       expiration,
    })

    tokenString, err := token.SignedString(secretKey)
    if err != nil {
        return "", 0, err
    }

    return tokenString, expiration, nil
}

// GenerateRefreshToken generates a long-lived refresh token
func GenerateRefreshToken(userID uint,role string) (string, int64, error) {
    expiration := time.Now().Add(time.Hour * 24 * 365).Unix() // Refresh token expires in 1 year
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":   userID,
        "role":      role,
        "exp":       expiration,
    })

    tokenString, err := token.SignedString(refreshSecretKey)
    if err != nil {
        return "", 0, err
    }

    return tokenString, expiration, nil
}

// ValidateToken validates access and refresh tokens against JWT_SECRET and JWT_REFRESH_SECRET
// and checks if the token is expired.
func ValidateToken(tokenString string, isRefreshToken bool) (*jwt.Token, jwt.MapClaims, error) {
    var key []byte
    if isRefreshToken {
        key = refreshSecretKey
    } else {
        key = secretKey
    }

    // Parse the token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return key, nil
    })

    // Handle parsing errors
    if err != nil {
        // Check for token expiration
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, nil, errors.New("token has expired")
        }
        return nil, nil, fmt.Errorf("failed to parse token: %v", err)
    }

    // Extract claims and check expiration manually
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return nil, nil, errors.New("invalid token")
    }

    // Manual expiration check
    if exp, ok := claims["exp"].(float64); ok {
        if time.Now().Unix() > int64(exp) {
            return nil, nil, errors.New("token has expired")
        }
    }

    return token, claims, nil
}