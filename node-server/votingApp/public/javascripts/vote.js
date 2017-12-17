$( document ).ready(function() {
    $('#send-vote').click(function(target){
      $('.form-check input').each(function(_, option){
        if($(option).is(':checked')){
            $.ajax({
              type: "POST",
              url: '/voting/place',
              data: {
                vote: $(option).val()
              }
            });
          };
          $('#placeVote').modal('hide');
      });
    });
    $('#showResultsButton').click(function(target){
      $.ajax({
        type: "GET",
        url: "voting/all"
      }).done(function(result){
        var datasets = [
          {
            data: [],
            backgroundColor: [
               'lightblue',
               'lightgreen',
               'lightcoral',
               'yellow'
            ],


          }
        ];
        var labels = [];
        $(result).each(function(_, option){
          datasets[0].data.push(option.Record.count);
          //labels.push(option.Record.name);
        });
        var ctx = document.getElementById("results");
        var formatedData = {
            datasets: datasets,
            // These labels appear in the legend and in the tooltips when hovering different arcs

            labels: [ "Michael Mayor", "Benjamin Bürgermeister", "Major Tom", "Petra Principal" ]
            //labels: labels,


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

      });

    });


});
