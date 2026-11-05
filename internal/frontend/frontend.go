package frontend

const HTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>System Verleih Dashboard</title>
    <style>
        body { font-family: sans-serif; margin: 0; padding: 0; background-color: #f4f4f4; }
        .navbar { background-color: #333; color: white; padding: 1rem; display: flex; gap: 1rem; }
        .navbar a { color: white; text-decoration: none; cursor: pointer; font-weight: bold; }
        .navbar a:hover { text-decoration: underline; color: #ffeb3b; }
        .container { padding: 2rem; max-width: 1200px; margin: 0 auto; }

        .metric-cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1rem; margin-bottom: 2rem; }
        .card { background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); text-align: center; border-bottom: 4px solid #4caf50; }
        .card h3 { margin: 0 0 0.5rem 0; color: #666; font-size: 1rem; }
        .card p { margin: 0; font-size: 2rem; font-weight: bold; color: #333; }

        .actions { display: flex; gap: 1rem; margin-bottom: 2rem; }
        button { padding: 0.5rem 1rem; background-color: #2196f3; color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold; box-shadow: 0 2px 4px rgba(0,0,0,0.2); transition: background-color 0.3s; }
        button:hover { background-color: #1976d2; }

        .promo-banner {
            background-color: #ff9800;
            color: white;
            padding: 1rem;
            text-align: center;
            font-weight: bold;
            border-radius: 8px;
            margin-bottom: 2rem;
            animation: pulse 2s infinite;
        }
        @keyframes pulse {
            0% { transform: scale(1); }
            50% { transform: scale(1.02); }
            100% { transform: scale(1); }
        }

        .content-section { background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); display: none; }
        .content-section.active { display: block; }

        .item-list { display: flex; flex-direction: column; gap: 1rem; }
        .item-row { display: flex; align-items: center; gap: 1rem; padding: 1rem; border: 1px solid #ddd; border-radius: 4px; }
        .item-row svg { width: 50px; height: 50px; }

        /* Modal */
        .modal { display: none; position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.5); justify-content: center; align-items: center; }
        .modal.active { display: flex; }
        .modal-content { background: white; padding: 2rem; border-radius: 8px; width: 100%; max-width: 400px; }
        .modal-content h2 { margin-top: 0; }
        .form-group { margin-bottom: 1rem; }
        .form-group label { display: block; margin-bottom: 0.5rem; }
        .form-group input, .form-group select { width: 100%; padding: 0.5rem; border: 1px solid #ccc; border-radius: 4px; box-sizing: border-box; }
        .modal-actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 1.5rem; }
    </style>
</head>
<body>

    <nav class="navbar">
        <a onclick="showSection('dashboard')">Dashboard</a>
        <a onclick="showSection('items')">Items</a>
        <a onclick="showSection('customers')">Customers</a>
        <a onclick="showSection('rentals')">Rentals</a>
    </nav>

    <div class="container">

        <div class="promo-banner">
            🔥 Use voucher code "HALF" to get 50% off any item! Verified users get an extra 10% off their entire cart! 🔥
        </div>

        <div id="dashboard" class="content-section active">
            <h2>Dashboard</h2>

            <div class="metric-cards">
                <div class="card">
                    <h3>Total Items</h3>
                    <p id="stat-total">0</p>
                </div>
                <div class="card">
                    <h3>Available</h3>
                    <p id="stat-available">0</p>
                </div>
                <div class="card">
                    <h3>Rented</h3>
                    <p id="stat-rented">0</p>
                </div>
                <div class="card">
                    <h3>Registered Users</h3>
                    <p id="stat-users">0</p>
                </div>
            </div>

            <div class="actions">
                <button onclick="openModal('modal-cart')">Shopping Cart (<span id="cart-count">0</span>)</button>
                <button onclick="showSection('items')">Manage Items</button>
                <button onclick="openModal('modal-verify-id')">Manage Persons (Verify ID)</button>
            </div>
        </div>

        <div id="items" class="content-section">
            <h2>Inventory Items</h2>
            <div class="item-list" id="dynamic-item-list">
                <!-- Items will be injected here via JS -->
            </div>
        </div>

        <div id="customers" class="content-section">
            <h2>Verified Customers</h2>
            <div id="customer-list" class="item-list"></div>
        </div>

        <div id="rentals" class="content-section">
            <h2>Active Rentals</h2>
            <p>Active rentals would appear here.</p>
        </div>

    </div>

    <!-- Modals -->
    <div id="modal-verify-id" class="modal">
        <div class="modal-content">
            <h2>Verify Customer ID</h2>
            <div class="form-group">
                <label>Customer ID</label>
                <input type="text" id="vid-customer" value="cust_123">
            </div>
            <div class="form-group">
                <label>ID Number</label>
                <input type="text" id="vid-number" value="AB123456">
            </div>
            <div class="form-group">
                <label>ID Type</label>
                <select id="vid-type">
                    <option value="PASSPORT">Passport</option>
                    <option value="ID_CARD">ID Card</option>
                </select>
            </div>
            <!-- Dummy Base64 image for the payload -->
            <input type="hidden" id="vid-doc" value="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxMDAiIGhlaWdodD0iMTAwIj48cmVjdCB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgZmlsbD0iI2NjYyIvPjx0ZXh0IHg9IjEwIiB5PSI1MCI+SUQgRE9DPC90ZXh0Pjwvc3ZnPg==">

            <div class="modal-actions">
                <button onclick="closeModal('modal-verify-id')" style="background:#ccc; color:#333;">Cancel</button>
                <button onclick="submitVerifyID()">Verify</button>
            </div>
        </div>
    </div>

    <div id="modal-cart" class="modal">
        <div class="modal-content">
            <h2>Shopping Cart</h2>
            <div id="cart-items-container" class="item-list" style="margin-bottom:1rem; max-height:200px; overflow-y:auto;"></div>

            <div class="form-group" style="margin-top: 1rem; border-top: 1px solid #eee; padding-top: 1rem;">
                <label>Verified Customer ID (Optional for 10% off)</label>
                <input type="text" id="cart-customer-id" placeholder="e.g. cust_123" style="border: 2px solid #4caf50;">
            </div>

            <div class="modal-actions">
                <button onclick="closeModal('modal-cart')" style="background:#ccc; color:#333;">Close</button>
                <button onclick="submitCheckout()">Checkout</button>
            </div>
        </div>
    </div>

    <script>
        let cart = [];

        function showSection(id) {
            document.querySelectorAll('.content-section').forEach(el => el.classList.remove('active'));
            document.getElementById(id).classList.add('active');
            if (id === 'dashboard') {
                fetchStats();
            } else if (id === 'items') {
                fetchItems();
            } else if (id === 'customers') {
                fetchCustomers();
            }
        }

        function openModal(id) {
            document.getElementById(id).classList.add('active');
        }

        function closeModal(id) {
            document.getElementById(id).classList.remove('active');
        }

        function fetchStats() {
            fetch('/api/v1/stats')
                .then(res => res.json())
                .then(data => {
                    document.getElementById('stat-total').innerText = data.total_items || 0;
                    document.getElementById('stat-available').innerText = data.available_items || 0;
                    document.getElementById('stat-rented').innerText = data.rented_items || 0;
                    document.getElementById('stat-users').innerText = data.registered_users || 0;
                })
                .catch(err => console.error('Error fetching stats:', err));
        }

        function fetchItems() {
            fetch('/api/v1/items')
                .then(res => res.json())
                .then(data => {
                    const list = document.getElementById('dynamic-item-list');
                    list.innerHTML = '';
                    data.forEach(item => {
                        const isPro = item.is_pro ? "(PRO) " : "";
                        const svgSki = '<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100" fill="#eee"/><path d="M 20 80 L 80 20" stroke="#333" stroke-width="5"/><text x="35" y="85" font-size="12">SKI</text></svg>';
                        const svgBike = '<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100" fill="#eee"/><circle cx="30" cy="70" r="15" stroke="#333" stroke-width="3" fill="none"/><circle cx="70" cy="70" r="15" stroke="#333" stroke-width="3" fill="none"/><path d="M 30 70 L 50 40 L 70 70 M 50 40 L 20 40" stroke="#333" stroke-width="3" fill="none"/><text x="30" y="25" font-size="12">BIKE</text></svg>';
                        let icon = svgBike;
                        if (item.type === 'SKI') icon = svgSki;
                        else if (item.type === 'HELMET') icon = '<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100" fill="#eee"/><path d="M 20 80 Q 50 20 80 80" fill="none" stroke="#333" stroke-width="5"/><text x="25" y="50" font-size="12">HELMET</text></svg>';
                        else if (item.type === 'BOOTS') icon = '<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100" fill="#eee"/><path d="M 40 20 L 40 80 L 80 80" fill="none" stroke="#333" stroke-width="5"/><text x="30" y="50" font-size="12">BOOTS</text></svg>';
                        else if (item.type === 'GOGGLES') icon = '<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100" fill="#eee"/><ellipse cx="50" cy="50" rx="30" ry="15" fill="none" stroke="#333" stroke-width="5"/><text x="20" y="30" font-size="12">GOGGLES</text></svg>';

                        const status = item.is_rented ? 'Rented' : 'Available';

                        let btn = item.is_rented ? '' : '<button onclick="addToCart(\'' + item.id + '\', \'' + item.description + '\')">Add to Cart</button>';

                        list.innerHTML += '<div class="item-row">' +
                            icon +
                            '<div style="flex-grow:1;">' +
                            '<strong>' + isPro + item.description + ' (ID: ' + item.id + ')</strong>' +
                            '<br>' + status + ' | Base Rate: €' + item.hourly_rate.toFixed(2) + '/hr' +
                            '</div>' + btn + '</div>';
                    });
                })
                .catch(err => console.error('Error fetching items:', err));
        }

        function addToCart(id, name) {
            cart.push({ id: id, name: name, voucher: '', duration_hours: 24 });
            document.getElementById('cart-count').innerText = cart.length;
            renderCart();
            alert(name + ' added to cart.');
        }

        function renderCart() {
            const container = document.getElementById('cart-items-container');
            container.innerHTML = '';
            cart.forEach((cItem, index) => {
                container.innerHTML += '<div style="margin-bottom:0.5rem; display:flex; flex-wrap:wrap; gap:0.5rem; align-items:center; border-bottom:1px solid #eee; padding-bottom:0.5rem;">' +
                    '<div style="flex-basis:100%;"><strong>' + cItem.name + '</strong></div>' +
                    '<input type="number" placeholder="Hours" style="width:70px; border:2px solid #2196f3;" onchange="updateDuration(' + index + ', this.value)" value="' + cItem.duration_hours + '" min="1">' +
                    '<input type="text" placeholder="Voucher" style="width:100px; border:2px solid #ff9800;" onchange="updateVoucher(' + index + ', this.value)" value="' + cItem.voucher + '">' +
                    '<button onclick="removeFromCart(' + index + ')" style="background:#f44336; margin-left:auto;">Remove</button>' +
                    '</div>';
            });
        }

        function updateVoucher(index, val) {
            cart[index].voucher = val;
        }

        function updateDuration(index, val) {
            cart[index].duration_hours = parseFloat(val);
        }

        function fetchCustomers() {
            fetch('/api/v1/customers')
                .then(res => res.json())
                .then(data => {
                    const list = document.getElementById('customer-list');
                    list.innerHTML = '';
                    if (data && data.length > 0) {
                        data.forEach(c => {
                            list.innerHTML += '<div class="item-row"><div><strong>' + c.customer_id + '</strong><br>ID: ' + c.id_number + ' (' + c.id_type + ')</div></div>';
                        });
                    } else {
                        list.innerHTML = '<p>No verified customers found.</p>';
                    }
                })
                .catch(err => console.error('Error fetching customers:', err));
        }

        function removeFromCart(index) {
            cart.splice(index, 1);
            document.getElementById('cart-count').innerText = cart.length;
            renderCart();
        }

        function submitCheckout() {
            if (cart.length === 0) {
                alert('Cart is empty');
                return;
            }

            const custId = document.getElementById('cart-customer-id').value;
            const payload = {
                customer_id: custId,
                cart_items: cart.map(i => ({ id: i.id, voucher: i.voucher, duration_hours: i.duration_hours }))
            };

            fetch('/api/v1/checkout', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            })
            .then(res => res.json().then(data => ({ status: res.status, body: data })))
            .then(res => {
                if (res.status === 201) {
                    alert('Checkout successful! Total Price: €' + res.body.total_price.toFixed(2));
                    cart = [];
                    document.getElementById('cart-count').innerText = cart.length;
                    closeModal('modal-cart');
                    fetchStats();
                } else {
                    alert('Checkout failed: ' + res.body);
                }
            })
            .catch(err => console.error('Error during checkout:', err));
        }

        function submitVerifyID() {
            const payload = {
                customer_id: document.getElementById('vid-customer').value,
                id_number: document.getElementById('vid-number').value,
                id_type: document.getElementById('vid-type').value,
                document_base64: document.getElementById('vid-doc').value
            };

            fetch('/api/v1/customers/verify-id', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            })
            .then(res => {
                if (res.ok) {
                    alert('Customer ID verified successfully!');
                    closeModal('modal-verify-id');
                    fetchStats();
                } else {
                    alert('Verification failed.');
                }
            })
            .catch(err => console.error('Error verifying ID:', err));
        }

        // Initial fetch
        fetchStats();
    </script>
</body>
</html>`
