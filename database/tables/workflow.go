package tables

const WorkflowsTable = `
    CREATE TABLE workflows (
		ID SERIAL PRIMARY KEY,
        CONSTRAINT fk_circuits
          FOREIGN KEY (circuits_id)
          REFERENCES circuits(id)
          ON DELETE CASCADE
        asn VARCHAR(225),
        divicder BOOLEAN,
        product VARCHAR(225),
        version VARCHAR(225),
        reporter VARCHAR(100),
        multicars BOOLEAN,
        insync BOOLEAN,
        status VARCHAR(50),
        end_date DATE,
        created_at TIMESTAMP WITH TIME ZONE
    )
`
