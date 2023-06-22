### About

This is an instant messaging server, written in go, and meant to be intergrated as a microservice in a bigger application.

### Stack:
- Go 1.20.3
- gorilla/mux
- gorilla/websocket

### Endoints:
- Websocket connexion: ``GET /ws``
- Send a message: ``POST /users/{id}/friends/{friendId}/message``

### Files
- **main.go :** Inits environment variables, db connexion, routing and starts server.
- **handlers/*:** Main functions, called from router.
- **data/sqlConn.go :** Db connexion.
- **data/messaging.go :** Db call functions.
- **bin/*:** Compiled binaries.

### Operation
1. When a user opens a front-end app, he would be connected to a websocket ``GET /ws``, the service then checks the validity of the authentication token sent, by making a request to the main backend server.

2) The websocket is added to a **list** containing all **active connections**.

3. When a user sends a message to a friend, it is written to the database (if there is valid authentication).

4) The program then checks whether the recipient currently has an **active connection**, in order to send the message in real time.

5. If a connection linked to the recipient's ID is not found, a push notification will be sent.
