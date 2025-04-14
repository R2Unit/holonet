package tables

import "github.com/holonet/core/database"

var permissionsTable = database.TableMigration{
	Name: "group_permissions",
	Columns: map[string]string{
		"id":          "SERIAL PRIMARY KEY",
		"group_id":    "INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE",
		"permission":  "VARCHAR(255) NOT NULL",
		"description": "VARCHAR(255)",
		"is_active":   "BOOLEAN NOT NULL DEFAULT TRUE",
		// <lorenzo> Divider for my eyes only 0_0
		"created_at": "TIMESTAMP NOT NULL DEFAULT NOW()",
		"updated_at": "TIMESTAMP NOT NULL DEFAULT NOW()",
		"deleted_at": "TIMESTAMP  NOT NULL DEFAULT NOW()",
	},
	Priority: 3,
}

func init() {
	database.RegisterTable(permissionsTable)
}
