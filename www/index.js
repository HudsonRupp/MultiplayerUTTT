const socket = new WebSocket("ws://localhost:8080/ws");

socket.addEventListener("open", (event) => {
    console.log(event)
});

socket.addEventListener("message", (event) => {
    console.log(atob(event.data))
})


function sendWs() {
    a = document.getElementById("msg").value
    console.log("sending message " + a)
    socket.send(a)
}