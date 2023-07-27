Feed Provider
This project is a feed provider that fetches news articles from a remote source and exposes them through a RESTful API. It provides three endpoints: 
/ping
, 
/articles
, and 
/articles/{id}
.

Endpoints
/ping
: Returns a "PONG" response to check if the server is running.
/articles
: Returns all articles in JSON format.
/articles/{id}
: Returns a specific article by its ID in JSON format.
Running Tests
To run the tests for this project, follow these steps:

Make sure you have Go installed on your machine.
Open a terminal and navigate to the project directory.
Run the following command to execute the tests:
   go test ./...
This will run all the tests in the project.

Running Locally
To run the project locally, follow these steps:

Make sure you have Go installed on your machine.
Open a terminal and navigate to the project directory.
Run the following command to build the project:
   go build -o feed-provider ./main.go
After the build is successful, run the following command to start the server:
   ./feed-provider
The server will start running on 
http://localhost:8080
.

Running with Docker Compose
To run the project with Docker Compose, follow these steps:

Make sure you have Docker and Docker Compose installed on your machine.
Open a terminal and navigate to the project directory.
Run the following command to build the Docker image:
   docker-compose build
After the build is successful, run the following command to start the containers:
   docker-compose up
The server will start running on 
http://localhost:8080
.

API Documentation
/ping
: GET request to check if the server is running.
/articles
: GET request to retrieve all articles.
/articles/{id}
: GET request to retrieve a specific article by its ID.
Dependencies
This project uses the following dependencies:

mux: A powerful HTTP router for building Go web applications.
gocron: A Golang library for cron scheduling.
Please refer to the respective documentation for more information on these dependencies.

That's it! You should now have a clear understanding of the project and how to run it locally and with Docker Compose. Feel free to modify the instructions and add any additional information specific to your project.