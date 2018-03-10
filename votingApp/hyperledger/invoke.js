/**
 * Created by stephde on 02.12.17.
 */

let util = require('util');
let HyperledgerUtils = require("./hyperledgerUtils");

let tx_Id = null;
let _onTransactionSubmitted = null;

/**
 * This function invokes a transaction on hypeledger with the given parameters.
 * Returns a promise which resolves to the transaction response as JSON object.
 *
 * @param fabricClient
 *      fabric client to execute transaction on
 * @param channel
 *      actual channel object to invoke transaction on
 * @param onTransactionSubmitted
 *      function which is to be called when a transaction was submitted
 * @param transactionFunc
 *      query function identifier as string, which refers to the chaincode method
 * @param args
 *      parameters of the transaction as array of string e.g. ['CAR1', 'user1']
 * @param userId
 *      id of the user who is trying to execute the transaction
 * @returns {Promise.<TResult>}
 */
exports.invokeTransaction = function (fabricClient, channel, onTransactionSubmitted, transactionFunc, args, userId) {
    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting

    _onTransactionSubmitted = onTransactionSubmitted;

    return HyperledgerUtils.createDefaultKeyValueStore().then((stateStore) => {
        // assign the store to the fabric client
        fabricClient.setStateStore(stateStore);

        // creates the default cryptoStore and adds it to the client
        HyperledgerUtils.createDefaultCryptoKeyStore(fabricClient);

        return fabricClient.getUserContext(userId, true);
    }).then((userFromStore) => {
        if (userFromStore && userFromStore.isEnrolled()) {
      		console.log('Successfully loaded user1 from persistence');
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

function sendTransaction(fabricClient, channel, request) {
    let transaction_id_string = tx_Id.getTransactionID(); //Get the transaction ID string to be used by the event processing
    let sendPromise = channel.sendTransaction(request);

    _onTransactionSubmitted(transaction_id_string);

    return sendPromise;
}
