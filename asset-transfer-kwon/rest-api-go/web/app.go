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
	http.HandleFunc("/query", setups.Query)
	http.HandleFunc("/invoke", setups.Invoke)
	http.HandleFunc("/verify", setups.Verify)

	http.HandleFunc("/invokeInit", setups.InvokeInit)
	//http.HandleFunc("/invokeUpdate", setups.InvokeUpdate)
	//http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request){fmt.Fprintf(w, "<html><title>Haccp</title><button>1</buton><button>2</buton><button>3</buton><button>4</buton></html>")})
	http.HandleFunc("/", LoadMain)
	http.HandleFunc("/main", LoadMain)

	fmt.Println("Listening (http://localhost:3001/)...")
	if err := http.ListenAndServe(":3001", nil); err != nil {
		fmt.Println(err)
	}
}
