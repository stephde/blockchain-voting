/**
 * Created by stephde on 02.12.17.
 */

let Hyperledger = require("./hyperledger.js");
let hyperledger;

const users = ["user1", "user2", "user3"]

function run() {
    hyperledger = new Hyperledger();

    //ToDo: handle async code execution properly

    let userId;
    for(userId in users) {
        hyperledger.registerUser({id: userId})
    }

    hyperledger.initVote();
    hyperledger.setEligible(users)
    hyperledger.beginSignUp("Do you like Blockchain?")

    for(userId in users) {
        hyperledger.registerForVote(userId)
    }

    hyperledger
        .computeTally()
        .then(console.log, console.log);
}

run();
