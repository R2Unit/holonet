package tables

import "github.com/holonet/core/database"

var tokensTable = database.TableMigration{
	Name: "tokens",
	Columns: map[string]string{
		"id":         "SERIAL PRIMARY KEY",
		"user_id":    "INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE",
		"token":      "VARCHAR(255) NOT NULL UNIQUE",
		"expires_at": "TIMESTAMP NOT NULL",
		"created_at": "TIMESTAMP NOT NULL DEFAULT NOW()",
		"updated_at": "TIMESTAMP NOT NULL DEFAULT NOW()",
		"deleted_at": "TIMESTAMP  NOT NULL DEFAULT NOW()",
	},
	Priority: 2,
}

func init() {
	database.RegisterTable(tokensTable)
}
