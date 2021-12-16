package sqlite

import (
	"database/sql"
	"os"
	"path"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

//
// Client is the database client
type Client struct {
	logger *zap.SugaredLogger
	db     *sql.DB
}

func New(settings map[string]string, subpath string) (*Client, error) {
	logger := zap.S().With("package", "store.sqlite")

	dbPath := path.Join(subpath, "db", settings["database.file"])
	if settings["database.auto_create"] == "true" {

		if _, err := os.Stat(dbPath); err == nil {
			os.Remove(dbPath)
		}
		logger.Debugf("%q database file removed", dbPath)
		file, err := os.Create(dbPath) // Create SQLite file
		if err != nil {
			return nil, err
		}

		logger.Debugf("%q database file created", dbPath)
		file.Close()
	}

	// Connect to database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	maxCon, err := strconv.Atoi(settings["database.max_connections"])
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxCon)

	client := &Client{
		logger: logger,
		db:     db,
	}

	err = client.createTable()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *Client) createTable() error {
	createMessageTableSQL := `CREATE TABLE messages (
		"id" VARCHAR(64) PRIMARY KEY, 
		"text" TEXT NOT NULL, 
		"date" DATETIME NOT NULL
	);` // SQL Statement for Create Table

	c.logger.Debugf("Create Messages table...")
	statement, err := c.db.Prepare(createMessageTableSQL) // Prepare SQL Statement
	if err != nil {
		return err
	}
	_, err = statement.Exec() // Execute SQL Statements
	c.logger.Debugf("Messages table created")
	return err
}
