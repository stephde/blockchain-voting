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
        .then(() => hyperledger.initVote())
        .then(() => hyperledger.setEligible(userIds))
        .then(() => hyperledger.beginSignUp("Do you like Blockchain?"))
        .then(() => runFuncParallelForUsers((userId) => hyperledger.registerForVote(userId), userIds))
        //ToDo: add voting step
        .then(() => hyperledger.computeTally())
        .then(console.log, console.log)
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
