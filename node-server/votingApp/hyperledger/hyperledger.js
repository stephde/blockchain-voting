
Hyperledger = function(){
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
    return this;
  }

  _this.queryAll = function(){
    return query.executeQuery(_this.client, _this.channel, 'queryVotes', [''], defaultUserId);
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
      return invoke.invokeTransaction(_this.client, _this.channel, 'initVote', [], defaultUserId)
  }

  // beginSignUp requires initVote to have been called before
  _this.beginSignUp = function (question) {
      console.log("Starting Sign-Up phase...")
      return invoke.invokeTransaction(_this.client, _this.channel, 'beginSignUp', [question], defaultUserId)
  }

  _this.finishRegistrationPhase = function () {
      console.log("Finishing registration phase, starting Vote phase...")
      return invoke.invokeTransaction(_this.client, _this.channel, 'finishRegistrationPhase', [], defaultUserId)
  }

  _this.setEligible = function (userIds) {
      console.log("Setting eligible voters to: \n" + userIds)
      return invoke.invokeTransaction(_this.hlAdapter.client, _this.channel, 'setEligible', [userIds], defaultUserId)
  }

  _this.registerForVote = function (userId) {
      //ToDo: is the userId implicit?
      console.log("Registering user - " + userId + " - for vote...")
      //ToDo: what is up with the arguments? and what is the 4th argument?
      return invoke.invokeTransaction(_this.client, _this.channel, 'register', [userId], defaultUserId)
  }

  _this.finishRegistrationPhase = function () {
      console.log("Finishing registration phase...")
      return invoke.invokeTransaction(_this.client, _this.channel, 'finishRegistrationPhase', [], defaultUserId)
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
      //ToDo: is this a query or an invocation?
      return invoke.invokeTransaction(_this.client, _this.channel, 'computeTally', [], defaultUserId)
  }

  _this.vote = function(userId, selectedOption) {
    invoke.invokeTransaction(_this.client,
      _this.channel,
      'submitVote', //transaction function
      [userId, selectedOption],
      defaultUserId);
  }

  init();
}

module.exports = Hyperledger;
