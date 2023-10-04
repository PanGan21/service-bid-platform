# service-bid-platform

Scalable Auction Platform, a comprehensive system designed and developed as part of a thesis project. This project focuses on building a transparent bidding platform for service auctions, implementing microservices architecture, and adhering to clean architecture principles.

## Overview

The platform leverages microservices architecture to achieve scalability, maintainability, and flexibility in handling service auctions. Clean architecture principles are followed to ensure separation of concerns, making the system modular and easy to extend.

## Getting Started

Run the Backend<br />
`docker-compose up --build`
Run the Frontend<br />

```
cd clients/frontend
yarn install
yarn start
```

Navigate to http://localhost:3000 to access the web frontend.

## Project Structure

`pkg`: Contains reusable packages shared across microservices.
`api-gateway`: Configuration for the API gateway.
clients: Houses different clients, with the current implementation being a React web frontend.
`demo`: Includes demo data for testing purposes.
`integration-tests`: Integration tests executed through Docker compose.
`scripts`: Utility scripts for the project.
`services`: Microservices, including user authentication, request handling, auction authorization, and bid management.
`ssl`: Testing SSL certificates used during development.

## Auction Mechanism

The platform employs [Vickrey auctions](https://en.wikipedia.org/wiki/Vickrey_auction) to enhance bidding competition, ensuring transparency and fairness. Role-Based Access Control (RBAC) is implemented to manage user roles effectively.

## License

This project is licensed under the MIT License, allowing you to use, modify, and distribute the software freely.
