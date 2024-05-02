package web

import (
	"fmt"
	"net/http"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"crypto/sha256"
)

func databaseOpen() (*sql.DB, error) {
	db, err := sql.Open("mysql", "ck@tcp(localhost:3306)/haccp")

	err = db.Ping()
	if err != nil {
		if err.Error() == "Error 1049 (42000): Unknown database 'haccp'" {
			db.Close()
			db, err = sql.Open("mysql", "ck@tcp(localhost:3306)/")

			createDB := "Create Database haccp"
			_, err = db.Exec(createDB)

			fmt.Println("haccp Database created.")

			db, err = sql.Open("mysql", "ck@tcp(localhost:3306)/haccp")
		}
	}
	return db, err
}

func databaseTable(db *sql.DB, tableName string) error {
	checkTable := "Select Factory ,Time, Data from haccp." + tableName
	_, err := db.Exec(checkTable)
	if err != nil {
		if err.Error() == "Error 1146 (42S02): Table 'haccp."+tableName+"' doesn't exist" {
			createTable := "create table " + tableName + "(Factory text, Time datetime, Data Blob)"
			_, err = db.Exec(createTable)
		}
	}
	return err
}

func InitDatabase() error {
	db, err := sql.Open("mysql", "ck@tcp(localhost:3306)/haccp")

	err = db.Ping()
	if err != nil {
		if err.Error() == "Error 1049 (42000): Unknown database 'haccp'" {
			db.Close()
			db, err = sql.Open("mysql", "ck@tcp(localhost:3306)/")

			createDB := "Create Database haccp"
			_, err = db.Exec(createDB)

			fmt.Println("haccp Database created.")

			db, err = sql.Open("mysql", "ck@tcp(localhost:3306)/haccp")
		}
	}
	fmt.Println("Mysql Database name: haccp has been initiated.")

	return err
}

// Query handles chaincode query requests.
func (setup OrgSetup) Inquery(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	fmt.Println(queryParams)
	data := queryParams.Get("data")
	fmt.Println(data)
	name := queryParams.Get("name")
	fmt.Println(name)

	db, err := databaseOpen()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = databaseTable(db, "leaf")
	if err != nil {
		panic(err.Error())
	}

	//Hash call
	var hash [32]byte = sha256.Sum256([]byte(data))
	s := hash[:]

	insertRecord := "insert into leaf (Factory, Time, Data) values (?, ?, ?)"
	_, err = db.Exec(insertRecord, name, time.Now().Format("2006-01-02 15:04:05"), s)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "Response: %s", data)
}
