# Hydra API
## Go Directories
### `/cmd`
Main application for this project.

### `/internal`
Private code that cannot be imported from outside 

### `/pkg`
Public code that can be imported

### `/vendor`
Application dependencies using `go mod` 
The `go mod vendor` command will create the `/vendor` directory for you. 
The `go mod sync` command will update the `/vendor`. 

## Internal Directory
Inspired by Clean Code architecture (https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

According to Uncle Bob, we can divide our code in 4 layers:

### `/usecase`
Contains application specific business logic

### `/controller`
Adapters to convert data format most convenient to external and usecases

### `/repository`
Interfaces with database

### `/entities`
Data structures to encapsulate enterprise wide business rules

## How to run the project
### Software Prerequisite
1. Install go https://golang.org/doc/install
2. Install Docker https://docs.docker.com/get-docker/

### Commands at `/api`
* This will start up postgres docker container and run database migration
```
make setup
```
* This will start the api
```
make run
```
* This will update your local vendor dependencies and clean up go.mod
```
make update-vendor
```
