package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/stub"
)

// Asset represents an asset with specific attributes
type Asset struct {
	DEALERID    string `json:"dealerId"`
	MSISDN      string `json:"msisdn"`
	MPIN        string `json:"mpin"`
	BALANCE     int    `json:"balance"`
	STATUS      string `json:"status"`
	TRANSAMOUNT int    `json:"transAmount"`
	TRANSTYPE   string `json:"transType"`
	REMARKS     string `json:"remarks"`
}

// Init initializes the chaincode
func (t *Asset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke invokes the chaincode
func (t *Asset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "createAsset":
		return t.createAsset(stub, args)
	case "updateAsset":
		return t.updateAsset(stub, args)
	case "getAsset":
		return t.getAsset(stub, args)
	case "getAssetHistory":
		return t.getAssetHistory(stub, args)
	default:
		return shim.Error("Invalid function name")
	}
}

// createAsset creates a new asset
func (t *Asset) createAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 8 {
		return shim.Error("Invalid number of arguments")
	}

	asset := Asset{
		DEALERID:    args[0],
		MSISDN:      args[1],
		MPIN:        args[2],
		BALANCE:     0,
		STATUS:      "active",
		TRANSAMOUNT: 0,
		TRANSTYPE:   "",
		REMARKS:     "",
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(args[0], assetJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// updateAsset updates an existing asset
func (t *Asset) updateAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 9 {
		return shim.Error("Invalid number of arguments")
	}

	assetJSON, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	asset := Asset{}
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return shim.Error(err.Error())
	}

	asset.BALANCE = asset.BALANCE + int(args[8])
	asset.TRANSAMOUNT = int(args[8])
	asset.TRANSTYPE = args[7]
	asset.REMARKS = args[6]

	assetJSON, err = json.Marshal(asset)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(args[0], assetJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// getAsset retrieves an asset
func (t *Asset) getAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Invalid number of arguments")
	}

	assetJSON, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(assetJSON)
}

// getAssetHistory retrieves the history of an asset
func (t *Asset) getAssetHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Invalid number of arguments")
	}

	history, err := stub.GetHistoryForKey(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	historyJSON, err := json.Marshal(history)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(historyJSON)
}

func main() {
	err := shim.Start(new(Asset))
	if err != nil {
		fmt.Printf("Error starting Asset chaincode: %s", err)
	}
}