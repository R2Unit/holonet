package tables

// Imports the "database" package from "github.com/holonet/core" for database-related operations.
import "github.com/holonet/core/database"

// exampleTable defines the schema for the "example" table including its name and columns with respective data types and constraints.
var exampleTable = database.TableMigration{
	Name: "example",
	Columns: map[string]string{
		"id":             "SERIAL PRIMARY KEY",
		"username":       "VARCHAR(50) NOT NULL",
		"email":          "VARCHAR(100) NOT NULL",
		"password":       "<PASSWORD>",
		"created_at":     "TIMESTAMP NOT NULL",
		"updated_at":     "TIMESTAMP NOT NULL",
		"deleted_at":     "TIMESTAMP",
		"deleted":        "BOOLEAN NOT NULL DEFAULT FALSE",
		"deleted_by":     "VARCHAR(50)",
		"deleted_reason": "VARCHAR(255)",
	},
}

// init initializes the application by registering the exampleTable schema with the database migration system.
func init() {
	database.RegisterTable(exampleTable)
}
