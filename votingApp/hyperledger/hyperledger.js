
Hyperledger = function() {
    let initClient = require("./initFabricClient.js"),
        query = require("./query.js"),
        registration = require("./registerUser.js"),
        enroll = require("./enrollAdmin.js"),
        invoke = require("./invoke.js"),
        hlAdapter,
        channel,
        client,
        _this = this;

  const host = 'grpc://localhost';
  const channelId = 'mychannel';
  const defaultUserId = 'user1';

  function init (){
    initClient.initFabricClient(host, channelId, defaultUserId, _this.onTransactionCommitted)
        .then((hlAdapter) => {
            _this.hlAdapter = hlAdapter;
            _this.channel = _this.hlAdapter.channel;
            _this.client = _this.hlAdapter.client;
            _this.eventHub = _this.hlAdapter.eventHub;
            _this.commitedTransactions = []
            _this.pendingTransactions = []
        });
    return this;
  }

  _this.onTransactionCommitted = function (transaction) {
      let header = transaction.payload.header; // the "header" object contains metadata of the transaction
      const tx_id = header.channel_header.tx_id;
      _this.commitedTransactions.push(tx_id)

      // remove tx_id from pending
      _this.pendingTransactions = _this.pendingTransactions.filter(val => val !== tx_id)
  }

  _this.onTransactionSubmitted = function (tx_id) {
      _this.pendingTransactions.push(tx_id)
  }

  _this.isTransactionPending = function () {
      return this.pendingTransactions.length > 0
  }

  _this.queryAll = function(){
    return query.executeQuery(_this.client, _this.channel, 'computeTally', [''], defaultUserId);
  }

  _this.registerUser = function(user) {
    return registration.registerUser(_this.client, user.id, defaultUserId);
  }

  _this.register = function (userID) {
      console.log("Registering " + userID + " for the vote")
      return invoke.invokeTransaction(_this.client, _this.channel, 'register', [userID], defaultUserId)
  }

  _this.getUser = function (userId) {
    //ToDo: actually get user
    return new Promise((resolve, reject) => resolve({id: userId}));
  }

  _this.enrollAdmin = function() {
    return enroll.enrollAdmin(_this.client, defaultUserId);
  }

  _this.initVote = function () {
      console.log("Initializing the vote...")
      return invoke.invokeTransaction(_this.client, _this.channel,
          _this.onTransactionSubmitted, 'initVote', [], defaultUserId)
  }

  // beginSignUp requires initVote to have been called before
  _this.beginSignUp = function (question) {
      console.log("Starting Sign-Up phase...")
      return invoke.invokeTransaction(_this.client, _this.channel, _this.onTransactionSubmitted, 'beginSignUp', [question], defaultUserId)
  }

  _this.finishRegistrationPhase = function () {
      console.log("Finishing registration phase, starting Vote phase...")
      return invoke.invokeTransaction(_this.client, _this.channel,
          _this.onTransactionSubmitted, 'finishRegistrationPhase', [], defaultUserId)
  }

  _this.setEligible = function (userIds) {
      console.log("Setting eligible voters to: \n" + userIds)
      return invoke.invokeTransaction(_this.client, _this.channel, _this.onTransactionSubmitted, 'setEligible',
          userIds, defaultUserId)
  }

  _this.registerForVote = function (userId) {
      console.log("Registering user - " + userId + " - for vote...")
      return invoke.invokeTransaction(_this.client, _this.channel, _this.onTransactionSubmitted, 'register',
          [userId], defaultUserId)
  }

  _this.question = function(){
      console.log("Getting the question ...")
      return query.executeQuery(_this.hlAdapter.client, _this.hlAdapter.channel, 'question', [''], defaultUserId);
  }

  _this.submitVote = function (userID, vote) {
      console.log("Submitting a vote for " + userID)
      return invoke.invokeTransaction(_this.client, _this.channel, 'submitVote', [userID, vote], defaultUserId)
  }

  _this.computeTally = function () {
      console.log("Computing the tally...")
      return query.executeQuery(_this.client, _this.channel, 'computeTally', [], defaultUserId);
  }

  _this.vote = function(userId, selectedOption) {
    return invoke.invokeTransaction(_this.client,
      _this.channel,
      _this.onTransactionSubmitted,
      'submitVote', //transaction function
      [userId, selectedOption],
      defaultUserId);
  }

  _this.close = function () {
      console.log("\nCommited " + _this.commitedTransactions.length +  " transactions: \n")
      console.log(JSON.stringify(_this.commitedTransactions))
      console.log("\nPending transactions: \n")
      console.log(JSON.stringify(_this.pendingTransactions))

      _this.eventHub.unregisterBlockEvent(_this.hlAdapter.blockListenerId)
      _this.eventHub.disconnect()
  }

  init()
}

module.exports = Hyperledger;
