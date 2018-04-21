$(document).ready(() => {
  function proper_size() {
    var document_height = $(document).height();
    var document_width = $(document).width();
    /*
      $("#content").css("height", document_height);
      $("#content").css("width", document_width);
      $("#board").width(document_width);
      $("#board").height(document_height);
      $("#game-board").width($("#board").width()-2);
      $("#game-board").height($("#board").height()-20);
    */
    $("#content").kinetic();
  }
  $(document).resize(() => {
    proper_size();
  });
  const GAME_SIZE = 30;
  console.log("Ready!");
  var game = $("<table id='game-board'></table>");
  for (var i = 0; i < GAME_SIZE; i++) {
    var new_row = $("<tr></tr>");
    for (var j = 0; j < GAME_SIZE; j++) {
      var new_column = $("<td></td>");
      //new_column.css("border-radius", "15%");
      new_row.append(new_column);
    }
    game.append(new_row);
  }
  $("#board").append(game);
  proper_size();
});
