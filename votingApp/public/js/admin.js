var overlay = $('#overlay');
function initVote() {
  overlay.show();
  $.ajax({
          type: "POST",
          url: '/voting/initVote',
          data: {}
  }).done(function(data, status, jqXHR){
    overlay.hide();
  }).fail(function(jqXHR, textStatus, e){
    showError(jqXHR.responseText);
  });;

  document.getElementById('initVoteBtn').disabled = true;
  document.getElementById('eligible').removeAttribute('hidden');
}

function setEligible() {
  var eligibleUsers = document.getElementById('addresses').value.split(',')
  overlay.show();
  $.ajax({
          type: "POST",
          url: '/voting/setEligible',
          data: {
            'eligibleUsers': eligibleUsers
          }
  }).done(function(){
    overlay.hide();
  }).fail(function(jqXHR, textStatus, e){
    showError(jqXHR.responseText);
  });;

  document.getElementById('registrationSetQuestion').removeAttribute('hidden');
}

function beginRegistration() {
  var votingQuestion = document.getElementById('questioninput').value;
  overlay.show();

  $.ajax({
          type: "POST",
          url: '/voting/beginSignUp',
          data: {
            votingQuestion: votingQuestion
          }
  }).done(function(){
    overlay.hide();
  }).fail(function(jqXHR, textStatus, e){
    showError(jqXHR.responseText);
  });;

  document.getElementById('setupfs').hidden = true;
  document.getElementById('registerfs').removeAttribute('hidden');
  $('#progressbar li').removeClass('active');
  $('#progressbar li:eq(1)').addClass('active');
}

function finishRegistration() {
  overlay.show();
  $.ajax({
          type: "POST",
          url: '/voting/finishRegistrationPhase',
          data: {}
  }).done(function(){
    overlay.hide();
  }).fail(function(jqXHR, textStatus, e){
    showError(jqXHR.responseText);
  });;

  document.getElementById('registerfs').hidden = true;
  document.getElementById('castfs').removeAttribute('hidden');
  $('#progressbar li').removeClass('active');
  $('#progressbar li:eq(2)').addClass('active');

}

function tally() {

  document.getElementById('castfs').hidden = true;
  document.getElementById('tallyfs').removeAttribute('hidden');
  $('#progressbar li').removeClass('active');
  $('#progressbar li:eq(3)').addClass('active');
  initResults();

}

function initResults(){
  overlay.show();
  $.ajax({
    type: "GET",
    url: "/voting/getTally"
  }).done(function(data, status, jqXHR){
    drawGraph(data);
    overlay.hide();
  }).fail(function(jqXHR, textStatus, e){
    showError(jqXHR.responseText);
  });
};

function showError(msg){
  $.notify(msg);
};

function drawGraph(result){
  var datasets = [
    {
      data: [],
      backgroundColor: [
         'lightblue',
         'lightgreen',
         'lightcoral',
         'lightyellow'
      ],
    }
  ];
  var labels = [ 'NO', 'YES' ];

  datasets[0].data = Object.values(result.Votes);
  var ctx = document.getElementById("results");
  var formatedData = {
      datasets: datasets,
      // These labels appear in the legend and in the tooltips when hovering different arcs
      labels: labels,
  };
  var options = {

  };
  if(typeof(myPieChart) !== 'undefined'){
    myPieChart.data = formatedData;
    myPieChart.update();
  } else{

    myPieChart = new Chart(ctx,{
      type: 'pie',
      data: formatedData,
      options: options
    });
  }
}
$('#showResultsButton').click(function(target){

});
