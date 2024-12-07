# CNAD-Assignment1
CNAD Assignment 1 - Joseph Wan

# Electric Car Sharing System Microservice

## Introduction

This project implements a microservices-based electric-sharing system application/platform. The main few features and functions include allowing users to reserve vehicles, figure out their billing based on their membership type, rental duration and view their invoices. This project emphasies modular design using Go and adheres to a microservices architecture for scalability and maintainability.

## Design Considerations
The microservices for this project have been organised into separate services, each contained within their own respective folder for clarity, organisation and ease of maintenance. Each service is responsible for a specific functionality of the system, ensuring that they are loosely coupled and can evolve independently. Should a particular service fail, the application should remain functional as long as all services do not fail.

<h3><strong>Breakdown of Each Service</strong></h3>

- User Service: Responsible for the management of user information, including their username, email, password and membership type.

- Vehicle Service: Handles the collection of vehicles on the platform, vehicle reservations

- Billing Service: Responsible for calculating the total amount for a reservation and applying discounts based on the type of membership a user has, and storing their respective invoices.

## Microservices Architecture Diagram

