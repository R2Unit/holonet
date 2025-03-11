package tables

const TasksTable = `
	CREATE TABLE tasks (
		ID SERIAL PRIMARY KEY,
        task_id TEXT UNIQUE,
	    command TEXT,
        args TEXT,         -- JSON array of strings
        files JSONB,       -- JSON object mapping filenames to contents
        reporter TEXT,
        hosts TEXT,
        task_template TEXT,
        status TEXT,       -- e.g. 'pending', 'running', 'completed', 'failed'
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
	)
`
