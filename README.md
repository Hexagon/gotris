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

## Start using local Redis instance

```bash
GOTRIS_ASSETS=./assets go run main.go
```

## Start using remote Redis instance
```bash
GOTRIS_REDIS_ADDR=172.18.0.12:6379 GOTRIS_ASSETS=./assets go run main.go
```

Game will now be available at

http://127.0.0.1:8080


## Available environment variables

GOTRIS_REDIS_ADDR

GOTRIS_REDIS_PASS

GOTRIS_PORT

GOTRIS_ASSETS
