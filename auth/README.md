# Authentication

## RegisterUser
Create a new user and assign a group.
```golang
package main

import (
	"log"

	"github.com/quanza/talos-core/auth"
)

func main() {
	// Initialize the database
	db, err := auth.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize authentication service
	authService := auth.AuthService{DB: db}

	// Register a new user
	err = authService.RegisterUser("newuser", "securepassword", "default")
	if err != nil {
		log.Println("User registration failed:", err)
	} else {
		log.Println("User registered successfully")
	}
}
```

## UpdateUserGroup
Modify a user's group.
```golang
package main

import (
	"log"

	"github.com/quanza/talos-core/auth"
)

func main() {
	// Initialize the database
	db, err := auth.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize authentication service
	authService := auth.AuthService{DB: db}

	// Update the user's group
	err = authService.UpdateUserGroup("newuser", "admin")
	if err != nil {
		log.Println("Failed to update user group:", err)
	} else {
		log.Println("User group updated successfully")
	}
}
```
## Update a User's Password
```golang
package main

import (
	"fmt"
	"log"

	"github.com/quanza/talos-core/auth"
)

func main() {
	// Initialize the database
	db, err := auth.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize authentication service
	authService := auth.AuthService{DB: db}

	// Update the user's password
	err = authService.UpdateUserPassword("newuser", "newsecurepassword")
	if err != nil {
		log.Println("Failed to update password:", err)
	} else {
		log.Println("Password updated successfully")
	}
}
```

## DeleteUser
Remove a user.
```golang
package main

import (
	"log"

	"github.com/quanza/talos-core/auth"
)

func main() {
	// Initialize the database
	db, err := auth.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize authentication service
	authService := auth.AuthService{DB: db}

	// Delete a user
	err = authService.DeleteUser("newuser")
	if err != nil {
		log.Println("Failed to delete user:", err)
	} else {
		log.Println("User deleted successfully")
	}
}
```

## AuthenticateUser
Validate login credentials.
```golang
package main

import (
	"fmt"
	"log"

	"github.com/quanza/talos-core/auth"
)

func main() {
	// Initialize the database
	db, err := auth.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize authentication service
	authService := auth.AuthService{DB: db}

	// Authenticate a user
	user, err := authService.AuthenticateUser("newuser", "securepassword")
	if err != nil {
		log.Println("Authentication failed:", err)
	} else {
		fmt.Println("Authenticated user:", user.Username, "Group:", user.Group)
	}
}
```