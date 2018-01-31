/**
 * Created by stephde on 02.12.17.
 */

let Hyperledger = require("./hyperledger.js");
let hyperledger = new Hyperledger();

function runElection() {
    const numOfUsers = 50;

    let userIds = []
    for (let i=0; i < numOfUsers; i++) {
        userIds.push("user" + i);
    }

    runFuncParallelForUsers((userId) => hyperledger.registerUser({id: userId}), userIds)
        .then(() => timedCall(hyperledger.initVote, [], 'Init Vote'))
        .then(() => timedCall(hyperledger.setEligible, userIds, 'Set Eligible'))
        .then(() => timedCall(hyperledger.beginSignUp, "Do you like Blockchain?", 'begin sign up'))
        .then(() => timedCall(() => runFuncParallelForUsers(
            (userId) => hyperledger.registerForVote(userId), userIds), [], 'register for vote'))
        .then(() => promisedTimeout(2000))
        .then(() => timedCall(hyperledger.finishRegistrationPhase, [], 'finishRegistrationPhase'))
        .then(() => timedCall(() =>
            runFuncParallelForUsers(
                (userId) => hyperledger.vote(userId, '1'), userIds), [], 'voting'))
        .then(() => promisedTimeout(5000))
        .then(() => timedCall(hyperledger.computeTally, [], "compute tally"))
        .then(console.log)
        .catch(console.log)
}

function timedCall(func, params, identifier){
    let start = new Date().getTime();
    let promise = func(params);
    let end = new Date().getTime();

    console.log("Time spend for " + identifier + ": " + end-start + "ms")
    return promise;
}

async function runFuncParallelForUsers(func, userIds) {
    let promises = [];

    let index;
    for (index in userIds) {
        await promisedTimeout(200)
        promises.push(func(userIds[index]))
    }

    return Promise.all(promises);
}

// promisified version of setTimeOut
function promisedTimeout(ms) {
    console.log("Waiting for " + ms + "ms ...")
    return new Promise(resolve => setTimeout(resolve, ms));
}

runElection();
