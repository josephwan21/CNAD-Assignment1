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
            document.getElementById('membership-tier').innerText = data.membership;

            // Fetch and display user's reservations
            fetchReservations(data.id);
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
function fetchReservations(userId) {
    fetch(`http://localhost:8080/reservations?user_id=${userId}`)
        .then(response => response.json())
        .then(data => {
            if (data.reservations && data.reservations.length > 0) {
                // Display reservations
                const reservationsList = data.reservations.map(reservation => {
                    return `<div class="reservation">
                                <p>Vehicle: ${reservation.vehicle.make} ${reservation.vehicle.model}</p>
                                <p>Reservation from ${reservation.start_time} to ${reservation.end_time}</p>
                            </div>`;
                }).join('');
                document.getElementById('reservations').innerHTML = reservationsList;
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

// Logout function: remove the token and redirect to the login page
document.getElementById('logout-btn').addEventListener('click', function() {
    localStorage.removeItem('jwt');  // Remove the token from local storage
    window.location.href = 'login.html';  // Redirect to login page
});

// Call checkAuth on page load to verify authentication and fetch data
window.onload = checkAuth;
