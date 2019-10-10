# Logging-Service
This is simple logging service gets an events as http requests and writes it to a database.

## Getting Started

### Prerequisites
You need to install MongoDB on your local machine or you should have MongoDB instance in your network before you can run the service. Please check [here](https://docs.mongodb.com/manual/installation/) for installation instructions.

### Installation
1. Run `dep ensure` in order to get all necessary packages
2. Compile `go build`

### Running
For starting the service use the following: </br>
`./logging-service --srvAddr "localhost" --srvPort 8080 --dbAddr "localhost" --dbPort 27017`</br>
It logs http request metrics in every 5 second interval.
### Load Test
1. Navigate to Load-Test directory: `cd tools/load-test/`
2. Compile: `go build main.go`
3. Run: `./main --time <duration>` duration should be specified in minutes, for instance: `./main --time 2` - will run the test for two minutes. 
