/**
 * Created by stephde on 02.12.17.
 */

let util = require('util');
let HyperledgerUtils = require("./hyperledergerUtils");

let tx_Id = null;
let globalEventHub = null;

/**
 * This function invokes a transaction on hypeledger with the given parameters.
 * Returns a promise which resolves to the transaction response as JSON object.
 *
 * @param fabricClient
 *      fabric client to execute transaction on
 * @param channel
 *      actual channel object to invoke transaction on
 * @param transactionFunc
 *      query function identifier as string, which refers to the chaincode method
 * @param args
 *      parameters of the transaction as array of string e.g. ['CAR1', 'user1']
 * @param userId
 *      id of the user who is trying to execute the transaction
 * @returns {Promise.<TResult>}
 */
exports.invokeTransaction = function (fabricClient, channel, eventHub, transactionFunc, args, userId) {
    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting

    globalEventHub = eventHub;

    return HyperledgerUtils.createDefaultKeyValueStore().then((stateStore) => {
        // assign the store to the fabric client
        fabricClient.setStateStore(stateStore);

        // creates the default cryptoStore and adds it to the client
        HyperledgerUtils.createDefaultCryptoKeyStore(fabricClient);

        // TODO replace string
        return fabricClient.getUserContext(userId, true);
    }).then((userFromStore) => {
        if (userFromStore && userFromStore.isEnrolled()) {
		// console.log('Successfully loaded user1 from persistence');
      		let member_user = userFromStore;
      	} else {
      		throw new Error('Failed to get user1.... run registerUser.js');
      	}

        // send the transaction proposal to the peers
        return proposeTransaction(fabricClient, channel, transactionFunc, args);
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

function proposeTransaction(fabricClient, channel, transactionFunc, args) {
    // get a transaction id object based on the current user assigned to fabric client
    tx_Id = fabricClient.newTransactionID();
    console.log("Assigning transaction_id: ", tx_Id._transaction_id);

    let chainConfig = HyperledgerUtils.getHyperledgerConfig();

    // must send the proposal to endorsing peers
    let request = {
        //targets: let default to the peer assigned to the client
        chaincodeId: chainConfig.chainCodeId, //'vote',
        fcn: transactionFunc, //'',
        args: args, //[''],
        chainId: chainConfig.chainId, //'vote',
        txId: tx_Id
    };

    console.log(request);
    // send the transaction proposal to the peers
    return channel.sendTransactionProposal(request);
}

function handleResponse(response) {
    // console.log('Send transaction promise and event listener promise have completed');
    // check the results in the order the promises were added to the promise all list
    if (response && response[0] && response[0].status === 'SUCCESS') {
        // console.log('Successfully sent transaction to the orderer.');
    } else {
        console.error('Failed to order the transaction. Error code: ' + response.status);
    }

    if(response && response[1] && response[1].event_status === 'VALID') {
        // console.log('Successfully committed the change to the ledger by the peer');
    } else {
        console.log('Transaction failed to be committed to the ledger due to ::'+response[1].event_status);
    }

    return response
}

function checkProposalResponse(proposalResponses) {
    let isProposalGood = false;
    if (proposalResponses && proposalResponses[0].response && proposalResponses[0].response.status === 200) {
        isProposalGood = true;

        // console.log(util.format(
        //     'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s"',
        //     proposalResponses[0].response.status, proposalResponses[0].response.message));
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
    let transaction_id_string = tx_Id.getTransactionID(); //Get the transaction ID string to be used by the event processing
    let promises = [];

    let sendPromise = channel.sendTransaction(request);
    promises.push(sendPromise); //we want the send transaction first, so that we know where to check status

    // get an eventhub once the fabric client has a user assigned. The user
    // is required bacause the event registration must be signed
    let event_hub = globalEventHub
    //fabricClient.getEventHub("peerio") --> we can only use this, of we use loadNetworkConfig in init client

    // using resolve the promise so that result status may be processed
    // under the then clause rather than having the catch clause process
    // the status
    let txPromise = new Promise((resolve, reject) => {
        let handle = setTimeout(() => {
            event_hub.disconnect();
            resolve({event_status : 'TIMEOUT'}); //we could use reject(new Error('Trnasaction did not complete within 30 seconds'));
        }, 10000);
        event_hub.connect();
        event_hub.registerTxEvent(transaction_id_string, (tx, code) => {
            // this is the callback for transaction event status
            // first some clean up of event listener
            console.log("##########################\n\nClearing timeout ...\n\n")
            clearTimeout(handle);

            event_hub.unregisterTxEvent(transaction_id_string);
            event_hub.disconnect();

            // now let the application know what happened
            let return_status = {event_status : code, tx_id : transaction_id_string};
            if (code !== 'VALID') {
                console.error('The transaction was invalid, code = ' + code);
                resolve(return_status); // we could use reject(new Error('Problem with the tranaction, event status ::'+code));
            } else {
                // console.log('The transaction has been committed on peer ' + event_hub._ep._endpoint.addr);
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
