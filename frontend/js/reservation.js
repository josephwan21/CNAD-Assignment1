// reservation.js

// Fetch available vehicles and display them
fetch('http://localhost:8082/vehicles')
    .then(response => response.json())
    .then(vehicles => {
        const vehiclesList = document.getElementById('vehicles-list');
        vehicles.forEach(vehicle => {
            const vehicleDiv = document.createElement('div');
            vehicleDiv.classList.add('vehicle');
            vehicleDiv.innerHTML = `<h3>${vehicle.make} ${vehicle.model}</h3><p>License Plate: ${vehicle.license_plate}</p>`;
            vehiclesList.appendChild(vehicleDiv);
        });
    })
    .catch(err => console.error('Error:', err));

// Reserve vehicle
document.getElementById('reserve-btn').addEventListener('click', function() {
    const selectedVehicleId = 1;  // Just for demo, should be dynamically selected
    const userId = 1;             // User ID from session/localStorage

    const reservationData = {
        user_id: userId,
        vehicle_id: selectedVehicleId,
        start_time: '2024-12-01T10:00:00Z',
        end_time: '2024-12-01T12:00:00Z'
    };

    fetch('http://localhost:8082/reserve', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(reservationData)
    })
    .then(response => response.json())
    .then(data => {
        if (data.status === 'Active') {
            alert('Reservation successful');
        } else {
            alert('Error reserving vehicle');
        }
    })
    .catch(err => console.error('Error:', err));
});
