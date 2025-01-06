# Cart Service
Part of [eshop](https://github.com/idoyudha/eshop) Microservices Architecture.

## Overview
This service handles user cart create, read, update, and delete (CRUD). Using redis as main database for user cart for high performance and low latency requirements, backed with MySQL if cache is missed ([Cache-Aside Pattern](https://learn.microsoft.com/en-us/azure/architecture/patterns/cache-aside)).

## Architecture
```
eshop-auth
├── .github/
│   └── workflows/
├── cmd/
│   └── app/
├── config/
├── internal/   
│   ├── app/
│   ├── controller/
│   │   ├── http/
│   │   |   └── v1/
│   │   └── kafka
│   │       └── v1/
│   ├── entity/
│   ├── usecase/
│   │   └── repo/
│   └── utils/
├── migrations/
└── pkg/
    ├── httpserver/
    ├── kafka/
    ├── logger/
    ├── mysql/
    └── redis/
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