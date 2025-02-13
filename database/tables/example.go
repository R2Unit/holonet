package tables

const ExampleTable = `
	CREATE TABLE example (
		ID SERIAL PRIMARY KEY,
		name VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
`
