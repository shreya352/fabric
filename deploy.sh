#!/bin/bash

# Create the network
docker-compose -f network/docker-compose.yml up -d

# Create the channel
docker exec -it peer0.org1.example.com peer channel create -o orderer.example.com:7050 -c mychannel -f ./network/channel-artifacts/channel.tx

# Join the channel
docker exec -it peer0.org1.example.com peer channel join -b mychannel.block
docker exec -it peer0.org2.example.com peer channel join -b mychannel.block

 # Install the chaincode
docker exec -it peer0.org1.example.com peer chaincode install -n asset -v 1.0 -p github.com/chaincode/asset
docker exec -it peer0.org2.example.com peer chaincode install -n asset -v 1.0 -p github.com/chaincode/asset

# Instantiate the chaincode
docker exec -it peer0.org1.example.com peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n asset -v 1.0 -c '{"Args":[]}'