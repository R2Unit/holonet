package tables

const WorkersTable = `
	CREATE TABLE workers (
		ID SERIAL PRIMARY KEY,
		worker VARCHAR(225),
		task VARCHAR(225),
		status VARCHAR(225),
		hosts VARCHAR(225),
		task_template VARCHAR(225),
		reporter VARCHAR(225),
		created_at TIMESTAMP DEFAULT NOW()
	)
`
