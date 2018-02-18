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

    let start;

    // runFuncParallelForUsers((userId) => hyperledger.registerUser({id: userId}), userIds)
    hyperledger.initVote()
        .then(() => waitForTransactions())
        .then(() => hyperledger.setEligible(userIds))
        .then(() => waitForTransactions())
        .then(() => hyperledger.beginSignUp("Do you like Blockchain?"))
        .then(() => waitForTransactions())
        .then(() => runFuncParallelForUsers(
                (userId) => hyperledger.registerForVote(userId), userIds))
        .then(() => waitForTransactions())
        .then(() => hyperledger.finishRegistrationPhase())
        .then(() => waitForTransactions())
        .then(() => start = new Date().getTime())
        .then(() => runFuncParallelForUsers(
                (userId) => hyperledger.vote(userId, '1'), userIds))
        .then(() => waitForTransactions())
        .then(() => printTimeSince(start, 'voting phase'))
        .then(() => hyperledger.computeTally())
        .then(() => waitForTransactions())
        .then(() => hyperledger.close())
        .catch(console.log)
}

async function waitForTransactions() {
    while(true) {
        if(hyperledger.isTransactionPending()){
            await promisedTimeout(500)
            continue;
        }
        break;
    }
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
        await promisedTimeout(100)
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
