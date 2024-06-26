package web

import (
	"fmt"
	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// OrgSetup contains organization's config to interact with the network.
type OrgSetup struct {
	OrgName      string
	MSPID        string
	CryptoPath   string
	CertPath     string
	KeyPath      string
	TLSCertPath  string
	PeerEndpoint string
	GatewayPeer  string
	Gateway      client.Gateway
}

// Serve starts http web server.
func Serve(setups OrgSetup) {

	// User verification handler
	http.HandleFunc("/login", setups.Login)

	// User Token verification handler
	http.HandleFunc("/loadverify", setups.Loadverify)

	// Main html page handler
	http.HandleFunc("/", setups.PageMain)

	// Protected routes
	http.Handle("/invoke", http.HandlerFunc(setups.Invoke))
	http.Handle("/invokeInit", http.HandlerFunc(setups.InvokeInit))
	http.Handle("/query", JWTAndRoleMiddleware(http.HandlerFunc(setups.Query)))
	http.Handle("/data", JWTAndRoleMiddleware(http.HandlerFunc(setups.Inquery)))
	http.Handle("/verify", JWTAndRoleMiddleware(http.HandlerFunc(setups.Verify)))
	http.Handle("/dailyInvoke", http.HandlerFunc(setups.DailyInvoke))
	//test

	// http.HandleFunc("/data", setups.Inquery)
	// http.HandleFunc("/dailyInvoke", setups.DailyInvoke)
	// Jmeter로 테스트시 중간 미들웨어에서 토큰 검증으로 인해 거부되기 때문에 검증을 거치지 않는 테스트 용도의 코드

	// Server start in background
	fmt.Println("Listening (http://localhost:3001/)...")

	if err := http.ListenAndServe(":3001", nil); err != nil {
		fmt.Println(err)
	}
}
