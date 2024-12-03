

// Handle registration
document.getElementById('register-form').addEventListener('submit', function(event) {
    event.preventDefault();
    const name = document.getElementById('name').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    fetch('http://localhost:8080/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, email, password })
    })
    .then(response => response.json())
    .then(data => {
        if (data.message === 'User registered successfully!') {
            alert('Registration successful')
            window.location.href = 'login.html';
        } else {
            alert('Registration failed');
        }
    })
    .catch(err => console.error('Error:', err));
});
