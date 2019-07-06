# GO RabbitMQ + MongoDB

Simplicity project to demo how combination messaging system with database.

### About

Create a consumer read message from queue, then use `x-request-id` for querying in remote service, if the result return as `undeliveryCount` = 0 save data to Mongodb, if not move message to `delayQueue` with TTL `5000ms` and try it again.

## Run

```
$ docker-compose up --build

$ go run main.go
```

## Publish Message
