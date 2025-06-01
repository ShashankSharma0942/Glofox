## Glofox

This is a backend service for managing fitness class bookings. The project is structured to handle various functionalities such as managing classes, creating bookings, and API request handling.

## Directory Breakdown

### `glofox/`

- **core**: Contains the core logic for data storage.
  - `mapstore.go`: Implements `MapStore` for storing class and booking data.

- **config**: Contains configuration loading logic.
  - `config.go`: Reads `config.json` and provides configuration to the application.

- **internal**: Holds the business logic and API request handling for classes and bookings.
  - **service**: Contains the service layer for business logic.
    - `booking_service.go`: Handles booking logic.
    - `class_service.go`: Handles class creation and management.
    - `service_test.go`: Contains unit tests for services.

  - **handler**: Contains the API handlers that manage the HTTP requests for classes and bookings.
    - `booking_handler.go`: API endpoint for booking-related requests.
    - `class_handler.go`: API endpoint for class-related requests.
    - `handler_test.go`: Unit tests for handlers.
  
- **models/dto**: Contains data transfer objects (DTOs) for class and booking data.
  - `booking.go`: Models related to bookings.
  - `class.go`: Models related to classes.

- **utils**: Utility functions for handling common tasks across the application.
  - `response.go`: Utility for generating standard API responses.

- **cmd/server**: Contains the server bootstrapping logic.

- **main.go**: The entry point of the application that sets up and starts the server.

### Root Files

- **`go.mod`**: Go module file that defines project dependencies.
- **`go.sum`**: Go checksum file for dependencies.
- **`config.json`**: JSON configuration file containing runtime parameters.
- **`README.md`**: This file, which provides an overview of the project and its structure.

## Configuration

### `config.json`

The application reads its runtime configuration from a `config.json` file located at the root of the project, if you want to change the application port or base route, please use config.json Here's a sample of config.json:
  ```json
  {
    "port": "7000",
    "baseRoute": "/glofox",
    "dateFormat": "2006-01-02"
  }
   ```
## How to Set Up the Project

1. Clone the repository:
   ```bash
   git clone <repository_url>
2. Navigate the Project Repositry
   ```bash
    cd glofox
3. Install the Dependencies
   ```bash
    go mod tidy 
4. RUN the Application
    ```bash
    go run main.go
## How to Run UT for Project

1. Run below Mentioned Command From cmd Directory
   ```bash
   go test ../internal/...
