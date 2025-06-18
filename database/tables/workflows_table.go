package tables

import "github.com/holonet/core/database"

var workflowsTable = database.TableMigration{
	Name: "workflows",
	Columns: map[string]string{
		"id":          "SERIAL PRIMARY KEY",
		"name":        "VARCHAR(255) NOT NULL",
		"description": "TEXT",
		"code":        "TEXT NOT NULL",
		"status":      "VARCHAR(50) NOT NULL",
		"created_at":  "TIMESTAMP NOT NULL DEFAULT NOW()",
		"updated_at":  "TIMESTAMP NOT NULL DEFAULT NOW()",
	},
	Priority: 5,
}

var workflowExecutionsTable = database.TableMigration{
	Name: "workflow_executions",
	Columns: map[string]string{
		"id":            "SERIAL PRIMARY KEY",
		"workflow_id":   "INTEGER NOT NULL",
		"status":        "VARCHAR(50) NOT NULL",
		"parameters":    "JSONB",
		"result":        "JSONB",
		"error_message": "TEXT",
		"scheduled_at":  "TIMESTAMP NOT NULL",
		"started_at":    "TIMESTAMP",
		"completed_at":  "TIMESTAMP",
		"created_at":    "TIMESTAMP NOT NULL DEFAULT NOW()",
		"updated_at":    "TIMESTAMP NOT NULL DEFAULT NOW()",
	},
	ForeignKeys: map[string]string{
		"workflow_id": "REFERENCES workflows(id)",
	},
	Priority: 6,
}

func init() {
	database.RegisterTable(workflowsTable)
	database.RegisterTable(workflowExecutionsTable)
}
