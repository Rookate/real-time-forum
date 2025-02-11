let ws;
let typingTimeout;

export function connect() {
    ws = new WebSocket("ws://localhost:8080/ws");

    ws.onopen = function () {
        const conversationUUID = window.location.pathname.split('/').pop();
        fetchMessages(conversationUUID);
        console.log("Connected to WebSocket server");
    };

    ws.onmessage = function (event) {
        const data = JSON.parse(event.data);
        if (data.type === "messages") {
            const messageContainer = document.getElementById("messages");
            //   messageContainer.innerHTML = "";

            data.messages.forEach(msg => {
                const message = document.createElement("div");
                message.classList.add("messageCard");
                console.log("Prout", data.messages);
                message.innerText = msg.message;
                messageContainer.appendChild(message);
            });

        } else if (data.type === "message") {
            const messageContainer = document.getElementById("messages")
            const message = document.createElement("div");
            message.classList.add("messageCard");
            message.innerText = data.content;
            messageContainer.appendChild(message);

        } else if (data.type === "typing") {
            const typingIndicator = document.getElementById("isTapping");
            if (data.isTyping) {
                typingIndicator.innerText = `${data.username} est en train d'écrire...`;
            } else {
                typingIndicator.innerText = "";
            }
        }
    };

    ws.onclose = function () {
        console.log("WebSocket connection closed, retrying...");
        setTimeout(connect, 1000); // Reconnect after 1 second
    };

    ws.onerror = function (error) {
        console.error("WebSocket error:", error);
    };
}

export async function sendMessage() {
    let input = document.getElementById("messageInput");
    let message = input.value;

    const conversationUUID = window.location.pathname.split('/').pop();

    console.log(conversationUUID)

    const data = {
        content: message,
        conversation_uuid: conversationUUID
    }

    if (message !== "") {
        ws.send(JSON.stringify({ type: "message", content: message }));

        try {
            const response = await fetch("http://localhost:8080/api/message/createMessage", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data)
            });

            if (response.ok) {
                input.value = "";
            } else {
                const error = await response.json();
                alert("Erreur lors de la création du message: " + error.message);
            }
        } catch (error) {
            console.error("Erreur lors de la création du message:", error);
        }

        clearInterval(typingTimeout);
        ws.send(JSON.stringify({ type: "typing", isTyping: false, username: "User" }));
        input.value = "";
    }
}

function fetchMessages(conversationUUID) {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        console.error("WebSocket not connected");
        return;
    }

    ws.send(JSON.stringify({ type: "getMessages", conversation_uuid: conversationUUID }));
}

export function isTapping() {
    ws.send(JSON.stringify({ type: "typing", isTyping: true, username: "User" }));

    clearTimeout(typingTimeout);
    typingTimeout = setTimeout(() => {
        ws.send(JSON.stringify({ type: "typing", isTyping: false, username: "User" }));
    }, 500);
}

document.getElementById("messageInput").addEventListener("input", isTapping);

connect();

