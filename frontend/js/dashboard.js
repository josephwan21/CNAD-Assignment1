// Check if the user is authenticated by checking for the auth token
function checkAuth() {
    const token = localStorage.getItem('jwt');
    if (!token) {
        alert('You need to log in first');
        window.location.href = 'login.html';  // Redirect to login page if no token found
    } else {
        // If logged in, fetch user data
        fetchUserData(token);
    }
}

// Decode JWT and extract user ID
const token = localStorage.getItem('jwt');
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
        if (data.length > 0) {
            // Display reservations
            const reservationsList = data.map(reservation => {
                return `<div class="reservation">
                            <p>Vehicle: <strong>${reservation.make} ${reservation.model}</strong></p>
                            <p>Reservation from ${reservation.start_time} to ${reservation.end_time}</p>
                            <button class="update-res-btn" data-id="${reservation.id}">Update Reservation</button>
                            <button class="delete-res-btn" data-id="${reservation.id}">Delete Reservation</button>
                            <div class="update-form" id="update-form-${reservation.id}" style="display: none;">
                                <label for="start-time-${reservation.id}">Start Time:</label>
                                <input type="datetime-local" id="start-time-${reservation.id}">
                                <label for="end-time-${reservation.id}">End Time:</label>
                                <input type="datetime-local" id="end-time-${reservation.id}">
                                <button class="save-update-btn" data-id="${reservation.id}">Save</button>
                                <button class="cancel-update-btn" data-id="${reservation.id}">Cancel</button>
                            </div>
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
    const updateForm = document.getElementById(`update-form-${reservationId}`);
    updateForm.style.display = 'block'; // Show the form to update reservation timings

    const updateButton = document.querySelector('.update-res-btn');
    updateButton.style.display = 'none';
}

// Handle save update button click
function handleSaveUpdate(event) {
    const reservationId = event.target.dataset.id;
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
function handleDelete(event) {
    const reservationId = event.target.dataset.id;

    // Confirm before deletion
    if (confirm('Are you sure you want to delete this reservation?')) {
        const token = localStorage.getItem('jwt');
        fetch(`http://localhost:8082/reservations/${reservationId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        })
        .then(response => response.json())
        .then(data => {
            console.log('Deleted reservation:', data);
            fetchReservations();  // Refresh reservations list
        })
        .catch(err => {
            console.error('Error deleting reservation:', err);
            alert('An error occurred while deleting your reservation.');
        });
    }
}

function fetchInvoicesByUser(userId) {
    fetch(`http://localhost:8083/invoices?user_id=${userId}`)
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
                invoiceItem.innerHTML = `<strong>Invoice #${invoice.id}</strong><br>Vehicle ID: ${invoice.vehicle_id}<br>Vehicle: ${invoice.make} ${invoice.model}<br>Total Amount: $${invoice.total_amount}<br>Discount: $${invoice.discount}`;
                
                invoicesList.appendChild(invoiceItem);
            });
        })
        .catch(err => {
            console.error('Error fetching invoices:', err);
            alert('Error fetching invoices');
        });
}

// Logout function: remove the token and redirect to the login page
document.getElementById('logout-btn').addEventListener('click', function() {
    localStorage.removeItem('jwt');  // Remove the token from local storage
    window.location.href = 'login.html';  // Redirect to login page
});

// Call checkAuth on page load to verify authentication and fetch data
window.onload = checkAuth;
