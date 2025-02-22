package auth

import "log"

func (a *AuthService) UpdateUserPassword(username, newPassword string) error {
	salt, err := generateSalt()
	if err != nil {
		return err
	}

	hashedPassword := hashPassword(newPassword, salt)

	_, err = a.DB.Exec("UPDATE users SET password = $1, salt = $2 WHERE username = $3",
		hashedPassword, salt, username)
	if err != nil {
		log.Println("Error updating user password:", err)
	}
	return err
}
