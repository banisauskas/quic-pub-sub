## QUIC Publisher/Subscriber

Server-client communication using **QUIC** protocol (instead of TCP or UDP) written in **Go**. Server forwards messages from publishers to subscribers.

### Server

At the beginning, run 1 instance:

```shell
cd server
go run .
```

At the end, terminate pressing `Ctrl+C`.

### Publisher Client

Run any number of instances:

```shell
cd client-publish
go run .
```

Terminate pressing `Ctrl+C`.

### Subscriber Client

Run any number of instances:

```shell
cd client-subscribe
go run .
```

Terminate pressing `Ctrl+C`.

## References

* https://en.wikipedia.org/wiki/QUIC
* https://github.com/quic-go/quic-go