package merkle

import (
	"crypto/sha256"
)

func testArrayCreate(quant int)[][32]byte{
	var testArray [][32]byte
	for i := 0; i < quant; i += 1{
		testArray = append(testArray, sha256.Sum256(make([]byte, i)))
	}
	return testArray
}

func flatData(data [][32]byte) []byte {
    var result []byte
    for _, item := range data {
        result = append(result, item[:]...)
    }
    return result
}

func root(data [][32]byte) [32]byte {
	var root [32]byte
	root = sha256.Sum256(flatData(data))
	return root
}

func calculateMerkleRoot_byte(data [][]byte) []byte {
	if len(data) == 0 {
		return nil
	}

	if len(data) == 1 {
		return data[0]
	}

	var newLevel [][]byte
	for i := 0; i < len(data); i += 2 {
		hash := sha256.New()
		if i+1 < len(data) {
			hash.Write(data[i])
			hash.Write(data[i+1])
			newLevel = append(newLevel, hash.Sum(nil))
		} else {
			// 홀수 개의 데이터인 경우, 마지막 데이터를 복제하여 처리
			hash.Write(data[i])
			hash.Write(data[i])
			newLevel = append(newLevel, hash.Sum(nil))
		}
	}

	return calculateMerkleRoot_byte(newLevel)
}