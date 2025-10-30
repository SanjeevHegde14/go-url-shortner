# Go URL Shortener API
*(ongoing)*

This project is a simple URL Shortener API built with Go.  
It demonstrates backend API design, HTTP routing, and containerization using Docker.  
The goal is to provide a practical learning example of building and deploying a RESTful Go service from scratch.


## Overview

The Go URL Shortener is a RESTful web service that allows users to:
- Shorten a long URL into a shorter unique code
- Redirect from the short URL back to the original one
- (Planned) Store URLs persistently using a lightweight SQLite database

At this stage, the project runs a basic Go HTTP server, which will be expanded to handle real API logic and Dockerized deployment.


## Tech Stack

- **Language:** Go (Golang)
- **Database:** SQLite (planned)
- **Containerization:** Docker (in progress)
- **Version Control:** Git + GitHub


## Running the Project Locally

1. Clone this repository:
   ```bash
   git clone https://github.com/<yourusername>/go-url-shortener.git
   cd go-url-shortener
   ````

2. Run the API:
    ```bash
    go run cmd/api/main.go
    ```
3. Open your browser and visit:
    http://localhost:8080
    
