package tables

const CircuitTables = `
	CREATE TABLE customers (
		ID SERIAL PRIMARY KEY,
		provider VARCHAR
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
`
