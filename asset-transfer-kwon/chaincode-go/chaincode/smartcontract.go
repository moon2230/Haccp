package chaincode

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	//"encoding/base64"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Cahincode by Kwon -start-
type Haccp struct {
	Fa         string `json:Fa`
	MerkleRoot string `json:MerkleRoot`
	Time       string `json:Time`
}

func (s *SmartContract) InitHaccp(ctx contractapi.TransactionContextInterface) error {
	haccps := []Haccp{
		{Fa: "Fa120230101"},
		{Fa: "Fa120230102"},
		{Fa: "Fa120230103"},
		{Fa: "Fa120230104"},
		{Fa: "Fa120230105"},
		{Fa: "Fa120230106"},
	}

	for _, haccp := range haccps {
		hash := sha256.New()
		hash.Write([]byte(haccp.Fa))
		hashSum := hash.Sum(nil)
		haccp.MerkleRoot = string(hashSum)
		//fmt.Println(base64.StdEncoding.EncodeToString(hash.Sum([]byte(haccp.Fa))))
		haccp.Time = time.Now().Format("2006-01-02 15:04:05")
		haccpJSON, err := json.Marshal(haccp)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(haccp.Fa, haccpJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

func (s *SmartContract) GetAllHaccp(ctx contractapi.TransactionContextInterface) ([]*Haccp, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var haccps []*Haccp
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var haccp Haccp
		err = json.Unmarshal(queryResponse.Value, &haccp)
		if err != nil {
			return nil, err
		}
		haccps = append(haccps, &haccp)
	}

	return haccps, nil
}

func (s *SmartContract) CreateHaccp(ctx contractapi.TransactionContextInterface, faid string) error {
	// exists, err := s.HaccpExists(ctx, faid)
	// if err != nil {
	// 	return err
	// }
	// if exists {
	// 	return fmt.Errorf("the asset %s already exists", faid)
	// }

	haccp := Haccp{Fa: faid}
	hash := sha256.New()
	hash.Write([]byte(haccp.Fa))
	hashSum := hash.Sum(nil)
	haccp.MerkleRoot = string(hashSum)
	haccp.Time = time.Now().Format("2006-01-02 15:04:05")
	haccpJSON, err := json.Marshal(haccp)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(haccp.Fa, haccpJSON)
}

func (s *SmartContract) HaccpExists(ctx contractapi.TransactionContextInterface, faid string) (bool, error) {
	haccpJSON, err := ctx.GetStub().GetState(faid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return haccpJSON != nil, nil
}

// func (s *SmartContract) UpdateHaccp(ctx contractapi.TransactionContextInterface, faid string, mkroot string) error {
// 	// -- This needs to be thought through
// 	/*exists, err := s.HaccpExists(ctx, faid)
// 	if err != nil {
// 		return err
// 	}
// 	if !exists {
// 		return fmt.Errorf("the asset %s does not existtt", mkroot)
// 	}*/

// 	// overwriting original asset with new asset
// 	haccp := Haccp{Fa: faid}
// 	haccp.MerkleRoot = mkroot
// 	haccp.Time = time.Now().Format("2006-01-02 15:04:05")
// 	haccpJSON, err := json.Marshal(haccp)
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.GetStub().PutState(haccp.Fa, haccpJSON)
// }

func (s *SmartContract) UpdateHaccp(ctx contractapi.TransactionContextInterface, faid string, mkroot string, time string) error {
	// -- This needs to be thought through
	/*exists, err := s.HaccpExists(ctx, faid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not existtt", mkroot)
	}*/

	// overwriting original asset with new asset
	haccp := Haccp{
		Fa:         faid,
		MerkleRoot: mkroot,
		Time:       time,
	}
	haccpJSON, err := json.Marshal(haccp)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(haccp.Fa, haccpJSON)
}

func (s *SmartContract) ReadHaccp(ctx contractapi.TransactionContextInterface, faid string) (*Haccp, error) {
	haccpJSON, err := ctx.GetStub().GetState(faid)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if haccpJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", faid)
	}

	var haccp Haccp
	err = json.Unmarshal(haccpJSON, &haccp)
	if err != nil {
		return nil, err
	}

	return &haccp, nil
}

//Cahincode by Kwon -end-
