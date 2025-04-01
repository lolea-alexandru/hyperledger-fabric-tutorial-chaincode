package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
)

// The Contract type makes sure our SmartContract type matches its definition
type SmartContract struct {
	contractapi.Contract
}

// Define the basic fields of an asset
type Asset struct {
	AppraisedValue int    `json:"AppraisedValue"`
	Color          string `json:"Color"`
	ID             string `json:"ID"`
	Owner          string `json:"Owner"`
	Size           int    `json:"Size"`
}

// Initialize the ledger with a few dummy assets
func (smartContractPointer *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// Define the dummy assets to be added
	assets := []Asset{
		{ID: "asset1", Color: "blue", Size: 5, Owner: "Tomoko", AppraisedValue: 300},
		{ID: "asset2", Color: "red", Size: 5, Owner: "Brad", AppraisedValue: 400},
		{ID: "asset3", Color: "green", Size: 10, Owner: "Jin Soo", AppraisedValue: 500},
		{ID: "asset4", Color: "yellow", Size: 10, Owner: "Max", AppraisedValue: 600},
		{ID: "asset5", Color: "black", Size: 15, Owner: "Adriana", AppraisedValue: 700},
		{ID: "asset6", Color: "white", Size: 15, Owner: "Michel", AppraisedValue: 800},
	}

	// Go through every asset in the "assets" slice
	for _, asset := range assets {
		// Try to marshall the current asset
		assetJSON, err := json.Marshal(asset)

		// Check if an error has occured
		if err != nil {
			return err
		}

		// Get a stub for the current context
		stub := ctx.GetStub()

		// Add the asset to the transaction's writeset as a data-write proposal
		// PutState affects the ledger only after the transaction is validated and submitted
		err = stub.PutState(asset.ID, assetJSON)

		// Check if an error has occured
		if err != nil {
			return fmt.Errorf("Failed to put to world state. %v", err)
		}
	}

	return nil
}
