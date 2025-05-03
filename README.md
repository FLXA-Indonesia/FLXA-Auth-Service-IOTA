# FLXA Auth Service for IOTA

Welcome to the official repository for the FLXA Auth Service, developed by FLXA Indonesia. This service provides authentication functionalities within the FLXA ecosystem.

## Overview
The FLXA Auth Service offers a backend solution for managing user authentication processes. Built with Go.

## Features
1. User Authentication: Secure login and registration processes.

2. Token Management: Issue and validate authentication tokens.

3. RESTful API: Exposes endpoints for frontend integration.

4. Deployment Ready: Configured for deployment on platforms like Vercel.

## Repository Structure
```
├── api/                  # API route handlers
├── controllers/          # Business logic controllers
├── initializers/         # Initialization scripts and configurations
├── migrate/              # Database migration files
├── models/               # Data models
├── utilities/            # Utility functions
├── .env.example          # Example environment variables
├── .gitignore            # Git ignore rules
├── go.mod                # Go module file
├── go.sum                # Go dependencies checksum file
└── vercel.json           # Vercel deployment configuration
```

## Getting Started
### Prerequisites

1. Go (version 1.16 or higher)

### Installation
1. Clone the repository:
```bash
git clone https://github.com/FLXA-Indonesia/FLXA-Auth-Service-IOTA.git
cd FLXA-Auth-Service-IOTA
```

2. Install dependencies:
```bash
go mod tidy
```

### Running the Application
To start the development server:

```bash
go run main.go
```

The server will start on the default port (e.g., http://localhost:3000). You can modify the port and other configurations as needed.

### Deployment
The project includes a vercel.json file, making it ready for deployment on Vercel. To deploy:

1. Install Vercel CLI:
```bash
npm install -g vercel
```

2. Deploy:
```bash
vercel
```

Follow the prompts to complete the deployment process.

## API Endpoints
The service exposes the RESTful API endpoints accessible at `/api`

Note: Authentication and authorization mechanisms should be implemented to secure these endpoints.

## Contributing
We welcome contributions to enhance the FLXA User Service. To contribute:
1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Commit your changes with clear messages.
4. Push your branch and open a pull request detailing your modifications.
5. Please ensure your code adheres to the project's coding standards and includes relevant tests.

## License
This project is licensed under the [GNU Affero General Public License V3](LICENSE)
