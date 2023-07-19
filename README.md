# Go Temporal Example

## Description
- The following is a description of a Go source code that demonstrates the integration of Temporal.io, 
a powerful workflow automation platform. 
Temporal.io simplifies the development of complex distributed systems by providing a programming model 
that enables developers to write fault-tolerant and scalable workflows using Go and other programming languages.
- This Go source-code showcases how to integrate Temporal.io into a Go application, 
enabling developers to build resilient, fault-tolerant, and scalable distributed systems. 
The combination of Temporal.io's workflow automation capabilities with the versatility of 
Go makes it a powerful solution for complex business processes and mission-critical applications.

## Key features
- Dependencies
- Workflow Definition
- Activity Implementations
- Workflow Registration
- Temporal.io Client Setup
- Main Function
- Error Handling and Logging

## Prerequisites

## How to run local

## How to run ci/cd

## Technologies
- Languages: Go
- SDK: temporal.io
- Network: gRPC, HTTP
- Docker
- Cloud vendor: AWS Secret Manager

## Directory structure
    📁 api // 
    |__ 📁 applications
        |__ 📁 commands
            |__ campaign.go
            |__ record.go
            |__ usage.go
        |__ 📁 queries
            |__ campaign.go
            |__ record.go
            |__ usage.go
    |__ 📁 domains
        |__ 📁 entities
            |__ campaign.go
            |__ campaign_draft.go
            |__ campaign_integrate.go
            |__ campaign_meta_field.go
            |__ ...
        |__ 📁 services
        |__ 📁 repositories
            |__ campaign_repository.go
            |__ record_repository.go
        |__ 📁 values
    |__ 📁 infrastructures
        |__ 📁 persistence
        |__ 📁 secrets
    |__ 📁 interfaces
        |__ 📁 middleware
        |__ campaign_handle.go
        |__ configure_handle.go
        |__ record_handle.go
        |__ usage_handle.go
    |__ 📁 cmd
        |__ main.go
    |__ 📁 configs
        |__ .env
        |__ app.go
    |__ 📁 docs
    |__ docker-compose.yml
    |__ go.mod
    |__ go.sum
    |__ README.md


## References
- https://dev.to/stevensunflash/using-domain-driven-design-ddd-in-golang-3ee5
- https://github.com/victorsteven/food-app-server/tree/master/domain
- DDD vs CQRS: https://github.com/authena-ru/courses-organization/tree/main
- Go GRPC middleware: https://github.com/grpc-ecosystem/go-grpc-middleware/tree/main/interceptors