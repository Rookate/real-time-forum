import { fetchAllUsers } from "./api/fetchAllUsers.js";
import { getConversation } from "./api/fetchConversation.js";

// Container
const container = document.querySelector('.chat[data-section="private-message"]');
const sidebar = document.getElementById('sidebar');

// Exemple de données
const conversationData = {
    username: "Luna",
    profile_picture: "https://ih1.redbubble.net/image.3546319896.8009/flat,750x,075,f-pad,750x1000,f8f8f8.jpg",
    // messages: [
    //     {
    //         username: "Rokat",
    //         profile_picture: "https://ih1.redbubble.net/image.3546319896.8009/flat,750x,075,f-pad,750x1000,f8f8f8.jpg",
    //         content: "C'est en 1901 dans Psychopathologie de la vie quotidienne que Sigmund Freud détaille le fonctionnement du lapsus.",
    //         date: "Wed, 05 February 2025 11:40"
    //     },
    //     {
    //         username: "Luna",
    //         profile_picture: "https://ih1.redbubble.net/image.3546319896.8009/flat,750x,075,f-pad,750x1000,f8f8f8.jpg",
    //         content: "Ce concept de lapsus révèle beaucoup sur la relation entre le conscient et l'inconscient.",
    //         date: "Wed, 05 February 2025 11:45"
    //     },
    //     {
    //         username: "Rokat",
    //         profile_picture: "https://ih1.redbubble.net/image.3546319896.8009/flat,750x,075,f-pad,750x1000,f8f8f8.jpg",
    //         content: "C'est en 1901 dans Psychopathologie de la vie quotidienne que Sigmund Freud détaille le fonctionnement du lapsus.",
    //         date: "Wed, 05 February 2025 11:40"
    //     },
    // {
    //     username: "Luna",
    //     profile_picture: "https://ih1.redbubble.net/image.3546319896.8009/flat,750x,075,f-pad,750x1000,f8f8f8.jpg",
    //     content: "Ce concept de lapsus révèle beaucoup sur la relation entre le conscient et l'inconscient.",
    //     date: "Wed, 05 February 2025 11:45"
    // },
    // {
    //     username: "Rokat",
    //     profile_picture: "https://ih1.redbubble.net/image.3546319896.8009/flat,750x,075,f-pad,750x1000,f8f8f8.jpg",
    //     content: "C'est en 1901 dans Psychopathologie de la vie quotidienne que Sigmund Freud détaille le fonctionnement du lapsus.",
    //     date: "Wed, 05 February 2025 11:40"
    // },
    // {
    //     username: "Luna",
    //     profile_picture: "https://ih1.redbubble.net/image.3546319896.8009/flat,750x,075,f-pad,750x1000,f8f8f8.jpg",
    //     content: "Ce concept de lapsus révèle beaucoup sur la relation entre le conscient et l'inconscient.",
    //     date: "Wed, 05 February 2025 11:45"
    // }
    // ]
};

function conversation(obj) {
    container.innerHTML = '';
    // Header
    const header = displayHeader(obj);
    container.appendChild(header);

    // Display the conversation messages
    displayMessage(obj.messages);

    // Display the input area
    const input = displayInput(obj);
    container.appendChild(input);
}

function displayMessage(messages) {
    if (!messages || messages.length === 0) {
        const noMessage = document.createElement('span');
        noMessage.classList.add('no-msg-content')
        noMessage.textContent = "Start a new conversation with @Rokat"

        container.appendChild(noMessage)
        return
    }

    const divContainer = document.createElement('div');
    divContainer.classList.add("message-container-div");

    messages.forEach(messageData => {
        const chat = document.createElement('div');
        chat.classList.add('message');

        // Message content
        const messageContent = displayContentMessage(messageData);

        // Append to container
        chat.appendChild(messageContent);
        divContainer.appendChild(chat);
    });

    container.appendChild(divContainer);
}

function displayContentMessage(content) {
    const userMessage = document.createElement('div');
    userMessage.classList.add('user-message');

    // Profile Picture
    const profilePicture = document.createElement('div');
    profilePicture.classList.add('profile-picture-user');

    const image = document.createElement('img');
    image.classList.add('pp-discord');
    image.src = content.profile_picture;
    image.alt = content.username;

    // User message container
    const userMessageContainer = document.createElement('div');
    userMessageContainer.classList.add('user-message-container');

    // User info
    const userInfo = document.createElement('div');
    userInfo.classList.add('user-info');

    const username = document.createElement('span');
    username.classList.add('username');
    username.textContent = content.username;

    const timestamp = document.createElement('span');
    timestamp.classList.add('timestamp');
    timestamp.textContent = content.date;

    const contentMessage = document.createElement('p');
    contentMessage.textContent = content.content;

    // Assemble the message
    profilePicture.appendChild(image);
    userInfo.appendChild(username);
    userInfo.appendChild(timestamp);
    userMessageContainer.appendChild(userInfo);
    userMessageContainer.appendChild(contentMessage);
    userMessage.appendChild(profilePicture);
    userMessage.appendChild(userMessageContainer);

    return userMessage;
}

function displayHeader(content) {
    const header = document.createElement('header');

    const image = document.createElement('img');
    image.classList.add('pp-header-discord');
    image.src = content.profile_picture;
    image.alt = content.username;

    const username = document.createElement('span');
    username.textContent = content.username;

    header.appendChild(image);
    header.appendChild(username);

    return header;
}

function displayInput(content) {
    const inputUser = document.createElement('div');
    inputUser.classList.add('input-user');

    const input = document.createElement('input');
    input.type = "text";
    input.placeholder = `Message @${content.username}`;

    const sendButton = document.createElement('button');
    sendButton.id = "send-chat";
    sendButton.textContent = "Send";

    inputUser.appendChild(input);
    inputUser.appendChild(sendButton);

    return inputUser;
}

// Appel de la fonction conversation avec les données
// conversation(conversationData);


function displayConversation(content) {
    // Conversation Item
    const conversationItem = document.createElement('div');
    conversationItem.classList.add('conversation-item')

    // Profile Picture
    const imageContainer = document.createElement('div');
    imageContainer.classList.add('image-container')


    const image = document.createElement('img');
    image.classList.add('pp-discord')
    image.src = content.profile_picture;
    image.alt = content.username;

    //Username
    const username = document.createElement('span');
    username.textContent = content.username;

    //Append all element
    imageContainer.appendChild(image);
    conversationItem.appendChild(imageContainer);
    conversationItem.appendChild(username);

    return conversationItem;
}

function createFriendList(friends) {
    const friendsContainer = document.createElement('div');
    friendsContainer.classList.add('friend-list');

    // Create search input
    const searchDiv = document.createElement('div');
    searchDiv.id = 'search-friend-div';

    const searchInput = document.createElement('input');
    searchInput.type = 'text';
    searchInput.id = 'input-friends-list';
    searchInput.placeholder = 'Search';

    searchDiv.appendChild(searchInput);
    container.appendChild(searchDiv);

    // Create friends list
    friends.forEach(friend => {
        const friendDiv = document.createElement('div');
        friendDiv.addEventListener('click', () => showConversation(friend.user_uuid))
        friendDiv.classList.add('friend');

        const profilePic = document.createElement('img');
        profilePic.src = friend.profile_picture || "https://upload.wikimedia.org/wikipedia/commons/8/87/Chimpanzee-Head.jpg?uselang=fr";
        profilePic.classList.add('pp-header-discord');

        const username = document.createElement('span');
        username.textContent = friend.username;

        friendDiv.appendChild(profilePic);
        friendDiv.appendChild(username);
        friendsContainer.appendChild(friendDiv);
    });

    return friendsContainer;
}

async function showConversation(user_uuid) {
    const conv = await getConversation(user_uuid);
    if (conv) {
        conversation(conv);
    }
}

export async function showFriendsList() {
    sidebar.classList.add('close');

    const users = await fetchAllUsers()
    // const friendsList = createFriendList(users)
    // container.appendChild(friendsList)
}

const messages = [
    {
        username: "Rokat", // Nom de l'utilisateur
        profile_picture: "https://example.com/pp-rokat.jpg", // URL de l'image de profil
        content: "C'est en 1901 dans Psychopathologie de la vie quotidienne que Sigmund Freud détaille le fonctionnement du lapsus.",
        date: "Wed, 05 February 2025 11:40", // Horodatage du message
    },
    {
        username: "Luna",
        profile_picture: "https://example.com/pp-luna.jpg",
        content: "C'est fascinant! Je vais devoir lire ça.",
        date: "Wed, 05 February 2025 11:45",
    }
    // Ajoute d'autres messages ici
];

//