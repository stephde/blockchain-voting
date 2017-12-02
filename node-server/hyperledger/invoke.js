/**
 * Created by stephde on 02.12.17.
 */

let path = require('path');
let util = require('util');
let os = require('os');

let store_path = path.join(__dirname, 'hfc-key-store');
console.log('Store path:'+store_path);
let tx_id = null;

/**
 * This function invokes a transaction on hypeledger with the given parameters.
 * Returns a promise which resolves to the transaction response as JSON object.
 *
 * @param fabricClient
 *      fabric client to execute transaction on
 * @param channel
 *      actual channel object to invoke transaction on
 * @param chaincodeId
 *      id of the chaincode as string
 * @param transactionFunc
 *      query function identifier as string, which refers to the chaincode method
 * @param chainId
 *      id of the chain / channel as string
 * @param args
 *      parameters of the transaction as array of string e.g. ['CAR1', 'user1']
 * @returns {Promise.<TResult>}
 */
exports.invokeTransaction = function (fabricClient, channel, chaincodeId, transactionFunc, chainId, args) {
    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
    return fabricClient.newDefaultKeyValueStore({ path: store_path
    }).then((state_store) => {
        return getUserContext('user1');
    }).then((userFromStore) => {

        // createCar chaincode function - requires 5 args, ex: args: ['CAR12', 'Honda', 'Accord', 'Black', 'Tom'],
        // changeCarOwner chaincode function - requires 2 args , ex: args: ['CAR10', 'Barry'],
        // must send the proposal to endorsing peers
        let request = {
            //targets: let default to the peer assigned to the client
            chaincodeId: chaincodeId, //'fabcar',
            fcn: transactionFunc, //'',
            args: args, //[''],
            chainId: chainId, //'mychannel',
            txId: tx_id
        };

        // send the transaction proposal to the peers
        return proposeTransaction(fabricClient, channel, userFromStore, request);
    }).then((results) => {
        let proposalResponses = results[0];
        let proposal = results[1];

        //throws an error if proposal was rejected
        checkProposalResponse(proposalResponses);

        // build up the request for the orderer to have the transaction committed
        let request = {
            proposalResponses: proposalResponses,
            proposal: proposal
        };

        return sendTransaction(fabricClient, channel, request);
    }).then((response) => {
        return handleResponse(response)
    }).catch((err) => {
        console.error('Failed to invoke successfully :: ' + err);
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

function proposeTransaction(fabricClient, channel, userFromStore, request) {
    if (userFromStore && userFromStore.isEnrolled()) {
        console.log('Successfully loaded user1 from persistence');
    } else {
        throw new Error('Failed to get user1.... run registerUser.js');
    }

    // get a transaction id object based on the current user assigned to fabric client
    tx_id = fabricClient.newTransactionID();
    console.log("Assigning transaction_id: ", tx_id._transaction_id);

    // send the transaction proposal to the peers
    return channel.sendTransactionProposal(request);
}

function handleResponse(response) {
    console.log('Send transaction promise and event listener promise have completed');
    // check the results in the order the promises were added to the promise all list
    if (response && response[0] && response[0].status === 'SUCCESS') {
        console.log('Successfully sent transaction to the orderer.');
    } else {
        console.error('Failed to order the transaction. Error code: ' + response.status);
    }

    if(response && response[1] && response[1].event_status === 'VALID') {
        console.log('Successfully committed the change to the ledger by the peer');
    } else {
        console.log('Transaction failed to be committed to the ledger due to ::'+response[1].event_status);
    }

    return response
}

function checkProposalResponse(proposalResponses) {
    let isProposalGood = false;
    if (proposalResponses && proposalResponses[0].response &&
        proposalResponses[0].response.status === 200) {
        isProposalGood = true;
        console.log('Transaction proposal was good');

        console.log(util.format(
            'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s"',
            proposalResponses[0].response.status, proposalResponses[0].response.message));
    } else {
        console.error('Transaction proposal was bad');
    }

    if( ! isProposalGood) {
        console.error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
        throw new Error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
    }
}

//ToDo: refactor this function
function sendTransaction(fabricClient, channel, request) {
    // set the transaction listener and set a timeout of 30 sec
    // if the transaction did not get committed within the timeout period,
    // report a TIMEOUT status
    let transaction_id_string = tx_id.getTransactionID(); //Get the transaction ID string to be used by the event processing
    var promises = [];

    let sendPromise = channel.sendTransaction(request);
    promises.push(sendPromise); //we want the send transaction first, so that we know where to check status

    // get an eventhub once the fabric client has a user assigned. The user
    // is required bacause the event registration must be signed
    let event_hub = fabricClient.newEventHub();
    event_hub.setPeerAddr('grpc://localhost:7053');

    // using resolve the promise so that result status may be processed
    // under the then clause rather than having the catch clause process
    // the status
    let txPromise = new Promise((resolve, reject) => {
        let handle = setTimeout(() => {
            event_hub.disconnect();
            resolve({event_status : 'TIMEOUT'}); //we could use reject(new Error('Trnasaction did not complete within 30 seconds'));
        }, 3000);
        event_hub.connect();
        event_hub.registerTxEvent(transaction_id_string, (tx, code) => {
            // this is the callback for transaction event status
            // first some clean up of event listener
            clearTimeout(handle);
            event_hub.unregisterTxEvent(transaction_id_string);
            event_hub.disconnect();

            // now let the application know what happened
            let return_status = {event_status : code, tx_id : transaction_id_string};
            if (code !== 'VALID') {
                console.error('The transaction was invalid, code = ' + code);
                resolve(return_status); // we could use reject(new Error('Problem with the tranaction, event status ::'+code));
            } else {
                console.log('The transaction has been committed on peer ' + event_hub._ep._endpoint.addr);
                resolve(return_status);
            }
        }, (err) => {
            //this is the callback if something goes wrong with the event registration or processing
            reject(new Error('There was a problem with the eventhub ::'+err));
        });
    });
    promises.push(txPromise);

    return Promise.all(promises);
}