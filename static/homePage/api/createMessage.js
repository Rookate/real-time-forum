import { UserInfo } from "../app.js";
import { ws } from "../direct-message.js";
import { typingTimeout } from "./isTipping.js";

export async function sendMessage(conversationUUID) {
    let input = document.getElementById("messageInput");
    let message = input.value;

    const data = {
        content: message,
        conversation_uuid: conversationUUID,
        sender_uuid: UserInfo.user_uuid,
        sender_username: UserInfo.username,
        sender_profile_picture: UserInfo.profile_picture,
    }

    if (message !== "") {
        ws.send(JSON.stringify({ type: "single_message", content: data }));

        try {
            const response = await fetch("http://localhost:8080/api/message/createMessage", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data)
            });

            if (response.ok) {
                const responseData = await response.json();
                // Récupérer le message_uuid
                const messageUUID = responseData.message_uuid;
                input.value = "";
                return messageUUID; // Retourne l'UUID du message si besoin
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