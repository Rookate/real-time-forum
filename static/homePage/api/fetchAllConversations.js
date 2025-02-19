export async function fetchAllConversations() {
    try {
        const response = await fetch("/api/conversations/getConversation", {
            method: "GET",
            headers: {
                "Content-Type": "application/json"
            },
        });

        if (response.ok) {
            const data = await response.json()
            console.log("Conversations", data)
            return data
        } else {
            const error = await response.json();
            console.error("Erreur lors de la récupération des users " + error.message);
        }
    } catch (error) {
        console.error("Erreur lors de l'envoi du commentaire:", error);
    }
}