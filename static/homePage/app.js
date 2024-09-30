document.addEventListener('DOMContentLoaded', () => {
    fetchPosts();
});

// Sélectionnez les éléments
const toggleButton = document.getElementById('toggle-menu-btn');
const sidebar = document.getElementById('sidebar');
const darkModeToggle = document.getElementById('dark-mode-toggle');

// Fonction pour appliquer le mode
function applyMode(mode) {
    const root = document.documentElement;

    if (mode === 'dark') {
        root.style.setProperty('--background-color', '#1C1C1C');
        root.style.setProperty('--text-color', '#000000');
        root.style.setProperty('--second-text-color', '#FFFFFF');
        root.style.setProperty('--border-color', '#5E5E5F');
        root.style.setProperty('--background-message-color', '#272727');
        darkModeToggle.textContent = 'Light Mode';
    } else {
        root.style.setProperty('--background-color', '#f5f5f5');
        root.style.setProperty('--text-color', '#FFFFFF');
        root.style.setProperty('--second-text-color', '#000000');
        root.style.setProperty('--border-color', '#9C9FA8');
        root.style.setProperty('--background-message-color', '#FFFFFF');
        darkModeToggle.textContent = 'Dark Mode';
    }

    // Enregistrer la préférence dans le Local Storage
    localStorage.setItem('theme', mode);
}

// Vérifie la préférence au chargement
const userPreference = localStorage.getItem('theme');
if (userPreference) {
    applyMode(userPreference);
} else {
    // Si aucune préférence n'est trouvée, définir le mode par défaut (par exemple, light)
    applyMode('light');
}

// Écouteur d'événement pour le bouton de basculement
darkModeToggle.addEventListener('click', () => {
    const currentMode = localStorage.getItem('theme') || 'light';
    const newMode = currentMode === 'dark' ? 'light' : 'dark';
    applyMode(newMode);
});

// Événement pour le menu
toggleButton.addEventListener('click', () => {
    sidebar.classList.toggle('close');
});


async function fetchPosts() {
    const messagesList = document.getElementById('users-post');
    messagesList.innerHTML = '<p>Loading...</p>';
    try {
        const response = await fetch("http://localhost:8080/api/post/fetchAllPost");
        if (!response.ok) {
            throw new Error("Error retrieving posts");
        }

        const posts = await response.json();
        messagesList.innerHTML = '';

        if (posts.length === 0) {
            messagesList.innerHTML = '<p>No posts available.</p>';
        } else {
            posts.forEach(post => {
                DisplayMessages(post);
            });
        }
    } catch (error) {
        messagesList.innerHTML = '<p>Error loading posts. Please try again.</p>';
        console.error(error);
    }
}

function DisplayMessages(post) {
    const displayTimeStamp = post.created_at ? new Date(post.created_at).toLocaleString() : new Date().toLocaleString();

    const messagesList = document.getElementById('users-post');

    const messageItem = document.createElement('div');
    messageItem.classList.add('message-item');
    messageItem.setAttribute('post_uuid', post.post_uuid);

    const messageHeader = document.createElement('div');
    messageHeader.classList.add('message-header');

    const userInfo = document.createElement('div');
    userInfo.classList.add('user-info');

    const profilePicture = document.createElement('img');
    profilePicture.src = post.profile_picture || 'default-profile-picture.jpg';
    profilePicture.alt = 'Profile Picture';
    profilePicture.classList.add('profile-picture');

    const userNameSpan = document.createElement('span');
    userNameSpan.classList.add('username');
    userNameSpan.textContent = post.username;

    userInfo.appendChild(profilePicture);
    userInfo.appendChild(userNameSpan);

    const timeStampSpan = document.createElement('span');
    timeStampSpan.classList.add('timestamp');
    timeStampSpan.textContent = displayTimeStamp;

    const deleteButton = document.createElement('button');
    deleteButton.classList.add('delete-button');
    deleteButton.textContent = 'Delete';
    deleteButton.addEventListener('click', () => {
        deletePost(post.post_uuid);
    });

    messageHeader.appendChild(userInfo);
    messageHeader.appendChild(timeStampSpan);
    messageHeader.appendChild(deleteButton);

    const messageContent = document.createElement('div');
    messageContent.classList.add('message-content');
    messageContent.textContent = post.content;

    messageItem.appendChild(messageHeader);
    messageItem.appendChild(messageContent);

    messagesList.appendChild(messageItem);
    messagesList.scrollTop = messagesList.scrollHeight;
}


async function deletePost(post_uuid) {
    const confirmDelete = confirm("Êtes-vous sûr de vouloir supprimer ce post ?");
    if (!confirmDelete) return;

    try {
        const response = await fetch("http://localhost:8080/api/post/deletePost", {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ post_uuid: post_uuid }),
        });

        if (!response.ok) {
            throw new Error("Erreur lors de la suppression du post");
        }

        // Recharger les posts après suppression
        fetchPosts();
    } catch (error) {
        console.error(error);
    }
}

async function createPost(event) {
    event.preventDefault();
    const userName = document.getElementById('user-name').textContent;
    const messageInput = document.getElementById("message");
    const messageContent = messageInput.value;

    const data = {
        userName: userName,
        content: messageContent,
    };

    const response = await fetch("http://localhost:8080/api/post/createPost", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    });

    if (response.ok) {
        alert("Post créé avec succès!");
        messageInput.value = '';
    } else {
        const error = await response.json();
        alert("Erreur lors de la création du post: " + error.message);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    const addButton = document.getElementById('add-button');
    const modalPost = document.getElementById('modal-post');
    const userPost = document.getElementById('users-post')

    let isModal = false;

    function NewPost() {
        CreatedModal();
        const newpost = document.getElementById('created-post');
        newpost.style.display = 'flex';
        modalPost.style.display = 'flex';
        isModal = true;

        // Ajouter un écouteur d'événement pour fermer le modal lorsqu'un clic se produit
        document.addEventListener('click', closeModal);
    }

    // Fonction pour récupérer un cookie
    function getCookie(name) {
        const value = `; ${document.cookie}`;
        const parts = value.split(`; ${name}=`);

        if (parts.length === 2) return parts.pop().split(';').shift();
    }

    const token = getCookie('token'); // Récupère le cookie

    if (token) {
        console.log('Token récupéré:', token);
    } else {
        console.log('Cookie non trouvé.');
    }

    function CreatedModal() {
        const createdPost = document.createElement('div');
        createdPost.classList.add('created-post');
        createdPost.id = 'created-post';

        const postHeader = document.createElement('div');
        postHeader.classList.add('post');
        postHeader.textContent = 'New Post';

        const userName = document.createElement('div');
        userName.classList.add('user-name');
        userName.id = 'user-name'
        userName.textContent = token;

        const form = document.createElement('form');
        form.classList.add('message-input');
        form.id = 'message-form';

        const inputField = document.createElement('input');
        inputField.type = 'text';
        inputField.id = 'message';
        inputField.placeholder = "What's new ?";

        const postButton = document.createElement('button');
        postButton.id = 'post-btn';
        postButton.textContent = 'Post';

        form.appendChild(inputField);
        createdPost.appendChild(postHeader);
        createdPost.appendChild(userName);
        createdPost.appendChild(form);
        createdPost.appendChild(postButton);

        // Ajout du formulaire dans le modal
        userPost.appendChild(createdPost);

        postButton.addEventListener('click', createPost)
    }


    // Fonction pour fermer le modal si on clique à l'extérieur de 'created-post'
    function closeModal(event) {
        const newpost = document.getElementById('created-post');
        const modalpost = document.getElementById('modal-post');
        // Vérifie si le clic a eu lieu en dehors de 'created-post'
        if (!newpost.contains(event.target) && event.target !== addButton && isModal) {
            newpost.style.display = 'none'; // Ferme le post
            modalpost.style.display = 'none'; // Ferme aussi le fond modal
            isModal = false;

            document.removeEventListener('click', closeModal);
        }
    }

    addButton.addEventListener('click', NewPost);

});