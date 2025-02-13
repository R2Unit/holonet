package auth

import "log"

// UpdateUserGroup updates the user's group (e.g., admin, user, etc.).
func (a *AuthService) UpdateUserGroup(username, user_group string) error {
	_, err := a.DB.Exec("UPDATE users SET user_group = $1 WHERE username = $2", user_group, username)
	if err != nil {
		log.Println("Error updating user user_group:", err)
	}
	return err
}
