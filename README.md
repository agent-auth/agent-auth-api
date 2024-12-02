## Dev Setup

### Getting the Keycloak API Keys for backend/API

```
curl -X POST 'http://localhost:8081/realms/master/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data 'client_id=admin-cli&grant_type=password&username=user&password=uoZfRaacg2'
```

### Getting the MongoDB running

```
docker pull mongodb/mongodb-community-server:latest
docker run --name mongodb -p 27017:27017 -d mongodb/mongodb-community-server:latest
```

### Saas UI/API Documentation

Generateting the API keys to access the API

```
curl -X POST 'http://localhost:8081/realms/saas-ui-api-users/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data 'client_id=saas-ui-api-users&grant_type=password&username=gufranmirza1@gmail.com&password=admin'
```

Making a request to the API

```
curl -X POST 'http://localhost:8002/workspaces' \
-H 'accept: application/json' \
-H 'Authorization: Bearer <access_token>'
```
