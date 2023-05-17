package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"my-first-blockchain/src"
	"net/http"
)

var blockchain = src.NewBlockchain()

var nodeAddress = uuid.NewString()

func mineBlock(w http.ResponseWriter, req *http.Request) {
	previousblock := blockchain.GetPreviousBlock()
	previousProof := previousblock.Proof

	proof := src.ProofOfWork(previousProof)

	previousHash := src.HashBlock(previousblock)
	_, transaction := blockchain.AddTransaction(nodeAddress, "Node 1", 1)
	blockchain.Transactions = append(blockchain.Transactions, transaction)

	newBlock := blockchain.CreateBlock(proof, previousHash)
	blockchain.Chain = append(blockchain.Chain, newBlock)

	var data []byte
	jsonBytes, _ := json.Marshal(newBlock)
	data = jsonBytes
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func getChain(w http.ResponseWriter, req *http.Request) {
	response := struct {
		Chain  []src.Block `json:"chain"`
		Length int         `json:"length"`
	}{
		Chain:  blockchain.Chain,
		Length: len(blockchain.Chain),
	}

	var data []byte
	jsonBytes, _ := json.Marshal(response)
	data = jsonBytes
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func isValid(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Status  int    `json:"status"`
		isValid bool   `json:"isValid"`
		Message string `json:"message"`
	}

	valid := src.IsChainValid(blockchain.Chain)

	var r response

	if valid {
		r = response{
			Status:  http.StatusOK,
			isValid: true,
			Message: "This blockchain is valid!",
		}
	} else {
		r = response{
			Status:  http.StatusUnprocessableEntity,
			isValid: false,
			Message: "This blockchain is Invalid!",
		}
	}

	var data []byte
	jsonBytes, _ := json.Marshal(r)
	data = jsonBytes
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func addTransaction(w http.ResponseWriter, req *http.Request) {
	type Transaction struct {
		Sender   string  `json:"sender"`
		Receiver string  `json:"receiver"`
		Amount   float64 `json:"amount"`
	}

	var transaction Transaction
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&transaction)

	type response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	var r response

	if err != nil {
		fmt.Println(err)
	} else {
		if transaction.Receiver == "" || transaction.Sender == "" || transaction.Amount == 0 {
			r = response{
				Status:  http.StatusBadRequest,
				Message: "Error: Empty Json",
			}
			w.WriteHeader(http.StatusBadRequest)
		} else {
			index, _ := blockchain.AddTransaction(transaction.Sender, transaction.Receiver, transaction.Amount)
			r = response{
				Status:  http.StatusOK,
				Message: "This transaction will be added to Block " + string(index),
			}
			w.WriteHeader(http.StatusOK)
		}
	}

	var data []byte
	jsonBytes, _ := json.Marshal(r)
	data = jsonBytes
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func connectNodes(w http.ResponseWriter, req *http.Request) {
	type bodyNodes struct {
		Nodes []string
	}
	var nodes bodyNodes
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&nodes)

	type response struct {
		Status     int      `json:"status"`
		Message    string   `json:"message"`
		TotalNodes []string `json:"totalNodes"`
	}

	var r response
	var totalNodes []string

	if err != nil {
		r = response{
			Status:  http.StatusBadRequest,
			Message: "Error: Invalid Json",
		}
		w.WriteHeader(http.StatusBadRequest)
	} else {
		if len(nodes.Nodes) == 0 {
			r = response{
				Status:  http.StatusBadRequest,
				Message: "Error: Empty Nodes List",
			}
			w.WriteHeader(http.StatusBadRequest)
		} else {
			for _, n := range nodes.Nodes {
				blockchain.AddNode(n)
				totalNodes = append(totalNodes, n)
			}

			r = response{
				Status:     http.StatusCreated,
				Message:    "All the nodes are now connected, The Blockchain now contains the following nodes: ",
				TotalNodes: totalNodes,
			}
			w.WriteHeader(http.StatusCreated)
		}
	}

	var data []byte
	jsonBytes, _ := json.Marshal(r)
	data = jsonBytes
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func replaceChain(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Status  int         `json:"status"`
		Message string      `json:"message"`
		Chain   []src.Block `json:"actualChain"`
	}

	isChainReplaced := blockchain.ReplaceChain()
	var r response

	if isChainReplaced {
		r = response{
			Status:  http.StatusOK,
			Message: "The nodes had different chains so the chain was replaced by the longest one.",
			Chain:   blockchain.Chain,
		}
	} else {
		r = response{
			Status:  http.StatusUnprocessableEntity,
			Message: "All good. The chain is the largest one.",
			Chain:   blockchain.Chain,
		}
	}

	var data []byte
	jsonBytes, _ := json.Marshal(r)
	data = jsonBytes
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func main() {

	http.HandleFunc("/mine-block", mineBlock)
	http.HandleFunc("/get-chain", getChain)
	http.HandleFunc("/is-valid", isValid)
	http.HandleFunc("/add-transaction", addTransaction)
	http.HandleFunc("/connect-nodes", connectNodes)
	http.HandleFunc("/replace-chain", replaceChain)

	http.ListenAndServe(":8091", nil)

}
