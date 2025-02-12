export async function getConversation(user_uid) {
    try {
        const response = await fetch("/api/conversations/createConversation", {
            headers: {
                "Content-Type": "application/json"
            },
            method: "POST",
            body: JSON.stringify(user_uid)
        });

        if (response.ok) {
            return await response.json();
        } else {
            const error = await response.json();
            alert("Erreur lors de la création de la conversation " + error.message);
        }
    } catch (error) {
        console.error("Erreur lors de la création d'une conversation:", error);
    }
}