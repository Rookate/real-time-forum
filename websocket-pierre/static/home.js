async function home() {
    const body = document.body;
    const buttonToConversation = document.createElement('button');
    buttonToConversation.classList.add("conversationButton")
    buttonToConversation.textContent = 'Go to Conversation';
    buttonToConversation.addEventListener('click', async () => {
        console.log("Conversation clicked")
        try {
            const response = await fetch("http://localhost:8080/api/conversations/createConversation", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body:
                    JSON.stringify({
                        title: "Nouvelle conversation",
                        participants: ["User1", "User2"]
                    })
            });
            if (response.ok) {
                const conversation = await response.json();
                console.log("Conversation created:", conversation.conversation_uuid);
                window.location.href = `/conversation/${conversation.conversation_uuid}`;  // Redirection avec l'UUID
            } else {
                const error = await response.json();
                alert("Erreur lors de la création du post: " + error.message);
            }
        } catch (error) {
            console.error("Erreur lors de la création du post:", error);
        }
    });
    body.appendChild(buttonToConversation);
}


document.addEventListener('DOMContentLoaded', async () => {
    home();
});