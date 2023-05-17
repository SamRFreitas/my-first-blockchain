package src

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Transaction struct {
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
}

type Block struct {
	Index        int           `json:"index"`
	Timestamp    time.Time     `json:"timestamp"`
	Proof        int64         `json:"proof"`
	PreviousHash string        `json:"previousHash"`
	Transactions []Transaction `json:"transactions"`
}

type Blockchain struct {
	Chain        []Block       `json:"chain"`
	Transactions []Transaction `json:"transaction"`
	Nodes        Network       `json:"network"`
}

var difficulty = "0000"
var network = NewNetwork()

type GetChainResponseData struct {
	Chain  []Block `json:"chain"`
	Length int     `json:"length"`
}

func NewBlockchain() Blockchain {
	b := Blockchain{
		Chain:        nil,
		Transactions: nil,
		Nodes:        *network,
	}

	block := Block{
		Index:        len(b.Chain) + 1,
		Timestamp:    time.Now(),
		Proof:        1,
		PreviousHash: "0",
		Transactions: nil,
	}

	b.Chain = append(b.Chain, block)

	return b
}

func (b Blockchain) CreateBlock(proof int64, previousHash string) Block {

	block := Block{
		Index:        len(b.Chain) + 1,
		Timestamp:    time.Now(),
		Proof:        proof,
		PreviousHash: previousHash,
		Transactions: b.Transactions,
	}

	return block
}

func (b Blockchain) GetPreviousBlock() Block {
	length := len(b.Chain)
	block := b.Chain[length-1]
	return block
}

func ProofOfWork(previousProof int64) int64 {

	newProof := 1
	checkedProof := false

	for {
		hashOperation := CheckProofs(int64(newProof), previousProof)
		if hashOperation[0:4] == difficulty {
			checkedProof = true
		} else {
			newProof += 1
		}
		if checkedProof == true {
			return int64(newProof)
		}
	}

}

func CheckProofs(proof int64, previousProof int64) string {

	x := math.Pow(float64(proof), 2)
	y := math.Pow(float64(previousProof), 2)

	calculation := int64(x - y)
	calculationHex := strconv.FormatInt(calculation, 16)

	h := fmt.Sprintf("%x", sha256.Sum256([]byte(calculationHex)))

	return h
}

func HashBlock(block Block) string {
	blockEncodedUint8, _ := json.Marshal(block)
	blockJson := string(blockEncodedUint8)

	h := fmt.Sprintf("%x", sha256.Sum256([]byte(blockJson)))
	return h
}

func IsChainValid(chain []Block) bool {
	previousBlock := chain[0]
	blockIndex := 1

	if len(chain) == 1 {
		return true
	}

	for {
		block := chain[blockIndex]
		if block.PreviousHash != HashBlock(previousBlock) {
			return false
		}

		previousProof := previousBlock.Proof
		proof := block.Proof
		hashOperation := CheckProofs(proof, previousProof)
		if hashOperation[:4] != difficulty {
			return false
		}

		if blockIndex == len(chain)-1 {
			return true
		}

		previousBlock = block
		blockIndex += 1

	}
}

func (b Blockchain) AddTransaction(sender string, receiver string, amount float64) (int, Transaction) {
	transaction := Transaction{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
	}

	previousBlock := b.GetPreviousBlock()

	return previousBlock.Index + 1, transaction
}

func (b Blockchain) AddNode(address string) {
	url, _ := url.Parse(address)
	b.Nodes.Add(*url)
}

func (b Blockchain) ReplaceChain() (bool, []Block) {
	var longestChain []Block
	maxLength := len(b.Chain)

	for node, _ := range b.Nodes.Nodes {

		address := "http://" + node.Host + "/get-chain"

		response, _ := http.Get(address)

		var data GetChainResponseData
		decoder := json.NewDecoder(response.Body)
		decoder.Decode(&data)

		if response.StatusCode == 200 {

			if data.Length > maxLength && IsChainValid(data.Chain) {

				maxLength = data.Length
				longestChain = data.Chain

			}
		}

	}

	if len(longestChain) == 0 {
		return false, b.Chain
	}

	return true, longestChain

}
