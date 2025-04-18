package tables

import "github.com/holonet/core/database"

var tasksTable = database.TableMigration{
	Name: "tasks",
	Columns: map[string]string{
		"id":          "SERIAL PRIMARY KEY",
		"task_name":   "VARCHAR(255) NOT NULL UNIQUE",
		"description": "VARCHAR(255)",
		"task_type":   "VARCHAR(255) NOT NULL",
		"task_value":  "VARCHAR(255) NOT NULL",
		// TODO: Add support for multiple workflow types
		//"workflow_code": "TEXT",
		"workflow_json": "JSONB",
		//"workflow_yaml": "TEXT",
		"priority":   "INTEGER NOT NULL DEFAULT 0",
		"user_id":    "INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE",
		"enabled":    "BOOLEAN NOT NULL DEFAULT TRUE",
		"created_at": "TIMESTAMP NOT NULL DEFAULT NOW()",
		"updated_at": "TIMESTAMP NOT NULL DEFAULT NOW()",
		"deleted_at": "TIMESTAMP NOT NULL DEFAULT NOW()",
	},
	Priority: 3,
}

func init() {
	database.RegisterTable(tasksTable)
}
