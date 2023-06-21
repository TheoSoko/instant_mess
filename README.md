## Zemus Messaging 

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

### Fichiers / dossiers
- **main.go :** Initialisation des variables d'environnement, de la connexion à la bdd, routing et lancement du serveur.
- **handlers/*:** Fonctions principales appellées depuis le router.
- **data/sqlConn.go :** Connexion à la bdd.
- **data/messaging.go :** Fonctions d'appels à la bdd.
- **env/.env :** Variables d'environnement.
- **go.mod :** Gère le nom du module, la version de go, les dépendances.
- **go.sum :** Liste toutes les dépendances.
- **bin/*:** Fichiers binaires compilés.


### Infos pratiques

Le service tourne sur le port **4000** d'une "instance" (un serveur virtuel) OVH cloud ; ```162.19.92.192```, avec Debian comme OS.
<br/>

Le service est relancé à chaque fois que le serveur l'est, grace à **crontab** et **nohup**.

Le fichier de logs, nommé "logs.out", est situé à la racine du projet. 

**Pour compiler le projet :**

```go build -o bin/server.exe```

**Pour lancer le serveur :** 

```./bin/service.exe``` (Windows)

```./bin/service``` (Linux)



### Fonctionnement
1. Quand un utilisateur ouvre l'application, il est connecté à une websocket ```GET /ws```, le service vérifie ensuite la validité du token d'authentification envoyé, en faisant une requête à l'API principale. 

2) La websocket est ajoutée à une **liste** contenant toutes les **connexions actives**.

3. Quand un utilisateur envoie un message à un ami, celui-ci est écrit dans la base de données (à condition d'une authentification valide).

4) Puis le programme cherche si le destinataire a actuellement une **connexion active** afin d'y envoyer le message en temps réel.

5. Si une connexion liée à l'identifiant du destinaire n'est pas trouvée, une   "notification push" sera envoyée. (Pas encore mise en place)

