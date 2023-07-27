
# Feed Provider

This project is a feed provider that fetches news articles from a remote source and exposes them through a RESTful API. It provides three endpoints:

- /ping
- /articles
- /articles/{id}

## Running Tests
To run the tests for this project, follow these steps:
- Make sure you have Go installed on your machine.
- Open a terminal and navigate to the project directory.
- Run the following command to execute the tests:
`go test ./...`


This will run all the tests in the project.

  

## Running locally

To run the project locally, follow these steps:  

 1. Make sure you have Go installed on your machine.
 2. Open a terminal and navigate to the project directory.
 3. Run the following command to build the project: `go build -o feed-provider ./main.go`
 4. After the build is successful, run the following command to start the server: `./feed-provider`
 
The server will start running on
http://localhost:8080
  

## Running with Docker Compose

To run the project with Docker Compose, follow these steps:

 1. Make sure you have Docker and Docker Compose installed on your machine.
 2. Open a terminal and navigate to the project directory.
 3. Run the following command to build the Docker image:`docker-compose build`
 4. After the build is successful, run the following command to start the containers:`docker-compose up`

The server will start running on
http://localhost:8080
  

## API Documentation

- `/ping`: GET request to check if the server is running.
- `/articles`: GET request to retrieve all articles.
- `/articles/{id}`:: GET request to retrieve a specific article by its ID.

## Dependencies

This project uses the following dependencies:  
- [mux](https://github.com/gorilla/mux): A powerful HTTP router for building Go web applications.
- [gocron](https://github.com/go-co-op/gocron): A Golang library for cron scheduling.
Please refer to the respective documentation for more information on these dependencies.


## Improvements
The test coverage in this projects needs to be increased, the tests should cover all functions and areas of code to ensure the code's outmost reliability.

 