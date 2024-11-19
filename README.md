# Receipt Processor

## How to run

There are two options to run the service - on the local machine if Go is installed locally or on a Docker container. A Makefile is provided with the following targets:

### Option 1: Local Machine

```
> make run_local
```

### Option 2: Docker

```
> make run_docker
```

---

## Comments about the implementation

1. The `gorilla/mux` router is used used for its lightweight nature and helpful abstractions over the standard `net/http` package.
2. The `.env` file is included for ease of local development but should be excluded in production setups to maintain security and environment isolation.
3. An in-memory Go map is used for data storage as persistence was not a requirement for this implementation.
4. Areas of improvement:
   1. Implement unit tests to ensure robustness and maintainability.
   2. Perform stricter validation of the receipt JSON payload to handle edge cases.
   3. Separate business logic from handlers into dedicated service layers for better modularity and reusability.
