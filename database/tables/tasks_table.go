package tables

import "github.com/holonet/core/database"

var tasksTable = database.TableMigration{
	Name: "tasks",
	Columns: map[string]string{
		"id":            "SERIAL PRIMARY KEY",
		"task_name":     "VARCHAR(255) NOT NULL UNIQUE",
		"description":   "VARCHAR(255)",
		"task_type":     "VARCHAR(255) NOT NULL",
		"task_value":    "VARCHAR(255) NOT NULL",
		"yaml_template": "TEXT NOT NULL",
		"enabled":       "BOOLEAN NOT NULL DEFAULT TRUE",
		// <lorenzo> Divider for my eyes only 0_0
		"created_at": "TIMESTAMP NOT NULL DEFAULT NOW()",
		"updated_at": "TIMESTAMP NOT NULL DEFAULT NOW()",
		"deleted_at": "TIMESTAMP  NOT NULL DEFAULT NOW()",
	},
	Priority: 3,
}

func init() {
	database.RegisterTable(tasksTable)
}
