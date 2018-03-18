/**
 * Created by stephde on 02.12.17.
 */

let Fabric_Client = require('fabric-client');
let HyperledgerUtils = require("./hyperledgerUtils");

const ordererPort = ':7050';
const peerPort = ':7051';
const eventHubPort = ':7053';

exports.initFabricClient = function (host, channelId, userId, onTransactionCommitted) {
    return new Promise(((resolve, reject) => {

        let fabricClient = new Fabric_Client();

        // setup the fabric network
        let channel = fabricClient.newChannel(channelId); //
        let peer = fabricClient.newPeer(host + peerPort);
        channel.addPeer(peer);
        let order = fabricClient.newOrderer(host + ordererPort);
        channel.addOrderer(order);


        initUserContext(fabricClient, userId)
            .then((userContext) => {
                // get an eventhub once the fabric client has a user assigned. The user
                // is required because the event registration must be signed
                let eventHub = fabricClient.newEventHub();
                eventHub.setPeerAddr(host + eventHubPort);

                let blockListenerId = initBlockEventListener(eventHub, channelId, onTransactionCommitted);
                eventHub.connect();

                let clientWrapper = {
                    client: fabricClient,
                    channel: channel,
                    peer: peer,
                    eventHub: eventHub,
                    blockListenerId: blockListenerId
                };

                console.log("Fabric client has been initialized.");
                resolve(clientWrapper);
            })

    }))
}

function initBlockEventListener(eventHub, myChannelId, onTransactionCommitted) {
    eventHub.registerBlockEvent((block) => {
        let first_tx = block.data.data[0]; // get the first transaction
        let header = first_tx.payload.header; // the "header" object contains metadata of the transaction
        let channel_id = header.channel_header.channel_id;
        if (myChannelId !== channel_id) return;

        console.log("Received - " + block.data.data.length + " - transactions in the block.")
        console.log("First header: " + JSON.stringify(first_tx))

        block.data.data.forEach(onTransactionCommitted)
    }, (err) => {
         console.log(err)
   });
}

function initUserContext(fabricClient, userId) {
    return HyperledgerUtils.createDefaultKeyValueStore().then((stateStore) => {
        // assign the store to the fabric client
        fabricClient.setStateStore(stateStore);

        // creates the default cryptoStore and adds it to the client
        HyperledgerUtils.createDefaultCryptoKeyStore(fabricClient);

        return fabricClient.getUserContext(userId, true);
    }).then((userFromStore) => {
        if (userFromStore && userFromStore.isEnrolled()) {
            // console.log('Successfully loaded user1 from persistence');
            let member_user = userFromStore;
        } else {
            throw new Error('Failed to get user.... run registerUser.js');
        }

        return userFromStore;
    })
}
