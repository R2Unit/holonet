package tables

import "github.com/holonet/core/database"

var tokenPoliciesTable = database.TableMigration{
	Name: "token_policies",
	Columns: map[string]string{
		"id":                   "SERIAL PRIMARY KEY",
		"name":                 "VARCHAR(255) NOT NULL UNIQUE",
		"description":          "TEXT",
		"rate_limit_per_min":   "INTEGER NOT NULL DEFAULT 60",   // Default: 60 requests per minute
		"max_requests_per_day": "INTEGER NOT NULL DEFAULT 1000", // Default: 1000 requests per day
		"active":               "BOOLEAN NOT NULL DEFAULT TRUE",
		"created_at":           "TIMESTAMP NOT NULL DEFAULT NOW()",
		"updated_at":           "TIMESTAMP NOT NULL DEFAULT NOW()",
	},
	Priority: 2,
}

func init() {
	database.RegisterTable(tokenPoliciesTable)
}
