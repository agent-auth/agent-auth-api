## Dev Setup

### Getting the Keycloak API Keys for backend/API

```
curl -X POST 'http://localhost:8080/realms/master/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data 'client_id=admin-cli&grant_type=password&username=admin&password=uoZfRaacg2'
```

### Saas UI/API Documentation

Generateting the API keys to access the API

```
curl -X POST 'http://localhost:8080/realms/saas-ui-api-users/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data 'client_id=saas-ui-api-users&grant_type=password&username=gufranmirza1@gmail.com&password=uoZfRaacg2'
```

Making a request to the API

```
curl -X POST 'http://localhost:8002/workspaces' \
-H 'accept: application/json' \
-H 'Authorization: Bearer <access_token>'
```
