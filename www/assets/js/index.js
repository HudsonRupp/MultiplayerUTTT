var socket

function sendWs() {
    a = document.getElementById("msg").value
    console.log("sending message " + a)
    socket.send(a)
}

function connectNewRoom(roomId) {
    if (!roomId) {
        roomId = "00000000-0000-0000-0000-000000000000"
    }
    window.location.href = "/game?roomId=" + roomId
}

function getRooms() {
    $.ajax({
        url : 'http://localhost:8000/rooms',
        type: 'GET',
        success : function(data) {
            console.log(data)
            setHtml(data)
            if(!socket) {
                setTimeout(getRooms, 500);
            }
        },
        cors: true,
        crossDomain: true,
        "headers": {
            "accept": "application/json",
        }
    });
}

function setHtml(data) {
    var string = ""
    for (var room in data) {
        string +=`
      <div class="rowElement" onclick="connectNewRoom('`+ data[room].id + `')">
        <h2>`+data[room].name+`</h2>
        <ul>
          <li>`+data[room].id+`</li>
          <li>`+data[room].status+`</li>
          <li>`+`:)`+`</li>
        </ul>
        <p class="capacity">`+data[room].occupants + `\\` + data[room].capacity +`</p>
      </div>
      `
    }
    string +=`<div class="newRoom" onclick="connectNewRoom()">
                <h1>+</h1>
            </div>`
    document.getElementsByClassName("rooms")[0].innerHTML = string;
}
//console.log(JSON.parse(demoData))
getRooms()
//setHtml(JSON.parse(demoData))