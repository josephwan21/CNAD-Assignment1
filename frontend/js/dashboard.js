const token = localStorage.getItem('jwt');

if (!token) {
    alert('You need to log in first');
    window.location.href = 'login.html';  // Redirect to login page if no token found
} else {
    fetchUserData(token);
}


// Decode JWT and extract user ID
const decoded = JSON.parse(atob(token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')));

// Directly assign the user ID
const userId = decoded.userid;

console.log("User ID:", userId);

// Fetch user details from the backend using the auth token
function fetchUserData(token) {
    fetch('http://localhost:8080/user-profile', {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`  // Send the token in the Authorization header
        }
    })
    .then(response => response.json())
    .then(data => {
        console.log(data)
        if (data) {
            // Display user info (name and membership tier)
            document.getElementById('welcome-message').innerText = `Hello, ${data.name}`;
            document.getElementById('membership-tier').innerHTML = `<strong>${data.membership}</strong>`;

            // Fetch and display user's reservations
            fetchReservations();
            fetchInvoicesByUser(userId);
            fetchRentalHistory(userId);
        } else {
            alert('Unable to fetch user data');
        }
    })
    .catch(err => {
        console.error('Error fetching user data:', err);
        alert('An error occurred while fetching your profile.');
    });
}

// Handle profile update submission
document.getElementById('update-profile-form').addEventListener('submit', function(event) {
    event.preventDefault();

    const name = document.getElementById('update-name').value;
    const email = document.getElementById('update-email').value;

    const token = localStorage.getItem('jwt');
    
    fetch('http://localhost:8080/update-profile', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ name, email })
    })
    .then(response => response.json())
    .then(data => {
        if (data.message === 'Profile updated successfully.') {
            alert('Profile updated successfully!');
            localStorage.setItem('jwt', data.token);
            window.location.reload();
        } else {
            alert('Profile update failed');
        }
    })
    .catch(err => console.error('Error:', err));
});

// Fetch the user's reservations
function fetchReservations() {
    const token = localStorage.getItem('jwt');
    fetch('http://localhost:8082/reservations', {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
    .then(response => response.json())
    .then(data => {
        console.log('Fetched reservations data:', data);
        if (data && data.length > 0) {
            // Display reservations
            const reservationsList = data.map(reservation => {
                return `<div class="reservation">
                            <p>Vehicle: <strong>${reservation.make} ${reservation.model}</strong></p>
                            <p>Reservation from ${reservation.start_time} to ${reservation.end_time}</p>
                            <button class="update-res-btn" data-vehicle_id="${reservation.vehicle_id}" data-id="${reservation.id}">Update Reservation</button>
                            <button class="delete-res-btn" data-vehicle_id="${reservation.vehicle_id}" data-start_time="${reservation.start_time}" data-end_time="${reservation.end_time}" data-id="${reservation.id}">Delete Reservation</button>
                            <form class="update-form" id="update-form-${reservation.id}" style="display: none;">
                                <label for="start-time-${reservation.id}">Start Time:</label>
                                <input type="datetime-local" id="start-time-${reservation.id}" required>
                                <label for="end-time-${reservation.id}">End Time:</label>
                                <input type="datetime-local" id="end-time-${reservation.id}" required>
                                <p class="cost-estimate" id="cost-estimate-${reservation.id}">Estimated Cost: $0.00</p>
                                <button type="submit" class="save-update-btn" data-vehicle_id="${reservation.vehicle_id}" data-id="${reservation.id}">Save</button>
                                <button class="cancel-update-btn" data-id="${reservation.id}">Cancel</button>
                            </form>
                        </div>`;
            }).join('');
            document.getElementById('reservations').innerHTML = reservationsList;
            // Add event listeners for Update buttons
            const updateButtons = document.querySelectorAll('.update-res-btn');
            updateButtons.forEach(button => {
                button.addEventListener('click', handleUpdate);
            });

            // Add event listeners for Delete buttons
            const deleteButtons = document.querySelectorAll('.delete-res-btn');
            deleteButtons.forEach(button => {
                button.addEventListener('click', handleDelete);
            });

            // Add event listeners for Save buttons in the form
            const saveUpdateButtons = document.querySelectorAll('.save-update-btn');
            saveUpdateButtons.forEach(button => {
                button.addEventListener('click', handleSaveUpdate);
            });

            // Add event listeners for Cancel buttons in the form
            const cancelUpdateButtons = document.querySelectorAll('.cancel-update-btn');
            cancelUpdateButtons.forEach(button => {
                button.addEventListener('click', handleCancelUpdate);
            });
        } else {
            // Show message if no reservations
            document.getElementById('reservations').innerHTML = '<p>You have no active reservations.</p>';
        }
    })
    .catch(err => {
        console.error('Error fetching reservations:', err);
        alert('An error occurred while fetching your reservations.');
    });
}

// Handle update button click
function handleUpdate(event) {
    const reservationId = event.target.dataset.id;
    const vehicleId = event.target.dataset.vehicle_id;
    const updateForm = document.getElementById(`update-form-${reservationId}`);
    updateForm.style.display = 'block'; // Show the form to update reservation timings

    const updateButton = document.querySelector('.update-res-btn');
    updateButton.style.display = 'none';

    const startTimeInput = document.getElementById(`start-time-${reservationId}`);
    const endTimeInput = document.getElementById(`end-time-${reservationId}`);
    const costEstimateElement = document.getElementById(`cost-estimate-${reservationId}`);

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
                costEstimateElement.textContent = `Estimated Cost: $${(data.total_amount).toFixed(2)}`;
            })
            .catch(err => {
                console.error('Error fetching cost estimate:', err);
                costEstimateElement.textContent = `Error fetching cost`;
            });
        });
    });
}

// Handle save update button click
function handleSaveUpdate(event) {
    const reservationId = event.target.dataset.id;
    const vehicleId = event.target.dataset.vehicle_id;
    const newStartTime = document.getElementById(`start-time-${reservationId}`).value;
    const newEndTime = document.getElementById(`end-time-${reservationId}`).value;

    const startDate = new Date(newStartTime);
    const endDate = new Date(newEndTime);

    const startTimeISO = startDate.toISOString();
    const endTimeISO = endDate.toISOString();

    // Call API to update the reservation
    const token = localStorage.getItem('jwt');
    fetch(`http://localhost:8082/reservations/${reservationId}`, {
        method: 'PUT',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            start_time: startTimeISO,
            end_time: endTimeISO
        })
    })
    .then(response => response.json())
    .then(data => {
        console.log('Updated reservation:', data);
        updateInvoice(parseInt(reservationId), parseInt(vehicleId), startTimeISO, endTimeISO)
        
        fetchReservations();  // Refresh reservations list
    })
    .catch(err => {
        console.error('Error updating reservation:', err);
        alert('An error occurred while updating your reservation.');
    });
}

// Handle cancel update button click
function handleCancelUpdate(event) {
    const reservationId = event.target.dataset.id;
    const updateForm = document.getElementById(`update-form-${reservationId}`);
    updateForm.style.display = 'none'; // Hide the form

    // Show the update button again
    const updateButton = document.querySelector('.update-res-btn');
    updateButton.style.display = 'block';
}

// Handle delete button click
async function handleDelete(event) {
    const reservationId = event.target.dataset.id;
    const vehicleId = event.target.dataset.vehicle_id;
    const startTime = event.target.dataset.start_time;
    const endTime = event.target.dataset.end_time;
    console.log("timings: ", reservationId, startTime, endTime)

    // Confirm before deletion
    if (confirm('Are you sure you want to delete this reservation?')) {
        const token = localStorage.getItem('jwt');
        await fetch(`http://localhost:8082/reservations/${reservationId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        })
        .then(response => response.json())
        .then(data => {
            console.log('Deleted reservation:', data);

            //deleteInvoice(reservationId, token);
            console.log("Deleted invoice by reservation ID:", reservationId);

            CompleteOrCancelReservationHandler(reservationId, userId, vehicleId, startTime, endTime, 'Canceled')
            
            fetchReservations();  // Refresh reservations list
            
            fetchInvoicesByUser(userId);
        })
        .catch(err => {
            console.error('Error deleting reservation:', err);
            alert('An error occurred while deleting your reservation.');
        });

    }
}

async function fetchInvoicesByUser(userId) {
    await fetch(`http://localhost:8083/invoices?user_id=${userId}`)
        .then(response => response.json())
        .then(invoices => {
            console.log('Invoices:', invoices);
            // Optionally render the invoices to the UI
            const invoicesList = document.getElementById('invoices');
            invoicesList.innerHTML = ''; // Clear existing invoices
            if (!invoices || invoices.length === 0) {
                const noInvoicesMessage = document.createElement('div');
                noInvoicesMessage.textContent = 'No invoices found.';
                invoicesList.appendChild(noInvoicesMessage);
                return; // Exit function if no invoices are found
            }
            invoices.forEach(invoice => {
                const invoiceItem = document.createElement('div');
                invoiceItem.innerHTML = `<strong>Invoice ID: ${invoice.id}</strong><br>Vehicle ID: ${invoice.vehicle_id}<br>Vehicle: ${invoice.make} ${invoice.model}<br>Final Amount: $${invoice.total_amount}<br>Discount: $${invoice.discount}`;
                
                invoicesList.appendChild(invoiceItem);
            });
        })
        .catch(err => {
            console.error('Error fetching invoices:', err);
            alert('Error fetching invoices');
        });
}

// Function to delete the invoice associated with the reservation
async function deleteInvoice(reservationId, token) {
    await fetch(`http://localhost:8083/invoices?reservation_id=${reservationId}`, {
        method: 'DELETE',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        console.log("Deleted Data: ", data)
        if (data.message === 'Invoice deleted successfully') {
            console.log('Deleted invoice:', data);
            fetchReservations();  // Refresh reservations list
            alert('Reservation and associated invoice deleted successfully!');
        } else {
            alert('Error deleting invoice');
        }
    })
    .catch(err => {
        console.error('Error deleting invoice:', err);
        alert('An error occurred while deleting the invoice.');
    });
}

// Logout function: remove the token and redirect to the login page
document.getElementById('logout-btn').addEventListener('click', function() {
    localStorage.removeItem('jwt');  // Remove the token from local storage
    window.location.href = 'login.html';  // Redirect to login page
});

// Fetch user rental history
async function fetchRentalHistory(userId) {
    try {
        const response = await fetch(`http://localhost:8080/getRentals?user_id=${userId}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Failed to fetch rental history');
        }
        
        const rentalHistory = await response.json();
        displayRentalHistory(rentalHistory); // Display the fetched rental history
    } catch (err) {
        console.error('Error fetching rental history:', err);
        alert('An error occurred while fetching your rental history.');
    }
}

// Display rental history on the page
function displayRentalHistory(history) {
    const rentalHistoryList = document.getElementById('rental-history');
    rentalHistoryList.innerHTML = ''; // Clear previous content

    console.log("History: ", history)
    
    if (!history || history.length === 0) {
        rentalHistoryList.innerHTML = '<p>Your rental history is empty.</p>';
        return;
    }

    history.forEach(entry => {
        console.log("Entry: ", entry);
        const rentalEntryDiv = document.createElement('div');
        rentalEntryDiv.classList.add('rental-history-entry');
        
        rentalEntryDiv.innerHTML = `
            <p><strong>Vehicle ID:</strong> ${entry.vehicle_id}</p>
            <p><strong>Vehicle:</strong> ${entry.make} ${entry.model}</p>
            <p><strong>Reservation:</strong> ${entry.start_time} to ${entry.end_time}</p>
            <p><strong>Rental Cost:</strong> $${entry.total_amount}</p>
        `;
        
        rentalHistoryList.appendChild(rentalEntryDiv);
    });
}


async function CompleteOrCancelReservationHandler(reservationId, userId, vehicleId, startTime, endTime, rentalStatus) {
    // Rental history data
    const rentalHistoryData = {
        reservation_id: parseInt(reservationId),
        user_id: parseInt(userId),
        vehicle_id: parseInt(vehicleId),
        start_time: startTime,
        end_time: endTime,
        rental_status: rentalStatus
    };

    try {
        const response = await fetch('http://localhost:8080/addRentalEntry', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(rentalHistoryData)
        });

        if (!response.ok) {
            throw new Error('Failed to add rental history');
        }

        const data = await response.json();
        console.log('Rental history added:', data);

    } catch (err) {
        console.error('Error adding rental history:', err);
        alert('An error occurred while adding your rental history.');
    }
}

function updateInvoice(reservationId, vehicleId, startTime, endTime) {
    const invoiceData = {
        reservation_id: reservationId,
        user_id: userId,
        vehicle_id: vehicleId,
        start_time: startTime,
        end_time: endTime
    };

    // Send request to backend to generate the invoice
    fetch(`http://localhost:8083/updateinvoice?reservation_id=${reservationId}`, {
        method: 'PUT',
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
            alert("Error updating invoice.");
        }
    })
    .catch(err => {
        console.error('Error updating invoice:', err);
        alert('Error updating invoice');
    });
}
