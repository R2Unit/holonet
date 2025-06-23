package tables

import (
	"github.com/holonet/core/database"
)

var netboxCredentialsTable = database.TableMigration{
	Name: "netbox_credentials",
	Columns: map[string]string{
		"id":               "SERIAL PRIMARY KEY",
		"user_id":          "INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE",
		"netbox_username":  "VARCHAR(255) NOT NULL",
		"netbox_password":  "VARCHAR(255) NOT NULL",
		"netbox_token":     "VARCHAR(255) NOT NULL",
		"netbox_group":     "VARCHAR(255) NOT NULL",
		"netbox_host":      "VARCHAR(255) NOT NULL",
		"is_encrypted":     "BOOLEAN NOT NULL DEFAULT TRUE",
		"last_verified_at": "TIMESTAMP",
		"created_at":       "TIMESTAMP NOT NULL DEFAULT NOW()",
		"updated_at":       "TIMESTAMP NOT NULL DEFAULT NOW()",
		"deleted_at":       "TIMESTAMP NOT NULL DEFAULT NOW()",
	},
	Priority: 4,
}

func init() {
	database.RegisterTable(netboxCredentialsTable)
}
