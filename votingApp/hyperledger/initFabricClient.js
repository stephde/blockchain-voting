/**
 * Created by stephde on 02.12.17.
 */

let Fabric_Client = require('fabric-client');
let HyperledgerUtils = require("./hyperledgerUtils");

exports.initFabricClient = function (host, channelId, userId, onTransactionCommitted) {
    let fabricClient = new Fabric_Client();

    // setup the fabric network
    let channel = fabricClient.newChannel(channelId); //
    let peer = fabricClient.newPeer(host);
    channel.addPeer(peer);
    let order = fabricClient.newOrderer('grpc://localhost:7050')
    channel.addOrderer(order);

    // get an eventhub once the fabric client has a user assigned. The user
    // is required bacause the event registration must be signed
    let eventHub = fabricClient.newEventHub();
    eventHub.setPeerAddr('grpc://localhost:7053');

    let blockListenerId = null;
    initUserContext(fabricClient, userId)
        .then((userContext) => {
            blockListenerId = initBlockEventListener(eventHub, channelId, onTransactionCommitted);
            eventHub.connect();
        })

    let clientWrapper = {
        client: fabricClient,
        channel: channel,
        peer: peer,
        eventHub: eventHub,
        blockListenerId: blockListenerId
    };

    //console.log("Fabric client has been initialized with : ", clientWrapper);
    console.log("Fabric client has been initialized.");

    return clientWrapper;
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
         console.log('\n\nOh snap!\n\n');
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
            throw new Error('Failed to get user1.... run registerUser.js');
        }

        return userFromStore;
    })
}
