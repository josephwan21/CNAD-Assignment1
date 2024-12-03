// billing.js

// Fetch invoice data
fetch('http://localhost:8082/billing')
    .then(response => response.json())
    .then(billing => {
        const invoiceDetails = document.getElementById('invoice-details');
        invoiceDetails.innerHTML = `
            <p>Total Amount: $${billing.total_amount}</p>
            <p>Discount: $${billing.discount}</p>
            <p>Final Amount: $${billing.total_amount - billing.discount}</p>
        `;
    })
    .catch(err => console.error('Error:', err));
