package tables

import "github.com/holonet/core/database"

var tokensTable = database.TableMigration{
	Name: "tokens",
	Columns: map[string]string{
		"id":                   "SERIAL PRIMARY KEY",
		"user_id":              "INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE",
		"token":                "VARCHAR(255) NOT NULL UNIQUE",
		"policy_id":            "INTEGER REFERENCES token_policies(id)",
		"expires_at":           "TIMESTAMP NOT NULL",
		"request_count":        "INTEGER NOT NULL DEFAULT 0",
		"last_request_at":      "TIMESTAMP",
		"requests_today":       "INTEGER NOT NULL DEFAULT 0",
		"requests_today_reset": "TIMESTAMP",
		"created_at":           "TIMESTAMP NOT NULL DEFAULT NOW()",
		"updated_at":           "TIMESTAMP NOT NULL DEFAULT NOW()",
		"deleted_at":           "TIMESTAMP NOT NULL DEFAULT NOW()",
	},
	Priority: 3,
}

func init() {
	database.RegisterTable(tokensTable)
}
