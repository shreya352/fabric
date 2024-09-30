package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func main() {
	// Create a new Fabric SDK
	sdk, err := fabsdk.New(config.FromFile("config.yaml"))
	if err != nil {
		log.Fatalf("Failed to create SDK: %v", err)
	}

	// Create a new gateway
	gateway, err := client.Connect(
		sdk,
		client.WithUser("Admin"),
		client.WithOrg("Org1"),
	)
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}

	// Create a new channel
	channel, err := gateway.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get channel: %v", err)
	}

	// Create a new contract
	contract := channel.GetContract("asset")

	http.HandleFunc("/createAsset", func(w http.ResponseWriter, r *http.Request) {
		// Create a new asset
		asset := Asset{
			DEALERID:    "dealer1",
			MSISDN:      "1234567890",
			MPIN:        "1234",
			BALANCE:     100,
			STATUS:      "active",
			TRANSAMOUNT: 0,
			TRANSTYPE:   "",
			REMARKS:     "",
		}

		// Invoke the createAsset function
		_, err := contract.SubmitTransaction("createAsset", asset.DEALERID, asset.MSISDN, asset.MPIN, fmt.Sprintf("%d", asset.BALANCE), asset.STATUS, fmt.Sprintf("%d", asset.TRANSAMOUNT), asset.TRANSTYPE, asset.REMARKS)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	http.HandleFunc("/getAsset", func(w http.ResponseWriter, r *http.Request) {
		// Get the asset ID from the request
		assetID := r.URL.Query().Get("assetID")

		// Invoke the getAsset function
		result, err := contract.EvaluateTransaction("getAsset", assetID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Unmarshal the result
		var asset Asset
		err = json.Unmarshal(result, &asset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the asset to the response
		json.NewEncoder(w).Encode(asset)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}