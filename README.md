# KahootQuizClone

KahootQuizClone is a real-time quiz game platform inspired by Kahoot. It allows users to create, host, and participate in interactive quizzes with multiple-choice questions. The project is built using a modern tech stack, featuring a Svelte frontend and a Go backend.

## Features

- Create and edit custom quizzes
- Host live quiz sessions
- Join quiz games using a unique game code
- Real-time gameplay with instant feedback
- Leaderboard to track player scores
- Responsive design for both desktop and mobile devices

## Tech Stack

### Frontend
- Svelte
- TypeScript
- Tailwind CSS
- SPA Router

### Backend
- Go
- Fiber (Web framework)
- WebSocket for real-time communication
- MongoDB for data storage

## Project Structure

The project is divided into two main parts:

1. Frontend (`/frontend`)
2. Backend (`/backend`)

### Frontend Structure

- `/src`: Main source code
  - `/views`: Different views/pages of the application
  - `/lib`: Reusable components
  - `/service`: Services for API calls and game logic
  - `/model`: TypeScript interfaces and types

### Backend Structure

- `/cmd`: Entry point for the application
- `/internal`: Core application logic
  - `/controller`: HTTP and WebSocket handlers
  - `/service`: Business logic and game management
  - `/entity`: Data models
  - `/collection`: Database operations

## Getting Started

### Prerequisites

- Node.js and npm
- Go 1.23 or later
- MongoDB

### Running the Frontend

1. Navigate to the `frontend` directory
2. Install dependencies:
   ```
   npm install
   ```
3. Start the development server:
   ```
   npm run dev
   ```

### Running the Backend

1. Navigate to the `backend` directory
2. Install Go dependencies:
   ```
   go mod download
   ```
3. Run the server:
   ```
   go run cmd/quiz/quiz.go
   ```

## API Endpoints

- `GET /api/quizzes`: Fetch all quizzes
- `GET /api/quizzes/:quizId`: Fetch a specific quiz
- `PUT /api/quizzes/:quizId`: Update a quiz
- `GET /ws`: WebSocket endpoint for real-time game communication

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).