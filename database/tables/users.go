package tables

const UsersTable = `
	CREATE TABLE users (
		ID SERIAL PRIMARY KEY,
		username VARCHAR(225),
		firstname VARCHAR(225),
		lastname VARCHAR(225),
		password VARCHAR(225),
		salt VARCHAR(225),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
`
