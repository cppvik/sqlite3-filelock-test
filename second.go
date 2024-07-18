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
	fmt.Println("Second: Opening the DB")
	db, err := sql.Open("sqlite3", "file:./test.db?_locking=EXCLUSIVE&_mode=rwc&_mutex=full&_busy_timeout=60000")
	if err != nil {
		log.Fatal(err)
	}

	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY AUTOINCREMENT, value TEXT)`)
	if err != nil {
		log.Fatal(err)
	}

	iterations := 10000

	// Write operations
	for i := 0; i < iterations; i++ {
		value := fmt.Sprintf("Read %d", i)
		_, err = db.Exec(`INSERT INTO test (value) VALUES (?)`, value)
		if err != nil {
			log.Fatalf("Second: Error inserting value %s: %v\n", value, err)
		}
		fmt.Printf("Second: Inserted %s\n", value)
	}
	db.Close()

	// Read operation
	time.Sleep(1 * time.Second)
	fmt.Println("Second: Reading from DB:")
	db, err = sql.Open("sqlite3", "file:./test.db?_mode=ro&_mutex=full&_busy_timeout=60000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Second: Opened the DB")
	defer db.Close()

	file, err := os.OpenFile("result.second.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := io.Writer(file)

	fmt.Println("Second: Running SELECT query")
	rows, err := db.Query(`SELECT id, value FROM test`)
	if err != nil {
		log.Fatalf("Second: Error reading data: %v\n", err)
	}
	fmt.Println("Second: Got response from SELECT query")
	for rows.Next() {
		var id int
		var value string
		err = rows.Scan(&id, &value)
		if err != nil {
			log.Fatalf("Second: Error scanning row: %v\n", err)
		}
		data := fmt.Sprintf("Second: id: %d, value: %s\n", id, value)
		// fmt.Printf(data)
		_, err := writer.Write([]byte(data))
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Second: Finished reading from DB")

	err = rows.Err()
	if err != nil {
		log.Fatalf("Second: Error iterating rows: %v\n", err)
	}
}
