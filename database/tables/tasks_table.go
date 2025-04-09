package tables

import "github.com/holonet/core/database"

var tasksTable = database.TableMigration{
	Name: "tasks",
	Columns: map[string]string{
		"id":         "SERIAL PRIMARY KEY",
		"user_id":    "INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE",
		"task_type":  "VARCHAR(255) NOT NULL",
		"task_value": "VARCHAR(255) NOT NULL",
		"created_at": "TIMESTAMP NOT NULL",
	},
	Priority: 3,
}

func init() {
	database.RegisterTable(tasksTable)
}
