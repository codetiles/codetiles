var zoomratio = 1;

$(document).ready(() => {


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
      if (Math.random() > .9) {
        new_column.text(function() {
          return Math.floor(Math.random() * 100 + 1)
        });
        var x = Math.random()
        if (x > .66) {
          new_column.addClass("blue");
        }
        else if (x < .33) {
          new_column.addClass("red")
        }
        else {
          new_column.addClass("yellow")
        }

      }

    }
    game.append(new_row);
  }
  $("#board").append(game);

  proper_size();

  $("#publish").on("click", function () {
    console.log("Hello there, you just clicked the submit button!");
  });
});

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
  // $("#code").css("position", "absolute");
  // $("#code").css("height", $(document).height()*2/5);
  // $("#code").css("width", $(document).width()*2/5);
  // $("#code").css("right", "0px");
  // $("#code").css("bottom", "0px");
  // $("#code").css("background-color", "blue"); // For the sake of identifying it
  // $("#editor").css("width", $("#code").width()*99/100);
  // $("#editor").css("height", $("#code").height());
  // $("#editor").css("resize", "none");
  // $("#editor").css("display", "block");
  // $("#editor").css("margin-left", "auto");
  // $("#editor").css("margin-right", "auto");
  // $("#editor").css("background-color", "#1A5569");
  // $("#editor").css("color", "white");
  $("#content").kinetic();
}

function zoomin() {
  var document_width = $(document).width();
  if (zoomratio >= 2) {
    return
  }

  zoomratio += .2;
  $("#game-board").css("width", document_width * zoomratio)
  $("#game-board").css("height", document_width * zoomratio)
  $("#game-board").css("font-size", String(16 * zoomratio) + "px"); 
}

function zoomout() {
  var document_width = $(document).width();
  if (zoomratio <= 1) {
    return;
  }

  zoomratio -= .2;
  $("#game-board").css("width", document_width * zoomratio)
  $("#game-board").css("height", document_width * zoomratio)
  $("#game-board").css("font-size", String(16 * zoomratio) + "px"); 
}
