package web

import (
	"fmt"
	"net/http"
)

// Query handles chaincode query requests.
func (setup OrgSetup) HashData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received hashing request")

	queryParams := r.URL.Query()
	fmt.Println("Params: %s\n", queryParams)

	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")
	function := queryParams.Get("function")
	args := r.URL.Query()["args"]

	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)

	evaluateResponse, err := contract.EvaluateTransaction(function, args...)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}