package tables

import "github.com/holonet/core/database"

var jobsTable = database.TableMigration{
	Name: "jobs",
	Columns: map[string]string{
		"id":           "SERIAL PRIMARY KEY",
		"queue_name":   "VARCHAR(255) NOT NULL",
		"payload":      "JSONB NOT NULL",
		"description":  "VARCHAR(255) NOT NULL",
		"status":       "VARCHAR(20) NOT NULL",
		"priority":     "INTEGER NOT NULL",
		"available_at": "TIMESTAMP NOT NULL",
		"attempts":     "INTEGER NOT NULL DEFAULT 0",
		"max_attempts": "INTEGER NOT NULL DEFAULT 1",
		"last_error":   "TEXT",
		"locked_by":    "VARCHAR(255)",
		"locked_at":    "TIMESTAMP",
		"created_at":   "TIMESTAMP NOT NULL DEFAULT NOW()",
		"updated_at":   "TIMESTAMP NOT NULL DEFAULT NOW()",
		"deleted_at":   "TIMESTAMP  NOT NULL DEFAULT NOW()",
	},
	Priority: 2,
}

func init() {
	database.RegisterTable(jobsTable)
}
