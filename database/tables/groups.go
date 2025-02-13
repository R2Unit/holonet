package tables

const GroupsTables = `
	CREATE TABLE groups (
		ID SERIAL PRIMARY KEY,
		user_id VARCHAR(255),
		group_name VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
`
