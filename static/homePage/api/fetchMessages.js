import { ws } from "../direct-message.js";

const MESSAGE_PER_PAGE = 10;
export function fetchMessages(conversationUUID, getMoreMessages = false, currentIndexMessages) {
    console.log('fetch message :', conversationUUID)
    const type = getMoreMessages ? "getMoreMessages" : "getMessages"
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        console.error("WebSocket not connected");
        return;
    }
    ws.send(JSON.stringify({
        type: type, conversation_uuid: conversationUUID, offset: currentIndexMessages, limit: MESSAGE_PER_PAGE
    }));

    return currentIndexMessages += MESSAGE_PER_PAGE;
}

