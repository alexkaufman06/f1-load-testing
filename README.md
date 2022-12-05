# F1 Load Testing

Example load using F1 and Docker. 

https://www.youtube.com/watch?v=yKJ_h0K6liM
18:40

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