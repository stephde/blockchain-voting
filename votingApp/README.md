# Liquid Voting Client

## Installation

Please Run: `npm install` to install all needed Dependencies

Install hyperledger and start the network. (Follow Installation instructions under /hyperledger)

## Running the App

To start the webserver run:
`sudo npm start`
Sudo ist required since the app uses port 80.

Open your http://localhost with your Browser to access the app. 

## Hyperledger Client

You can find the hyperledger js adapter under /hyperledger.
The entry point of the adapter is the hyperledger.js file, which defines the public interface.
It delegates all transaction to invoke.js and all queries to query.js.
In order to submit transactions an admin and a Uuer have to be enrolled.

When the app is started the adapter tries to connect to hyperledger on the route specified in hyperledger.js (per default localhost).
It connects to a peer (per default on port 7051) and an orderer (per default on port 7051).
Moreover ot creates an event which subscribes to all block events on the chain (per default on port 7053).
Every time a block is commited the `onTransactionCommitted(tx_id)` function is called for every transaction in the block.

Everytime a transaction is submitted its id is stored in the hyperledger client as pending transaction.
Once the callback function is executed with the same transactionId, we know that the transaction has been committed successfully.

Additionally to the hyperledger client there is a main.js file in the hyperledger directory.
By runnning this file, as complete election (from user registration until tally computation) can be simulated.
The number of participants can be specified in the file.