var socket

function sendWs() {
    a = document.getElementById("msg").value
    console.log("sending message " + a)
    socket.send(a)
}

function connectNewRoom() {
    socket = new WebSocket("ws://localhost:8080/ws")

    socket.addEventListener("open", (event) => {
        console.log(event)
    });
    
    socket.addEventListener("message", (event) => {
        console.log(atob(event.data))
    })
}

function getRooms() {
    $.ajax({
        url : 'http://localhost:8080/rooms',
        type: 'GET',
        success : function(data) {
            console.log(data)
            if (data.length != 0) {
                setHtml(data)
            }
            if(!socket) {
                setTimeout(getRooms, 3000);
            }
        },
        cors: true,
        crossDomain: true,
        "headers": {
            "accept": "application/json",
        }
    });
}

demoData = `[
    {
        "name": "Room1",
        "id": "uuid-uuid-uuid-uuid",
        "occupants": 1,
        "capacity": 2,
        "status": "Waiting to start"
    },
    {
        "name": "Room2",
        "id": "uuid-uuid-uuid-uuid",
        "occupants": 2,
        "capacity": 2,
        "status": "In Progress"
    },
    {
        "name": "Room3",
        "id": "uuid-uuid-uuid-uuid",
        "occupants": 1,
        "capacity": 2,
        "status": "Finished"
    }
]`

function setHtml(data) {
    var string = ""
    for (var room in data) {
        console.log(data[room])
        string +=`
      <div class="rowElement">
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