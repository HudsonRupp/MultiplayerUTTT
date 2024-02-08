var socket

const params =  new URLSearchParams(window.location.search);
const roomId = params.get("roomId")

socket = new WebSocket("ws://localhost:8000/ws/" + roomId)

socket.addEventListener("open", (event) => {
    console.log(event)
});

socket.addEventListener("message", (event) => {
    data = JSON.parse(event.data)
    console.log(data)
    switch(data.messageType) {
        case "RoomInfo":
            updateRoom(data.content)
    }
})

function sendWs() {
    a = document.getElementById("msg").value
    console.log("sending message " + a)
    socket.send(a)
}

function updateRoom(info) {
    console.log(info)
    roomInfo = JSON.parse(info)

    document.getElementById("name").innerHTML = roomInfo.name
    document.getElementById("capacity").innerHTML = roomInfo.occupants + "/" + roomInfo.capacity
    document.getElementById("status").innerHTML = roomInfo.status
}

function makeMove(i, j) {
    //send ws

    changeSquare(i, j, "X")

    console.log(i, j)
}

function checkSquare() {
    
}

function changeSquare(i, j, html) {
    outI = Math.floor(i / 3)
    outJ = Math.floor(j / 3)
    inI = i % 3
    inJ = j % 3
    var childNodes = document.getElementById("c" + outI.toString() + outJ.toString()).childNodes
    innerId = "i" + inI.toString() + inJ.toString()
    for (var k = 0; k < childNodes.length; k++) {
        var rowNodes = childNodes[k].childNodes
        for (var i = 0; i < rowNodes.length; i++) {
            if (rowNodes[i].id == innerId) {
                rowNodes[i].innerHTML = html
            }
        }
    }
}

function makeBoard() {
    boardDiv = document.getElementById("board")
    board = "";
    for (let i = 0; i < 3; i++) {
      board += "<div class=\"row\">"
        for (let j = 0; j < 3; j++) {
            board += `
                <div class="square" id="c` + i.toString() + j.toString() + `">
                  ` + getInnerBoard(i, j) + `
                </div>
            `
        }
      board += "</div>"
    }
  
    boardDiv.innerHTML = board
}

function getInnerBoard(outI, outJ) {
  board = "";
  for (let i = 0; i < 3; i++) {
      board += "<div class=\"innerRow\">"
        for (j = 0; j < 3; j++) {
            board += `
                <div class="innerSquare" id="i` + i.toString() + j.toString() + `" onclick=makeMove(`+((outI * 3) + i)+`,`+((outJ * 3)+j)+`)>
                    
                </div>
            `
        }
      board += "</div>"
  }
  return board
}
makeBoard()