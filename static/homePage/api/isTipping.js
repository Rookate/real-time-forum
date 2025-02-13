import { UserInfo } from "../app.js";
import { ws } from "../direct-message.js";

export let typingTimeout;

export function isTipping() {
    ws.send(JSON.stringify({ type: "typing", isTyping: true, user_uuid: UserInfo.user_uuid }));

    clearTimeout(typingTimeout)
    typingTimeout = setTimeout(() => {
        ws.send(JSON.stringify({ type: "typing", isTyping: false, user_uuid: UserInfo.user_uuid }));
    }, 2000)

}