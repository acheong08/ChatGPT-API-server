# ChatGPT API Server
[![Release Go Binaries](https://github.com/ChatGPT-Hackers/ChatGPT-API-server/actions/workflows/release.yml/badge.svg)](https://github.com/ChatGPT-Hackers/ChatGPT-API-server/actions/workflows/release.yml)
# Quickstart

## Setup

1. Install Go
2. `go install github.com/ChatGPT-Hackers/ChatGPT-API-server@latest`

If the latest commit fails, try using one of the release binaries

# Build

1. `git clone https://github.com/ChatGPT-Hackers/ChatGPT-API-server/`
2. `cd ChatGPT-API-server`
3. `go install .`

# Usage

`ChatGPT-API-server <port> <API_KEY>`

The admin key can be anything you want. It's just for authenticating yourself.

# Connect agents

Take note of your IP address or domain name. This could be `localhost` or a remote IP address. The default port is `8080`

Check out our [firefox agent](https://github.com/ChatGPT-Hackers/ChatGPT-API-agent). More versions in the works.

There is also a [Python based client](https://github.com/ahmetkca/chatgpt-unofficial-api-docker/tree/ChatGPT-API-agent) by @ahmetkca (WIP)

# Usage

## Quickstart

(After connecting agents)

```bash
 $ curl "http://localhost:8080/api/ask" -X POST --header 'Authorization: <API_KEY>' -d '{"content": "Hello world", "conversation_id": "<optional>", "parent_id": "<optional>"}'
```
Note: if you want to use `conversation_id`, you also need to use `parent_id`!

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

#### /client/register (GET)

N/A. Used for websocket

#### /api/ask (POST)

Headers: `Authorization: <USER_TOKEN>`

_The user token can be set by the admin via /admin/users/add. You can also use the api key as the token. Both work by default_

Data:

```json
{
  "content": "Hello world",
  "conversation_id": "<optional>",
  "parent_id": "<optional>"
}
```

Do not enter conversation or parent id if not available.
If you want to use either of these, you need to specify both! i.e. `request.parent_id=response.response_id` and `request.conversation_id=response.conversation_id`

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

#### /api/connections (GET)

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

#### /admin/users/add (POST)

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

#### /admin/users/delete (POST)

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

#### /admin/users (GET)

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

open `docker-compose.yml` and add your own custom api-key in `<API_KEY>` section

```yaml
version: "3"

services:
  chatgpt-api-server:
    build: .
    ports:
      - "8080:8080"
    command: ["ChatGPT-API-server", "8080", "<API_KEY>", "-listen", "0.0.0.0"]
```

then run:

`docker-compose up` or `docker-compose up -d` (if you want a persistent instance)
