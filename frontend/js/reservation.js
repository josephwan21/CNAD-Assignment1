
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
                <form class="reservation-form" id="reservation-form-${vehicle.id}" style="display:none;">
                    <label for="start-time-${vehicle.id}">Start Time:</label>
                    <input type="datetime-local" id="start-time-${vehicle.id}" required>
                    <label for="end-time-${vehicle.id}">End Time:</label>
                    <input type="datetime-local" id="end-time-${vehicle.id}" required>
                    <p class="cost-estimate" id="cost-estimate-${vehicle.id}">Estimated Cost: $0.00</p>
                    <button type="submit" class="submit-reservation" data-vehicle-id="${vehicle.id}">Submit Reservation</button>
                </form>
            `;
            vehiclesList.appendChild(vehicleDiv);
        });

        // Handle reserve button click
        document.querySelectorAll('.reserve-btn').forEach(button => {
            button.addEventListener('click', function() {
                const vehicleId = this.getAttribute('data-vehicle-id');
                document.getElementById(`reservation-form-${vehicleId}`).style.display = 'block';

                // Attach listeners to the start and end time inputs
                const startTimeInput = document.getElementById(`start-time-${vehicleId}`);
                const endTimeInput = document.getElementById(`end-time-${vehicleId}`);
                const costEstimateElement = document.getElementById(`cost-estimate-${vehicleId}`);
        
                [startTimeInput, endTimeInput].forEach(input => {
                    input.addEventListener('change', () => {
                        const startTime = startTimeInput.value;
                        const endTime = endTimeInput.value;
        
                        if (!startTime || !endTime) {
                            costEstimateElement.textContent = `Estimated Cost: $0.00`;
                            return;
                        }

                        if (endTime <= startTime) {
                            costEstimateElement.textContent = `End time must be after start time.`;
                            startTimeInput.value = "";
                            endTimeInput.value = "";
                            return;
                        }
        
                        const estimateData = {
                            user_id: parseInt(userId),
                            vehicle_id: parseInt(vehicleId),
                            start_time: new Date(startTime).toISOString(),
                            end_time: new Date(endTime).toISOString()
                        };
        
                        fetch('http://localhost:8083/calculatebilling', {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify(estimateData)
                        })
                        .then(response => response.json())
                        .then(data => {
                            console.log(data);
                            costEstimateElement.textContent = `Estimated Cost: $${data.total_amount}`;
                        })
                        .catch(err => {
                            console.error('Error fetching cost estimate:', err);
                            costEstimateElement.textContent = `Error fetching cost`;
                        });
                    });
                });
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
                    console.log("Reservation data:", data);
                    if (data.status === 'Active') {
                        alert('Reservation successful');
                        // Hide the form after reservation
                        createInvoice(data.id, parseInt(vehicleId), startTimeISO, endTimeISO);
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

// Function to create the invoice after a reservation is successful
function createInvoice(reservationId, vehicleId, startTime, endTime) {
    const invoiceData = {
        reservation_id: reservationId,
        user_id: userId,
        vehicle_id: vehicleId,
        start_time: startTime,
        end_time: endTime
    };

    // Send request to backend to generate the invoice
    fetch('http://localhost:8083/createinvoice', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(invoiceData)
    })
    .then(response => response.json())
    .then(invoice => {
        console.log(invoice)
        // Handle successful invoice creation
        if (invoice) {
            // Optionally, display more invoice details here
            console.log("Invoice details:", invoice);
        } else {
            alert("Error generating invoice.");
        }
    })
    .catch(err => {
        console.error('Error creating invoice:', err);
        alert('Error generating invoice');
    });
}



                


