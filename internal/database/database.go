package database

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB() {
    var err error
    db, err = sql.Open("sqlite", "./data.db")
    if err != nil {
        log.Fatal(err)
    }

    createTable := `
    CREATE TABLE IF NOT EXISTS plant_data (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
        temperature REAL,
        humidity REAL,
        soil REAL,
        light REAL
    );`
    _, err = db.Exec(createTable)
    if err != nil {
        log.Fatal(err)
    }
}

// Inserts the parsed MQTT JSON message into the database
func InsertMessageToDB(message string) {
    if db == nil {
        log.Println("Database not initialized")
        return
    }

    // Parse JSON string to map
    var data map[string]float32
    err := json.Unmarshal([]byte(message), &data)
    if err != nil {
        log.Println("Error parsing JSON:", err)
        return
    }

    // Prepare insert statement
    stmt, err := db.Prepare(`INSERT INTO plant_data(temperature, humidity, soil, light) VALUES (?, ?, ?, ?)`)
    if err != nil {
        log.Println("Error preparing SQL statement:", err)
        return
    }
    defer stmt.Close()

    _, err = stmt.Exec(data["temperature"], data["humidity"], data["moisture"], data["light"])
    if err != nil {
        log.Println("Error inserting data into DB:", err)
        return
    }

    log.Println("Inserted message into DB:", data)
}
