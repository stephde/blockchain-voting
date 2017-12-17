'use strict';
/*
* Copyright IBM Corp All Rights Reserved
*
* SPDX-License-Identifier: Apache-2.0
*/
/*
 * Register and Enroll a user
 */

let Fabric_CA_Client = require('fabric-ca-client');
let HyperledgerUtils = require("./hyperledergerUtils");
let fabric_ca_client = null;
let admin_user = null;
let member_user = null;

exports.registerUser = function(fabric_client, user){
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
        } else {
            throw new Error('Failed to get admin.... run enrollAdmin.js');
        }

        // at this point we should have the admin user
        // first need to register the user with the CA server
        return fabric_ca_client.register({enrollmentID: user, affiliation: 'org1.department1', role: 'client'}, admin_user);
    }).then((secret) => {
        // next we need to enroll the user with CA server
        console.log('Successfully registered user1 - secret:'+ secret);

        return fabric_ca_client.enroll({enrollmentID: user, enrollmentSecret: secret});
    }).then((enrollment) => {
      console.log('Successfully enrolled member user "user1" ');
      return fabric_client.createUser(
         {username: 'user1',
         mspid: 'Org1MSP',
         cryptoContent: { privateKeyPEM: enrollment.key.toBytes(), signedCertPEM: enrollment.certificate }
         });
    }).then((user) => {
         member_user = user;

         return fabric_client.setUserContext(member_user);
    }).then(()=>{
         console.log('User1 was successfully registered and enrolled and is ready to intreact with the fabric network');

    }).catch((err) => {
        console.error('Failed to register: ' + err);
    	if(err.toString().indexOf('Authorization') > -1) {
    		console.error('Authorization failures may be caused by having admin credentials from a previous CA instance.\n' +
    		'Try again after deleting the contents of the store directory '+store_path);
    	}
    });
  }

// create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
