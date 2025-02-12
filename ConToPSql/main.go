package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

func getTables(db *sql.DB) ([]string, error) {
	query := `
	SELECT table_name 
	FROM information_schema.tables 
	WHERE table_schema='bookings'
`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	//проверяем наличие ошибок
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

// Функция для получения списка баз данных
func getDatabases(db *sql.DB) ([]string, error) {
	// SQL-запрос для получения баз данных
	query := `SELECT datname FROM pg_database;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}
		databases = append(databases, dbName)
	}

	// Проверяем наличие ошибок при переборе строк
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return databases, nil
}

func main() {
	// Настройка строки подключения
	connStr := "user=postgres password=5421 dbname=demo host=localhost port=5432 sslmode=disable"

	//Открываем подключение к базе данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Пинг базы данных
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connetcted to the database!")

	// Получение списка баз данных
	databases, err := getDatabases(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Databases in the server:")
	for _, database := range databases {
		fmt.Println(database)
	}

	// Получение списка таблиц
	tables, err := getTables(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tables list:")
	for _, table := range tables {
		fmt.Println(table)
	}

	fmt.Println("\n GETING AIRCRAFT SEATS COUNT BY SEAT CLASS")
	getSeatsCountByClassReq := "SELECT aircraft_code, fare_conditions, count(*) FROM seats_data GROUP BY aircraft_code, fare_conditions ORDER BY aircraft_code;"
	rows, err := db.Query(getSeatsCountByClassReq)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer rows.Close()

	type SeatCount struct {
		aircraft_code   string
		fare_conditions string
		count           int
	}

	var foundedRows []SeatCount
	for rows.Next() {
		var row SeatCount
		if err := rows.Scan(&row.aircraft_code, &row.fare_conditions, &row.count); err != nil {
			log.Fatal(err)
		}
		foundedRows = append(foundedRows, row)
	}

	headerStr := "| Aircraft code | Fare Condition | Seats Count |"
	fmt.Println(strings.Repeat("-", len(headerStr)))
	fmt.Println(headerStr)
	fmt.Println(strings.Repeat("-", len(headerStr)))
	for _, row := range foundedRows {
		fmt.Printf("| %-13s | %-14s | %-11d |\n", row.aircraft_code, row.fare_conditions, row.count)
	}
	fmt.Println(strings.Repeat("-", len(headerStr)))

}
