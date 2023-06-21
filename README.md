Zemus Messaging 
----------
### Présentation
Zemus messaging est un microservice de messagerie instantanée écrit en Go, et visant à être intégré dans le système Zemus.

### Stack:
- Go 1.20.3
- gorilla/mux
- gorilla/websocket

### Endoints:
- API: ``http://api.zemus.info:4000``
- Connexion websocket: ``GET /ws``
- Envoi message: ``POST /users/{id}/friends/{friendId}/message``
