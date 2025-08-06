let currentOrder = null;

async function getOrder() {
    const orderUid = document.getElementById('orderUid').value;
    if (!orderUid) {
        showError('Please enter Order UID');
        return;
    }
    
    showLoading(true);
    hideError();
    hideOrderDetails();
    hideJsonViewer();
    
    try {
        const response = await fetch(`/order/${orderUid}`);
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`HTTP error! Status: ${response.status}, Message: ${errorText}`);
        }
        
        currentOrder = await response.json();
        displayOrder(currentOrder);
        showOrderDetails();
        showLoading(false);
    } catch (error) {
        showError(`Error: ${error.message}`);
        showLoading(false);
    }
}

function displayOrder(order) {
    
    document.getElementById('orderInfo').innerHTML = `
        <div class="info-item">
            <strong>Order UID</strong>
            <span>${order.order_uid}</span>
        </div>
        <div class="info-item">
            <strong>Track Number</strong>
            <span>${order.track_number}</span>
        </div>
        <div class="info-item">
            <strong>Entry</strong>
            <span>${order.entry}</span>
        </div>
        <div class="info-item">
            <strong>Customer ID</strong>
            <span>${order.customer_id}</span>
        </div>
        <div class="info-item">
            <strong>Date Created</strong>
            <span>${new Date(order.date_created).toLocaleString()}</span>
        </div>
        <div class="info-item">
            <strong>Delivery Service</strong>
            <span>${order.delivery_service}</span>
        </div>
    `;
    
    const delivery = order.delivery;
    document.getElementById('deliveryInfo').innerHTML = `
        <div class="info-item">
            <strong>Name</strong>
            <span>${delivery.name}</span>
        </div>
        <div class="info-item">
            <strong>Phone</strong>
            <span>${delivery.phone}</span>
        </div>
        <div class="info-item">
            <strong>Email</strong>
            <span>${delivery.email}</span>
        </div>
        <div class="info-item">
            <strong>Address</strong>
            <span>${delivery.city}, ${delivery.address}</span>
        </div>
        <div class="info-item">
            <strong>Region</strong>
            <span>${delivery.region}</span>
        </div>
        <div class="info-item">
            <strong>ZIP Code</strong>
            <span>${delivery.zip}</span>
        </div>
    `;
    
    const payment = order.payment;
    document.getElementById('paymentInfo').innerHTML = `
        <div class="info-item">
            <strong>Transaction</strong>
            <span>${payment.transaction}</span>
        </div>
        <div class="info-item">
            <strong>Amount</strong>
            <span>$${(payment.amount / 100).toFixed(2)}</span>
        </div>
        <div class="info-item">
            <strong>Currency</strong>
            <span>${payment.currency}</span>
        </div>
        <div class="info-item">
            <strong>Provider</strong>
            <span>${payment.provider}</span>
        </div>
        <div class="info-item">
            <strong>Bank</strong>
            <span>${payment.bank}</span>
        </div>
        <div class="info-item">
            <strong>Payment Date</strong>
            <span>${new Date(payment.payment_dt * 1000).toLocaleString()}</span>
        </div>
    `;
    
    let itemsHtml = '';
    order.items.forEach(item => {
        itemsHtml += `
            <div class="item-card">
                <h3>${item.name} <span class="brand">(${item.brand})</span></h3>
                <div class="info-grid">
                    <div class="info-item">
                        <strong>Price</strong>
                        <span>$${(item.price / 100).toFixed(2)}</span>
                    </div>
                    <div class="info-item">
                        <strong>Sale</strong>
                        <span>${item.sale}%</span>
                    </div>
                    <div class="info-item">
                        <strong>Total</strong>
                        <span>$${(item.total_price / 100).toFixed(2)}</span>
                    </div>
                    <div class="info-item">
                        <strong>Size</strong>
                        <span>${item.size}</span>
                    </div>
                    <div class="info-item">
                        <strong>Status</strong>
                        <span class="status status-${item.status}">${item.status}</span>
                    </div>
                    <div class="info-item">
                        <strong>Brand</strong>
                        <span>${item.brand}</span>
                    </div>
                </div>
            </div>
        `;
    });
    
    document.getElementById('itemsList').innerHTML = itemsHtml;
    document.getElementById('itemsCount').textContent = order.items.length;
    
    document.getElementById('jsonViewer').textContent = 
        JSON.stringify(order, null, 2);
}

function toggleJson() {
    const viewer = document.getElementById('jsonViewer');
    viewer.style.display = viewer.style.display === 'none' ? 'block' : 'none';
}

function showLoading(show) {
    document.getElementById('loading').style.display = show ? 'flex' : 'none';
}

function showOrderDetails() {
    document.getElementById('orderDetails').classList.add('active');
}

function hideOrderDetails() {
    document.getElementById('orderDetails').classList.remove('active');
}

function showError(message) {
    const errorEl = document.getElementById('error');
    errorEl.textContent = message;
    errorEl.style.display = 'block';
}

function hideError() {
    document.getElementById('error').style.display = 'none';
}

function showJsonViewer() {
    document.getElementById('jsonViewer').style.display = 'block';
}

function hideJsonViewer() {
    document.getElementById('jsonViewer').style.display = 'none';
}

document.addEventListener('DOMContentLoaded', () => {
    getOrder();
    
    document.getElementById('orderUid').addEventListener('keyup', event => {
        if (event.key === 'Enter') {
            getOrder();
        }
    });
});