package tables

import "github.com/holonet/core/database"

var usersTable = database.TableMigration{
	Name: "users",
	Columns: map[string]string{
		"id":            "SERIAL PRIMARY KEY",
		"username":      "VARCHAR(255) NOT NULL UNIQUE",
		"email":         "VARCHAR(255) NOT NULL UNIQUE",
		"password_hash": "VARCHAR(255) NOT NULL",
		"created_at":    "TIMESTAMP NOT NULL",
		"updated_at":    "TIMESTAMP",
		"deleted_at":    "TIMESTAMP",
	},
	Priority: 1,
}

func init() {
	database.RegisterTable(usersTable)
}
