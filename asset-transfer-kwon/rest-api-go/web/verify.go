package web

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
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

//중앙 집중식 방식 즉 DB에서 호출해서 비교하고 이를 전달해주는 경우 검증과정 어차피 DB에 추가하는건 이미 코드로 구현되어 있으니 위에 코드만 바꾸면 블록 추가 즉 데이터를 추가하는 과정에서 일어나는 상황과 데이터를 검증하는 과정에서 상황에서 비교가 가능하고 이를 또한 어차피 DB에 저장한다고 하면 mkRoot의
//를 계산하고 추가 용량이 들어가는데 이는 수식을통해 이미 증명한봐 충분한 가짓수를 늘리면 어차피 큰 차이가 없다는것 또한 시나리오에 따라 주기는 길게 가져간다는것 이런것을 보면 증명가능할듯 결국 실제 중앙집중식 서버와 블록체인 네트워크를 하이브리드로 운용했을때 얻는 이점 및 단점에 대해 비교한다고 생각하묜
// 어떨까 정확히 오프체인 기반으로는 쓸만하지않을까? 이런느낌 또한 구조에 대한 경우는 참고 문서를 통해서
