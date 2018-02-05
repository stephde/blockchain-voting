let express = require('express');
let path = require('path');
let favicon = require('serve-favicon');
let logger = require('morgan');
let bodyParser = require('body-parser');

let Hyperledger = require('./hyperledger/hyperledger.js');
let hyperledger = new Hyperledger();

let app = express();

// uncomment after placing your favicon in /public
//app.use(favicon(path.join(__dirname, 'public', 'favicon.ico')));
app.use(logger('dev'));
app.use(express.static(path.join(__dirname, 'public')));

app.use( bodyParser.json() );       // to support JSON-encoded bodies
app.use(bodyParser.urlencoded({     // to support URL-encoded bodies
  extended: true
}));

app.get('/:name', function (req, res, next) {
  var options = {
    root: __dirname + '/public/',
    dotfiles: 'deny',
    headers: {
        'x-timestamp': Date.now(),
        'x-sent': true
    }
  };

  var fileName = req.params.name;
  res.sendFile(fileName, options, function (err) {
    if (err) {
      next(err);
    } else {
      console.log('Sent:', fileName);
    }
  });
})

app.post('/voting/initVote', function(req, res, next) {
  hyperledger.initVote();
});

app.post('/voting/setEligible', function(req, res, next) {
  hyperledger.setEligible(req.body.eligibleUsers);
});

app.post('/voting/beginSignUp', function(req, res, next) {
  hyperledger.beginSignUp(req.body.votingQuestion);
});

app.post('/voting/registerUser', function(req, res, next) {
  // todo: check if eligible
  hyperledger.register(req.body.userID);
});

app.get('/voting/question', function(req, res, next) {
  hyperledger.question().then(function (results) {
    res.json({  question: results.toString('utf8') });
  });
});

app.post('/voting/finishRegistrationPhase', function(req, res, next) {
  hyperledger.finishRegistrationPhase();
})

// catch 404 and forward to error handler
app.use(function(req, res, next) {
  let err = new Error('Not Found');
  err.status = 404;
  next(err);
});

// error handler
app.use(function(err, req, res, next) {
  // set locals, only providing error in development
  res.locals.message = err.message;
  res.locals.error = req.app.get('env') === 'development' ? err : {};

  // render the error page
  res.status(err.status || 500);
});

module.exports = app;
