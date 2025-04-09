package tables

import "github.com/holonet/core/database"

var usersTable = database.TableMigration{
	Name: "users",
	Columns: map[string]string{
		"id":            "SERIAL PRIMARY KEY",
		"username":      "VARCHAR(255) NOT NULL UNIQUE",
		"email":         "VARCHAR(255) NOT NULL UNIQUE",
		"password_hash": "VARCHAR(255) NOT NULL",
		//"first_name":    "VARCHAR(255)",
		//"last_name":     "VARCHAR(255)",
		"created_at": "TIMESTAMP NOT NULL",
		"updated_at": "TIMESTAMP",
		"deleted_at": "TIMESTAMP",
	},
}

func init() {
	database.RegisterTable(usersTable)
}
