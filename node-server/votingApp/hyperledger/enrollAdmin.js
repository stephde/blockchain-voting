'use strict';

let Fabric_CA_Client = require('fabric-ca-client');
let HyperledgerUtils = require("./hyperledergerUtils");

let fabric_ca_client = null;
let admin_user = null;

exports.enrollAdmin = function(fabric_client){
  HyperledgerUtils.createDefaultKeyValueStore().then((state_store) => {
      // assign the store to the fabric client
      fabric_client.setStateStore(state_store);
      let crypto_suite = HyperledgerUtils.createDefaultCryptoKeyStore(fabric_client);
      let tlsOptions = {
      	trustedRoots: [],
      	verify: false
      };
      // be sure to change the http to https when the CA is running TLS enabled
      fabric_ca_client = new Fabric_CA_Client('http://localhost:7054', tlsOptions , 'ca.example.com', crypto_suite);

      // first check to see if the admin is already enrolled
      return fabric_client.getUserContext('admin', true);
  }).then((user_from_store) => {
      if (user_from_store && user_from_store.isEnrolled()) {
          console.log('Successfully loaded admin from persistence');
          admin_user = user_from_store;
          return null;
      } else {
          // need to enroll it with CA server
          return fabric_ca_client.enroll({
            enrollmentID: 'admin',
            enrollmentSecret: 'adminpw'
          }).then((enrollment) => {
            console.log('Successfully enrolled admin user "admin"');
            return fabric_client.createUser({
                username: 'admin',
                mspid: 'Org1MSP',
                cryptoContent: { privateKeyPEM: enrollment.key.toBytes(), signedCertPEM: enrollment.certificate }
            });
          }).then((user) => {
            admin_user = user;
            return fabric_client.setUserContext(admin_user);
          }).catch((err) => {
            console.error('Failed to enroll and persist admin. Error: ' + err.stack ? err.stack : err);
            throw new Error('Failed to enroll admin');
          });
      }
  }).then(() => {
      console.log('Assigned the admin user to the fabric client ::' + admin_user.toString());
  }).catch((err) => {
      console.error('Failed to enroll admin: ' + err);
  });


}
// create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting