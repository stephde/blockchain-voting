'use strict';

// ToDo: should this go into the init method ?
var store_path = path.join(__dirname, 'hfc-key-store');
console.log('Store path:'+store_path);


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
export function executeQuery(fabricClient, chainCodeId, queryFunc, args) {
    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
    return fabricClient.newDefaultKeyValueStore({
        path: store_path
    }).then((stateStore) => {
        return getUserContext('user1')
    }).then((userFromStore) => {
        // queryCar chaincode function - requires 1 argument, ex: args: ['CAR4'],
        // queryAllCars chaincode function - requires no arguments , ex: args: [''],
        const request = {
            //targets : --- letting this default to the peers assigned to the channel
            chaincodeId: chainCodeId, //'fabcar',
            fcn: queryFunc, //'queryAllCars',
            args: args //['']
        };

        executeQueryFor(userFromStore, request)
    }).then((response) => {
        return handleResponse(response)
    }).catch((err) => {
        console.error('Failed to query successfully :: ' + err);
    });
}


// -------------------- private functions --------------------- //

function getUserContext(userID) {
    // assign the store to the fabric client
    fabric_client.setStateStore(state_store);
    let crypto_suite = Fabric_Client.newCryptoSuite();
    // use the same location for the state store (where the users' certificate are kept)
    // and the crypto store (where the users' keys are kept)
    let crypto_store = Fabric_Client.newCryptoKeyStore({path: store_path});
    crypto_suite.setCryptoKeyStore(crypto_store);
    fabric_client.setCryptoSuite(crypto_suite);

    // get the enrolled user from persistence, this user will sign all requests
    return fabric_client.getUserContext(userID, true);
}

function executeQueryFor(userFromStore, request) {
    if (userFromStore && userFromStore.isEnrolled()) {
        console.log('Successfully loaded user from persistence', userFromStore);
    } else {
        throw new Error('Failed to get user.... run registerUser.js');
    }

    // send the query proposal to the peer
    return channel.queryByChaincode(request);
}

function handleResponse(queryResponse) {
    console.log("Query has completed, checking results");
    // query_responses could have more than one  results if there multiple peers were used as targets
    
    if (queryResponses && queryResponses.length == 1) {
        if (queryResponses[0] instanceof Error) {
            console.error("error from query = ", queryResponses[0]);
        } else {
            console.log("Response is ", queryResponses[0].toString());
            return queryResponse[0];
        }
    } else {
        console.log("No payloads were returned from query");
    }

    return {};
}