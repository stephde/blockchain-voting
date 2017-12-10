'use strict';

let Fabric = require('fabric-client');
let path = require('path');
let util = require('util');
let os = require('os');
// ToDo: should this go into the init method ?
var store_path = path.join(__dirname, 'hfc-key-store');

/**
 * This function executes a query against the given hyperledger client and returns a promise which
 * resolve to the response as a JSON object.
 *
 * @param fabricClient
 *      already initialized client to execute function on
 * @param chainCodeId
 *      id of the chain which should be queried
 * @param queryFunc
 *      query function identifier as string, which refers to the chaincode method
 * @param args
 *      arguments for the query
 *
 * @returns {Promise.<TResult>}
 */
exports.executeQuery = function (fabricClient, channel, chainCodeId, queryFunc, args) {
    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
    return Fabric.newDefaultKeyValueStore({
        path: store_path
    }).then((stateStore) => {
        // assign the store to the fabric client
        fabricClient.setStateStore(stateStore);

        return getUserContext(fabricClient, 'user1')
    }).then((userFromStore) => {
        // queryCar chaincode function - requires 1 argument, ex: args: ['CAR4'],
        // queryAllCars chaincode function - requires no arguments , ex: args: [''],
        const request = {
            //targets : --- letting this default to the peers assigned to the channel
            chaincodeId: chainCodeId, //'fabcar',
            fcn: queryFunc, //'queryAllCars',
            args: args //['']
        };

        return executeQueryFor(userFromStore, channel, request)
    }).then((response) => {
        return handleResponse(response)
    }).catch((err) => {
        console.error('Failed to query successfully :: ' + err);
    });
}


// -------------------- private functions --------------------- //

function getUserContext(fabricClient, userID) {
    let crypto_suite = Fabric.newCryptoSuite();
    // use the same location for the state store (where the users' certificate are kept)
    // and the crypto store (where the users' keys are kept)
    let crypto_store = Fabric.newCryptoKeyStore({path: store_path});
    crypto_suite.setCryptoKeyStore(crypto_store);
    fabricClient.setCryptoSuite(crypto_suite);

    // get the enrolled user from persistence, this user will sign all requests
    return fabricClient.getUserContext(userID, true);
}

function executeQueryFor(userFromStore, channel, request) {
    if (userFromStore && userFromStore.isEnrolled()) {
        //console.log('Successfully loaded user from persistence', userFromStore);
    } else {
        throw new Error('Failed to get user.... run registerUser.js');
    }

    // send the query proposal to the peer
    return channel.queryByChaincode(request);
}

function handleResponse(queryResponses) {
    console.log("Query has completed, checking results");
    // query_responses could have more than one  results if there multiple peers were used as targets
    if (queryResponses && queryResponses.length == 1) {
        if (queryResponses[0] instanceof Error) {
            console.error("error from query = ", queryResponses[0]);
        } else {
            console.log("Response is ", queryResponses[0].toString());
            return queryResponses[0];
        }
    } else {
        console.log("No payloads were returned from query");
    }

    return {};
}
