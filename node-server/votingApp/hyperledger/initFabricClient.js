/**
 * Created by stephde on 02.12.17.
 */

let Fabric_Client = require('fabric-client');

exports.initFabricClient = function (host, channelId) {
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


    let clientWrapper = {
        client: fabricClient,
        channel: channel,
        peer: peer,
        eventHub: eventHub
    };

    console.log("Fabric client has been initialized with : ", clientWrapper);

    return clientWrapper;
}
//
