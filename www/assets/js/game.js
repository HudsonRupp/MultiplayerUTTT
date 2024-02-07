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