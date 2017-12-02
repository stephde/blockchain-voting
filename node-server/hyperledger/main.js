/**
 * Created by stephde on 02.12.17.
 */

import { initFabricClient } from "./initFabricClient"
import { executeQuery } from "./query"
import { invokeTransaction } from "./invoke"

const host = 'grpc://localhost:7051';
const channel = 'mychannel';

function run() {
    let hlAdaper = initFabricClient(host, channel);

    //sample usage of query execution
    let queryPromise = executeQuery(hlAdaper.client, hlAdaper.channel, 'fabcar', 'queryAllCars', [''])
    //sample usage of transaction invocation
    let transactionPromise = invokeTransaction(hlAdaper.client, hlAdaper.channel, 'fabcar', '', 'mychannel', [''])
}

run();