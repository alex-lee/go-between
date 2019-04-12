# go-between

`go-between` is a simple UDP proxy for testing purposes, forwarding
traffic sent to a frontend address to some backend address.

This was originally written to watch DNS traffic and simulate failures.

Sessions are tracked from frontend to backend, so that return traffic
can be sent to the proper address.

## Usage

Run the proxy:
```sh
$ make cmd
$ ./go-between -front=127.0.0.1:8053 -back=8.8.8.8:53
```

Send some traffic:
```
$ dig +short @127.0.0.1 -p 8053 www.google.com
```

## Todo

Sessions shouldn't live forever. After a certain amount of idle time
they should be removed.
