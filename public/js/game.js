var zoomratio = 1;

$(document).ready(() => {
  // SETUP
  // scroll wheel zoomin
  $('#content').on('mousewheel', function(event) {
    event.deltaY === 1 ? zoomin() : zoomout();
  });
  // set game size properly on window resizing
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
    }
    game.append(new_row);
    $(document).delegate('#editor', 'keydown', function(e) {
    var keyCode = e.keyCode || e.which;

    if (keyCode == 9) {
      e.preventDefault();
      var start = this.selectionStart;
      var end = this.selectionEnd;

      // set textarea value to: text before caret + tab + text after caret
      $(this).val($(this).val().substring(0, start)
        + "\t"
        + $(this).val().substring(end));

      // put caret at right position again
      this.selectionStart =
      this.selectionEnd = start + 1;
    }
  });
  }
  $("#board").append(game);

  proper_size();

  $("#publish").on("click", function () {
    submitCode();
  });

  // Websocket stuffs
  const socket = new WebSocket('ws://localhost:8080/api/v1/ws/gameboard');

  socket.addEventListener('message', function (event) {
    console.log(event.data);
    if (event.data === "User is not in a game" || event.data === "User does not exist") {
      window.location.href = '/';
    } else {
      determineLoad(event.data);
    }
  });

  socket.addEventListener('open', function (event) {
    socket.send(localStorage.getItem('user_id'));
  });
});

function loadBoard(boardString) {
  let lines = boardString.split(`\n`);
  let tiles = [];
  for (let i=0;i<lines.length;i++) {
    for (let j=0;j<lines[i].length;j+=3) {
      let number = lines[i].substring(j+1, j+3);
      let color = lines[i][j];
      let el = $('#game-board tr:eq('+String(i)+') td:eq('+String(j/3)+')')
      if (number != "00") {
        el.text(number);
      }
      else {
        el.html("&nbsp")
      }
      switch (color) {
        case "b":
          el.removeClass();
          el.addClass('blue');
          break;
        case "r":
          el.removeClass();
          el.addClass('red');
          break;
        default:
          el.removeClass();
      }
    }
  }
}

function loadDifference(boardString) {
  for (let i = 0; i < boardString.length; i += 7) {
    let x = boardString.substring(i, i + 2);
    let y = boardString.substring(i + 2, i + 4);
    let c = boardString[4]
    let v = boardString.substring(i + 5, i + 7);
    let el = $("#game-board tr:eq(" + x + ") td:eq(" + y + ")")
    el.text(v)
    switch (c) {
      case "b":
        el.removeClass();
        el.addClass('blue');
        break;
      case "r":
        el.removeClass();
        el.addClass('red');
        break;
      default:
        el.removeClass();
    }
  }
}

function determineLoad(str) {
  if (str[0] == "d") {
    loadDifference(str.substring(1, str.length));
    return;
  }
  loadBoard(str);
}

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
    return;
  }
  zoomratio += .1;
  $("#game-board").css("width", document_width * zoomratio)
  $("#game-board").css("height", document_width * zoomratio)
  $("#game-board").css("font-size", String(16 * zoomratio) + "px");
}

function zoomout() {
  var document_width = $(document).width();
  if (zoomratio <= 1) {
    return;
  }
  zoomratio -= .1;
  $("#game-board").css("width", document_width * zoomratio)
  $("#game-board").css("height", document_width * zoomratio)
  $("#game-board").css("font-size", String(16 * zoomratio) + "px");
}
