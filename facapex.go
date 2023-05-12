/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a capex transactions
type SmartContract struct {
	contractapi.Contract
}

// capex describes basic details of asset transactions
type Capex struct {
	BU   string `json:"bu"`
	COCD  string `json:"cocd"`
	DOCNO string `json:"docno"`
	MRU  string `json:"mru"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Capex
}

// InitLedger adds a base set of capex  to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	transactions := []Capex{
		Capex{BU: "B1", COCD: "C1", DOCNO: "D1", MRU: "M1"},
		Capex{BU: "B2", COCD: "C2", DOCNO: "D2", MRU: "M2"},
		Capex{BU: "B3", COCD: "C3", DOCNO: "D3", MRU: "M3"},
		Capex{BU: "B4", COCD: "C4", DOCNO: "D4", MRU: "M4"},
		Capex{BU: "B5", COCD: "C5", DOCNO: "D5", MRU: "M5"},
	}

	for i, transaction := range transactions {
		transactionAsBytes, _ := json.Marshal(transaction)
		err := ctx.GetStub().PutState("Transaction"+strconv.Itoa(i), transactionAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// Createtransaction adds a new transaction to the world state with given details
func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, transactionNumber string, bu string, cocd string, docno string, mru string) error {
	transaction := Capex{
		BU:   bu,
		COCD:  cocd,
		DOCNO: docno,
		MRU:  mru,
	}

	transactionAsBytes, _ := json.Marshal(transaction)

	return ctx.GetStub().PutState(transactionNumber,transactionAsBytes)
}

// Querytransaction returns the transaction stored in the world state with given id
func (s *SmartContract) QueryTransaction(ctx contractapi.TransactionContextInterface, transactionNumber string) (*Capex, error) {
	transactionAsBytes, err := ctx.GetStub().GetState(transactionNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if transactionAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", transactionNumber)
	}

	transaction := new(Capex)
	_ = json.Unmarshal(transactionAsBytes, transaction)

	return transaction, nil
}

// QueryAlltransaction returns all transaction found in world state
func (s *SmartContract) QueryAllTransactions(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		transaction := new(Capex)
		_ = json.Unmarshal(queryResponse.Value, transaction)

		queryResult := QueryResult{Key: queryResponse.Key, Record: transaction}
		results = append(results, queryResult)
	}

	return results, nil
}

// Changetransaction updates the owner field of mru with given id in world state
func (s *SmartContract) ChangeTransaction(ctx contractapi.TransactionContextInterface, transactionNumber string, newmru string) error {
	transaction, err := s.QueryTransaction(ctx, transactionNumber)

	if err != nil {
		return err
	}

	transaction.MRU = newmru

	transactionAsBytes, _ := json.Marshal(transaction)

	return ctx.GetStub().PutState(transactionNumber, transactionAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fa capex chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fa capex chaincode: %s", err.Error())
	}
}
