#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1
starttime=$(date +%s)
LANGUAGE=${1:-"golang"}
CC_SRC_PATH=hpi.de/go

# clean the keystore
rm -rf ./hfc-key-store
docker rm -f $(docker ps -aq) || true

# Start Docker containers
docker-compose -f docker-compose.yml down
docker-compose -f docker-compose.yml up -d ca.example.com orderer.example.com peer0.org1.example.com couchdb

# wait for Hyperledger Fabric to start
# in case of errors when running later commands, issue export FABRIC_START_TIMEOUT=<larger number>
export FABRIC_START_TIMEOUT=10
sleep ${FABRIC_START_TIMEOUT}

# Create the channel
printf "\n############### Creating channel\n"
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.example.com/msp" peer0.org1.example.com peer channel create -o orderer.example.com:7050 -c mychannel -f /etc/hyperledger/configtx/channel.tx
# Join peer0.org1.example.com to the channel.
printf "\n############### Joining peer to the channel\n"
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.example.com/msp" peer0.org1.example.com peer channel join -b mychannel.block

# Now launch the CLI container in order to install, instantiate chaincode
docker-compose -f ./docker-compose.yml up -d cli

# Install go dependencies
# docker exec cli /opt/gopath/src/hpi.de/go/upgradeGo.sh
docker exec cli /opt/gopath/src/hpi.de/go/loadDependencies.sh
docker exec peer0.org1.example.com /opt/gopath/src/hpi.de/go/upgradeGo.sh
# docker exec peer0.org1.example.com /opt/gopath/src/hpi.de/go/loadDependencies.sh

printf "\n############### Installing chaincode\n"
docker exec cli peer chaincode install -n vote -v 1.0 -p "$CC_SRC_PATH" -l "$LANGUAGE"

printf "\n############### Instantiate chaincode\n"
docker exec cli peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n vote -l "$LANGUAGE" -v 1.0 -c '{"Args":["Foo"]}'

printf "\n############### Total setup execution time : $(($(date +%s) - starttime)) secs ...\n\n\n"
