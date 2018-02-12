#!/bin/bash

function join_by { local IFS="$1"; shift; echo "$*"; }

run_register() {
  for (( i = 1; i <= $1; ++i));
  do
    docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"Function":"register", "Args":["'$i'"]}' &> out.log &
  done
}

run_voting() {
  for ((i=1;i<=$1;i++));
  do
    docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"Function":"submitVote", "Args":["'$i'", "1"]}' &> out.log &
  done
}

run_benchmark() {
  numberOfUsers=$1

  printf "############### Init vote \n"
  docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"Function":"initVote", "Args":[]}' &> out.log
  sleep 5

  userIDs=$(echo '"'$(seq -s '","' $numberOfUsers))
  userIDs=${userIDs%??};

  printf "############### Set eligible \n"
  docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"Function":"setEligible", "Args":['$userIDs']}' &> out.log
  sleep 3

  printf "############### Begin Signup \n"
  docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"Function":"beginSignUp", "Args":["Do you like Blockchain?"]}' &> out.log
  sleep 3

  printf "############### Register \n"
  # gdate is date for MacOS
  START=$(gdate +%s%N)
  run_register numberOfUsers
  END=$(gdate +%s%N)
  DIFF=$(echo "$END - $START" | bc)
  echo $DIFF
  sleep 10

  printf "############### Finish registration phase \n"
  docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"Function":"finishRegistrationPhase", "Args":[]}' &> out.log
  sleep 3

  printf "############### Vote \n"
  # gdate is date for MacOS
  START=$(gdate +%s%N)
  run_voting numberOfUsers
  END=$(gdate +%s%N)
  DIFF=$(echo "$END - $START" | bc)
  echo $DIFF
  sleep 2

  printf "############### Compute tally \n"
  docker exec cli peer chaincode query -C mychannel -n vote -c '{"Function":"computeTally","Args":[]}'
}

for i in 10 20 30 40 50 60 70 80 90 100 200
do
  printf "#### Running benchmark for "$i" ####\n"
  run_benchmark $i
done
