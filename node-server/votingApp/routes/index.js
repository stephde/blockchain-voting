let express = require('express');
let router = express.Router();

let Hyperledger = require('../hyperledger/hyperledger.js')

/* GET home page. */
router.get('/', function(req, res, next) {
  res.render('index', { title: 'Express' });
});

router.post('/voting/place', function(req, res, next) {
  let hyperledger = new Hyperledger();
  hyperledger.vote(req.body.vote);
  console.log(req.body.vote);

  res.render('index', { title: 'Express' });
});

router.get('/voting/all', function(req, res, next) {
  let hyperledger = new Hyperledger();

  hyperledger.queryAll().then(function(results){
    console.log(results.toString('utf8'));
    debugger
    res.setHeader('Content-Type', 'application/json');
    res.send(results.toString('utf8'));
  });

});

module.exports = router;
