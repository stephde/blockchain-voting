'use strict';

let HyperledgerUtils = require("./hyperledergerUtils");

/**
 * This function executes a query against the given hyperledger client and returns a promise which
 * resolve to the response as a JSON object.
 *
 * @param fabricClient
 *      already initialized client to execute function on
 * @param chainCodeId
 *      id of the chain which should be queried
 * @param channel
 *      the channel which should be queried
 * @param queryFunc
 *      query function identifier as string, which refers to the chaincode method
 * @param args
 *      arguments for the query
 *
 * @returns {Promise.<TResult>}
 */
exports.executeQuery = function (fabricClient, channel, chainCodeId, queryFunc, args) {
    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
    return HyperledgerUtils.createDefaultKeyValueStore()
    .then((stateStore) => {
        // assign the store to the fabric client
        fabricClient.setStateStore(stateStore);

        // creates the default cryptoStore and adds it to the client
        HyperledgerUtils.createDefaultCryptoKeyStore(fabricClient);

        return fabricClient.getUserContext("user1", true);
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
        return handleResponse(response);
    }).catch((err) => {
        console.error('Failed to query successfully :: ' + err);
    });
}


// -------------------- private functions --------------------- //

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

    return queryResponses;
}
