-- Drop any existing databases as shown below
DROP DATABASE IF EXISTS CarSharingUserService;
DROP DATABASE IF EXISTS CarSharingVehicleService;
DROP DATABASE IF EXISTS CarSharingBillingService;

-- Creation of the User Service Database
CREATE DATABASE CarSharingUserService;
USE CarSharingUserService;

CREATE TABLE Users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    membership VARCHAR(50) DEFAULT 'Basic',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    verified BOOLEAN DEFAULT FALSE,
    verification_code VARCHAR(255)
);

CREATE TABLE Rentals (
    id INT AUTO_INCREMENT PRIMARY KEY,
    reservation_id INT,
    user_id INT,
    vehicle_id INT,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    rental_status VARCHAR(50) CHECK (rental_status IN ('Completed', 'Canceled'))
);

-- Creation of the Vehicle Service Database
CREATE DATABASE CarSharingVehicleService;
USE CarSharingVehicleService;

CREATE TABLE Vehicles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    make VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    license_plate VARCHAR(50) UNIQUE NOT NULL,
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insertion of sample data
INSERT INTO Vehicles (make, model, license_plate, is_available)
VALUES 
('Tesla', 'Model 3', 'ABC123', TRUE),
('Ford', 'Mustang', 'XYZ789', TRUE),
('Chevrolet', 'Camaro', 'DEF456', FALSE),
('BMW', 'X5', 'GHI321', TRUE),
('Audi', 'A4', 'JKL654', TRUE);

CREATE TABLE Reservations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    vehicle_id INT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status VARCHAR(20) DEFAULT 'Active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insertion of sample vehicle reservations
INSERT INTO Reservations (user_id, vehicle_id, start_time, end_time, status)
VALUES
(2, 2, '2024-12-03 13:00:00', '2024-12-03 15:00:00', 'Active'),
(3, 3, '2024-12-03 16:00:00', '2024-12-03 18:00:00', 'Completed'),
(4, 4, '2024-12-03 19:00:00', '2024-12-03 21:00:00', 'Active'),
(5, 5, '2024-12-03 22:00:00', '2024-12-03 23:59:59', 'Canceled');

-- Creation of the Billing Service Database
CREATE DATABASE CarSharingBillingService;
USE CarSharingBillingService;

CREATE TABLE Billing (
    id INT AUTO_INCREMENT PRIMARY KEY,
    reservation_id INT NOT NULL,
    user_id INT NOT NULL,
    vehicle_id INT NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    discount DECIMAL(10, 2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
 
