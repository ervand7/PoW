# Word of Wisdom TCP Server

Welcome to the **Word of Wisdom** TCP server project! This server provides inspirational quotes to clients but requires them to solve a Proof of Work (PoW) challenge first. This mechanism helps protect the server from DDoS attacks by ensuring that each client expends computational effort before receiving a quote.

## Features

- **Proof of Work Protection**: Mitigates DDoS attacks using a hash-based PoW challenge-response protocol.
- **Random Quotes**: Delivers a random quote upon successful PoW verification.
- **Dockerized Setup**: Easily deployable using Docker with multi-stage builds and Docker Compose.

## Quick Start

### Prerequisites

- **Docker** and **Docker Compose** installed on your machine.

### Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/ervand7/PoW.git
   cd word-of-wisdom
   ```

2. **Configure Environment Variables**

   Create a `.env` file in both the `server` and `client` directories.

   **For `server/.env`:**

   ```env
   PORT=:8000
   DIFFICULTY=4
   QUOTES="To inspire love, is to merit love.;Habit is given to us from above: a substitute for happiness.;I am not in the habit of flattering myself with hopes that are not destined to be realized.;The less we show our love to a woman, the more we draw her to us.;It’s better to have loved and lost than never to have loved at all."
   ```

   - **PORT**: The port on which the server listens.
   - **DIFFICULTY**: The number of leading zeros required in the hash (adjusts the PoW difficulty).
   - **QUOTES**: A semicolon-separated list of quotes.

   **For `client/.env`:**

   ```env
   SERVER_ADDRESS=server:8000
   ```

   - **SERVER_ADDRESS**: The address and port of the server.

3. **Build and Run with Docker Compose**

   From the root directory:

   ```bash
   docker-compose up --build
   ```

   This command builds the Docker images and starts both the server and client containers.

## Proof of Work Algorithm Choice

### Overview

The server uses a hash-based PoW algorithm where the client must find a nonce that, when combined with a challenge string and hashed using SHA-256, results in a hash with a specified number of leading zeros.

### Reasons for This Choice

- **Simplicity**: Easy to implement and understand, facilitating maintenance and auditing.
- **Adjustable Difficulty**: The difficulty level can be tuned by changing the number of leading zeros required, allowing dynamic response to server load.
- **Stateless Verification**: The server does not need to store state between client requests, enhancing scalability and reducing resource consumption.
- **Security**: Hash functions like SHA-256 are cryptographically secure, making the PoW robust against attacks.

### Why Hash-Based PoW?

Hash-based PoW algorithms are well-established in blockchain technologies like Bitcoin due to their effectiveness in requiring computational effort (work) that is verifiable with minimal overhead. This makes them ideal for preventing abuse such as DDoS attacks.
