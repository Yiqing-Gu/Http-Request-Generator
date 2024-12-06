
HTTP Server and Request Simulator

## Overview

This is a Go program that serves as both an HTTP server and a client for sending HTTP POST requests based on configurations defined in a JSON file. `serv/serv.go` can:

1. Start an HTTP server listening on a specified port.
2. Handle incoming requests, responding with a simple "hi" message.
3. Read a JSON configuration file to register dynamic routes and send HTTP POST requests.
4. Parse and print data from the JSON configuration file, including one-time data arrays.

This program can also act as a **MOCK SERVER** to simulate an HTTP server with configurable paths and behavior.

---

## Features

- **HTTP Server**:
  - Starts an HTTP server based on the configuration file or default settings.
  - Registers dynamic routes specified in the configuration file.

- **Configurable HTTP Client**:
  - Sends HTTP POST requests to specified servers with configurable paths and body.

- **MOCK SERVER**:
  - Simulates an HTTP server with user-defined routes and responses.

- **JSON Parsing**:
  - Uses the `gjson` library to extract specific values from the JSON file.

- **Data Handling**:
  - Prints arrays of one-time data from the JSON configuration file.

---

## Requirements

- Go 1.18+ installed on your system.
- A valid JSON configuration file (default name: `serverConfig.json`).

---

## Installation

1. Clone the repository or download `serv.go` to your local machine.
2. Install required dependencies:
   ```bash
   go get github.com/tidwall/gjson
   ```

---

## Usage

### Run the Program
```bash
go run serv.go -f <json-config-file>
```

- The `-f` flag allows you to specify the JSON configuration file. If omitted, it defaults to `mock-http.json`.

### Example JSON Configuration File

#### File: `serverConfig.json`

```json
{
  "config": {
    "port": "8080",
    "path": ["/api/test", "/api/demo"],
    "baseUrl": "http://localhost:8080",
    "body": "{"message": "hello"}"
  },
  "data": {
    "once": [1, 2, 3],
    "repeat": [4, 5, 6]
  }
}
```

---

### Using MOCK SERVER

#### Steps:
1. Ensure the `serverConfig.json` file exists in the execution directory.
2. Run the program:
   - Default: Automatically loads `serverConfig.json`.
   - Specify a different configuration file:
     ```bash
     go run serv.go -f customConfig.json
     ```

#### Configuration Details:
1. **`port`**:
   - Specifies the port on which the server will listen.
   - Directly use the port number (e.g., `"8080"`), without adding a colon.

2. **`path`**:
   - Specifies the routes to register.
   - Duplicate routes will only be registered once.

#### Example:
- Start the server:
  ```bash
  go run serv.go -f serverConfig.json
  ```

- Access registered routes (e.g., `/api/test`):
  ```bash
  curl http://localhost:8080/api/test
  ```

---

### HTTP POST Requests
When running with a valid configuration file, the program sends an HTTP POST request to the server specified in `baseUrl` and `path`.

Example:
- JSON Config:
  ```json
  {
    "config": {
      "baseUrl": "http://localhost:8080",
      "path": "/api/demo",
      "body": "{"key":"value"}"
    }
  }
  ```

- Program output:
  ```
  [1, 2, 3]
  body
  ```

---

## Code Explanation

1. **HTTP Server**:
   - Starts an HTTP server and registers routes specified in `path`.

2. **JSON Parsing**:
   - Reads the JSON configuration file to configure routes and requests.

3. **HTTP POST**:
   - Sends a POST request with the specified body to `baseUrl` + `path`.

4. **MOCK SERVER**:
   - Registers and handles incoming requests on specified routes.

---

## Test Program Directory (`testProgram`)

The `testProgram` directory contains **JavaScript files** used for testing the functionality of systems, such as subway control systems. These programs are designed to:

1. Simulate various operational scenarios.
2. Test system responses.
3. Validate data processing workflows.

### Usage:
- These scripts can be executed in a controlled environment to ensure the system behaves as expected.
- Before running, ensure your environment is configured to meet the test requirements.

The directory serves as a valuable resource for developers working on system validation and functionality testing. Refer to the included documentation or code comments in the `testProgram` directory for more details.
