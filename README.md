# whoknows

## Overview
This monorepo contains a comprehensive web application modernization project, developed as part of the DevOps elective course at KEA (Copenhagen School of Design and Technology). 

The project demonstrates the transformation of a legacy Python/Flask application into a modern, distributed system using Go and React, while implementing contemporary DevOps practices and tools.

The project showcases real-world application of DevOps principles including containerization with Docker, CI/CD through GitHub Actions, automated testing, security scanning (CodeQL, ZAP), dependency management (Dependabot), monitoring (Prometheus, Loki, Grafana etc.) and
Container orchestration is handled through Docker Compose, providing configuration-as-code for our development, testing, and production environments.

 The modernized architecture emphasizes scalability, security, and maintainability, serving as a practical example of modern software development and deployment practices.

Our solution includes frontend (React/TypeScript), backend (Go), automated content scraping, and comprehensive monitoring, all orchestrated through as close to professional-grade DevOps workflowsas we can get. 

This project represents both a technical transformation and an educational journey in applying DevOps methodologies to legacy system modernization.

## Legacy System Description

Refactored Python/Flask monolith from an abandoned Python2 codebase.

#### Current version features:

- Python 3.x with Flask 3.0.3
- SQLite database with users/pages tables
- Bcrypt password hashing
- Bilingual content (EN/DA)
- Session-based authentication
- REST API endpoints
- Jinja2 templating

#### Key limitations:

- SQLite scalability constraints
- Limited search capabilities
- Technical debt from original Python2 version

## Modernized System

### Frontend

#### Modern React/TypeScript SPA with:

- `Vite` for fast development and optimized builds
- `Docker` configurations for dev, test, and production
- `TailwindCSS` for styling
- `Playwright` for E2E testing
- Components structure:
  - Reusable UI components (Nav, LoadingSpinner, etc.)
  - Type-safe helpers and utilities
  - Organized view components
- Environment-based API configuration
- React Router for navigation
- Toast notifications
- Cookie-based authentication
- Weather integration
- Nginx for production serving

#### Infrastructure:

- Multi-stage Docker builds
- Development hot-reloading
- Production optimization
- Automated testing setup

### Backend
#### Modern Go service with:

- `RESTful API` with `Gorilla Mux` router
- JWT-based authentication
- `Swagger` API documentation
- `Prometheus` metrics and monitoring
- Structured logging with `Logrus`
- Request `sanitization` and `validation`
- `CORS` support
- `Docker` configurations for dev/test/prod

#### Key Features

- Hot reloading in development
- Middleware for metrics, caching and authentication
- Input validation and sanitization
- Comprehensive error handling
- API versioning

#### Database
`PostgreSQL` with `GORM` ORM featuring:

- Migration management
- Transaction support
- Four main models:
  - User (auth, profile)
  - Page (content storage)
  - JWT (token management)
  - SearchLog (query tracking)
- Test support with SQLite

## Wikipedia Content Scraper Lambda
### Overview
Lambda function that automatically scrapes Wikipedia content based on popular user queries.

### Features

- Scrapes top 20 queries from query_log table
- 48-hour cooldown between scrapes for each query
- Stores content in pages table
- Runs daily at 22:00 UTC
- Records metrics for monitoring and debugging

### Technical Details

- Go 1.21
- Uses Colly for scraping
- PostgreSQL for storage
- AWS Lambda with EventBridge trigger
- Retry logic for DB connections
- Request timeouts for scraping
- Transaction-based content storage

### Testing

1. Build: GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
2. Package: zip function.zip bootstrap
3. Deploy and test via AWS Lambda console
4. Monitor execution in CloudWatch logs

### Deployment
EventBridge cron expression: cron(0 22 ? * * *)

## Monitor System
See [The Monitor Repository](https://github.com/CEM-KEA/whoknows_monitor)

## DevOps Evolution
### Legacy System
Basic Python Flask application with:

- Manual deployments
- No containerization
- Simple requirements.txt
- No automated testing
- No monitoring
- Basic error logging
- No security scanning
- Manual dependency updates

### Modernized System
Comprehensive DevOps pipeline with:

**Infrastructure & Containerization**
- Docker multi-stage builds
- Environment-specific configurations
- Docker Compose for service orchestration
- Azure VM deployment
- AWS Monitor deployment

**CI/CD Pipeline**

- Automated GitHub Actions workflows
- Semantic versioning
- Container registry integration
- Automated deployment
- Post-deployment verification

**Testing & Quality**

- Unit testing
- E2E testing with Playwright
- ZAP security scanning
- SonarQube code quality
- Comprehensive test environments

**Security**

- CodeQL analysis
- Dependabot dependency management
- Security vulnerability scanning
- Authentication monitoring
- CORS and input validation and sanitazion

**Monitoring & Logging**

- Prometheus metrics
- Structured logging
- Performance monitoring
- Error tracking
- Health checks

**Development Process**

- Standardized PR process
- Issue templates
- Automated code reviews
- Documentation requirements
- Security disclosure process

This transformation represents a shift from manual, error-prone processes to automated, secure, and monitored development operations.