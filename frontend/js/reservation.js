
// Decode JWT and extract user ID
const token = localStorage.getItem('jwt');
const decoded = JSON.parse(atob(token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')));

// Directly assign the user ID
const userId = decoded.userid;

console.log("User ID:", userId);

// Fetch available vehicles and display them
fetch('http://localhost:8082/vehicles')
    .then(response => response.json())
    .then(vehicles => {
        const vehiclesList = document.getElementById('vehicles-list');
        vehicles.forEach(vehicle => {
            const vehicleDiv = document.createElement('div');
            vehicleDiv.classList.add('vehicle');
            vehicleDiv.innerHTML = `
                <h3>${vehicle.make} ${vehicle.model}</h3>
                <p>License Plate: ${vehicle.license_plate}</p>
                <button class="reserve-btn" data-vehicle-id="${vehicle.id}">Reserve Vehicle</button>
                <div class="reservation-form" id="reservation-form-${vehicle.id}" style="display:none;">
                    <label for="start-time-${vehicle.id}">Start Time:</label>
                    <input type="datetime-local" id="start-time-${vehicle.id}" required>
                    <label for="end-time-${vehicle.id}">End Time:</label>
                    <input type="datetime-local" id="end-time-${vehicle.id}" required>
                    <button class="submit-reservation" data-vehicle-id="${vehicle.id}">Submit Reservation</button>
                </div>
            `;
            vehiclesList.appendChild(vehicleDiv);
        });

        // Handle reserve button click
        document.querySelectorAll('.reserve-btn').forEach(button => {
            button.addEventListener('click', function() {
                const vehicleId = this.getAttribute('data-vehicle-id');
                document.getElementById(`reservation-form-${vehicleId}`).style.display = 'block';
            });
        });

        // Handle form submission for reservation
        document.querySelectorAll('.submit-reservation').forEach(button => {
            button.addEventListener('click', function() {
                const vehicleId = this.getAttribute('data-vehicle-id');
                console.log("Vehicle ID:", vehicleId);
                console.log("Logged User ID:", userId);
                const startTime = document.getElementById(`start-time-${vehicleId}`).value;
                const endTime = document.getElementById(`end-time-${vehicleId}`).value;

                const startDate = new Date(startTime);
                const endDate = new Date(endTime);

                // Convert the Date objects to ISO 8601 format (which is the standard format expected by most APIs)
                const startTimeISO = startDate.toISOString();
                const endTimeISO = endDate.toISOString();

                //console.log(startTimeISO, endTimeISO);

                const reservationData = {
                    user_id: parseInt(userId),
                    vehicle_id: parseInt(vehicleId),
                    start_time: startTimeISO,
                    end_time: endTimeISO
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
                        // Hide the form after reservation
                        document.getElementById(`reservation-form-${vehicleId}`).style.display = 'none';
                    } else {
                        alert('Error reserving vehicle');
                    }
                })
                .catch(e => {
                console.error("Error:", e);
            });
        });
    });
});

                


