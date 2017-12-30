
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
  }

  _this.queryAll = function(){
    return query.executeQuery(_this.hlAdapter.client, _this.hlAdapter.channel, 'queryVotes', [''], defaultUserId);
  }

  _this.registerUser = function(user){
    registration.registerUser(_this.client, user);
  }

  _this.enrollAdmin = function(){
    enroll.enrollAdmin(_this.client);
  }

  _this.vote = function(selectedOption){
    invoke.invokeTransaction(_this.hlAdapter.client,
      _this.channel,
      'vote', //transaction function
      [selectedOption],
      defaultUserId);
  }

  init();
}

module.exports = Hyperledger;
