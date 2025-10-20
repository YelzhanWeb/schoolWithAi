// API Configuration
const API_URL = 'http://localhost:8080/api';

// Get auth token from localStorage
function getToken() {
    return localStorage.getItem('token');
}

// Get current user from localStorage
function getCurrentUser() {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
}

// Create axios instance with auth header
// frontend/js/api.js - ДОБАВИТЬ:
function createAxiosInstance() {
    const token = getToken();
    const instance = axios.create({
        baseURL: API_URL,
        headers: token ? { 'Authorization': `Bearer ${token}` } : {},
        timeout: 30000  // 30 секунд
    });
    
    // ДОБАВИТЬ INTERCEPTOR ДЛЯ ОШИБОК:
    instance.interceptors.response.use(
        response => response,
        error => {
            if (error.response?.status === 401) {
                // Токен истек - редирект на логин
                localStorage.removeItem('token');
                localStorage.removeItem('user');
                window.location.href = './login.html';
            } else if (error.code === 'ECONNABORTED') {
                showToast('Превышено время ожидания', 'error');
            } else if (!error.response) {
                showToast('Нет связи с сервером', 'error');
            }
            return Promise.reject(error);
        }
    );
    
    return instance;
}

// API Functions
const API = {
    // Auth
    auth: {
        login: (email, password) => 
            axios.post(`${API_URL}/auth/login`, { email, password }),
        
        register: (data) => 
            axios.post(`${API_URL}/auth/register`, data),
        
        me: () => 
            createAxiosInstance().get('/auth/me')
    },

    // Profile
    profile: {
        get: () => 
            createAxiosInstance().get('/profile'),
        
        create: (data) => 
            createAxiosInstance().post('/profile', data),
        
        update: (data) => 
            createAxiosInstance().put('/profile', data)
    },

    // Recommendations
recommendations: {
    get: () => 
        createAxiosInstance().get('/recommendations'),
    
    refresh: () => 
        createAxiosInstance().post('/recommendations/refresh'),
},

// Resources - ДОБАВИТЬ:
resources: {
    getById: (id) =>
        axios.get(`${API_URL}/resources/${id}`)
},

    // Courses
    courses: {
        getAll: () => 
            axios.get(`${API_URL}/courses`),
        
        getById: (id) => 
            axios.get(`${API_URL}/courses/${id}`),
        
        getModuleResources: (moduleId) => 
            axios.get(`${API_URL}/modules/${moduleId}/resources`)
    },

    // Progress
    progress: {
        get: () => 
            createAxiosInstance().get('/progress'),
        
        getStatistics: () => 
            createAxiosInstance().get('/progress/statistics'),
        
        update: (data) => 
            createAxiosInstance().post('/progress', data)
    }
};