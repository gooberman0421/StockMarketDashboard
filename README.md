# Real-Time Stock Market Dashboard

## Description
A real-time stock market dashboard that displays live prices, trends, and historical data of various stocks. The project utilizes TypeScript for the frontend, Python for the backend, and Go for real-time price updates and WebSocket communication.

## Requirements

- Node.js
- Python 3.8+
- Go 1.16+
- Flask

## Setup Instructions

### Frontend

1. Navigate to the frontend directory:
    ```bash
    cd frontend
    ```
2. Install dependencies:
    ```bash
    npm install
    ```
3. Build the project:
    ```bash
    npm run build
    ```
4. Start the development server:
    ```bash
    npm start
    ```

### Backend

1. Navigate to the backend directory:
    ```bash
    cd backend
    ```
2. Create a virtual environment:
    ```bash
    python -m venv venv
    ```
3. Activate the virtual environment:
    ```bash
    # On Windows
    venv\Scripts\activate

    # On macOS/Linux
    source venv/bin/activate
    ```
4. Install dependencies:
    ```bash
    pip install -r requirements.txt
    ```
5. Run the Flask application:
    ```bash
    flask run
    ```

### WebSocket Server

1. Navigate to the websocket directory:
    ```bash
    cd websocket
    ```
2. Install Go dependencies:
    ```bash
    go mod tidy
    ```
3. Start the WebSocket server:
    ```bash
    go run main.go
    ```

## Usage

1. **Access the frontend application:**
   Open a web browser and go to `http://localhost:3000`.

2. **Interact with the backend API:**
   The backend API is available at `http://localhost:5000/api`.

3. **Receive real-time stock price updates:**
   The WebSocket server runs on `ws://localhost:8080/ws`.
