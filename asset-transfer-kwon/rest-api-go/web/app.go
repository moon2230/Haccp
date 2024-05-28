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
	http.Handle("/query", JWTMiddleware(http.HandlerFunc(setups.Query)))
	http.Handle("/invoke", JWTMiddleware(http.HandlerFunc(setups.Invoke)))
	http.Handle("/invokeInit", JWTMiddleware(http.HandlerFunc(setups.InvokeInit)))
	http.Handle("/data", JWTMiddleware(http.HandlerFunc(setups.Inquery)))
	http.Handle("/verify", JWTMiddleware(http.HandlerFunc(setups.Verify)))
	http.Handle("/dailyInvoke", JWTMiddleware(http.HandlerFunc(setups.DailyInvoke)))

	// Server start in background
	fmt.Println("Listening (http://localhost:3001/)...")

	if err := http.ListenAndServe(":3001", nil); err != nil {
		fmt.Println(err)
	}
}
