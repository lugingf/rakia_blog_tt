# Simple Blog API platform with minimal functionality

Test Task

## Requirements

### Golang:
Ensure you have Go installed. You can download and install it from Go's official site.
Verify the installation by running:

```sh
go version
```

### golangci-lint:

Install golangci-lint for linting your Go code. You can install it using the following command:

```sh
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

#### Verify the installation by running:

```sh
golangci-lint --version
```

#### Swagger Codegen:

Install swagger for generating models from the API specification. You can install the swagger command-line tool using the following command:

```sh
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
```

#### Verify the installation by running:

```sh
swagger version
```

#### Make:

Ensure you have make installed on your system. Most Unix-like systems come with make pre-installed. On Debian-based systems, you can install it using:

```sh
sudo apt-get install build-essential
```

On MacOS, you can install it using:

```sh
xcode-select --install
```

### Environment Variables:
Ensure you have a .env.local file with the necessary environment variables. This file is sourced when running the application.

### Summary of Makefile Commands

#### lint:
Runs golangci-lint to lint the Go code.

```sh
make lint
```

#### check:
Runs lint and tests.

```sh
make check
```

#### run:
Builds the application and runs it, sourcing environment variables from .env.local.

```sh
make run
```

#### build:
 Builds the application with optimization flags.

```sh
make build
```

#### models:
 Generates models using Swagger from the provided API specification.

```sh
make models
```

## Endpoints

### Create a New Blog Post

```sh
curl -X POST http://localhost:8080/posts \
-H "Content-Type: application/json" \
-d '{
  "title": "Title 1",
  "content": "Quaerat sit dolorem velit. Ipsum non tempora magnam neque tempora. Tempora dolorem adipisci tempora neque labore. Dolorem sed dolore sed. Voluptatem consectetur dolor voluptatem. Quiquia adipisci voluptatem modi dolore. Dolor etincidunt neque consectetur dolor. Numquam etincidunt voluptatem sit amet tempora. Modi dolorem sed magnam consectetur. Dolor dolorem est amet magnam velit.",
  "author": "Author 1"
}'
```

Retrieve All Blog Posts

```sh
curl -X GET http://localhost:8080/posts
```

Retrieve a Specific Blog Post

```sh
curl -X GET http://localhost:8080/posts/1
```

Update an Existing Blog Post

```sh
curl -X PUT http://localhost:8080/posts/1 \
-H "Content-Type: application/json" \
-d '{
  "title": "Updated Title 1",
  "content": "Updated content for the first post.",
  "author": "Updated Author 1"
}'
```

Delete a Blog Post

```sh
curl -X DELETE http://localhost:8080/posts/1
```

#### Running the Server

To run the server, use the following command:

```sh
go run main.go
```

Or you can build and run as executable file
```sh
make build
make run
```

Replace http://localhost:8080 with the actual URL of your running server if it's different.