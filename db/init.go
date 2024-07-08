package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitializeDatabase() (*sqlx.DB, error) {
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	dbName := os.Getenv("DATABASE_NAME")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	sslMode := os.Getenv("DATABASE_SSL_MODE")

	//postgresql://[user[:password]@][netloc][:port][/dbname]

	fmt.Printf("Establishing connection to db: %s on %s:%s\n", dbName, host, port)
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbName, sslMode)
	connection, err := sqlx.Connect("postgres", connectionString)

	if err != nil {
		fmt.Printf("Error connecting to database - %v\n", err)
		return nil, err
	}
	fmt.Println("Initializing database schema.")
	runInitializationScripts(connection)
	fmt.Println("Database initialized")

	return connection, nil
}

func runInitializationScripts(connection *sqlx.DB) error {
	scriptsDirectory := os.Getenv("DATABSE_SCRIPTS_DIRECTORY")
	scripts, err := os.ReadDir(scriptsDirectory)
	if err != nil {
		fmt.Printf("Could not read scripts directory %v\n", err)
	}

	transaction := connection.MustBegin()
	for _, script := range scripts {
		filePath := scriptsDirectory + string(os.PathSeparator) + script.Name()
		file, err := os.Open(filePath)

		if err != nil {
			fmt.Printf("Error reading script %s\n", filePath)
			return err
		}

		defer file.Close()

		content, err := os.ReadFile(filePath)

		if err != nil {
			fmt.Printf("Error reading file content %s - %v", content, err)
			return err
		}

		fmt.Printf("Executing script: %s\n", filePath)

		transaction.MustExec(string(content))

	}
	err = transaction.Commit()
	if err != nil {
		fmt.Printf("Error initializing DB - %v", err)
		return err
	}
	return nil
}
