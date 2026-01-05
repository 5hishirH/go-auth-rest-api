project structure

```
project/
├── cmd/
│   └── api/
│       └── main.go
├── config/
│   └── local.yaml
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── storage/
│   │   ├── sqlite.go
│   │   ├── redis.go
│   │   └── minio.go
│   ├── utils/
│   │   └── response/
│   │       └── response.go
│   ├── auth/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   └── routes.go
│   └── server/
│       ├── server.go
│       └── router.go
└── go.mod
```

init project

```sh
go mod init github.com/5hishirH/go-auth-rest-api.git
```

run the project

```sh
go run cmd/api/main.go
```
