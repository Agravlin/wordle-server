# wordle-server

Multiplayer Wordle backend in Go

## Getting Started

You'll need Go 1.25+

```bash
git clone https://github.com/agravlin/wordle-server.git
cd wordle-server
go run ./cmd/server
```

The server listens on `:8080`.

Tests:

```bash
go test ./...
```

## API

### HTTP

| Method | Path          | Body                                  | Returns                 |
|--------|---------------|---------------------------------------|-------------------------|
| POST   | `/api/create` | `{"nickname": "alice"}`               | `{"room_id": "ABC123"}`, or `400` if the body is invalid |
| POST   | `/api/join`   | `{"room_id": "ABC123", "nickname": "bob"}` | `{"status": "ok"}`, `400` if the body is invalid or the nick is empty, `404` if the room doesn't exist, or `409` if the nick is taken |

The player who creates the room is its host.

### WebSocket

Connect to `ws://localhost:8080/ws?room=ABC123&nick=alice`. The handshake fails with `404` if the room doesn't exist, or `409` if the nick is taken.

**Client → server**

```jsonc
{ "action": "START_GAME" }                          // host only
{ "action": "SYNC_ROW", "row_state": [1,1,1,0,0] }  // what you're currently typing
{ "action": "GUESS", "guess": "crane" }
```

`row_state` must have exactly 5 cells, each `0` (empty) or `1` (typed). Anything else is ignored.

**Server → client**

Every message looks like `{ "type": "...", "payload": ... }`.

| `type`           | When                                              | `payload`                               |
|------------------|---------------------------------------------------|-----------------------------------------|
| `PLAYER_LIST`    | someone joins or leaves                           | `["alice", "bob"]`                      |
| `FULL_STATE`     | game starts or ends                               | `{ "alice": <board>, "bob": <board> }`  |
| `BOARD_UPDATE`   | a player typed or submitted a row                 | `{ "nick": "alice", "board": <board> }` |
| `GUESS_REJECTED` | your guess was rejected (sent only to you)        | `{ "message": "Not a word" }`           |

A guess is rejected if it isn't a word, has the wrong length, contains non-ASCII characters, or the game is already over.

A `<board>` looks like this:

```json
{ "grid": [[2,3,2,4,2], [4,4,4,4,4]], "CurrentRow": 2 }
```

Board cells are numbers: `0` empty, `1` typed, `2` gray, `3` yellow, `4` green.

## Contributing

The core game works, but there are still some improvements I want to make. If you have ideas or want to help out, feel free to open an issue or PR :)
