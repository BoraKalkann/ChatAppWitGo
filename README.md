# Go-ChatApp 🚀

Go-ChatApp is a real-time, modern chat application backend and web frontend built with **Go (Golang)**, **WebSockets**, and **MongoDB**. It features user authentication, real-time messaging, and multimedia (image) upload support.

## Features ✨

- **Real-time Messaging:** Lightning-fast communication powered by `gorilla/websocket`.
- **User Authentication:** Secure user registration and login system.
- **Password Hashing:** Passwords are fully secured using `bcrypt`.
- **JWT Authorization:** Stateless session management via JSON Web Tokens for API and WebSocket handshakes.
- **Image Uploads:** Users can upload and share `.jpg`, `.png`, and `.gif` files seamlessly in the chat.
- **Persistent Chat History:** All messages and users are safely stored in **MongoDB**, allowing users to see past messages when they log in.
- **Modern UI:** Clean, responsive, dark-mode CSS frontend with message bubbles and self/other distinctions.

---

## Tech Stack 🛠

- **Backend:** Go (Golang)
- **Database:** MongoDB (via `mongo-driver`)
- **WebSockets:** Gorilla WebSocket (`github.com/gorilla/websocket`)
- **Security:** `golang.org/x/crypto/bcrypt`, `github.com/golang-jwt/jwt/v5`
- **Environment Management:** `github.com/joho/godotenv`
- **Frontend:** Vanilla HTML, CSS, JavaScript (Fetch API, WebSockets API)

---

## Installation & Setup ⚙️

### 1. Prerequisites
- **Go:** `v1.20` or higher installed. ([Download Go](https://go.dev/))
- **MongoDB:** A running instance of MongoDB (Local or Atlas). Default configuration expects it at `localhost:27017`.

### 2. Clone the Repository
```bash
git clone https://github.com/BoraKalkann/ChatAppWitGo.git
cd ChatAppWitGo
```

### 3. Environment Variables
Create a `.env` file in the root directory. You can copy the template:
```bash
cp .env.example .env
```
Open `.env` and configure your secrets:
```env
MONGO_URI=mongodb://localhost:27017
JWT_SECRET=your_super_secret_jwt_key_here
```

### 4. Install Dependencies
```bash
go mod tidy
```

### 5. Run the Application
```bash
go run main.go
```
The server will start at `http://localhost:8080`.

---

## Usage 💻

1. Open your browser and navigate to `http://localhost:8080`.
2. Enter a username and password. 
   - *If the username doesn't exist, an account will be automatically created.*
   - *If the username exists, you will be logged in securely.*
3. Start chatting! You can open an **Incognito Window** or a secondary browser to test chatting with yourself as another user.
4. Click the **📎 (Attachment)** button next to the chat bar to upload images to the conversation.

---

## Folder Structure 📁

```text
├── config/              # Database connection logic
├── internal/
│   ├── chat/            # WebSocket Client, Hub, and Message structs
│   │   ├── auth/        # Login/Register Handlers and JWT generation
│   │   ├── models/      # Database models (User)
│   │   └── upload/      # Multipart form file uploader handler
├── public/              # Static Frontend Assets
│   ├── css/style.css    # UI Styling
│   ├── js/app.js        # Frontend logic (WS, fetch, DOM)
│   ├── uploads/         # Directory where uploaded images are saved
│   └── index.html       # The main chat interface
├── .env.example         # Template for environment variables
├── .gitignore           # Ignored files (includes .env and /uploads/)
├── main.go              # Application Entry Point & HTTP Routes
├── go.mod               # Go Modules dependencies
└── README.md            # You are here
```

---

## License 📜
This project is open-source and free to use.
