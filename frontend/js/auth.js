// Check if user is authenticated
function isAuthenticated() {
    return localStorage.getItem('token') !== null;
}

// Check if user is student
function isStudent() {
    const user = getCurrentUser();
    return user && user.role === 'student';
}

// Check if user is teacher
function isTeacher() {
    const user = getCurrentUser();
    return user && user.role === 'teacher';
}

// Logout
function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/login.html';
}

// Protect page (redirect to login if not authenticated)
function protectPage() {
    if (!isAuthenticated()) {
        window.location.href = '/login.html';
        return false;
    }
    return true;
}

// Protect student page
function protectStudentPage() {
    if (!protectPage()) return false;
    
    if (!isStudent()) {
        window.location.href = '/teacher.html';
        return false;
    }
    return true;
}

// Check profile existence
async function checkProfile() {
    try {
        const response = await API.profile.get();
        return response.data !== null;
    } catch (error) {
        if (error.response && error.response.status === 404) {
            return false;
        }
        console.error('Error checking profile:', error);
        return false;
    }
}

// Initialize auth on page load
async function initAuth() {
    if (!protectStudentPage()) return;

    const user = getCurrentUser();
    
    // Update navbar with user info
    const userNameElement = document.getElementById('userName');
    if (userNameElement) {
        userNameElement.textContent = user.full_name;
    }

    // Check if profile exists
    const hasProfile = await checkProfile();
    
    // Redirect to profile setup if needed (except on profile-setup page)
    if (!hasProfile && !window.location.pathname.includes('profile-setup')) {
        window.location.href = '/profile-setup.html';
    }
}