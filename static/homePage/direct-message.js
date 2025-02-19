import { fetchAllUsers } from "./api/fetchAllUsers.js";
import { createConversation } from "./api/createConversation.js";
import { formatTimestamp, throttle, updateURL } from "./utils.js";
import { fetchMessages } from "./api/fetchMessages.js";
import { sendMessage } from "./api/createMessage.js";
import { isTipping } from "./api/isTipping.js";
import { fetchAllConversations } from "./api/fetchAllConversations.js";


export let ws;
let convId
let activeUser
let currentIndexMessages = 0
// Container
const container = document.querySelector('.chat[data-section="private-message"]');
const sidebar = document.getElementById('sidebar');

const divContainer = document.createElement('div');
divContainer.classList.add("message-container-div");



ws = new WebSocket("ws://localhost:8080/ws");
ws.onopen = function () {
    console.log('Open')
}
ws.onmessage = function (event) {
    const data = JSON.parse(event.data);
    console.log('data recu', data)

    if (data.type === "single_message") {
        displaySingleMessage(data);
    } else if (data.type === "moreMessages") {
        displayMessage(data, true)
    } else if (data.type === "notification") {
        const messageDot = document.getElementById('private-message-dot')
        messageDot.style.display = "block"
    }
    else if (data.type === "messages") {
        displayMessage(data);
    } else if (data.type === "typing") {
        const typingSpan = document.getElementById('typing-span')
        if (data.isTyping) {
            typingSpan.style.visibility = "visible"
        } else {
            typingSpan.style.visibility = "hidden"
        }
    } else if (data.type === "active_users") {
        activeUser = data.active_users
    }
};




// const handleScroll = async () => {
//     if (
//         !isFetching &&
//         window.innerHeight + window.scrollY >= document.body.offsetHeight - 100
//     ) {
//         isFetching = true;
//         console.log("OUI")
//         displayMessage();
//         isFetching = false;
//     }
// };

// const throttleScroll = throttle(handleScroll, 300);
// window.addEventListener("scroll", throttleScroll);

function conversation(obj) {
    const conversationUUID = obj.conversation_uuid;
    updateURL(conversationUUID)

    container.innerHTML = '';
    divContainer.innerHTML = '';
    // Header
    const header = displayHeader(obj);
    container.appendChild(header);

    // Display the conversation messages
    currentIndexMessages = 0
    currentIndexMessages = fetchMessages(conversationUUID, null, currentIndexMessages);
    convId = conversationUUID

    // Message container
    container.appendChild(divContainer)

    // Display the input area
    const input = displayInput(obj);
    container.appendChild(input);

    const isTypingSpan = document.createElement('span');
    isTypingSpan.classList.add('typing-span');
    isTypingSpan.id = 'typing-span'
    isTypingSpan.innerHTML = `
    <span class="typing-dots">
        <span></span>
        <span></span>
        <span></span>
    </span> 
    ${obj.receiver_username} is typing...
`;

    container.appendChild(isTypingSpan);
}

function displayMessage(data, prepend = false) {
    const scrollPosition = divContainer.scrollTop;
    const initialScrollHeight = divContainer.scrollHeight;

    data.messages.forEach(messageData => {
        const chat = document.createElement('div');
        chat.classList.add('message');

        // Message content
        const messageContent = displayContentMessage(messageData);
        chat.appendChild(messageContent);

        divContainer.prepend(chat);
    });

    if (prepend) {
        // Maintenir la position relative du scroll lors de l'ajout en haut
        const newScrollHeight = divContainer.scrollHeight;
        const scrollDiff = newScrollHeight - initialScrollHeight;
        divContainer.scrollTop = scrollPosition + scrollDiff;
    } else {
        // Scroll tout en bas pour les nouveaux messages
        divContainer.scrollTop = divContainer.scrollHeight;
    }
}

function displaySingleMessage(message) {
    const chat = document.createElement('div');
    chat.classList.add('message');

    // Message content
    const messageContent = displayContentMessage(message);

    // Append to container
    chat.appendChild(messageContent);
    divContainer.appendChild(chat);
}

function displayContentMessage(content) {

    const userMessage = document.createElement('div');
    userMessage.classList.add('user-message');

    // Profile Picture
    const profilePicture = document.createElement('div');
    profilePicture.classList.add('profile-picture-user');

    const image = document.createElement('img');
    image.classList.add('pp-discord');
    image.src = content.sender_profile_picture;
    image.alt = content.sender_username;

    // User message container
    const userMessageContainer = document.createElement('div');
    userMessageContainer.classList.add('user-message-container');

    // User info
    const userInfo = document.createElement('div');
    userInfo.classList.add('user-info');

    const username = document.createElement('span');
    username.classList.add('username');
    username.textContent = content.sender_username;

    const timestamp = document.createElement('span');
    timestamp.classList.add('timestamp');
    timestamp.textContent = formatTimestamp(content.created_at);

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
    image.src = content.receiver_profile_picture || "https://upload.wikimedia.org/wikipedia/commons/8/87/Chimpanzee-Head.jpg?uselang=fr";
    image.alt = content.receiver_username;

    const username = document.createElement('span');
    username.textContent = content.receiver_username;

    header.appendChild(image);
    header.appendChild(username);

    return header;
}

function displayInput(content) {
    const inputUser = document.createElement('div');
    inputUser.classList.add('input-user');

    const input = document.createElement('input');
    input.type = "text";
    input.id = 'messageInput'
    input.placeholder = `Message @${content.receiver_username}`;
    input.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') {
            sendMessage(content.conversation_uuid, content.reciever, content.sender, content.sender_username, content.sen);
        }
    })
    input.addEventListener("input", isTipping);
    const sendButton = document.createElement('button');
    sendButton.id = "send-chat";
    sendButton.textContent = "Send";

    inputUser.appendChild(input);
    inputUser.appendChild(sendButton);

    return inputUser;
}

// Appel de la fonction conversation avec les données
// conversation(conversationData);
const conversationsList = document.getElementById("conversation-container"); // Conteneur principal


function displayAllConversations(conversations) {

    const sortedConv = conversations.sort((a, b) => b.update_at.localeCompare(a.update_at))
    sortedConv.forEach(content => {
        const convItem = displayConversation(content);
        conversationsList.appendChild(convItem)
    });
}

async function displayConversationHandler() {
    const conv = await fetchAllConversations();

    // list conversation
    conversationsList.innerHTML = ""; // Nettoyer avant d'ajouter les nouvelles conversations
    displayAllConversations(conv)

}

function displayConversation(content) {
    // Conversation Item
    const conversationItem = document.createElement('div');
    conversationItem.classList.add('conversation-item')
    conversationItem.addEventListener('click', () => {
        showConversation(content.receiver)
    })

    // Profile Picture
    const imageContainer = document.createElement('div');
    imageContainer.classList.add('image-container')


    const image = document.createElement('img');
    image.classList.add('pp-discord')
    image.src = content.receiver_profile_picture || "https://koreus.cdn.li/media/201404/90-insolite-34.jpg";
    image.alt = content.receiver_username;

    //Username
    const username = document.createElement('span');
    username.textContent = content.receiver_username;

    //Append all element
    imageContainer.appendChild(image);
    conversationItem.appendChild(imageContainer);
    conversationItem.appendChild(username);

    return conversationItem
}

function createFriendList(friends, clients) {
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
        friendDiv.addEventListener('click', () => {
            showConversation(friend.user_uuid)
        })
        friendDiv.classList.add('friend');

        const imageWrapper = document.createElement('div');
        imageWrapper.classList.add('image-wrapper');
        if (clients.includes(friend.user_uuid)) {
            const online = document.createElement('div');
            online.classList.add('online');
            imageWrapper.appendChild(online);
        }

        const profilePic = document.createElement('img');
        profilePic.src = friend.profile_picture || "https://upload.wikimedia.org/wikipedia/commons/8/87/Chimpanzee-Head.jpg?uselang=fr";
        profilePic.classList.add('pp-header-discord');

        const username = document.createElement('span');
        username.textContent = friend.username;

        imageWrapper.appendChild(profilePic);
        friendDiv.appendChild(imageWrapper);
        friendDiv.appendChild(username);
        friendsContainer.appendChild(friendDiv);
    });

    return friendsContainer;
}

async function showConversation(user_uuid) {
    container.innerHTML = "";
    const conv = await createConversation(user_uuid);
    console.log("conv recu", conv)
    if (conv) {
        conversation(conv);
    }
}

export async function showFriendsList() {
    container.innerHTML = "";
    sidebar.classList.add('close');

    displayConversationHandler()

    const privateMessageDot = document.getElementById("private-message-dot")
    privateMessageDot.style.display = "none"

    const users = await fetchAllUsers()
    const sortedUser = users.sort((a, b) => a.username.localeCompare(b.username));
    const friendsList = createFriendList(sortedUser, activeUser)
    container.appendChild(friendsList)
}

const OFFSET_TRIGGER = 50;
divContainer.addEventListener("scroll", throttle(async () => {
    if (divContainer.scrollTop <= OFFSET_TRIGGER) {
        currentIndexMessages = fetchMessages(convId, true, currentIndexMessages)
    }
}, 500))

ws.onclose = function () {
    console.log("WebSocket connection closed, retrying...");
    setTimeout(connect, 1000);
};

ws.onerror = function (error) {
    console.error("WebSocket error:", error);
};

// const privateMessageLink = document.getElementById('private-message-link');

// const moveButton = function () {
//     const randomX = Math.random() * (window.innerWidth - privateMessageLink.offsetWidth);
//     const randomY = Math.random() * (window.innerHeight - privateMessageLink.offsetHeight);
//     privateMessageLink.style.position = 'absolute'
//     privateMessageLink.style.zIndex = '99999'
//     privateMessageLink.style.left = randomX + "px";
//     privateMessageLink.style.top = randomY + "px";
// };
// privateMessageLink.addEventListener('mouseenter', moveButton)
// Object.defineProperty(window, "console", { value: null })
