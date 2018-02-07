#!/bin/bash

function join_by { local IFS="$1"; shift; echo "$*"; }

run_register() {
  for (( i = 1; i <= $1; ++i));
  do
    docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"Function":"register", "Args":["'$i'"]}' &
  done
}

run_voting() {
  for ((i=1;i<=$1;i++));
  do
    docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"Function":"submitVote", "Args":["'$i'", "1"]}'
  done
}

run_benchmark() {
  numberOfUsers=$1

  printf "############### Init vote \n"
  docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"Function":"initVote", "Args":[]}' &
  sleep 5

  userIDs=$(echo '"'$(seq -s '","' $numberOfUsers))
  userIDs=${userIDs%??};

  printf $userIDs"\n"

  printf "############### Set eligible \n"
  docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"Function":"setEligible", "Args":['$userIDs']}' &
  sleep 3

  printf "############### Begin Signup \n"
  docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"Function":"beginSignUp", "Args":["Do you like Blockchain?"]}' &
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
  docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vote -c '{"Function":"finishRegistrationPhase", "Args":[]}'
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

run_benchmark 100
