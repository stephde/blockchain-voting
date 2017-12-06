#!/bin/sh

# must be run in cli docker container

docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"Function":"vote", "Args":["Blue"]}'
docker exec cli peer chaincode query -C mychannel -n vote -c '{"Function":"queryVotes","Args":[]}'
docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"Function":"vote", "Args":["Red"]}'
sleep 2
docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"Function":"vote", "Args":["Red"]}'
sleep 2
docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"Function":"vote", "Args":["Red"]}'
sleep 10
docker exec cli peer chaincode query -C mychannel -n vote -c '{"Function":"queryVotes","Args":[]}'
