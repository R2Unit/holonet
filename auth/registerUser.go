package auth

import "log"

func (a *AuthService) RegisterUser(username, password, user_group string) error {
	salt, err := generateSalt()
	if err != nil {
		return err
	}

	hashedPassword := hashPassword(password, salt)

	_, err = a.DB.Exec("INSERT INTO users (username, password, salt, user_group) VALUES ($1, $2, $3, $4)",
		username, hashedPassword, salt, user_group)
	if err != nil {
		log.Println("Error registering user:", err)
	}
	return err
}
