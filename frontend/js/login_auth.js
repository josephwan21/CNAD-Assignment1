// Handle login
document.getElementById('login-form').addEventListener('submit', function(event) {
    event.preventDefault();
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    fetch('http://localhost:8080/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
    })
    .then(response => {
        if (!response.ok) {
            if (response.status === 401) {
                alert('Invalid credentials');
            } else {
                alert('Server error, please try again later.');
            }
        }
        return response.json();
    })
    .then(data => {
        if (data.token) {
            localStorage.setItem('jwt', data.token);
            window.location.href = 'dashboard.html';
        } else {
            alert('Invalid credentials');
        }
    })
    .catch(err => console.error('Error:', err));
});