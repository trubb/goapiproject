# goapiproject

A very basic REST API utilizing [gin-gonic/gin](https://github.com/gin-gonic/gin) for routing, [etcd-io/bbolt](https://github.com/etcd-io/bbolt) instead of a traditional database, and [https://github.com/asdine/storm](asdine/storm) as a wrapper around bbolt.

I chose to go for a less comprehensive API mainly due to a lack of time to sit down and work through the assignment.

## Building

- To build the project, run `make apiserver` or `make build` in the project root.  
- To clean up after a previous build, run `make clean`.  

## Running the server

To run the project from the project's root folder after building, run the following:

```md
./apiserver -api <port to communicate on>
```

The `apiserver` uses [urfave/cli/v2](https://github.com/urfave/cli/v2) to deal with flags, run `./apiserver -h` to display the help message shown below.

```md
$ ./apiserver -h
NAME:
   apiserver - A very basic REST API with basic database backing

USAGE:
   apiserver [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h               show help (default: false)
   --port PORT, --api PORT  assign PORT to send API calls to (default: "8080")
```

## Running tests

Run `make test` in the project's root folder to run tests for the handler functions.

`make test` runs the following commands:

```makefile
test:
   go test -v; \
   rm *_test.db
```

## API endpoints

The Achievement return type is defined as:

```go
type Achievement struct {
  ID          int `storm:"id,increment"`
  Name        string
  Description string
}
```

### Get one achievement

Returns the achievement with the specified ID.

Path: `/achievement/one`

Parameters:

- `id - integer (required)`

Response components:

- `achievement - Achievement`
  - `ID - positive integer`
  - `Name - string`
  - `Description - string`

Sample call with parameters: `<host>:<port>/achievement/one?id=1`

Returned message:

```json
{
    "achievement": {
        "ID": 1,
        "Name": "Explore",
        "Description": "Explored many parts of the world"
    }
}
```

### Get all achievements

Returns all achievements in the database.

Path: `/achievement/all`

Parameters:

- `none`

Response components:

- `achievements - []Achievement`
  - `ID - positive integer`
  - `Name - string`
  - `Description - string`

Sample call, no parameters needed: `<host>:<port>/achievement/all`

Returned message:

```json
{
  "achievements": [
      {
          "ID": 1,
          "Name": "Explore",
          "Description": "Explored many parts of the world"
      },
      {
          "ID": 2,
          "Name": "Produce",
          "Description": "Produced more than 200 items across all scenarios"
      },
      {
          "ID": 3,
          "Name": "Showtime",
          "Description": "Put on a show with the combine harvester"
      }
  ]
}
```

### Create a achievement

Adds an achievement to the database.

Path: `/achievement/create`

Parameters:

- `id - integer (required)`
- `Name - string (required)`
- `Description - string (required)`

Response components:

- `achievement - Achievement`
  - `ID - positive integer`
  - `Name - string`
  - `Description - string`
- `message - string`

Sample call with parameters: `<host>:<port>/achievement/create?name=Hello&description=Waved hello to another person`

Returned message:

```json
{
    "achievement": {
        "ID": 4,
        "Name": "Hello",
        "Description": "Waved hello to another person"
    },
    "message": "successfully saved new achievement 'Hello' as ID: 4"
}
```

### Update existing achievement

Updates an achievement that already exists in the database.

Path: `/achievement/update`

Parameters:

- `id - integer (required)`
- `Name - string (optional)`
- `Description - string (optional)`

Please note that at least one of `Name` and `Description` must be provided to result in an actual update.

Response components:

- `achievement`
  - `ID - positive integer`
  - `Name - String`
  - `Description - String`
- `message - string`

Sample call with parameters: `<host>:<port>/achievement/update?id=1&name=Answer&description=Hi hello yes`

Returned message:

```json
{
    "achievement": {
        "ID": 1,
        "Name": "Answer",
        "Description": "Hi hello yes"
    },
    "message": "successfully updated achievement to ID: '1 Name: 'Answer' Description: Hi hello yes"
}
```

### Delete existing achievement

Deletes an existing achievement from the database.

Path: `/achievement/delete`

Parameters:

- `id - integer (required)`

Response components:

- `message - string`

Sample call with parameters: `<host>:<port>/achievement/delete?id=1`

Returned message:

```json
{
    "message": "successfully deleted achievement with ID: '1"
}
```

### Reset database to initial state

Resets the database to a default starting state.

Path: `/achievement/reset/really`

Parameters:

- `none`

Response components:

- `message - string`

Sample call with parameters: `<host>:<port>/achievement/reset/really`

Returned message:

```json
{
    "message": "Successfully reset the database to its initial state"
}
```

### Availability check

A rudimentary way to let requestees know that the server is online and reachable.

Path: `/online`

Parameters:

- `none`

Response components:

- `message - string`

Sample call, no parameters needed: `<host>:<port>/online`

Returned message:

```json
{
    "message": "I'm alive!"
}
```
