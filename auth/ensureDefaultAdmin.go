package auth

import (
	"database/sql"
	"log"

	"github.com/r2unit/go-colours"
)

func ensureDefaultAdmin(db *sql.DB) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = 'admin')").Scan(&exists)
	if err != nil {
		log.Println("Error checking for admin user:", err)
		return
	}

	if !exists {
		password, err := generateRandomPassword()
		if err != nil {
			log.Println("Error generating random password:", err)
			return
		}

		salt, err := generateSalt()
		if err != nil {
			log.Println("Error generating salt:", err)
			return
		}

		hashedPassword := hashPassword(password, salt)

		_, err = db.Exec(
			"INSERT INTO users (username, password, salt, user_group) VALUES ($1, $2, $3, $4)",
			"admin",
			hashedPassword,
			salt,
			"admin",
		)
		if err != nil {
			log.Println("Error creating admin user:", err)
			return
		}
		log.Println("Default admin user created with password:", password)
	} else {
		log.Println(colours.Success("Admin user already exists"))
	}
}
