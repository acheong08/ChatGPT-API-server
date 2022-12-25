# ChatGPT API Server

# Quickstart 
## Setup
1. Install Go
2. `go install github.com/ChatGPT-Hackers/ChatGPT-API-server@latest`

# Build
1. `git clone https://github.com/ChatGPT-Hackers/ChatGPT-API-server/`
2. `cd ChatGPT-API-server`
3. `go install .`

# Usage
`ChatGPT-API-server <port> <API key>`

# Connect agents
Take note of your IP address or domain name. This could be `localhost` or a remote IP address. The default port is `8080`

Check out our [firefox agent](https://github.com/ChatGPT-Hackers/ChatGPT-API-agent). More versions in the works.

# Usage
```bash
 $ curl "http://localhost:8080/api/ask" -X POST --header 'Authorization: <API_KEY>' -d '{"content": "Hello world", "conversation_id": "<optional>", "parent_id": "<optional>"}'
 ```

# Docker

open `docker-compose.yml` and add your own custom api-key in `<api-key>` section

 ```yaml
version: "3"

services:
  chatgpt-api-server:
    build: .
    ports:
      - "8080:8080"
    command: ["ChatGPT-API-server", "8080", "<api-key>"]
 ```

then run:

`docker-compose up` or `docker-compose up -d` (if you want a persistent instance)
