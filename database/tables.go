package database

import "github.com/r2unit/holonet-core/database/tables"

func GetTables() map[string]string {
	return map[string]string{
		"example": tables.ExampleTable,
		"users":   tables.UsersTable,
		"groups":  tables.GroupsTables,
		//	"custmomers": tables.CustomersTable,
		//  "wofkflow": tables.WorkflowsTable,
		"workers": tables.WorkersTable,
	}
}
