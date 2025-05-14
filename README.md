# Deck Service

The Deck Service is a microservice that provides a gRPC API for managing decks of cards. It allows users to create and manipulate decks and cards.

This service is designed to be used as part of a larger application that requires card deck functionality (e.g., blackjack, poker).

## Running the Service

The service is hard coded to run on port 8080, and requires no parameters or arguments to run.

To start the server, simple run the following command:

```bash
go run cmd/deck/main.go
```

## API Documentation

The API is defined using protobuf, the contract implemented here is located at https://github.com/chn555/schemas/blob/main/proto/deck/v1/deck.proto