package web

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	//"strings"
)

type Haccp struct {
	Fa         string `json:Fa`
	MerkleRoot string `json:MerkleRoot`
	Time       string `json:Time`
}

// Query handles chaincode query requests.

func (setup OrgSetup) Verify(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Verification request")

	queryParams := r.URL.Query()
	factoryName := queryParams.Get("factoryName")
	date := queryParams.Get("date")
	data := queryParams.Get("data")
	//테스트용
	//originalDatas := strings.Split(data, ",")
	// var layerZero [][]byte
	// for _, value := range originalDatas {
	// 	hash := sha256.New()
	// 	hash.Write([]byte(value))
	// 	hashSum := hash.Sum(nil)
	// 	layerZero = append(layerZero, hashSum[:])
	// }

	// root := merkleTree(layerZero, len(layerZero))
	// mkroot := string(root[0])
	// 나중에 Create -> Updata로 바꿀시 위코드로 수정 필요
	hash := sha256.New()
	hash.Write([]byte(data))
	hashSum := hash.Sum(nil)
	mkroot := string(hashSum)
	//

	chainCodeName := "basic"
	channelID := "mychannel"
	function := "ReadHaccp"

	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	Bc_data, err := contract.EvaluateTransaction(function, factoryName+date)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	var Bc_haccp Haccp
	err = json.Unmarshal([]byte(Bc_data), &Bc_haccp)
	if err != nil {
		fmt.Fprintf(w, "Error Unmarshalling Haccp object: %s", err)
		return
	}

	Verify_haccp := Haccp{Fa: factoryName}
	Verify_haccp.MerkleRoot = mkroot
	Verify_haccp.Time = date
	Verify_haccpJSON, err := json.Marshal(Verify_haccp)
	if err != nil {
		fmt.Fprintf(w, "Error marshalling Haccp object: %s", err)
		return
	}
	err = json.Unmarshal(Verify_haccpJSON, &Verify_haccp)
	if err != nil {
		fmt.Fprintf(w, "Error Unmarshalling Haccp object: %s", err)
		return
	}
	if Bc_haccp.MerkleRoot == Verify_haccp.MerkleRoot {
		fmt.Fprintf(w, "Verification passed")
		fmt.Println("Verification passed")
	} else {
		fmt.Fprintf(w, "Merkle Root does not match, verification failed")
		fmt.Println("Verification failed")
	}
}
