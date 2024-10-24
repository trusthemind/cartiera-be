# Cars Sales App Back-end

## This application is created for Cartiera app, make market operations and car trades more comfortable for customers.

### Install all dev dependencies

```bash
go get all
```

---

### How to run server with terminal

run with Compile Daemon

```bash
CompileDaemon -command="./go-cars-app"
```

---

or run with command

```bash
go run main.go
```

## Build and run with Docker

build

```bash
docker build -t cartiera-be-go:local .
```

run application service

```bash
docker run cartiera-be-go:local .
```

## Helpers

- [License dependency checker](licenses.csv)
- [Swagger](https://car-sales-app-v2.up.railway.app/docs/index.html)
