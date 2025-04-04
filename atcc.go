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

// Creates a new asset with the provided details
func (s *SmartContract) CreateAsset(
	ctx contractapi.TransactionContextInterface,
	id string,
	color string,
	size int,
	owner string,
	appraisedValue int) error {
	// Check if the asset already exists on the ledger
	exists, err := s.AssetExists(ctx, id)

	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("The asset with id %s already exists", id)
	}

	// Create the new asset
	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}

	// Marshall the asset
	assetJSON, err := json.Marshal(asset)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// Reads the asset with the given id
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	// Try and retrieve the asset with the given id from the ledger
	assetJSON, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to read world state for asset with id %v", id)
	}

	if assetJSON == nil {
		return nil, fmt.Errorf("The asset with id %v does not exist", id)
	}

	// Unmarshal the found asset
	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)

	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// Updates an asset with the given parameters
func (s *SmartContract) UpdateAsset(
	ctx contractapi.TransactionContextInterface,
	id string,
	color string,
	size int,
	owner string,
	appraisedValue int) error {
	// Check if the asset exists on the ledger
	exists, err := s.AssetExists(ctx, id)

	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("The asset with id %s does not exist", id)
	}

	// Create the new asset
	updatedAsset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}

	// Marshall the updated asset
	updatedAssetJSON, err := json.Marshal(updatedAsset)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, updatedAssetJSON)
}

// Deletes an asset with the given id
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	// Check if the asset to be deleted is on the ledger or not
	exists, err := s.AssetExists(ctx, id)

	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("The asset with id %s does not exist", id)
	}

	// Delete the asset with the given id
	return ctx.GetStub().DelState(id)

	// Important note, DelState will record the id on the writestate and will delete the key once the transaction
	// has been validated and comitted
}

func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	// Read the asset from the ledger
	transferedAsset, err := s.ReadAsset(ctx, id)

	if err != nil {
		return err
	}

	// Change the owner
	transferedAsset.Owner = newOwner

	// Marshall the updated asset
	transferedAssetJSON, err := json.Marshal(transferedAsset)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, transferedAssetJSON)
}

func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// Get an iterator for all the keys
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}
	// Close the iterator when done with reading data
	defer resultsIterator.Close()

	// Go through all the key-value pairs in the ledger
	var assets []*Asset
	for resultsIterator.HasNext() {
		// Get the next KV pair
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		// Unmarshall the read asset
		var asset Asset

		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		// Apped the pointer to the read asset
		assets = append(assets, &asset)
	}

	return assets, nil
}

// Checks if the asset with the given id exists in the ledger
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	// Try and retrieve the asset with the given id from the ledger
	assetJSON, err := ctx.GetStub().GetState(id)

	if err != nil {
		return false, fmt.Errorf("Failed to read world state for asset with id %v", id)
	}

	// Return the asset (if found) together with no error
	return assetJSON != nil, nil
}

func main() {
	// Create new chaincode
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating the chaincode: %v", err)
	}

	// Start the chaincode
	err = assetChaincode.Start()
	if err != nil {
		log.Panicf("Error starting the chaincode: %v", err)
	}

}
