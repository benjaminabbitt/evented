#Binaries for executing integration tests
###Not useful in a real application, test only.

To Run:

### Start Mongo
`docker run mongo -p 27017:27017`

### Start Rabbit
`docker run -p 4369:4369 -p 5671:5671 -p 5672:5672 -p 25672:25672 rabbitmq`

## Domain A
###Command Handler
####Configuration:

```
---
port: 8080
domain: "domainA"
business:
    address: "localhost:8081"
eventStore:
    type: mongodb
    mongodb:
        url: mongodb://localhost:27017
        database: event
snapshotStore:
    type: mongodb
    mongodb:
        url: mongodb://localhost:27017
        database: snapshot
transport:
    type: amqp
    amqp:
        url: amqp://guest:guest@localhost:5672/
        exchange: evented

```

####Executable
``
