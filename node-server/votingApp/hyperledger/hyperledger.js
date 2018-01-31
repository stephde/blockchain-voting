
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

  const host = 'grpc://localhost:7051';
  const channelId = 'mychannel';
  const defaultUserId = 'user1';

  function init(){
    _this.hlAdapter = initClient.initFabricClient(host, channelId);
    _this.channel = _this.hlAdapter.channel;
    _this.client = _this.hlAdapter.client;
    _this.eventHub = _this.hlAdapter.eventHub;
    return this;
  }

  _this.queryAll = function(){
    return query.executeQuery(_this.client, _this.channel, 'queryVotes', [''], defaultUserId);
  }

  _this.registerUser = function(user) {
    return registration.registerUser(_this.client, user.id, defaultUserId);
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
          _this.eventHub, 'initVote', [], defaultUserId)
  }

  // beginSignUp requires initVote to have been called before
  _this.beginSignUp = function (question) {
      console.log("Starting Sign-Up phase...")
      return invoke.invokeTransaction(_this.client, _this.channel, _this.eventHub, 'beginSignUp', [question], defaultUserId)
  }

  _this.finishRegistrationPhase = function () {
      console.log("Finishing registration phase, starting Vote phase...")
      return invoke.invokeTransaction(_this.client, _this.channel,
          _this.eventHub, 'finishRegistrationPhase', [], defaultUserId)
  }

  _this.setEligible = function (userIds) {
      console.log("Setting eligible voters to: \n" + userIds)
      return invoke.invokeTransaction(_this.client, _this.channel,
          _this.eventHub, 'setEligible', userIds, defaultUserId)
  }

  _this.registerForVote = function (userId) {
      //ToDo: is the userId implicit?
      console.log("Registering user - " + userId + " - for vote...")
      return invoke.invokeTransaction(_this.client, _this.channel,
          _this.eventHub, 'register', [userId], defaultUserId)
  }

  _this.computeTally = function () {
      console.log("Computing the tally...")
      return query.executeQuery(_this.client, _this.channel, 'computeTally', [], defaultUserId);
  }

  _this.vote = function(userId, selectedOption) {
    return invoke.invokeTransaction(_this.client,
      _this.channel,
      _this.eventHub,
      'submitVote', //transaction function
      [userId, selectedOption],
      defaultUserId);
  }

  _this.close = function () {
      _this.eventHub.disconnect();
  }

  init();
}

module.exports = Hyperledger;
