# Cart Service
Part of [eshop](https://github.com/idoyudha/eshop) Microservices Architecture.

## Overview
This service handles user cart create, read, update, and delete (CRUD). Using redis as main database for user cart for high performance and low latency requirements, backed with MySQL if cache is missed ([Cache-Aside Pattern](https://learn.microsoft.com/en-us/azure/architecture/patterns/cache-aside)).

## Architecture
```
eshop-auth
├── .github/
│   └── workflows/  # github workflows to automatically test, build, and push
├── cmd/
│   └── app/        # configuration and log initialization
├── config/         # configuration
├── internal/   
│   ├── app/        # one run function in the `app.go`
│   ├── controller/ # serve handler layer
│   │   ├── http/
│   │   |   └── v1/ # rest http
│   │   └── kafka
│   │       └── v1/ # kafka subscriber
│   ├── entity/     # entities of business logic (models) can be used in any layer
│   ├── usecase/    # business logic
│   │   └── repo/   # abstract stirage (database) that business logic works with
│   └── utils/      # helpers function
├── migrations/     # sql migration
└── pkg/
    ├── httpserver/ # http server initialization
    ├── kafka/      # kafka initialization
    ├── logger/     # logger initialization
    ├── mysql/      # mysql initialization
    └── redis/      # redis initialization
```

## Tech Stack
- Backend: Go
- Authorization: AWS Cognito
- Database: Redis and MySQL
- CI/CD: Github Actions
- Message Broker: Apache Kafka
- Container: Docker

## API Documentation
tbd