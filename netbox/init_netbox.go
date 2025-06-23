package netbox

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/holonet/core/database"
	"github.com/holonet/core/logger"
)

type Group struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Permissions []int  `json:"permissions,omitempty"`
	URL         string `json:"url,omitempty"`
	Display     string `json:"display,omitempty"`
	UserCount   int    `json:"user_count,omitempty"`
}

type User struct {
	ID           int                      `json:"id,omitempty"`
	Username     string                   `json:"username"`
	Password     string                   `json:"password,omitempty"`
	Email        string                   `json:"email"`
	IsStaff      bool                     `json:"is_staff"`
	IsSuperuser  bool                     `json:"is_superuser"`
	Groups       []int                    `json:"groups,omitempty"`
	URL          string                   `json:"url,omitempty"`
	Display      string                   `json:"display,omitempty"`
	FirstName    string                   `json:"first_name,omitempty"`
	LastName     string                   `json:"last_name,omitempty"`
	DateJoined   string                   `json:"date_joined,omitempty"`
	LastLogin    string                   `json:"last_login,omitempty"`
	IsActive     bool                     `json:"is_active,omitempty"`
	GroupsDetail []map[string]interface{} `json:"groups_detail,omitempty"`
}

func InitNetboxAuth(client *Client, gatekeeper *Gatekeeper, db *sql.DB) error {
	logger.Info("Initializing NetBox authentication...")

	credentials, err := database.GetNetboxCredentials(db, 1)
	if err != nil {
		logger.Error("Failed to check for existing NetBox credentials: %v", err)
	} else if credentials != nil && credentials.NetboxGroup == "Robot" {
		logger.Info("NetBox authentication already initialized with Robot group, skipping")
		return nil
	}

	groups, err := createAuthGroups(gatekeeper)
	if err != nil {
		return fmt.Errorf("failed to create authentication groups: %v", err)
	}

	var superuserGroupID, robotGroupID int
	for _, group := range groups {
		if group.Name == "Superuser" {
			superuserGroupID = group.ID
		} else if group.Name == "Robot" {
			robotGroupID = group.ID
		}
	}

	if superuserGroupID == 0 {
		return fmt.Errorf("failed to find Superuser group ID")
	}

	if robotGroupID == 0 {
		return fmt.Errorf("failed to find Robot group ID")
	}

	username, password, token, err := createRobotUser(gatekeeper, superuserGroupID, robotGroupID)
	if err != nil {
		return fmt.Errorf("failed to create robot.holonet user: %v", err)
	}

	credentials = &database.NetboxCredentials{
		UserID:         1,
		NetboxUsername: username,
		NetboxPassword: password,
		NetboxToken:    token,
		NetboxGroup:    "Robot",
		NetboxHost:     client.Host,
		IsEncrypted:    false,
		LastVerifiedAt: time.Now(),
	}

	if err := database.StoreNetboxCredentials(db, *credentials); err != nil {
		logger.Error("Failed to store NetBox credentials: %v", err)
	} else {
		logger.Info("NetBox credentials stored in database")
	}

	logger.Info("NetBox authentication initialized successfully")
	return nil
}

func createAuthGroups(gatekeeper *Gatekeeper) ([]Group, error) {
	logger.Info("Creating NetBox authentication groups...")

	groups := []Group{
		{
			Name: "Superuser",
		},
		{
			Name: "Operator",
		},
		{
			Name: "Member",
		},
		{
			Name: "Read-Only",
		},
		{
			Name: "Robot",
		},
	}

	existingGroups, err := getExistingGroups(gatekeeper)
	if err != nil {
		return nil, err
	}

	var createdGroups []Group
	for _, group := range groups {
		var existingGroup *Group
		for i := range existingGroups {
			if existingGroups[i].Name == group.Name {
				existingGroup = &existingGroups[i]
				break
			}
		}

		if existingGroup != nil {
			logger.Info("Updating existing group: %s", group.Name)
			updatedGroup, err := updateGroup(gatekeeper, *existingGroup)
			if err != nil {
				return nil, fmt.Errorf("failed to update group %s: %v", group.Name, err)
			}
			createdGroups = append(createdGroups, updatedGroup)
		} else {
			logger.Info("Creating new group: %s", group.Name)
			newGroup, err := createGroup(gatekeeper, group)
			if err != nil {
				return nil, fmt.Errorf("failed to create group %s: %v", group.Name, err)
			}
			createdGroups = append(createdGroups, newGroup)
		}
	}

	logger.Info("Created/updated %d authentication groups", len(createdGroups))
	return createdGroups, nil
}

func getExistingGroups(gatekeeper *Gatekeeper) ([]Group, error) {
	data, err := gatekeeper.Request(http.MethodGet, "users/groups/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing groups: %v", err)
	}

	var response struct {
		Results []Group `json:"results"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		logger.Error("Failed to parse groups response: %v. Response data: %s", err, string(data))
		return nil, fmt.Errorf("failed to parse groups response: %v", err)
	}

	logger.Debug("Retrieved %d existing groups from NetBox", len(response.Results))
	return response.Results, nil
}

func createGroup(gatekeeper *Gatekeeper, group Group) (Group, error) {
	logger.Info("Creating new group in NetBox: %s", group.Name)

	groupData := map[string]interface{}{
		"name": group.Name,
	}
	data, err := gatekeeper.Request(http.MethodPost, "users/groups/", groupData)
	if err != nil {
		logger.Error("Failed to create group %s: %v", group.Name, err)
		return Group{}, fmt.Errorf("failed to create group: %v", err)
	}

	var createdGroup Group
	if err := json.Unmarshal(data, &createdGroup); err != nil {
		logger.Error("Failed to parse created group response: %v. Response data: %s", err, string(data))
		return Group{}, fmt.Errorf("failed to parse created group response: %v", err)
	}

	logger.Info("Successfully created group %s with ID %d", createdGroup.Name, createdGroup.ID)
	return createdGroup, nil
}

func updateGroup(gatekeeper *Gatekeeper, existingGroup Group) (Group, error) {
	logger.Info("Updating existing group in NetBox: %s (ID: %d)", existingGroup.Name, existingGroup.ID)

	groupData := map[string]interface{}{
		"name": existingGroup.Name,
	}

	endpoint := fmt.Sprintf("users/groups/%d/", existingGroup.ID)
	data, err := gatekeeper.Request(http.MethodPatch, endpoint, groupData)
	if err != nil {
		logger.Error("Failed to update group %s (ID: %d): %v", existingGroup.Name, existingGroup.ID, err)
		return Group{}, fmt.Errorf("failed to update group: %v", err)
	}

	var updatedGroup Group
	if err := json.Unmarshal(data, &updatedGroup); err != nil {
		logger.Error("Failed to parse updated group response: %v. Response data: %s", err, string(data))
		return Group{}, fmt.Errorf("failed to parse updated group response: %v", err)
	}

	logger.Info("Successfully updated group %s (ID: %d)", updatedGroup.Name, updatedGroup.ID)
	return updatedGroup, nil
}

func createRobotUser(gatekeeper *Gatekeeper, superuserGroupID, robotGroupID int) (string, string, string, error) {
	logger.Info("Creating robot.holonet user...")

	existingUsers, err := getExistingUsers(gatekeeper)
	if err != nil {
		logger.Error("Failed to get existing users: %v", err)
		return "", "", "", err
	}

	var existingUser *User
	for i := range existingUsers {
		if existingUsers[i].Username == "robot.holonet" {
			existingUser = &existingUsers[i]
			logger.Debug("Found existing robot.holonet user with ID: %d", existingUser.ID)
			break
		}
	}

	password, err := generateRandomPassword(16)
	if err != nil {
		logger.Error("Failed to generate password: %v", err)
		return "", "", "", fmt.Errorf("failed to generate password: %v", err)
	}
	logger.Debug("Generated random password for robot.holonet user")

	token := ""
	if existingUser != nil {
		logger.Info("Updating existing robot.holonet user (ID: %d)", existingUser.ID)
		existingUser.Password = password
		existingUser.IsStaff = true
		existingUser.IsSuperuser = true
		existingUser.Groups = []int{superuserGroupID, robotGroupID}

		endpoint := fmt.Sprintf("users/users/%d/", existingUser.ID)
		data, err := gatekeeper.Request(http.MethodPatch, endpoint, existingUser)
		if err != nil {
			logger.Error("Failed to update robot.holonet user: %v", err)
			return "", "", "", fmt.Errorf("failed to update robot.holonet user: %v", err)
		}

		logger.Debug("User update response: %s", string(data))

		token = gatekeeper.client.Token
	} else {
		logger.Info("Creating new robot.holonet user")
		newUser := User{
			Username:    "robot.holonet",
			Password:    password,
			Email:       "robot.holonet@example.com",
			IsStaff:     true,
			IsSuperuser: true,
			Groups:      []int{superuserGroupID, robotGroupID},
		}

		data, err := gatekeeper.Request(http.MethodPost, "users/users/", newUser)
		if err != nil {
			logger.Error("Failed to create robot.holonet user: %v", err)
			return "", "", "", fmt.Errorf("failed to create robot.holonet user: %v", err)
		}

		logger.Debug("User creation response: %s", string(data))

		token = gatekeeper.client.Token
	}

	logger.Info("robot.holonet user created/updated successfully with password: %s", password)
	return "robot.holonet", password, token, nil
}

func getExistingUsers(gatekeeper *Gatekeeper) ([]User, error) {
	data, err := gatekeeper.Request(http.MethodGet, "users/users/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing users: %v", err)
	}

	var response struct {
		Results []User `json:"results"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		logger.Error("Failed to parse users response: %v. Response data: %s", err, string(data))
		return nil, fmt.Errorf("failed to parse users response: %v", err)
	}

	logger.Debug("Retrieved %d existing users from NetBox", len(response.Results))
	return response.Results, nil
}

func generateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?"
	charsetLength := big.NewInt(int64(len(charset)))

	var password strings.Builder
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		password.WriteByte(charset[randomIndex.Int64()])
	}

	return password.String(), nil
}
