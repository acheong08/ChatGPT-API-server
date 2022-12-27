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

`ChatGPT-API-server <port> <ADMIN_KEY>`

The admin key can be anything you want. It's just for authenticating yourself.

# Connect agents

Take note of your IP address or domain name. This could be `localhost` or a remote IP address. The default port is `8080`

Check out our [firefox agent](https://github.com/ChatGPT-Hackers/ChatGPT-API-agent). More versions in the works.

# Usage

## Quickstart

(After connecting agents)

```bash
 $ curl "http://localhost:8080/api/ask" -X POST --header 'Authorization: <SECRET_KEY>' -d '{"content": "Hello world", "conversation_id": "<optional>", "parent_id": "<optional>"}'
```

## Routes

```go
	router.GET("/client/register", handlers.Client_register) // Used by agent
	router.POST("/api/ask", handlers.API_ask) // For making ChatGPT requests
	router.GET("/api/connections", handlers.API_getConnections) // For debugging
	router.POST("/admin/users/add", handlers.Admin_userAdd) // Adds an API token
	router.POST("/admin/users/delete", handlers.Admin_userDel) // Invalidates a token (based on user_id)
	router.GET("/admin/users", handlers.Admin_usersGet) // Get all users
```

### Parameters for each route

#### /client/register

N/A. Used for websocket

#### /api/ask

Headers: `Authorization: <USER_TOKEN>`
Data:

```json
{
  "content": "Hello world",
  "conversation_id": "<optional>",
  "parent_id": "<optional>"
}
```

Do not enter conversation or parent id if not available.

Response:

```json
{
  "id": "",
  "response_id": "<UUID>",
  "conversation_id": "<UUID>",
  "content": "<string>",
  "error": ""
}
```

#### /api/connections

Headers: None

Data: None

Response:

```json
{
  "connections": [
    {
      "Ws": {},
      "Id": "<UUID>",
      "Heartbeat": "<Time string>",
      "LastMessageTime": "<Time string>"
    }
  ]
}
```

#### /admin/users/add

Headers: None

Data:

```json
{
  "admin_key": "<string>"
}
```

Response:

```json
{
  "user_id": "<UUID>",
  "token": "<UUID>"
}
```

#### /admin/users/delete

Headers: None

Data:

```json
{
  "admin_key": "<string>",
  "user_id": "<UUID>"
}
```

Response:

```json
{ "message": "User deleted" }
```

#### /admin/users

Parameters: `?admin_key=<string>`

Example usage: `curl "http://localhost:8080/admin/users?admin_key=some_random_key"`

Response:

```json
{
  "users": [
    {
      "user_id": "<UUID>",
      "token": "<UUID>"
    },
    {
      "user_id": "<UUID>",
      "token": "<UUID>"
    },
    {
      "user_id": "<UUID>",
      "token": "<UUID>"
    },
    {
      "user_id": "<UUID>",
      "token": "<UUID>"
    },
    ...
  ]
}
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
