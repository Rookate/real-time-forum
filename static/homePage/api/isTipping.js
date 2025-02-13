import { UserInfo } from "../app.js";
import { ws } from "../direct-message.js";

export let typingTimeout;

export function isTipping() {
    console.log("IsTipping ==========>")
    ws.send(JSON.stringify({ type: "typing", isTyping: true, user_uuid: UserInfo.user_uuid }));

    clearTimeout(typingTimeout);
    typingTimeout = setTimeout(() => {
        document.getElementById("typing-span").style.visibility = "hidden"
    }, 2000)

}