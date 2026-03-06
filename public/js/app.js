const chatBox = document.getElementById("chat-box");
const msgInput = document.getElementById("msg-input");
const sendBtn = document.getElementById("send-btn");
const connStatus = document.getElementById("conn-status");
const statusDot = document.querySelector(".status-dot");
const attachBtn = document.getElementById("attach-btn");
const imageUpload = document.getElementById("image-upload");
const loginOverlay = document.getElementById("login-overlay");
const passwordInput = document.getElementById("password-input");
const loginError = document.getElementById("login-error");
const usernameInput = document.getElementById("username-input");

let token = "";
let socket;
let username = "";

usernameInput.addEventListener("keypress", function(event) {
    if (event.key === "Enter") {
        passwordInput.focus();
    }
});

passwordInput.addEventListener("keypress", function(event) {
    if (event.key === "Enter") {
        startChat();
    }
});

async function startChat() {
    const valUser = usernameInput.value.trim();
    const valPass = passwordInput.value.trim();
    
    if (valUser === "" || valPass === "") {
        loginError.innerText = "Kullanıcı adı ve şifre gereklidir.";
        loginError.style.display = "block";
        return;
    }

    try {
        const response = await fetch("/api/auth", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                username: valUser,
                password: valPass
            })
        });

        const data = await response.json();

        if (!response.ok) {
            loginError.innerText = data.error || "Giriş yapılamadı. Tekrar deneyin.";
            loginError.style.display = "block";
            return;
        }

        token = data.token;
        username = valUser;
        loginOverlay.classList.add("hidden");
        
        connectWebSocket();
    } catch (error) {
        console.error("Auth hatası:", error);
        loginError.innerText = "Sunucuya bağlanılamadı.";
        loginError.style.display = "block";
    }
}

function connectWebSocket() {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const host = window.location.host || "localhost:8080";
    const wsUrl = protocol + "//" + host + "/ws?token=" + encodeURIComponent(token);
    
    socket = new WebSocket(wsUrl);

    socket.onopen = function() {
        connStatus.innerText = "Çevrimiçi (" + username + ")";
        statusDot.style.backgroundColor = "#10b981"; 
        msgInput.disabled = false;
        sendBtn.disabled = false;
        attachBtn.disabled = false;
        msgInput.focus();
        
        appendSystemMessage("Sohbete katıldınız.");
    };

    socket.onmessage = function(event) {
        try {
            const data = JSON.parse(event.data);
            
            if (data.sender && data.message) {
                const isSelf = data.sender === username;
                appendMessage(data.sender, data.message, isSelf, data.type);
            }
        } catch (e) {
            console.error("Geçersiz veri formatı:", event.data);
        }
    };

    socket.onclose = function(event) {
        connStatus.innerText = "Bağlantı koptu";
        statusDot.style.backgroundColor = "#ef4444"; 
        msgInput.disabled = true;
        sendBtn.disabled = true;
        attachBtn.disabled = true;
        appendSystemMessage("Sunucu bağlantısı kesildi.");
    };
    
    socket.onerror = function() {
        connStatus.innerText = "Bağlantı Hatası";
        statusDot.style.backgroundColor = "#ef4444"; 
    };
}

function sendMessage() {
    const message = msgInput.value.trim();
    if (message !== "") {
        socket.send(JSON.stringify({ type: "text", message: message }));
        msgInput.value = ""; 
        msgInput.focus();
    }
}

imageUpload.addEventListener("change", async function() {
    if (this.files && this.files.length > 0) {
        const file = this.files[0];
        const formData = new FormData();
        formData.append("file", file);

        try {
            const response = await fetch("/api/upload", {
                method: "POST",
                body: formData
            });

            if (response.ok) {
                const data = await response.json();
                socket.send(JSON.stringify({ type: "image", message: data.url }));
            } else {
                alert("Resim yüklenemedi.");
            }
        } catch (error) {
            console.error("Upload error:", error);
            alert("Resim yükleme hatası.");
        }
        
        this.value = "";
    }
}
);

msgInput.addEventListener("keypress", function(event) {
    if (event.key === "Enter") {
        sendMessage();
    }
});

function appendMessage(sender, text, isSelf, type = "text") {
    const msgElement = document.createElement("div");
    msgElement.className = "message " + (isSelf ? "self" : "other");
    
    let html = "";
    if (!isSelf) {
        html += `<div class="sender-name">${escapeHtml(sender)}</div>`;
    }
    
    if (type === "image") {
        html += `<div class="bubble"><img src="${escapeHtml(text)}" class="chat-image" alt="Gönderilen Resim" onclick="window.open(this.src, '_blank')"></div>`;
    } else {
        html += `<div class="bubble">${escapeHtml(text)}</div>`;
    }
    
    msgElement.innerHTML = html;
    chatBox.appendChild(msgElement);
    
    if (type === "image") {
        const img = msgElement.querySelector('img');
        if (img) {
            img.onload = scrollToBottom;
        }
    }
    
    scrollToBottom();
}

function appendSystemMessage(text) {
    const msgElement = document.createElement("div");
    msgElement.className = "system-message";
    msgElement.innerText = text;
    chatBox.appendChild(msgElement);
    scrollToBottom();
}

function scrollToBottom() {
    chatBox.scrollTop = chatBox.scrollHeight;
}

function escapeHtml(unsafe) {
    return String(unsafe)
         .replace(/&/g, "&amp;")
         .replace(/</g, "&lt;")
         .replace(/>/g, "&gt;")
         .replace(/"/g, "&quot;")
         .replace(/'/g, "&#039;");
}
