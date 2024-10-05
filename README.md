# goproxy

# Inicializar o servidor

```sh
go run ./cmd/server https://www.youtube.com/
```
O argumento esperado é a url do backend

# Bloqueie algum ip

```sh
go run ./cmd/cli block ::1
```

# Remover ip da lista de bloqueio

```sh
go run ./cmd/cli remove ::1
```
