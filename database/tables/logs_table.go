package tables

import "github.com/holonet/core/database"

var logsTable = database.TableMigration{
	Name: "logs",
	Columns: map[string]string{
		"id":         "SERIAL PRIMARY KEY",
		"log_level":  "VARCHAR(20) NOT NULL",
		"message":    "TEXT NOT NULL",
		"file":       "VARCHAR(255)",
		"function":   "VARCHAR(255)",
		"line":       "INTEGER",
		"created_at": "TIMESTAMP NOT NULL",
	},
}

func init() {
	database.RegisterTable(logsTable)
}
