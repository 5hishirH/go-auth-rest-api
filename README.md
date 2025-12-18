project structure

```
project/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── user/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── model.go
│   │   └── routes.go
│   ├── shared/
│   │   ├── db/
│   │   ├── middleware/
│   │   └── response/
├── pkg/
├── go.mod
```

init project
```sh
go mod init github.com/5hishirH/go-auth-rest-api.git
```

run the project
```sh
go run cmd/api/main.go
```