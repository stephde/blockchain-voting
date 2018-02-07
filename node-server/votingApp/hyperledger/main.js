/**
 * Created by stephde on 02.12.17.
 */

let Hyperledger = require("./hyperledger.js");
let hyperledger = new Hyperledger();

function runElection() {
    const numOfUsers = 100;
    let userIds = []
    for (let i=0; i < numOfUsers; i++) {
        userIds.push("user" + i);
    }

    let start;

    runFuncParallelForUsers((userId) => hyperledger.registerUser({id: userId}), userIds)
        .then(() => hyperledger.initVote())
        .then(() => hyperledger.setEligible(userIds))
        .then(() => hyperledger.beginSignUp("Do you like Blockchain?"))
        .then(() => runFuncParallelForUsers(
                (userId) => hyperledger.registerForVote(userId), userIds))
        .then(() => promisedTimeout(5000))
        .then(() => hyperledger.finishRegistrationPhase())
        .then(() => start = new Date().getTime())
        .then(() => runFuncParallelForUsers(
                (userId) => hyperledger.vote(userId, '1'), userIds))
        .then(() => promisedTimeout(3000))
        .then(() => hyperledger.computeTally())
        .then(() => printTimeSince(start, 'voting phase'))
        .catch(console.log)
}

function printTimeSince(start, identifier){
    let end = new Date().getTime();
    console.log("\n\n###########################\n\n")
    console.log("Total time for " + identifier + " is " + (end-start) + " ms")
    console.log("\n\n###########################\n\n")

    return Promise.resolve();
}

async function runFuncParallelForUsers(func, userIds) {
    let promises = [];

    let index;
    for (index in userIds) {
        await promisedTimeout(50)
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
