/**
 * Created by stephde on 02.12.17.
 */

let Hyperledger = require("./hyperledger.js");
let hyperledger = new Hyperledger();


function runElection() {
    const numOfUsers = 3;

    let userIds = []
    for (let i=0; i < numOfUsers; i++) {
        userIds.push("user" + i);
    }

    runFuncParallelForUsers((userId) => hyperledger.registerUser({id: userId}), userIds)
        .then(() => timedCall(hyperledger.initVote, [], 'Init Vote'))
        .then(() => timedCall(hyperledger.setEligible, userIds, 'Set Eligible'))
        .then(() => timedCall(hyperledger.beginSignUp, "Do you like Blockchain?", 'begin sign up'))
        .then(() => timedCall(()=>runFuncParallelForUsers((userId) => hyperledger.registerForVote(userId), userIds), [], 'register for vote'))
        //ToDo: add voting step
        .then(() => timedCall(hyperledger.computeTally, [], "compute tally"))
        .then(console.log, console.log)
}

function timedCall(func, params, identifier){
  var start = new Date().getTime();
  var promise = func(params);
  var end = new Date().getTime();
  var totalTime = end-start;
  console.log("Time spend for "+identifier + ": "+ totalTime+ "ms")
  return promise;
}

function runFuncParallelForUsers(func, userIds) {
    let promises = [];

    let index;
    for (index in userIds) {
        promises.push(func(userIds[index]))
    }

    return Promise.all(promises);
}

runElection();
