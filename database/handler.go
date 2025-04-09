package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type DBHandler struct {
	DB *sql.DB
}

type TableMigration struct {
	Name    string
	Columns map[string]string
}

var tableMigrations []TableMigration

func RegisterTable(tm TableMigration) {
	tableMigrations = append(tableMigrations, tm)
}

func RegisteredTableCount() int {
	return len(tableMigrations)
}

// NewDBHandler initializes a new DBHandler by setting up a connection to the database using environment variables.
// It attempts to connect with retries and returns an error if the connection cannot be established.
func NewDBHandler() (*DBHandler, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "holonet"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "insecure"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "holonet"
	}
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	log.Printf("Attempting to open a connection to the database at %s:%s...", host, port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	const maxRetries = 5
	var attempt int
	for attempt = 1; attempt <= maxRetries; attempt++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("Attempt %d: failed to connect to database: %v", attempt, err)
		if attempt < maxRetries {
			time.Sleep(5 * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	log.Printf("Connected to DB at %s:%s, using database: %s", host, port, dbname)
	return &DBHandler{DB: db}, nil
}

// EnsureTable ensures the specified table exists in the database with the required schema and adds missing columns if needed.
func (handler *DBHandler) EnsureTable(tableName string, columns map[string]string) error {
	log.Printf("Ensuring table '%s' exists and has the required schema...", tableName)

	checkTableQuery := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = $1
		);`
	var exists bool
	log.Printf("Executing query to check if table '%s' exists.", tableName)
	if err := handler.DB.QueryRow(checkTableQuery, tableName).Scan(&exists); err != nil {
		return fmt.Errorf("error checking existence of table %s: %w", tableName, err)
	}

	if !exists {
		log.Printf("Table '%s' does not exist. Preparing to create table...", tableName)
		colDefs := ""
		first := true
		for col, def := range columns {
			if !first {
				colDefs += ", "
			}
			colDefs += fmt.Sprintf("%s %s", col, def)
			first = false
		}
		createTableQuery := fmt.Sprintf("CREATE TABLE %s (%s);", tableName, colDefs)
		log.Printf("Built CREATE TABLE query: %s", createTableQuery)
		log.Printf("Executing CREATE TABLE query for '%s'...", tableName)
		if _, err := handler.DB.Exec(createTableQuery); err != nil {
			return fmt.Errorf("error creating table %s: %w", tableName, err)
		}
		log.Printf("Table '%s' created successfully.", tableName)
	} else {
		log.Printf("Table '%s' already exists. Verifying schema for missing columns...", tableName)
		columnQuery := `
			SELECT column_name 
			FROM information_schema.columns 
			WHERE table_schema = 'public' AND table_name = $1;`
		log.Printf("Executing query to retrieve existing columns for table '%s'.", tableName)
		rows, err := handler.DB.Query(columnQuery, tableName)
		if err != nil {
			return fmt.Errorf("error retrieving columns for table %s: %w", tableName, err)
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {

			}
		}(rows)

		existingColumns := make(map[string]bool)
		for rows.Next() {
			var colName string
			if err := rows.Scan(&colName); err != nil {
				return fmt.Errorf("error scanning column name for table %s: %w", tableName, err)
			}
			existingColumns[colName] = true
			log.Printf("Found existing column '%s' in table '%s'.", colName, tableName)
		}

		for col, def := range columns {
			if !existingColumns[col] {
				alterTableQuery := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", tableName, col, def)
				log.Printf("Column '%s' is missing in table '%s'. Built ALTER TABLE query: %s", col, tableName, alterTableQuery)
				log.Printf("Executing ALTER TABLE query for table '%s'...", tableName)
				if _, err := handler.DB.Exec(alterTableQuery); err != nil {
					return fmt.Errorf("error adding column %s to table %s: %w", col, tableName, err)
				}
				log.Printf("Added missing column '%s' to table '%s'.", col, tableName)
			} else {
				log.Printf("Column '%s' exists in table '%s'. No action needed.", col, tableName)
			}
		}
		log.Printf("Table '%s' schema verified and up-to-date.", tableName)
	}

	return nil
}

// Migrate performs the database migration by ensuring all registered tables exist and are updated with the correct schema.
func (handler *DBHandler) Migrate() error {
	log.Printf("Starting database migration for %d table(s)...", len(tableMigrations))
	for _, tm := range tableMigrations {
		log.Printf("Migrating table: %s", tm.Name)
		if err := handler.EnsureTable(tm.Name, tm.Columns); err != nil {
			log.Printf("Migration error in table '%s': %v", tm.Name, err)
			return err
		}
		log.Printf("Finished migrating table: %s", tm.Name)
	}
	log.Printf("Database migration completed successfully.")
	return nil
}
