function getOrder() {
    const orderId = document.getElementById('orderId').value.trim();
    if (!orderId) {
        showError('Please enter an Order UID');
        return;
    }

    // Reset UI
    hideError();
    hideSuccess();
    showLoading();
    hideOrderInfo();

    // Make API request
    fetch(`/api/order/${orderId}`)
        .then(response => {
            if (!response.ok) {
                return response.json().then(err => {
                    throw new Error(err.error || `Error: ${response.status} ${response.statusText}`);
                });
            }
            return response.json();
        })
        .then(order => {
            hideLoading();
            renderOrder(order);
            showSuccess(`Order ${orderId} loaded successfully`);
        })
        .catch(error => {
            hideLoading();
            showError(error.message);
        });
}

function renderOrder(order) {
    const container = document.getElementById('orderInfo');
    
    let html = `
        <div class="order-section">
            <h3>Order Information</h3>
            <div class="field">
                <strong>Order UID:</strong>
                <span class="field-value">${order.order_uid}</span>
            </div>
            <div class="field">
                <strong>Track Number:</strong>
                <span class="field-value">${order.track_number}</span>
            </div>
            <div class="field">
                <strong>Entry:</strong>
                <span class="field-value">${order.entry}</span>
            </div>
            <div class="field">
                <strong>Locale:</strong>
                <span class="field-value">${order.locale}</span>
            </div>
            <div class="field">
                <strong>Customer ID:</strong>
                <span class="field-value">${order.customer_id}</span>
            </div>
            <div class="field">
                <strong>Delivery Service:</strong>
                <span class="field-value">${order.delivery_service}</span>
            </div>
            <div class="field">
                <strong>Date Created:</strong>
                <span class="field-value">${order.date_created}</span>
            </div>
        </div>
        
        <div class="order-section">
            <h3>Delivery Information</h3>
            <div class="field">
                <strong>Name:</strong>
                <span class="field-value">${order.delivery.name}</span>
            </div>
            <div class="field">
                <strong>Phone:</strong>
                <span class="field-value">${order.delivery.phone}</span>
            </div>
            <div class="field">
                <strong>Zip Code:</strong>
                <span class="field-value">${order.delivery.zip}</span>
            </div>
            <div class="field">
                <strong>City:</strong>
                <span class="field-value">${order.delivery.city}</span>
            </div>
            <div class="field">
                <strong>Address:</strong>
                <span class="field-value">${order.delivery.address}</span>
            </div>
            <div class="field">
                <strong>Region:</strong>
                <span class="field-value">${order.delivery.region}</span>
            </div>
            <div class="field">
                <strong>Email:</strong>
                <span class="field-value">${order.delivery.email}</span>
            </div>
        </div>
        
        <div class="order-section">
            <h3>Payment Information</h3>
            <div class="field">
                <strong>Transaction ID:</strong>
                <span class="field-value">${order.payment.transaction}</span>
            </div>
            <div class="field">
                <strong>Currency:</strong>
                <span class="field-value">${order.payment.currency}</span>
            </div>
            <div class="field">
                <strong>Provider:</strong>
                <span class="field-value">${order.payment.provider}</span>
            </div>
            <div class="field">
                <strong>Amount:</strong>
                <span class="field-value">$${(order.payment.amount / 100).toFixed(2)}</span>
            </div>
            <div class="field">
                <strong>Payment Date:</strong>
                <span class="field-value">${new Date(order.payment.payment_dt * 1000).toLocaleString()}</span>
            </div>
            <div class="field">
                <strong>Bank:</strong>
                <span class="field-value">${order.payment.bank}</span>
            </div>
        </div>
        
        <div class="order-section">
            <h3>Items (${order.items.length})</h3>
            <div class="items-container">
    `;

    order.items.forEach(item => {
        html += `
            <div class="item-card">
                <div class="field">
                    <strong>Name:</strong>
                    <span class="field-value">${item.name}</span>
                </div>
                <div class="field">
                    <strong>Brand:</strong>
                    <span class="field-value">${item.brand}</span>
                </div>
                <div class="field">
                    <strong>Price:</strong>
                    <span class="field-value">$${(item.price / 100).toFixed(2)}</span>
                </div>
                <div class="field">
                    <strong>Sale:</strong>
                    <span class="field-value">${item.sale}%</span>
                </div>
                <div class="field">
                    <strong>Size:</strong>
                    <span class="field-value">${item.size}</span>
                </div>
                <div class="field">
                    <strong>Status:</strong>
                    <span class="field-value">${item.status}</span>
                </div>
                <div class="field">
                    <strong>Total Price:</strong>
                    <span class="field-value">$${(item.total_price / 100).toFixed(2)}</span>
                </div>
            </div>
        `;
    });

    html += `
            </div>
        </div>
    `;

    container.innerHTML = html;
    showOrderInfo();
}

// UI helper functions
function showLoading() {
    document.getElementById('loading').style.display = 'block';
}

function hideLoading() {
    document.getElementById('loading').style.display = 'none';
}

function showOrderInfo() {
    document.getElementById('orderInfo').style.display = 'block';
}

function hideOrderInfo() {
    document.getElementById('orderInfo').style.display = 'none';
}

function showError(message) {
    const errorContainer = document.getElementById('errorContainer');
    errorContainer.textContent = message;
    errorContainer.style.display = 'block';
}

function hideError() {
    document.getElementById('errorContainer').style.display = 'none';
}

function showSuccess(message) {
    const successContainer = document.getElementById('successMessage');
    successContainer.textContent = message;
    successContainer.style.display = 'block';
}

function hideSuccess() {
    document.getElementById('successMessage').style.display = 'none';
}

// Initialize the page
document.addEventListener('DOMContentLoaded', () => {
    // Example of a valid order ID for quick testing
    document.getElementById('orderId').value = 'b563feb7b2b84b6test';
});