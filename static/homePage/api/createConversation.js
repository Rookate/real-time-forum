export async function createConversation(user_uid) {
    console.log("<Data-User-UUID", user_uid)
    try {
        const response = await fetch("http://localhost:8080/api/conversations/createConversation", {
            headers: { "Content-Type": "application/json" },
            method: "POST",
            body: JSON.stringify({ user_uuid: user_uid })
        });

        if (response.ok) {
            const data = await response.json();
            console.log("<Data>", data)
            return data // Tenter de parser si c'est bien un JSON
        } else {
            console.error("Erreur lors de la création de la conversation:");
            alert("Erreur lors de la création de la conversation: ");
        }
    } catch (error) {
        console.error("Erreur lors de la création d'une conversation:", error);
    }
}
