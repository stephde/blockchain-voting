function initVote() {
  $.ajax({
          type: "POST",
          url: '/voting/initVote',
          data: {}
  });

  document.getElementById('initVoteBtn').disabled = true;
  document.getElementById('eligible').removeAttribute('hidden');
}

function setEligible() {
  var eligibleUsers = document.getElementById('addresses').value.split(',')

  $.ajax({
          type: "POST",
          url: '/voting/setEligible',
          data: {
            'eligibleUsers': eligibleUsers
          }
  });

  document.getElementById('registrationSetQuestion').removeAttribute('hidden');
}

function beginRegistration() {
  var votingQuestion = document.getElementById('questioninput').value;

  $.ajax({
          type: "POST",
          url: '/voting/beginSignUp',
          data: {
            votingQuestion: votingQuestion
          }
  });

  document.getElementById('setupfs').hidden = true;
  document.getElementById('registerfs').removeAttribute('hidden');
  $('#progressbar li').removeClass('active');
  $('#progressbar li:eq(1)').addClass('active');
}

function finishRegistration() {
  $.ajax({
          type: "POST",
          url: '/voting/finishRegistrationPhase',
          data: {}
  });

  document.getElementById('registerfs').hidden = true;
  document.getElementById('castfs').removeAttribute('hidden');
  $('#progressbar li').removeClass('active');
  $('#progressbar li:eq(2)').addClass('active');

}

function tally() {
  $.ajax({
          type: "POST",
          url: '/voting/tallyVote',
          data: {}
  });

  document.getElementById('castfs').hidden = true;
  document.getElementById('tallyfs').removeAttribute('hidden');
  $('#progressbar li').removeClass('active');
  $('#progressbar li:eq(3)').addClass('active');
  initResults();
  var result = [
    { 'Record': {
        'count': 20,
        'name': 'Yes'
    }},
    { 'Record': {
        'count': 5,
        'name': 'No'
    }}
  ];
  drawGraph(result);
}

function initResults(){
  $.ajax({
    type: "GET",
    url: "/voting/getTally"
  }).done(function(result){
    drawGraph(result);
  });
}

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
  var labels = [];
  $(result).each(function(_, option){
    datasets[0].data.push(option.Record.count);
    labels.push(option.Record.name);
  });
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
