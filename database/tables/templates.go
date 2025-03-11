package tables

const TemplatesTable = `
	CREATE TABLE templates (
		ID SERIAL PRIMARY KEY,
		name VARCHAR(225) NOT NULL,
		description TEXT,
		content TEXT,
        created_by INT REFERENCES users(id),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	)
`
