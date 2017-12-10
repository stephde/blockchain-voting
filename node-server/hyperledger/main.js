/**
 * Created by stephde on 02.12.17.
 */

let Hyperledger = require("./hyperledger.js");
let hyperledger;

function run() {
    hyperledger = new Hyperledger();
    hyperledger.queryAll();
    hyperledger.vote('Yellow');
    setTimeout(hyperledger.queryAll, 3000);



}

run();
