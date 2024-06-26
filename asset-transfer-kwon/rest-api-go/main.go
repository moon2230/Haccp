package main

import (
	"fmt"
	"rest-api-go/web"
)

func main() {

	//Database init
	err := web.InitDatabase()

	//Initialize setup for Org1
	cryptoPath := "../../test-network/organizations/peerOrganizations/org1.example.com"
	orgConfig := web.OrgSetup{
		OrgName:      "Org1",
		MSPID:        "Org1MSP",
		CertPath:     cryptoPath + "/users/User1@org1.example.com/msp/signcerts/cert.pem",
		KeyPath:      cryptoPath + "/users/User1@org1.example.com/msp/keystore/",
		TLSCertPath:  cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt",
		PeerEndpoint: "localhost:7051",
		GatewayPeer:  "peer0.org1.example.com",
	}
	orgSetup, err := web.Initialize(orgConfig)
	if err != nil {
		fmt.Println("Error initializing setup for Org1: ", err)
	}

	orgSetup.InvokeInit2()

	go web.SaveDataPeriodically()
	//DB 데이터 저장시 성능 향상을 위해 초당 1000개이하의 데이터가 들어오면 1초에 한번씩 Insert 실행
	//만약 초당 1000개 이상이 들어오면 리스트에 1000개를 저장했다가 1000개일때 저장

	//Server start
	web.Serve(web.OrgSetup(*orgSetup))
}
