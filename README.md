<!-- [alt text](https://www.eyeonicon.xyz/media/logos/logo_256.png "Eye On Icon") -->

<p align="center">

  <a href="https://godoc.org/github.com/eyeonicon/get-omm-votes-go">
    <img src="https://godoc.org/github.com/eyeonicon/get-omm-votes-go?status.svg" alt="GoDoc">
  </a>

  <a href="./LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT">
  </a>

  <!-- make one for go ref stuff -->
  <a href="https://goreportcard.com/report/github.com/eyeonicon/get-omm-votes-go">
    <img src="https://goreportcard.com/badge/github.com/eyeonicon/get-omm-votes-go" alt="Go Report Card">
  </a>


</p>

# Get OMM delegation on specific validator

This is a simple Go package to get the OMM delegation on a specific validator.

## How to use

First clone repository. 

The main package is located at cmd/main.go. You can build it with the following command:

```bash
go build -o get-omm-votes-go cmd/main.go
```

Then you can run it with the following command:

```bash
./get-omm-votes-go <validator_address>
```

There will be a report exported to the reports folder with a filename that is a timestamp.

The calls package contains all the data fetching stuff. If needed there is also 'GetOMMTotalVotes' function that returns the total votes on the given validator.


