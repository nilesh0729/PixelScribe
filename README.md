# PixelScribe 🎙️

[![Go Tests](https://github.com/nilesh0729/PixelScribe/actions/workflows/test.yml/badge.svg)](https://github.com/nilesh0729/PixelScribe/actions/workflows/test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

PixelScribe is a modern, high-performance web application designed to help users improve their typing speed and accuracy through dictation practice. It features real-time audio generation using OpenAI's TTS API, server-side accuracy verification, and comprehensive performance analytics.

## 📋 Table of Contents

- [Screenshots](#-screenshots)
- [Features](#-features)
- [Tech Stack](#️-tech-stack)
- [Installation & Setup](#-installation--setup)
- [Project Structure](#-project-structure)
- [API Documentation](#️-api-documentation)
- [Contributing](#-contributing)
- [License](#-license)

## 📸 Screenshots

![Sign In Page](screenshots/signin.png)
*Sign in to your account*

![Registration Page](screenshots/register.png)
*Create a new account*

![Dashboard](screenshots/dashboard.png)
*Overview of your dictation performance with stats and recent activity*

![New Dictation](screenshots/new-dictation.png)
*Create a new dictation with customizable title, language, and content*

![Attempt History](screenshots/history.png)
*Track your progress with detailed attempt history*

## 🚀 Features

-   **Dictation Engine**:
    -   High-quality Text-to-Speech (TTS) using **OpenAI API**.
    -   Split "Listen then Type" workflow for focused learning.
    -   Real-time audio playback control.
-   **Smart Analysis**:
    -   **Visual Diffing**: Highlights missed, incorrect, and extra words (Green/Red highlighting).
    -   **Server-Side Verification**: Secure and accurate WPM and accuracy calculation.
-   **Performance Tracking**:
    -   Comprehensive Dashboard with charts and recent activity.
    -   Detailed **Attempt History** to track progress over time.
-   **User System**:
    -   Secure Authentication (JWT-based).
    -   Personalized Dictation Library.
-   **Modern Tech Stack**:
    -   Built with **Go (Golang)** for a robust, high-speed backend.
    -   **React + Tailwind CSS** for a responsive, clean UI.
    -   **Dockerized** for easy deployment.

## 🛠️ Tech Stack

### Backend
-   **Language**: Go (Golang)
-   **Framework**: Gin Gonic
-   **Database**: PostgreSQL
-   **ORM/Data Layer**: SQLC (Type-safe SQL generation)
-   **Auth**: JWT (JSON Web Tokens)
-   **TTS**: OpenAI API

### Frontend
-   **Framework**: React (Vite)
-   **Styling**: Tailwind CSS v4
-   **State/Data**: Axios, React Router

### DevOps
-   **Containerization**: Docker & Docker Compose
-   **CI**: GitHub Actions (Go Test Workflow)
-   **Config**: Viper (Environment variables)

## 📦 Installation & Setup

### Prerequisites
-   Docker & Docker Compose
-   (Optional) Go 1.23+ and Node.js 18+ for local dev without Docker.

### Quick Start (Docker)

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/nilesh0729/PixelScribe.git
    cd PixelScribe
    ```

2.  **Configure Environment**:
    Create a `app.env` file in the root directory (or use the example):
    ```env
    DB_SOURCE=postgresql://root:secret@postgres:5432/pixelscribe?sslmode=disable
    SERVER_ADDRESS=0.0.0.0:8080
    TOKEN_SYMMETRIC_KEY=12345678901234567890123456789012
    ACCESS_TOKEN_DURATION=15m
    OPENAI_API_KEY=your_openai_api_key_here
    ```

3.  **Run with Docker Compose**:
    ```bash
    docker compose up --build
    ```
    -   **Frontend**: http://localhost:5173
    -   **Backend**: http://localhost:8080

### Local Development (Manual)

#### Backend
1.  Navigate to `/api` or root.
2.  Install dependencies: `go mod download`.
3.  Run migrations (using `migrate` CLI or via Make).
4.  Start server: `go run main.go`.

#### Frontend
1.  Navigate to `/web`.
2.  Install dependencies: `npm install`.
3.  Start dev server: `npm run dev`.

## 📂 Project Structure

```
PixelScribe/
├── api/                # Go Backend Handlers & Logic
├── Result/             # SQLC Generated Code (DB Layer)
├── db/                 # SQL Migration Files
├── web/                # React Frontend Application
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── services/   # API Clients
│   │   └── ...
├── app.env             # Environment Secrets (GitIgnored)
├── docker-compose.yaml # Docker Orchestration
└── ...
```

## 🛡️ API Documentation

The API is RESTful and communicates via JSON. Key endpoints include:

-   `POST /users/login`: Authenticate user.
-   `POST /tts/generate`: Proxy to OpenAI TTS (Secure).
-   `POST /attempts`: Submit a dictation attempt for grading.
-   `GET /performance`: Fetch user stats.

## 🤝 Contributing

1.  Fork the repo.
2.  Create a feature branch (`git checkout -b feature/amazing-feature`).
3.  Commit your changes (`git commit -m 'Add amazing feature'`).
4.  Push to the branch (`git push origin feature/amazing-feature`).
5.  Open a Pull Request.

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

## 🙏 Acknowledgments

- **OpenAI** - For providing the high-quality TTS API
- **Go Community** - For the excellent tooling and libraries
- **React & Tailwind CSS** - For making frontend development enjoyable

---

<div align="center">
  Made with ❤️ by <a href="https://github.com/nilesh0729">Nilesh Khodiyar</a>
</div>
