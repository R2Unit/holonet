package auth

import "log"

func (a *AuthService) DeleteUser(username string) error {
	_, err := a.DB.Exec("DELETE FROM users WHERE username = $1", username)
	if err != nil {
		log.Println("Error deleting user:", err)
	}
	return err
}
