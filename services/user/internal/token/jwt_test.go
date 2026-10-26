package token_test

import (
	"log"
	"os"
	"shopik/services/user/internal/token"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestIssuser_Success(t *testing.T) {
	mustLoad()
	secretStr := os.Getenv("JWT_SECRET")
	secret := []byte(secretStr)
	TTLStr := os.Getenv("JWT_ACCESS_TTL")
	TTL, err := strconv.Atoi(TTLStr)
	assert.NoError(t, err)
	issuer, err := token.NewJWTIssuer(secret, time.Duration(TTL)*time.Minute)
	assert.NoError(t, err)
	tokenJWT, err := issuer.Issuer("4")
	assert.NoError(t, err)

}

func mustLoad() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("err .env loading, using enviroment OS")
	}
	envList := []string{
		"JWT_SECRET",
		"JWT_ACCESS_TTL",
	}
	var envNotExist []string
	for _, envName := range envList {
		if os.Getenv(envName) == "" {
			envNotExist = append(envNotExist, envName)
		}
	}
	if len(envNotExist) > 0 {
		log.Fatalf("env not exist: %v", envNotExist)
	}
}
