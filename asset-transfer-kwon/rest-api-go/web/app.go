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
	http.HandleFunc("/verifytoken", setups.verifyToken)

	//Main html page handler
	http.HandleFunc("/", setups.PageMain)

	//Blockchain Query handler - gettAll
	http.HandleFunc("/query", setups.Query)

	//Blockchain Invoke handler
	http.HandleFunc("/invoke", setups.Invoke)

	//Blockchain init handler -Init-
	http.HandleFunc("/invokeInit", setups.InvokeInit)

	//Data stream input handler -data-
	http.HandleFunc("/data", setups.Inquery)

	//Data integrity verification handler -verify-
	http.HandleFunc("/verify", setups.Verify)

	http.HandleFunc("/dailyInvoke", setups.DailyInvoke)

	//Server start in background
	fmt.Println("Listening (http://localhost:3001/)...")

	if err := http.ListenAndServe(":3001", nil); err != nil {
		fmt.Println(err)
	}
}
