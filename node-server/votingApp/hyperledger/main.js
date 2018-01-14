/**
 * Created by stephde on 02.12.17.
 */

let Hyperledger = require("./hyperledger.js");
let hyperledger;

const users = ["user1", "user2", "user3"]

function run() {
    hyperledger = new Hyperledger();

    hyperledger.registerUser({id: users[0]})
        .then(() => hyperledger.registerUser({id: users[1]}))
        .then(() => hyperledger.registerUser({id: users[2]}))
        .then(() => hyperledger.initVote())
        .then(() => hyperledger.setEligible(users))
        .then(() => hyperledger.beginSignUp("Do you like Blockchain?"))
        .then(() => hyperledger.registerForVote(users[0]))
        .then(() => hyperledger.registerForVote(users[1]))
        .then(() => hyperledger.registerForVote(users[2]))
        .then(() => hyperledger.computeTally())
        .then(console.log, console.log)
}

run();
