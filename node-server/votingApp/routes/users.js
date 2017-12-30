var express = require('express');
var router = express.Router();

/* GET users listing. */
router.get('/', function(req, res, next) {
  res.send('respond with a resource');
});

// get user by id
router.get('/:id', function (req, res, next) {
  let userId = req.params.id;


  console.log("Fetching user with id - " + userId + " from hyperledger...")
  let hyperledger = new Hyperledger();
  hyperledger.getUser(userId).then((user) => {
    console.log(user);
    res.json(user).send();
  })
});

// create user
router.post('/', function (req, res, next) {
  let user = req.body;

  console.log(user);
  let hyperledger = new Hyperledger();
  hyperledger.registerUser(user).then((result) => {
    res.json(result).send();
  })
})

router.post('/admin', function (req, res, next) {
    let hyperledger = new Hyperledger();
    hyperledger.enrollAdmin().then(() => {
        res.send(200);
    })
})

module.exports = router;
