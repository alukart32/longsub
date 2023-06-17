# Longsub app

This project is an example of a search service for finding the longest substring without repeating characters.

## How to run

The HTTP server can be started using Docker:

```shell
make docker-up
```

To find the largest substring, you need to execute the CLI application:

```shell
longsub your_string http://localhost:8080/api/substring
```
