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
CC_SRC_PATH=lehmann.com/voting/go

# clean the keystore
rm -rf ./hfc-key-store

# launch network; create channel and join peer to channel
cd ../../fabric-test/fabric-samples/basic-network/
./start.sh

# Now launch the CLI container in order to install, instantiate chaincode
# and prime the ledger with our 10 cars
docker-compose -f ./docker-compose.yml up -d cli

docker exec cli peer chaincode install -n vote -v 1.0 -p "$CC_SRC_PATH" -l "$LANGUAGE"
docker exec cli peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n vote -l "$LANGUAGE" -v 1.0 -c '{"Args":[""]}'
sleep 10
docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"function":"initLedger","Args":[""]}'

printf "\nTotal setup execution time : $(($(date +%s) - starttime)) secs ...\n\n\n"
printf "Start by installing required packages run 'npm install'\n"
printf "Then run 'node enrollAdmin.js', then 'node registerUser'\n\n"
printf "The 'node invoke.js' will fail until it has been updated with valid arguments\n"
printf "The 'node query.js' may be run at anytime once the user has been registered\n\n"
