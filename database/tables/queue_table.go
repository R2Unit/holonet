package tables

import "github.com/holonet/core/database"

var queueTable = database.TableMigration{
	Name: "queue",
	Columns: map[string]string{
		"id":           "SERIAL PRIMARY KEY",
		"task_id":      "INTEGER NOT NULL",
		"user_id":      "INTEGER NOT NULL",
		"state":        "VARCHAR(50) NOT NULL",
		"reporter":     "VARCHAR(255)",
		"queued_at":    "TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP",
		"started_at":   "TIMESTAMP",
		"completed_at": "TIMESTAMP",
	},
	ForeignKeys: map[string]string{
		"task_id": "REFERENCES tasks(id)",
		"user_id": "REFERENCES users(id)",
	},
	Priority: 4,
}

func init() {
	database.RegisterTable(queueTable)
}
