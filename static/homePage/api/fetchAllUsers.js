export async function fetchAllUsers() {
    try {
        const response = await fetch("/api/users/fetchAllUsers", {
            method: "GET",
            headers: {
                "Content-Type": "application/json"
            },
        });

        if (response.ok) {
            return await response.json()
        } else {
            const error = await response.json();
            alert("Erreur lors de la récupération des users " + error.message);
        }
    } catch (error) {
        console.error("Erreur lors de l'envoi du commentaire:", error);
        alert("Une erreur s'est produite lors de l'envoi du commentaire. Veuillez réessayer.");
    }
}