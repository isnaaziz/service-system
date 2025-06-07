// Global Variables
let token = "";
let currentUser = "";

// API Base URL
const API_BASE_URL = 'http://localhost:8800';

// Utility Functions
function showNotification(message, type = 'success') {
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.textContent = message;
    document.body.appendChild(notification);

    setTimeout(() => notification.classList.add('show'), 100);
    setTimeout(() => {
        notification.classList.remove('show');
        setTimeout(() => document.body.removeChild(notification), 300);
    }, 3000);
}

function updateStatus(connected, user = '') {
    const statusDot = document.getElementById('status-dot');
    const statusText = document.getElementById('status-text');
    const userInfo = document.getElementById('current-user');

    if (connected) {
        statusDot.classList.add('connected');
        statusText.textContent = 'Connected';
        userInfo.textContent = user;
    } else {
        statusDot.classList.remove('connected');
        statusText.textContent = 'Not Connected';
        userInfo.textContent = 'Not logged in';
    }
}

function showSection(sectionName) {
    // Hide all sections
    document.querySelectorAll('.section').forEach(section => {
        section.classList.remove('active');
    });

    // Remove active class from all nav buttons
    document.querySelectorAll('.nav-btn').forEach(btn => {
        btn.classList.remove('active');
    });

    // Show selected section
    document.getElementById(sectionName + '-section').classList.add('active');

    // Add active class to clicked button
    const activeButton = document.querySelector(`[data-section="${sectionName}"]`);
    if (activeButton) {
        activeButton.classList.add('active');
    }
}

function showResult(elementId, data, isError = false) {
    const element = document.getElementById(elementId);
    element.style.display = 'block';
    element.innerHTML = `<pre>${JSON.stringify(data, null, 2)}</pre>`;

    if (isError) {
        element.style.borderLeftColor = '#e74c3c';
    } else {
        element.style.borderLeftColor = '#27ae60';
    }
}

function clearForm(formPrefix) {
    document.querySelectorAll(`#${formPrefix}-section input`).forEach(input => {
        input.value = '';
    });
}

function validateEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

// API Functions
async function signup() {
    const username = document.getElementById('signup-username').value.trim();
    const email = document.getElementById('signup-email').value.trim();
    const password = document.getElementById('signup-password').value;

    // Validation
    if (!username || !email || !password) {
        showNotification('Please fill in all fields', 'error');
        return;
    }

    if (!validateEmail(email)) {
        showNotification('Please enter a valid email address', 'error');
        return;
    }

    if (password.length < 6) {
        showNotification('Password must be at least 6 characters long', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/signup`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, email, password })
        });

        const data = await response.json();
        showResult('signup-result', data, !response.ok);

        if (response.ok) {
            showNotification('Account created successfully!');
            clearForm('signup');
            // Automatically switch to login section
            setTimeout(() => showSection('login'), 1500);
        } else {
            showNotification(data.message || 'Signup failed', 'error');
        }
    } catch (error) {
        console.error('Signup error:', error);
        showResult('signup-result', { error: error.message }, true);
        showNotification('Connection error', 'error');
    }
}

async function login() {
    const username = document.getElementById('login-username').value.trim();
    const password = document.getElementById('login-password').value;

    if (!username || !password) {
        showNotification('Please fill in all fields', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });

        const data = await response.json();
        showResult('login-result', data, !response.ok);

        if (response.ok && data.token) {
            token = data.token;
            currentUser = username;
            updateStatus(true, username);
            showNotification('Login successful!');
            clearForm('login');

            // Auto-load users after successful login
            setTimeout(() => {
                showSection('users');
                getUsers();
            }, 1000);
        } else {
            showNotification(data.message || 'Login failed', 'error');
        }
    } catch (error) {
        console.error('Login error:', error);
        showResult('login-result', { error: error.message }, true);
        showNotification('Connection error', 'error');
    }
}

async function logout() {
    if (!token) {
        showNotification('Not logged in', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/logout`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token
            },
            body: JSON.stringify({ token })
        });

        // Reset state regardless of response
        token = "";
        currentUser = "";
        updateStatus(false);
        showNotification('Logged out successfully!');
        showSection('login');

        // Clear all result sections
        document.querySelectorAll('.result-container').forEach(container => {
            container.style.display = 'none';
        });

        // Clear users grid
        document.getElementById('users-grid').innerHTML = '';

    } catch (error) {
        console.error('Logout error:', error);
        showNotification('Logout error', 'error');

        // Force logout even if request fails
        token = "";
        currentUser = "";
        updateStatus(false);
        showSection('login');
    }
}

async function getUsers() {
    if (!token) {
        showNotification('Please login first', 'error');
        showSection('login');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/users`, {
            headers: { 'Authorization': 'Bearer ' + token }
        });

        const data = await response.json();

        if (response.ok) {
            if (Array.isArray(data)) {
                displayUsers(data);
                showNotification(`Loaded ${data.length} users successfully!`);
            } else {
                showResult('users-result', data, false);
                showNotification('Users loaded successfully!');
            }
        } else {
            showResult('users-result', data, true);
            showNotification(data.message || 'Failed to load users', 'error');

            // If unauthorized, redirect to login
            if (response.status === 401) {
                token = "";
                currentUser = "";
                updateStatus(false);
                showSection('login');
            }
        }
    } catch (error) {
        console.error('Get users error:', error);
        showResult('users-result', { error: error.message }, true);
        showNotification('Connection error', 'error');
    }
}

function displayUsers(users) {
    const usersGrid = document.getElementById('users-grid');
    const usersResult = document.getElementById('users-result');

    // Hide raw JSON result when displaying cards
    usersResult.style.display = 'none';

    usersGrid.innerHTML = '';

    if (users.length === 0) {
        usersGrid.innerHTML = '<div class="user-card"><p>No users found</p></div>';
        return;
    }

    users.forEach((user, index) => {
        const userCard = document.createElement('div');
        userCard.className = 'user-card';
        userCard.style.animationDelay = `${index * 0.1}s`;

        userCard.innerHTML = `
            <h3>${user.username || 'N/A'}</h3>
            <p><strong>Email:</strong> ${user.email || 'N/A'}</p>
            <p><strong>ID:</strong> ${user.id || 'N/A'}</p>
            <p><strong>Created:</strong> ${user.created_at ? new Date(user.created_at).toLocaleDateString() : 'N/A'}</p>
            ${user.username === currentUser ? '<p><strong>🌟 Current User</strong></p>' : ''}
        `;

        usersGrid.appendChild(userCard);
    });
}

async function changePassword() {
    if (!token) {
        showNotification('Please login first', 'error');
        showSection('login');
        return;
    }

    const oldPassword = document.getElementById('old-password').value;
    const newPassword = document.getElementById('new-password').value;

    if (!oldPassword || !newPassword) {
        showNotification('Please fill in all fields', 'error');
        return;
    }

    if (newPassword.length < 6) {
        showNotification('New password must be at least 6 characters long', 'error');
        return;
    }

    if (oldPassword === newPassword) {
        showNotification('New password must be different from old password', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/users/change-password`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token
            },
            body: JSON.stringify({
                old_password: oldPassword,
                new_password: newPassword
            })
        });

        const data = await response.json();
        showResult('change-password-result', data, !response.ok);

        if (response.ok) {
            showNotification('Password changed successfully!');
            clearForm('password');
        } else {
            showNotification(data.message || 'Password change failed', 'error');

            // If unauthorized, redirect to login
            if (response.status === 401) {
                token = "";
                currentUser = "";
                updateStatus(false);
                showSection('login');
            }
        }
    } catch (error) {
        console.error('Change password error:', error);
        showResult('change-password-result', { error: error.message }, true);
        showNotification('Connection error', 'error');
    }
}

// Event Listeners
document.addEventListener('DOMContentLoaded', function () {
    // Initialize dashboard
    updateStatus(false);

    // Add Enter key support for forms
    document.addEventListener('keypress', function (e) {
        if (e.key === 'Enter') {
            const activeSection = document.querySelector('.section.active');
            const target = e.target;

            // Prevent if target is a button to avoid double execution
            if (target.tagName.toLowerCase() === 'button') {
                return;
            }

            if (activeSection.id === 'login-section') {
                login();
            } else if (activeSection.id === 'signup-section') {
                signup();
            } else if (activeSection.id === 'password-section') {
                changePassword();
            }
        }
    });

    // Add input validation feedback
    const inputs = document.querySelectorAll('input');
    inputs.forEach(input => {
        input.addEventListener('blur', function () {
            if (this.type === 'email' && this.value && !validateEmail(this.value)) {
                this.style.borderColor = '#e74c3c';
            } else if (this.type === 'password' && this.value && this.value.length < 6) {
                this.style.borderColor = '#e74c3c';
            } else if (this.value) {
                this.style.borderColor = '#27ae60';
            } else {
                this.style.borderColor = '#ddd';
            }
        });

        input.addEventListener('focus', function () {
            this.style.borderColor = '#3498db';
        });
    });
});

// Auto-refresh users every 30 seconds if logged in and on users section
setInterval(() => {
    const activeSection = document.querySelector('.section.active');
    if (token && activeSection && activeSection.id === 'users-section') {
        getUsers();
    }
}, 30000);

// Handle connection errors globally
window.addEventListener('online', () => {
    showNotification('Connection restored', 'success');
});

window.addEventListener('offline', () => {
    showNotification('Connection lost', 'error');
}); q