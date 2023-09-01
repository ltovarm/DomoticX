package query

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

// JSONResponse is a struct to handle JSON responses containing data.
type JSONResponse struct {
	Data []map[string]interface{} `json:"data"`
}

// User is a structure containing the keys username and password.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Database is a struct to manage the database connection and operations.
type Database struct {
	DB         *sql.DB
	sqlStament *string
}

// NewDb creates a new instance of Database.
func NewDb() *Database {
	return &Database{nil, nil}
}

// ConnectToDatabaseFromEnvVar connects to the database using environment variables.
func (db *Database) ConnectToDatabaseFromEnvVar() error {
	sqlServerURL := os.Getenv("DATABASE_URL")
	// Connect to database
	var err error
	db.DB, err = sql.Open("postgres", sqlServerURL)
	if err != nil {
		panic(err)
	}

	log.Println(" > Connection to database successfully.")

	return err
}

// ConnectToDatabase connects to the database using provided parameters.
func (db *Database) ConnectToDatabase(user, pw, dbName, dbContainerName, port string) error {
	// Connect to database
	sqlServerURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pw, dbContainerName, port, dbName)
	var err error
	db.DB, err = sql.Open("postgres", sqlServerURL)
	if err != nil {
		panic(err)
	}

	log.Printf("Connection to database %s successfully completed.", dbName)

	return err
}

// CloseDatabase closes the database connection.
func (db *Database) CloseDatabase() {
	db.DB.Close()
}

// SendData inserts data into the specified SQL table.
func (db *Database) SendData(sqlTable string, parameters []string, values []string) (id int) {
	// Process parameters for statements in SQL
	myFormattedParameters := strings.Join(parameters, ", ")

	// Process values for statements in SQL
	myFormattedValues := strings.Join(values, ", ")

	sqlStatement := fmt.Sprintf(`
	INSERT INTO %s (%s)
	VALUES (%s)
	RETURNING id`, sqlTable, myFormattedParameters, myFormattedValues)

	err := db.DB.QueryRow(sqlStatement).Scan(&id)
	if err != nil {
		panic(err)
	}

	log.Println(" > New record ID is:", id)

	return id
}

// SendDataAsJSON inserts JSON data into the specified SQL table.
func (db *Database) SendDataAsJSON(jsonData []byte, sqlTable string) error {
	// Insert into table
	query := fmt.Sprintf("INSERT INTO %s (data) VALUES ($1)", sqlTable)
	log.Println(query)
	_, err := db.DB.Exec(query, jsonData)
	if err != nil {
		log.Fatalf(" > Error inserting into table: %s", err)
		return err
	}

	return nil
}

// GetXData retrieves the last n records from the specified SQL table.
func (db *Database) GetXData(sqlTable string, numberOfData int) ([]map[string]interface{}, error) {
	// Query to get last n rows from the table
	rows, err := db.DB.Query(fmt.Sprintf("SELECT data FROM %s ORDER BY id DESC LIMIT %d", sqlTable, numberOfData))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Array to store the result
	result := []map[string]interface{}{}

	// Loop through each row and scan its data into a map
	for rows.Next() {
		var jsonData []byte
		data := make(map[string]interface{})

		err := rows.Scan(&jsonData)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			return nil, err
		}

		result = append(result, data)
	}

	log.Println(" > Query successfully finished")

	return result, nil
}

// GetLastData retrieves the last record from the specified SQL table.
func (db *Database) GetLastData(sqlTable string) (map[string]interface{}, error) {
	data, err := db.GetXData(sqlTable, 1)
	if err != nil {
		return nil, err
	}

	return data[0], err
}

// GetAllData retrieves all records from the specified SQL table.
func (db *Database) GetAllData(parameters []string, sqlTable string) ([]map[string]interface{}, error) {
	// Query to get all rows from the table
	rows, err := db.DB.Query(fmt.Sprintf("SELECT COUNT(id) FROM %s;", sqlTable))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Array to store the result
	result := []map[string]interface{}{}

	// Loop through each row and scan its data into a map
	for rows.Next() {
		var jsonData []byte
		data := make(map[string]interface{})

		err := rows.Scan(&jsonData)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			return nil, err
		}

		result = append(result, data)
	}

	log.Println(" > Query successfully finished")

	return result, nil
}

// InsertUser adds a new user to the users table.
func (db *Database) InsertUser(user User) error {
	sqlQuery := `INSERT INTO users ("username", password) VALUES ($1, $2) ON CONFLICT ("username") DO NOTHING;`
	_, err := db.DB.Exec(sqlQuery, user.Username, user.Password)
	if err != nil {
		log.Fatalf(" > Error inserting into table: %s", err)
		return err
	}

	return nil
}

// CheckUser check if the user is register in the base data
func (db *Database) CheckUser(user User) (bool, error) {
	sqlQuery := "SELECT * FROM users WHERE \"username\" = $1 AND password = $2"
	result, err := db.DB.Exec(sqlQuery, user.Username, user.Password)
	if err != nil {
		log.Fatalf(" > Error query: %s", err)
		return false, err
	}
	if out, err := result.RowsAffected(); err != nil {
		log.Fatalf(" > Error get data from table: %s", err)
		return false, err
	} else if out >= 1 && err == nil {
		log.Println(" > User and password are correct.")
		return true, nil
	} else {
		log.Println(" > User or password not found.")
		return false, nil
	}

}
