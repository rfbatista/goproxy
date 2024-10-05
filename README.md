# goproxy

# Inicializar o servidor

```sh
go run ./cmd/server https://www.youtube.com/
```

O argumento esperado é a url do backend

# Bloqueie algum ip

`go run ./cmd/cli block <ip>`

Exemplo:
```sh
go run ./cmd/cli block ::1
```

# Remover ip da lista de bloqueio

`go run ./cmd/cli remove <ip>`

Exemplo:
```sh
go run ./cmd/cli remove ::1
```

# Agora você pode testar o proxy com curl

```sh
curl --location 'http://localhost:8080/?col'
```
