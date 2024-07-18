package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("Opening the DB")
	db, err := sql.Open("sqlite3", "file:./test.db?_locking=EXCLUSIVE&_mode=rwc&_mutex=full&_busy_timeout=60000")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY AUTOINCREMENT, value TEXT)`)
	if err != nil {
		log.Fatal(err)
	}

	iterations := 10000

	// Write operations
	for i := 0; i < iterations; i++ {
		value := fmt.Sprintf("Write %d", i)
		_, err := db.Exec(`INSERT INTO test (value) VALUES (?)`, value)
		if err != nil {
			log.Fatalf("Error inserting value %s: %v\n", value, err)
		}
		fmt.Printf("Writer: Inserted %s\n", value)
	}
	db.Close()

	// Read operation
	time.Sleep(1 * time.Second)
	fmt.Println("Reader: Reading from DB:")
	db, err = sql.Open("sqlite3", "file:./test.db?_mode=ro&_mutex=full&_busy_timeout=60000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Opened the DB")
	defer db.Close()

	file, err := os.OpenFile("result.first.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := io.Writer(file)

	fmt.Println("Running SELECT query")
	rows, err := db.Query(`SELECT id, value FROM test`)
	if err != nil {
		log.Fatalf("Error reading data: %v\n", err)
	}
	fmt.Println("Got response from SELECT query")

	for rows.Next() {
		var id int
		var value string
		err = rows.Scan(&id, &value)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		data := fmt.Sprintf("Reader: id: %d, value: %s\n", id, value)
		// fmt.Printf(data)
		_, err := writer.Write([]byte(data))
		if err != nil {
			panic(err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatalf("Error iterating rows: %v\n", err)
	}
}
