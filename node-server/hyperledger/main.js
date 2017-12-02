/**
 * Created by stephde on 02.12.17.
 */

import { initFabricClient } from "./initFabricClient"
import { executeQuery } from "./query"

const host = 'grpc://localhost:7051';
const channel = 'mychannel';

function run() {
    let fabricClient = initFabricClient(host, channel);
    
    let result = executeQuery(fabricClient, 'fabcar', 'queryAllCars', [''])
}

run();