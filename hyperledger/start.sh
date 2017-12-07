#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -ev

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1
starttime=$(date +%s)
LANGUAGE=${1:-"golang"}
CC_SRC_PATH=hpi.de/go

# clean the keystore
rm -rf ./hfc-key-store
docker rm -f $(docker ps -aq) || true
#docker network prune
#CHAINCODE_DOCKER_IMAGE=dev-peer0.org1.example.com-vote-1.0-833e0675ada610923f4bff0345ef4e64eb6947b413ae79ea5b9539f38d148627
#docker stop CHAINCODE_DOCKER_IMAGE || true && docker rm CHAINCODE_DOCKER_IMAGE || true && docker rmi CHAINCODE_DOCKER_IMAGE || true


# Start Docker containers
docker-compose -f docker-compose.yml down
docker-compose -f docker-compose.yml up -d ca.example.com orderer.example.com peer0.org1.example.com couchdb

# wait for Hyperledger Fabric to start
# in case of errors when running later commands, issue export FABRIC_START_TIMEOUT=<larger number>
export FABRIC_START_TIMEOUT=10
#echo ${FABRIC_START_TIMEOUT}
sleep ${FABRIC_START_TIMEOUT}

# Create the channel
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.example.com/msp" peer0.org1.example.com peer channel create -o orderer.example.com:7050 -c mychannel -f /etc/hyperledger/configtx/channel.tx
# Join peer0.org1.example.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.example.com/msp" peer0.org1.example.com peer channel join -b mychannel.block

# Now launch the CLI container in order to install, instantiate chaincode
docker-compose -f ./docker-compose.yml up -d cli

printf "Installing chaincode"
docker exec cli peer chaincode install -n vote -v 1.0 -p "$CC_SRC_PATH" -l "$LANGUAGE"

printf "Instantiate chaincode"
docker exec cli peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n vote -l "$LANGUAGE" -v 1.0 -c '{"Args":[""]}'
sleep 10

printf "Invoke chaincode"
docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"function":"initLedger","Args":[""]}'

printf "\nTotal setup execution time : $(($(date +%s) - starttime)) secs ...\n\n\n"
