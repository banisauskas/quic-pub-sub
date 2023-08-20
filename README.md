### Run Server

One instance:

```shell
cd server
go run server.go publisher.go subscriber.go
```

### Run Publisher Client

Any number of instances:

```shell
cd client-publish
go run client.go random.go
```