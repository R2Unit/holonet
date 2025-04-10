package tables

import "github.com/holonet/core/database"

var groupsTable = database.TableMigration{
	Name: "groups",
	Columns: map[string]string{
		"id":          "SERIAL PRIMARY KEY",
		"user_id":     "INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE",
		"name":        "VARCHAR(255) NOT NULL UNIQUE",
		"description": "VARCHAR(255)",
		"permissions": "TEXT NOT NULL",
		"created_at":  "TIMESTAMP NOT NULL",
		"updated_at":  "TIMESTAMP",
		"deleted_at":  "TIMESTAMP",
	},
	Priority: 3,
}

func init() {
	database.RegisterTable(groupsTable)
}
