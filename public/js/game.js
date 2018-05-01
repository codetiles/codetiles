var zoomratio = 1;

$(document).ready(() => {
  // scroll wheel zoomin
  $('#content').on('mousewheel', function(event) {
    event.deltaY === 1 ? zoomin() : zoomout();
  });

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
    submitCode();
  });
});

function submitCode() {
  let submitted_code = $('textarea#code').val();
  // verify that code contains characters
  if (submitted_code !== '' && /\s/.test(submitted_code) !== true) {
    // TODO: CHECK FOR POSSIBLE BUG WITH CODE NOT EXECUTING BECAUSE OF ENCODED CHARS
    submitted_code = encodeURI(submitted_code);
    // TODO: CHECK IF USER IS LOGGED IN BEFORE ALLOWING /GAME ACCESSED OR THIS CODE TO WORK
    let user_id = localStorage.getItem('user_id');
    let data = `{"id" : "` + user_id + `", "code" : "` + submitted_code + `"}`
    $.post("api/v1/uploadcode", data, function(data) {
      let response = JSON.parse(data);
      console.log(response)
    })
  }
}

function proper_size() {
  var document_height = $(document).height();
  var document_width = $(document).width();
  $("#content").kinetic();
}

function zoomin() {
  var document_width = $(document).width();
  if (zoomratio >= 2) {
    //return;
  }
  zoomratio += .1;
  $("#game-board").css("width", document_width * zoomratio)
  $("#game-board").css("height", document_width * zoomratio)
  $("#game-board").css("font-size", String(16 * zoomratio) + "px");
}

function zoomout() {
  var document_width = $(document).width();
  if (zoomratio <= 1) {
    //return;
  }
  zoomratio -= .1;
  $("#game-board").css("width", document_width * zoomratio)
  $("#game-board").css("height", document_width * zoomratio)
  $("#game-board").css("font-size", String(16 * zoomratio) + "px");
}
