package tables

const CustomersTable = `
	CREATE TABLE customers (
		ID SERIAL PRIMARY KEY,
		name VARCHAR(255),
		slug VARCHAR(25),
		description VARCHAR(255),
		customer_number VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
`
