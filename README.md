### Rinha de Backend 2024 Q1

[![Go](https://github.com/MarcusAdriano/rinha-de-backend-2024-q1-go/actions/workflows/go.yml/badge.svg)](https://github.com/MarcusAdriano/rinha-de-backend-2024-q1-go/actions/workflows/go.yml)

Mais detalhes sobre a rinha nesse link --> [rinha-de-backend-2024-q1](https://github.com/zanfranceschi/rinha-de-backend-2024-q1).

### Stack

Lang: Golang 1.21
- Nginx
- Fiber
- Postgres

Requisitos mínimos:
- Docker
- Go v1.21

### Variáveis de ambiente

- `SERVER_PORT` - Porta que a aplicação irá rodar;
- `DATABASE_URL` - Host do banco de dados (postgres);

### Build e teste da aplicação

Será realizado os testes e o build da app.

```bash
make
```

### Rodando a aplicação

```bash
export DATABASE_URL=postgres://postgres:mysecretpassword@localhost:5432/rinhabackend
```

```bash
go run cmd/rinha/main.go
```

Detalhes da API acesse o [swagger](http://localhost:8080/swagger/).

### Rodando um container postgres local

BD em docker com as configurações necessárias para rodar a aplicação.

```bash
docker run --rm --name rinhabackend-db -p 5432:5432 -v ./sql/init:/docker-entrypoint-initdb.d/ -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_DB=rinhabackend -d postgres:16.1
```
