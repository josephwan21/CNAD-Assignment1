function showForm(formType) {
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');
    const responseMessage = document.getElementById('responseMessage');
    responseMessage.textContent = ''; // Clear any previous messages

    if (formType === 'login') {
        loginForm.classList.add('active');
        registerForm.classList.remove('active');
    } else if (formType === 'register') {
        registerForm.classList.add('active');
        loginForm.classList.remove('active');
    }
}

document.getElementById('login').addEventListener('submit', async (event) => {
    event.preventDefault();
    const email = document.getElementById('loginEmail').value;
    const password = document.getElementById('loginPassword').value;

    try {
        const response = await fetch('http://localhost:8080/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password }),
        });

        const data = await response.json();
        if (response.ok) {
            alert('Login successful! Token: ' + data.token); // Replace with redirection or token handling
        } else {
            document.getElementById('responseMessage').textContent = data.error || 'Login failed';
        }
    } catch (error) {
        document.getElementById('responseMessage').textContent = 'An error occurred. Please try again.';
    }
});

document.getElementById('register').addEventListener('submit', async (event) => {
    event.preventDefault();
    const name = document.getElementById('registerName').value;
    const email = document.getElementById('registerEmail').value;
    const password = document.getElementById('registerPassword').value;

    try {
        const response = await fetch('http://localhost:8080/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, email, password }),
        });

        const data = await response.json();
        if (response.ok) {
            document.getElementById('responseMessage').textContent = data.message;
            document.getElementById('responseMessage').style.color = 'green';
        } else {
            document.getElementById('responseMessage').textContent = data.error || 'Registration failed';
            document.getElementById('responseMessage').style.color = 'red';
        }
    } catch (error) {
        document.getElementById('responseMessage').textContent = 'An error occurred. Please try again.';
    }
});
