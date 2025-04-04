import { DisplayMessages } from "./displayMessage.js";
import { initEventListeners } from "./comment.js";
import { getPPFromID, throttle } from "./utils.js";
import { NewPost } from "./newPost.js";
import { handleLogout } from "./logout.js";
import { toggleCommentReaction, toggleReaction } from "./reaction.js";
import { FetchMostUseCategories } from "./tendance.js";
import { fetchNotifications } from "./notifs.js";
import { initSectionEvents, showDashboard } from "./section.js";


export let UserInfo = null
export async function fetchUserInfo() {

    // Vérifie si le cookie 'session_token' existe
    const cookies = document.cookie.split('; ').reduce((acc, cookie) => {
        const [name, value] = cookie.split('=');
        acc[name] = value;
        return acc;
    }, {});

    if (!cookies['session_token']) {
        return;
    }

    try {
        const response = await fetch("http://localhost:8080/api/getSession");
        if (!response.ok) {
            document.cookie = "session_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
            throw new Error("Error retrieving user data");
        }

        const data = await response.json();
        UserInfo = data;

    } catch (error) {
        console.error(error);
    }
}

// Sélectionnez les éléments
const toggleButton = document.getElementById('toggle-menu-btn');
const sidebar = document.getElementById('sidebar');
const addButton = document.getElementById('add-button');

const darkModeToggles = document.querySelectorAll('.dark-mode-toggle');

darkModeToggles.forEach(button => {
    button.addEventListener('click', () => {
        console.log('Dark mode toggled');
    });
});


// Fonction pour appliquer le mode
function applyMode(mode) {
    const root = document.documentElement;

    if (mode === 'dark') {
        root.style.setProperty('--background-color', '#36393F');
        root.style.setProperty('--text-color', '#DCDDDE');
        root.style.setProperty('--second-text-color', '#FFFFFF');
        root.style.setProperty('--border-color', '#51545c');
        root.style.setProperty('--background-message-color', '#2F3136');
        root.style.setProperty('--accent-color', '#5865F2'); // Bleu Discord
        root.style.setProperty('--online-color', '#3BA55C'); // Vert (en ligne)
        root.style.setProperty('--offline-color', '#ED4245'); // Rouge (hors ligne)
        root.style.setProperty('--busy-color', '#FAA61A'); // Orange (occupé)
        root.style.setProperty('--idle-color', '#F0B126'); // Jaune (AFK)
        darkModeToggles.textContent = 'Light Mode';
    } else {
        root.style.setProperty('--background-color', '#FFFFFF');
        root.style.setProperty('--text-color', '#202225'); // Texte sombre
        root.style.setProperty('--second-text-color', '#5865F2'); // Bleu Discord
        root.style.setProperty('--border-color', '#E3E5E8'); // Bordure claire
        root.style.setProperty('--background-message-color', '#F6F6F7'); // Fond clair pour les messages
        root.style.setProperty('--accent-color', '#5865F2'); // Bleu Discord
        root.style.setProperty('--online-color', '#3BA55C'); // Vert (en ligne)
        root.style.setProperty('--offline-color', '#ED4245'); // Rouge (hors ligne)
        root.style.setProperty('--busy-color', '#FAA61A'); // Orange (occupé)
        root.style.setProperty('--idle-color', '#F0B126'); // Jaune (AFK)
        darkModeToggles.textContent = 'Dark Mode';
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
darkModeToggles.forEach(button => {
    button.addEventListener('click', () => {
        const currentMode = localStorage.getItem('theme') || 'light';
        const newMode = currentMode === 'dark' ? 'light' : 'dark';
        applyMode(newMode);
        // Mettre à jour le texte du bouton
        button.textContent = newMode === 'dark' ? 'Light Mode' : 'Dark Mode';
    });
});


document.addEventListener('DOMContentLoaded', async () => {
    await fetchUserInfo();

    const loginButton = document.getElementById('login-btn');
    const profilMenu = document.querySelector('.profil-menu');
    const logOut = document.getElementById('logout-button');

    // Vérifiez si les informations utilisateur sont valides
    if (UserInfo) {
        logOut.style.display = "flex"
        // Créer la div qui remplacera le bouton "Login"
        const profileDiv = document.createElement('div');
        profileDiv.classList.add('profile-container');

        const profileImage = document.createElement('img');
        profileImage.src = await getPPFromID(UserInfo.user_uuid);
        profileImage.alt = 'User profile';
        profileImage.classList.add('profile-image');

        // Créer le menu contextuel
        const menu = document.createElement('div');
        menu.classList.add('profile-menu');
        menu.style.display = 'none';

        const logoutButton = document.createElement('button');
        logoutButton.id = "logout-btn";
        logoutButton.textContent = 'Log Out';
        logoutButton.addEventListener('click', handleLogout);
        logOut.addEventListener('click', handleLogout)

        menu.appendChild(logoutButton);

        profileDiv.appendChild(profileImage);
        profileDiv.appendChild(menu);

        // Événement pour afficher/masquer le menu lorsque l'image est cliquée
        profileImage.addEventListener('click', () => {
            // Afficher ou masquer le menu
            menu.style.display = (menu.style.display === 'none') ? 'block' : 'none';
        });

        document.addEventListener('click', (event) => {
            if (!profileDiv.contains(event.target)) {
                menu.style.display = 'none';
            }
        });
        addButton.style.display = 'block';

        // Remplacer le bouton "Login" par la div
        profilMenu.replaceChild(profileDiv, loginButton);


        const moderationLink = document.getElementById('moderation-link')

        if (UserInfo.role !== "admin" && UserInfo.role !== "GOAT") {
            moderationLink.remove();
        }
        return;
    }

    const elementsToHide = [
        document.getElementById('notifications-link'),
        document.getElementById('request-link'),
        document.getElementById('moderation-link'),
        document.getElementById('private-message-link'),
    ];

    elementsToHide.forEach(element => {
        if (element) {
            const parentLi = element.closest('li');
            if (parentLi) parentLi.remove();
        }
    });

    // Si le bouton "Login" n'est pas déjà dans le menu
    if (!profilMenu.contains(loginButton)) {
        logOut.style.display = 'none';
        // Créer un nouveau bouton "Login" si nécessaire
        const newLoginButton = document.createElement('button');
        newLoginButton.id = 'login-btn';
        newLoginButton.textContent = 'Log in';

        // Ajouter à nouveau le bouton à la place de la div
        profilMenu.appendChild(newLoginButton);
    }

});


function toggleSidebar() {
    sidebar.classList.toggle('close');
}
// Ajouter l'événement
if (window.screen.width >= 768) {
    toggleButton.addEventListener('click', toggleSidebar);
}


export async function fetchPosts() {
    const messagesList = document.querySelector(`.users-post[data-section="home"]`);
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
            posts.sort((b, a) => new Date(b.created_at) - new Date(a.created_at));
            posts.forEach(post => {
                DisplayMessages(post, "home");
            });
        }
    } catch (error) {
        messagesList.innerHTML = '<p>Error loading posts. Please try again.</p>';
        console.error(error);
    }
    initEventListeners();
    fetchNotifications();
}


export function Reaction(event) {
    const likeButton = event.target.closest('.like-btn');
    const dislikeButton = event.target.closest('.dislike-btn');

    if (likeButton || dislikeButton) {
        const messageItem = (likeButton || dislikeButton).closest('.message-item');
        if (messageItem) {
            const postUuid = messageItem.getAttribute('post_uuid');
            toggleReaction(event, postUuid);
        }
    }
}

export function CommentReaction(event) {
    console.log("ça passe dans la fonction commentReact");
    const likeButton = event.target.closest('.like-comment-btn');
    const dislikeButton = event.target.closest('.dislike-comment-btn');

    console.log("ça trouve like button :", likeButton);
    console.log("ça trouve dislike button :", dislikeButton);

    if (likeButton || dislikeButton) {
        const messageItem = (likeButton || dislikeButton).closest('.message-item');
        if (messageItem) {
            const postUuid = messageItem.getAttribute('post_uuid');
            console.log("post_uuid :", postUuid);
            toggleCommentReaction(event, postUuid);
        }
    }
}

document.addEventListener('DOMContentLoaded', async () => {
    fetchPosts();
    addButton.addEventListener('click', NewPost);
    FetchMostUseCategories();

    // Gérer uniquement les événements sur les boutons de post
    document.body.addEventListener('click', (event) => {
        const isPostReaction = event.target.closest('.like-btn') || event.target.closest('.dislike-btn');
        const isCommentReaction = event.target.closest('.like-comment-btn') || event.target.closest('.dislike-comment-btn');

        if (isPostReaction) {
            Reaction(event);
        } else if (isCommentReaction) {
            CommentReaction(event);
        }
    });

    await fetchUserInfo();
    initSectionEvents(UserInfo);
    document.getElementById('dashboard-link').addEventListener('click', () => showDashboard(UserInfo));
});

function handleResize() {
    const elements = ['request-link', 'moderation-link'];

    if (window.innerWidth <= 480) {
        for (const section of elements) {
            const element = document.getElementById(section);
            if (element) {
                element.style.display = 'none';
            }
        }
        sidebar.classList.remove('close');
    } else {
        for (const section of elements) {
            const element = document.getElementById(section);
            if (element) {
                element.style.display = 'flex';
            }
        }
    }
}

const friendList = document.getElementById('close-friend-list');
const convList = document.getElementById('conversations-list');
const openFriendList = document.getElementById('open-friend-list');

function toggleCloseFriend() {
    console.log('click')
    friendList.classList.toggle('close');
    convList.classList.toggle('close');


    if (convList.classList.contains('close')) {
        // Ajouter un écouteur d'événement pour détecter la fin de la transition
        convList.addEventListener('transitionend', function onTransitionEnd() {
            convList.style.display = 'none';
            openFriendList.style.display = 'block'
            convList.removeEventListener('transitionend', onTransitionEnd);
        }, { once: true });
    } else {
        convList.style.display = 'flex';
        openFriendList.style.display = 'none'
    }
}
friendList.addEventListener('click', toggleCloseFriend)
openFriendList.addEventListener('click', toggleCloseFriend)


const throttleResize = throttle(handleResize, 500)

window.addEventListener('resize', throttleResize)
handleResize();
