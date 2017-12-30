
let Fabric = require('fabric-client');
let path = require('path');

/**
 * @type {{keyStoreName: string, chainCodeId: string, chainId: string}}
 *
 * @param chaincodeId
 *      id of the chaincode as string
 * @param chainId
 *      id of the chain / channel as string
 */
const hyperledgerConfig = {
    keyStoreName: 'hfc-key-store',
    chainCodeId: 'vote',
    chainId: 'vote'
}

const storePath = path.join(__dirname, hyperledgerConfig.keyStoreName);


let createCryptoKeyStore = function (fabricClient, storePath) {
    let crypto_suite = Fabric.newCryptoSuite();
    // use the same location for the state store (where the users' certificate are kept)
    // and the crypto store (where the users' keys are kept)
    let crypto_store = Fabric.newCryptoKeyStore({path: storePath});
    crypto_suite.setCryptoKeyStore(crypto_store);
    fabricClient.setCryptoSuite(crypto_suite);

    return crypto_suite;
}

/**
 * Create a CryptoKey store with default path. And add to client.
 *
 * @param client
 *      fabric client, in which the cryptoKeyStore will be attached
 *
 * @return CryptoSuite
 */
exports.createDefaultCryptoKeyStore = function (client) {

    return createCryptoKeyStore(client, storePath)
}

exports.createCryptoKeyStore = createCryptoKeyStore;

/**
 * Create a key value store with default path.
 *
 * @returns {Promise<IKeyValueStore>}
 */
exports.createDefaultKeyValueStore = function () {
    return Fabric.newDefaultKeyValueStore({
        path: storePath
    })
}

exports.getHyperledgerConfig = function () {
    return hyperledgerConfig;
}