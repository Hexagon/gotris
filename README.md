# gotris

Source code of http://tetris.56k.guru. 

Server side Tetris implementation written in go with lazy HTML5 client.


# Development setup

Make sure you have a running local instance of Redis, and that Go is installed

```bash
cd <your workspace path>
mkdir -p src/github.com/hexagon
cd src/github.com/hexagon
git clone https://github.com/Hexagon/gotris.git
cd gotris
go get -v
```

## Start with local MongoDB instance

```bash
GOTRIS_ASSETS=./assets go run main.go
```

## Start with remote MongoDB instance
```bash
GOTRIS_MONGO_ADDR=172.18.0.16 GOTRIS_ASSETS=./assets go run main.go
```

Game will now be available at

http://127.0.0.1:8080


## Available environment variables

Variable | Default
--- | ---
GOTRIS_MONGO_ADDR | 127.0.0.1:27017
GOTRIS_MONGO_USER | -
GOTRIS_MONGO_PASS | -
GOTRIS_MONGO_DB | gotris
GOTRIS_PORT | 8080
GOTRIS_ASSETS | -
