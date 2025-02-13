package auth

import (
	"database/sql"
	"errors"
)

// AuthenticateUser checks if the given username and password match a user in the database.
func (a *AuthService) AuthenticateUser(username, password string) (*User, error) {
	var user User
	err := a.DB.QueryRow("SELECT id, username, password, salt FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Salt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if hashPassword(password, user.Salt) != user.Password {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}
