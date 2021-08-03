# GoMiner

## Intro
This is a Web Server that acts as a [Stratum](https://braiins.com/stratum-v1/docs#developers) server. Currently it only supports two commands:

- [mining.authorize](https://en.bitcoin.it/wiki/Stratum_mining_protocol#mining.authorize): as long as at least the username (first param) is provided, it will always return true
- [mining.subscribe](https://en.bitcoin.it/wiki/Stratum_mining_protocol#mining.subscribe): it will create a new subscription or resume an existing one depending on the provided params. 

## Instructions
The following instructions are useful to Build, Test and Run the server.

## Building and Running
### Generating mocks
Before building the project, the `mocks` need to be created, so that the Unit Tests can run successfully. In order to create them:
```
moq -out controller/mock_service.go -pkg controller ./service Service
```
Or even simpler:
```
make gen
```

### Build
After the mocks are created, the project can be easily built with:
```
go build
```

Or even simpler:
```
make build
```
This last command generates the mocks before building the code. Coverage report can also be created with the command:
```
make cover
```


### Run
#### Environment Variables
The following environment variables need to be present or must be provided in a `.env` file in the same directory as the executable:
```
HTTP_PORT=
POSTGRES_HOST=
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
POSTGRES_PORT=
POSTGRES_USERS_TABLE_SCHEMA=
POSTGRES_USERS_TABLE_NAME=
POSTGRES_SUBSCRIPTIONS_TABLE_SCHEMA=
POSTGRES_SUBSCRIPTIONS_TABLE_NAME=
```

#### Database
This server uses a PostgreSQL DB. A `docker-compose.yaml` is included in order to spin it up. In order to do it:
```
docker-compose up
```

#### Execution
Many different ways to do it:
```
go run main.go
```
Or directly:
```
./stratum-server
```

Or even simpler
```
make run
```

## Technologies
### Go-Chi
[Go-Chi](https://github.com/go-chi/chi) has been used as the HTTP Router. It's lightweight, idiomatic and composable, therefore no further dependencies are added.

### Gorilla Mux
[Gorrila Mux](https://github.com/gorilla/mux) has been used for the WebSocket development.

### Moq
[Moq](https://github.com/matryer/moq) has been used in order to create mocks automatically.

## Architecture
This basic WebServer has been divided in:
- **config**: contains all the logic to retrieve environment variables
- **controller**: contains all APIs, router, decoding and encoding.
- **repository**: contains the interface to perform Insert/Update/Query operations on the PostgreSQL DB.
- **service**: contains all the specific business logic, including the websocket logic and the orchestration of the different pieces.

### Assumptions
There are a few things that are not 100% clear about the protocol. Therefore, I'll list all the assumptions I've made and each one of them could be easily modified if it's required:
- **Subscription IDs**: according to the documentation, the subscription IDs should be unique. Therefore, I used a UUID generator for it, since it didn't specify that the Subscription IDs need to be unique across all miners.
- **Subscribing with ExtraNonce1 as param**: according to the documentation, if the optional parameter `ExtraNonce1` is provided, a previous existing subscription should be resumed. In order to do it, I'm saving the status from the subscription and marking it as "active" when it's created and "inactive" when connection is lost. When subscribing, the `ExtraNonce1` must corresspond to a valid and "inactive" existing subscription. I'm assuming that it's forbidden to have multiple connections subscribing for the same `ExtraNonce1` simultaneously.
- **Only one `[mining.subscribe]` per connection**: I couldn't find in the documentation is this is indeed a condition, but it seems that the `[mining.subscribed]` method is performed at the beginning of the process.

### Improvements
There are many things that could be improved in the overall solution with the proper time:
- **Coverage**: improve coverage on every module and raise it to the maximum. I've included a few UTs to show how to structure them and how to use Mocks to test the modules independently.
- **websocket module**: it'd be great to move all the specific logic from the websocket into a separate module.
- **[mining.notify]**: it'd be great to add this piece of logic with more time. I'd essentially store the WebSocket connection in a pool, and with a Job I'd:
  - Periodically iterate through the connections and call `WebSocket.WriteMsg()` function.
  - When the WebSocket is shutdown, report it to the pool using a channel so that it can be deleted from the pool.

## CI
The project is configured in Gitlab with CI. The code is built and tested every time a new commit is pushed.

## Examples
A few examples are provided, including both commands and different kind of errors. [websocat](https://github.com/vi/websocat) has been used to test the code

### New Subscription specifying subscriber
```
▶ websocat ws://127.0.01:8080/api/v1/ws
{"id":1,"method":"mining.subscribe","params":["cgminer/4.10.0"]}
{"id":1,"result":[[["mining.set_difficulty","a00e3334-5b8e-41ba-9fac-e0a26b1fd000"],["mining.notify","828b75d3-bcce-4f8c-a40a-cea154fb880c"]],"00000011",4]}
{"id":2,"method":"mining.authorize","params":["user","pass"]}
{"id":2,"result":true}
```

### New Subscription without specifying subscriber
```
▶ websocat ws://127.0.01:8080/api/v1/ws
{"id":1,"method":"mining.subscribe"}
{"id":1,"result":[[["mining.set_difficulty","ac318a35-6093-4200-9860-dfc8e5a0acf7"],["mining.notify","d121b4c7-bf3b-4705-8f89-38df8e20c90d"]],"00000012",4]}
{"id":2,"method":"mining.authorize","params":["user","pass"]}
{"id":2,"result":true}
```

### Resuming Previous Subscription
```
▶ websocat ws://127.0.01:8080/api/v1/ws
{"id":1,"method":"mining.subscribe","params":["cgminer/4.10.0","0000000f"]}
{"id":1,"result":[[["mining.set_difficulty","aa326c83-273d-4efb-a0d2-cb4168e763ba"],["mining.notify","07cad3ce-ea62-4a15-b1b5-a3f6ce48e965"]],"0000000f",4]}
{"id":2,"method":"mining.authorize","params":["user","pass"]}
{"id":2,"result":true}
```

### Errors
#### Error with Invalid request
```
▶ websocat ws://127.0.01:8080/api/v1/ws
{}
{"error":{"code":-32600,"message":"Invalid Request"}}
```

#### Error with Invalid or unsupported method
```
▶ websocat ws://127.0.01:8080/api/v1/ws
{"id":1,"method":"mining.something"}
{"id":1,"error":{"code":-32601,"message":"Method not found"}}
```

#### Error with Invalid params
```
▶ websocat ws://127.0.01:8080/api/v1/ws
{"id":1,"method":"mining.authorize"}
{"error":{"code":-32602,"message":"Invalid params"}}
```

#### Error with Internal Error
```
▶ websocat ws://127.0.01:8080/api/v1/ws
{"error":{"code":-32603,"message":"Internal error"}}
```

#### Error with Invalid JSON-RPC format
```
▶ websocat ws://127.0.01:8080/api/v1/ws
[]
{"error":{"code":-32700,"message":"Parse error"}}
```
