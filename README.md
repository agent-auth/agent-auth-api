# recruiter-api

# Go Restful API Boilerplate

Easily extendible RESTful API boilerplate aiming to follow idiomatic go and best practice.

The goal of this boiler is to have a solid and structured foundation to build upon on.

Any feedback and pull requests are welcome and highly appreciated. Feel free to open issues just for comments and discussions.

## Features

The following feature set is a minimal selection of typical Web API requirements:

- Configuration using [viper](https://github.com/spf13/viper)
- CLI features using [cobra](https://github.com/spf13/cobra)
- PostgreSQL support including migrations using [go-pg](https://github.com/go-pg/pg)
- Structured logging with [Logrus](https://github.com/sirupsen/logrus)
- Routing with [chi router](https://github.com/go-chi/chi) and middleware
- JWT Authentication using [jwt-go](https://github.com/dgrijalva/jwt-go) with example passwordless email authentication
- Request data validation using [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)
- HTML emails with [gomail](https://github.com/go-gomail/gomail)

## Start Application

- Clone this repository
- Create a postgres database and set environment variables for your database accordingly if not using same as default
- Run the application to see available commands: `go run main.go`
- First initialize the database running all migrations found in ./database/migrate at once with command _migrate_: `go run main.go migrate`
- Run the application with command _serve_: `go run main.go serve`
- Go to http://127.0.0.1:8001/recruiter-api/v1/swagger to view the swagger API docs

## API Routes

chechout `src/web/docs` folder for swagger API documentation

### Testing

Package auth/pwdless contains example api tests using a mocked database.

---

![Screenshot from 2021-11-16 23-24-49](https://user-images.githubusercontent.com/17959487/142039740-5f5a6b5d-5210-403b-9e9f-54ea18f420bd.png)

```
curl -X POST 'http://localhost:8081/realms/master/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data 'client_id=admin-cli&grant_type=password&username=user&password=uoZfRaacg2'
```

```
curl -X POST 'http://localhost:8081/realms/saas-ui-api-users/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data 'client_id=saas-ui-api-users&grant_type=password&username=gufranmirza1@gmail.com&password=admin'
```

eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ1Q2R0bE52YnFxbE5TWDg3VUNBQ2M4NlBFcDltVjRZNm44Tzh6QUstcmswIn0.

curl -X 'GET' \
 'http://localhost:8002/master/0d37d66f-dffd-4c87-9301-28b49abc9c7a/oauth/resources' \
 -H 'accept: application/json' \
 -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ1Q2R0bE52YnFxbE5TWDg3VUNBQ2M4NlBFcDltVjRZNm44Tzh6QUstcmswIn0.eyJleHAiOjE3MzMwOTYyNzIsImlhdCI6MTczMzA5NTk3MiwianRpIjoiYzMwOWU0ZTktMzRmOS00N2QzLWE1ZmYtNjAyNjZjY2RmYWUxIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgxL3JlYWxtcy9zYWFzLXVpLWFwaS11c2VycyIsInN1YiI6ImE5NjYzZmUwLTIxYjQtNDhlNS1iODU5LTlkNzZkMTliNGIwNSIsInR5cCI6IkJlYXJlciIsImF6cCI6InNhYXMtdWktYXBpLXVzZXJzIiwic2lkIjoiYzI0MzIzMjUtZmFkMC00NDA3LTk1NTUtOWE4YWI1ZTZiYTg3IiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyIvKiJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsiYXBwX2RldmVsb3BlciIsImFwcF92aWV3ZXIiLCJhcHBfYWRtaW4iXX0sInNjb3BlIjoiZW1haWwgcHJvZmlsZSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwicm9sZXMiOlsiYXBwX2RldmVsb3BlciIsImFwcF92aWV3ZXIiLCJhcHBfYWRtaW4iXSwibmFtZSI6IkdVRlJBTiBCQUlHIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiZ3VmcmFubWlyemExQGdtYWlsLmNvbSIsImdpdmVuX25hbWUiOiJHVUZSQU4iLCJmYW1pbHlfbmFtZSI6IkJBSUciLCJlbWFpbCI6Imd1ZnJhbm1pcnphMUBnbWFpbC5jb20ifQ.cX7HTWW50OtVdwub_spBtU0wI7PUAwrQkU-lmWmy2g-vxd190KBy2lsOEP_lX03_PspqmAL4JqOJ3TQVQE184j_AqLd2Mp9IMbqK4iTLEOLzIqzSquuEjAVRIXSxZjIsB8lSp7n1IsgJgoTfdYjuh24tuVZdbX0oqKE7oWSOMPOsuJM5FCRVRRIL4k57lCx0SIlmEVyAAoszrwjwkbe2N6MDbIRyqy55vlWPzzPX12q3kSLegc2L_HgW_s64L8JqbzsA-lXa0j9qHxsAVD7zNLS46zFECA2fNIc5oeg2tlxAZlU2Y65oQXbUYiHAdTIFoCsW7xdG_ya70b2rUGx_mA'

docker pull mongodb/mongodb-community-server:latest
docker run --name mongodb -p 27017:27017 -d mongodb/mongodb-community-server:latest
