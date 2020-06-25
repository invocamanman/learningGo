package main

import (
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"

	//"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/syndtr/goleveldb/leveldb"
)

type tx struct {
	From   string //important capitalletters?¿
	To     string
	Amount int
	Hash   string
}

func (tx *tx) toBytes() []byte {
	b, err := json.Marshal(tx)
	if err != nil {
		fmt.Println("error:", err)
	}
	return []byte(b)
}

func fromBytes(bytes []byte) tx {
	var transaction tx
	err := json.Unmarshal(bytes, &transaction)
	if err != nil {
		fmt.Println("error:", err)
	}
	return transaction
}

//account := common.HexToAddress("0x71c7656ec7ab88b098defb751b7401b5f6d8976f")

func main() {
	db, err := leveldb.OpenFile("./db", nil)
	defer db.Close()
	if err != nil {
		fmt.Println("error? ", err)
	}

	addressA := generateRandomAddress()
	addressB := generateRandomAddress()

	const amount = 1000

	amountBin := make([]byte, 8)
	binary.LittleEndian.PutUint64(amountBin, amount)

	hashKek := crypto.Keccak256Hash(addressA.Bytes(), addressB.Bytes(), amountBin)

	transaction := tx{From: addressA.Hex(), To: addressB.Hex(), Amount: 1000, Hash: hashKek.Hex()}
	fmt.Println("transaction ", transaction)

	transactionByte := transaction.toBytes()
	fmt.Println("transactionByte ", transactionByte)

	err = db.Put([]byte(transaction.Hash), transactionByte, nil)
	if err != nil {
		fmt.Println("error? ", err)
	}

	data, err := db.Get([]byte(transaction.Hash), nil)
	if err != nil {
		fmt.Println("error? ", err)
	}

	transactionFromByte := fromBytes(data)
	fmt.Println("transactionFromByte ", transactionFromByte)
	Hi()
}

func generateRandomAddress() common.Address {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	//privateKeyBytes := crypto.FromECDSA(privateKey)
	//fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // 0xfad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	// publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	// fmt.Println(hexutil.Encode(publicKeyBytes))     // 0xa97df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05
	// fmt.Println(hex.EncodeToString(publicKeyBytes)) //a97df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05

	address := crypto.PubkeyToAddress(*publicKeyECDSA) //.Hex() --> string
	return address

	// account := common.HexToAddress(address)string-->hex
	// fmt.Println(account)
	// string(bytes) work!:D
}
