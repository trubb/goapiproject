# Instructions for building and running the api server

## Building

TODO fix this document:
To build the project, run `make apiserver` or `make build` in the project root.  
To clean up after a previous build, run `make clean`.  

## Running the server

To run the respective component from the project root, run the following:

```md
./apiserver --db <port to communicate on>
```

The `apiserver` uses the package `cli`[(urfave/cli/v2 on pkg.go.dev)](https://pkg.go.dev/github.com/urfave/cli/v2) to deal with flags, run `./apiserver -h`message shown below.

```md
$ ./apiserver -h
NAME:
   apiserver - Hi hello yes here should be information

USAGE:
   apiserver [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --apiPort PORT, --api PORT  assign PORT to send API calls to (default: "8080")
   --dbPort PORT, --db PORT    assign PORT to communicate with the database on (default: "4711")
   --help, -h                  show help (default: false)
```

## Running tests

TODO this section  
Run `go test [-v]` in `/countd/` to run tests for the `countd` functions `handleIncomingWord`, `updateWordAggregate`, `flush`, and `readHotList`.

I'm frankly not entirely sure how to test the `main` and `countd` functions when using `urfave/cli`, and have therefore not added tests for those two functions.

## sending requests (from outside WSL2)

access <http://wsl-eth0-IP:8080/ping>
eg <http://172.19.213.164:8080/ping>
