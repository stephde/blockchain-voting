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

    return {
        client: fabricClient,
        channel: channel,
        peer: peer
    };
}
//

