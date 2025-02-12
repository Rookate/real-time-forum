import { ws } from "../direct-message.js";

export function fetchMessages(conversationUUID) {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        console.error("WebSocket not connected");
        return;
    }

    ws.send(JSON.stringify({ type: "getMessages", conversation_uuid: conversationUUID }));
}
