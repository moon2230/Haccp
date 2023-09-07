package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type Haccp struct {
	Fa         string `json:Fa`
	MerkleRoot string `json:MerkleRoot`
	Time       string `json:Time`
}

func main() {
	// MySQL 데이터베이스 연결 설정
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/fabric_db")
	if err != nil {
		fmt.Println("MySQL 연결 실패:", err)
		return
	}
	defer db.Close()

	//조회하는 데이터 날짜 시간을 제외하고 날짜만 기준
	//data := time.Now().Format("20060102")

	//시간을 설정해서 그에따라 반복저으로 코드 실행할수있게 제작
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {

		//머클루투값 조회하는는 부분
		var mekleroot string
		err = db.QueryRow("SELECT hash FROM merkle_root WHERE data = ?", "20230906").Scan(&mekleroot)
		//잘 조회가 됬는지 확인하는 부분 나중에 제거 가능//
		if err == sql.ErrNoRows {
			fmt.Println("데이터가 MySQL에 존재하지 않습니다.")
		} else if err != nil {
			fmt.Println("데이터를 MySQL에서 조회하는 중 오류 발생:", err)
		} else {
			fmt.Println("MySQL에서 조회한 hasg:", mekleroot)
		}

		// 인식번호 조회하는부분
		var fa string
		err = db.QueryRow("SELECT num FROM merkle_root WHERE data = ?", "20230906").Scan(&fa)
		//잘 조회가 됬는지 확인하는 부분 나중에 제거 가능//
		if err == sql.ErrNoRows {
			fmt.Println("데이터가 MySQL에 존재하지 않습니다.")
		} else if err != nil {
			fmt.Println("데이터를 MySQL에서 조회하는 중 오류 발생:", err)
		} else {
			fmt.Println("MySQL에서 조회한 fa:", fa)
		}

		//mysql에서 불러온 데이터를 원장에 업데이트 한다 이때
		//cli에서 테스트 네트우크 폴더로 이동후 쉘을 실행해서 hash값을 배포한다 이때  UpdateHaccp함수를 사용한다
		//기본적으로 변수 설정이랑 그다음 설정해주면 자동으로 바로 배포가 가능하다 org1에 다가
		//CLI에서 체인코드를 사용해서 배포한다
		log.Println("============ application-mysql update starts ============")

		err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
		if err != nil {
			log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environment variable: %v", err)
		}

		walletPath := "wallet"
		// remove any existing wallet from prior runs
		os.RemoveAll(walletPath)
		wallet, err := gateway.NewFileSystemWallet(walletPath)
		if err != nil {
			log.Fatalf("Failed to create wallet: %v", err)
		}

		if !wallet.Exists("appUser") {
			err = populateWallet(wallet)
			if err != nil {
				log.Fatalf("Failed to populate wallet contents: %v", err)
			}
		}

		ccpPath := filepath.Join(
			"..",
			"..",
			"test-network",
			"organizations",
			"peerOrganizations",
			"org1.example.com",
			"connection-org1.yaml",
		)

		gw, err := gateway.Connect(
			gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
			gateway.WithIdentity(wallet, "appUser"),
		)
		if err != nil {
			log.Fatalf("Failed to connect to gateway: %v", err)
		}
		defer gw.Close()

		channelName := "mychannel"
		if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
			channelName = cname
		}

		log.Println("--> Connecting to channel", channelName)
		network, err := gw.GetNetwork(channelName)
		if err != nil {
			log.Fatalf("Failed to get network: %v", err)
		}

		chaincodeName := "basic"
		if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
			chaincodeName = ccname
		}

		log.Println("--> Using chaincode", chaincodeName)
		contract := network.GetContract(chaincodeName)
		//초기화 부분 필요하면 활성화
		log.Println("--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger")
		result, err := contract.SubmitTransaction("InitHaccp")
		if err != nil {
			log.Fatalf("Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("--> Submit Transaction: CreateAsset, creates new asset with ID, color, owner, size, and appraisedValue arguments")
		result, err = contract.SubmitTransaction("CreateHaccp", fa)
		if err != nil {
			log.Fatalf("Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("--> Submit Transaction: UpdateHaccp")
		result, err = contract.SubmitTransaction("UpdateHaccp", fa, mekleroot)
		if err != nil {
			log.Fatalf("Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")

	}

	//duration, _ := time.ParseDuration("1s")

}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := os.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := os.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := os.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}
