package web

import (
	"fmt"
	"net/http"

	"crypto/sha256"

	"github.com/hyperledger/fabric-gateway/pkg/client"

	"time"
)

// Invoke handles chaincode invoke requests.
func (setup *OrgSetup) Invoke(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := r.FormValue("chaincodeid")
	channelID := r.FormValue("channelid")
	function := r.FormValue("function")
	args := r.Form["args"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)

	hash := sha256.New()
	hash.Write([]byte(args[0]))
	hashSum := hash.Sum(nil)
	MerkleRoot := string(hashSum)
	Time := time.Now().Format("2006-01-02 15:04:05")
	args = append(args, MerkleRoot, Time)

	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		fmt.Fprintf(w, "Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		fmt.Fprintf(w, "Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		fmt.Fprintf(w, "Error submitting transaction: %s", err)
		return
	}
	fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
}

func (setup *OrgSetup) InvokeInit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "InitHaccp"
	args := r.Form["args"]
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)

	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		fmt.Fprintf(w, "Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		fmt.Fprintf(w, "Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		fmt.Fprintf(w, "Error submitting transaction: %s", err)
		return
	}

	fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
}

func (setup *OrgSetup) InvokeInit2() {
	fmt.Println("Received Invoke request")

	chainCodeName := "basic"
	channelID := "mychannel"
	function := "InitHaccp"
	args := "args"
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args))
	if err != nil {
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		return
	}
	fmt.Printf("Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
}

func merkleTree(layerZero [][]byte, jump int) [][]byte {
	if len(layerZero) == 1 {
		fmt.Println(layerZero)
		return layerZero
	}
	var newLayer [][]byte
	for i := 0; i < len(layerZero); i += jump {
		var con []byte
		for ii := 0; ii < jump; ii += 1 {
			if i+ii < len(layerZero) {
				con = append(layerZero[i+ii])
			} else {
				con = append(layerZero[i])
			}
		}
		hash := sha256.New()
		hash.Write([]byte(con))
		hashSum := hash.Sum(nil)
		newLayer = append(newLayer, hashSum)
	}
	if jump == 1 {
		return layerZero
	}
	fmt.Println(len(layerZero))

	return merkleTree(newLayer, jump)
}

func (setup *OrgSetup) DailyInvoke(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	args := queryParams.Get("args")

	db, err := databaseOpen()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var res string

	checkTable := "Select data from haccp.leaf where Factory = '" + args + "' and Date(Time) = curdate()"

	rows, err := db.Query(checkTable)
	if err != nil {
		fmt.Println("Query error.")
	}
	defer rows.Close()

	var newLayer [][]byte
	for rows.Next() {
		err := rows.Scan(&res)
		if err != nil {
			fmt.Println("Error")
		}
		raw := []byte(res)
		newLayer = append(newLayer, raw)
	}

	root := merkleTree(newLayer, len(newLayer))
	mkroot := string(root[0])

	//Total hash

	channelID := "mychannel"
	chainCodeName := "basic"
	function := "UpdateHaccp"

	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)

	faName := args
	//Time := time.Now().Format("20060102")//타입포맷 변경
	Time := time.Now().Format("2006-01-02 15:04:05")

	txn_proposal, err := contract.NewProposal(function, client.WithArguments(faName, Time, mkroot))
	if err != nil {
		fmt.Fprintf(w, "Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		fmt.Fprintf(w, "Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		fmt.Fprintf(w, "Error submitting transaction: %s", err)
		return
	}

	fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
}
