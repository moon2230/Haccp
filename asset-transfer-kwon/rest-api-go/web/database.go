package web

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
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
// func (setup OrgSetup) Inquery(w http.ResponseWriter, r *http.Request) {
// 	queryParams := r.URL.Query()
// 	fmt.Println(queryParams)
// 	data := queryParams.Get("data")
// 	fmt.Println(data)
// 	name := queryParams.Get("name")
// 	fmt.Println(name)

// 	db, err := databaseOpen()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	defer db.Close()

// 	err = databaseTable(db, "leaf")
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	//Hash call
// 	var hash [32]byte = sha256.Sum256([]byte(data))
// 	s := hash[:]

// 	insertRecord := "insert into leaf (Factory, Time, Data) values (?, ?, ?)"
// 	_, err = db.Exec(insertRecord, name, time.Now().Format("2006-01-02 15:04:05"), s)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
// 	fmt.Fprintf(w, "Response: %s", data)
// }

var mu sync.Mutex
var records [][]interface{}

func (setup OrgSetup) Inquery(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	data := queryParams.Get("data")
	name := queryParams.Get("name")

	var hash [32]byte = sha256.Sum256([]byte(data))
	s := hash[:]

	// 데이터를 슬라이스에 추가
	mu.Lock()
	defer mu.Unlock()
	records = append(records, []interface{}{name, time.Now().Format("2006-01-02 15:04:05"), s})

	fmt.Fprintf(w, "Response: %s", data)
}

const batchSize = 1000 // 일괄 삽입할 배치 크기

func SaveDataPeriodically() {
	// 1초 동안의 시간 측정용 변수
	var startTime time.Time
	for {
		// 1초마다 시간 측정을 재시작
		startTime = time.Now()

		// 1초 동안 데이터를 확인하고, 1000개 이상일 때는 바로 삽입
		for time.Since(startTime) < time.Second {
			mu.Lock()
			if len(records) >= batchSize {
				performBatchInsert()
				mu.Unlock()
				break // 1000개 이상이면 바로 삽입하고 루프 종료
			}
			mu.Unlock()
			time.Sleep(time.Millisecond) // 데이터를 확인하는 간격 설정
		}

		// 1초가 지나면 삽입 수행
		mu.Lock()
		if len(records) > 0 {
			performBatchInsert()
		}
		mu.Unlock()
	}
}

// 일괄 삽입 수행 함수
func performBatchInsert() {
	// DB 연결
	db, err := databaseOpen()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	// 쿼리 및 매개변수 준비
	var valueStrings []string
	var valueArgs []interface{}
	for _, record := range records {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, record...)
	}
	query := fmt.Sprintf("INSERT INTO leaf (Factory, Time, Data) VALUES %s", strings.Join(valueStrings, ","))

	// 일괄 삽입 실행
	_, err = db.Exec(query, valueArgs...)
	if err != nil {
		log.Println(err)
	}

	// 데이터 슬라이스 초기화
	records = [][]interface{}{}
}
