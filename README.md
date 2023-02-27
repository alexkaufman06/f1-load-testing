# F1 Load Testing

Example load using F1 and Docker.

Launch docker services:

```
docker-compose up -d
```

View running containers:

```
docker ps
```

Hit local endpoint:

```
curl -X POST http://localhost:8080/payments -v
```

Configure AWS: (region `us-west-2`)

```
aws configure
```

List queues:

```
aws --endpoint-url http://localhost:4100 sqs list-queues
```

Get Queue attributes:

```
aws --endpoint-url http://localhost:4100 sqs get-queue-attributes --queue-url "http://us-west-2.goaws/4100/100010001000/test-queue"
```

F1 help:

```
go run ./cmd/f1/main.go --help
```

F1 list scenarios:

```
go run ./cmd/f1/main.go scenarios ls
```

Compile app into binary:

```
go build -o f1 ./cmd/f1/main.go
```

Run binary in constant mode:
(This mode gives you more control over the rate compared to user mode)

```
./f1 run constant testScenario -r 100/s -d 5s
```

Run binary in user mode:

```
./f1 run users testScenario -c 1 -d 10s
```