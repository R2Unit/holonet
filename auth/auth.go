package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/quanza/talos-core/config"
	"github.com/r2unit/colours"

	_ "github.com/lib/pq"
)

type AuthService struct {
	DB *sql.DB
}

type User struct {
	ID       int
	Username string
	Password string
	Salt     string
}

func InitDB() (*sql.DB, error) {
	cfg := config.LoadConfig()

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(colours.Danger("Error opening database:"), err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Println("Error connecting to database:", err)
		return nil, err
	}

	log.Println(colours.Success("Database, Pool AUTH connection established"))

	ensureDefaultAdmin(db)

	return db, nil
}

func generateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

func generateRandomPassword() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password := make([]byte, 50)
	for i := range password {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[n.Int64()]
	}
	return string(password), nil
}

func hashPassword(password, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(password + salt))
	return hex.EncodeToString(hash.Sum(nil))
}
