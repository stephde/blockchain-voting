# must be run in cli docker container

peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n fabcar -c '{"Function":"changeCarOwner", "Args":["CAR4", "Max"]}'
peer chaincode query -C mychannel -n fabcar -c '{"Function":"queryCar","Args":["CAR4"]}'
peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n fabcar -c '{"Function":"changeCarOwner", "Args":["CAR4", "user1"]}'
peer chaincode query -C mychannel -n fabcar -c '{"Function":"queryCar","Args":["CAR4"]}'
