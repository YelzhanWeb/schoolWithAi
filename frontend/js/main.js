// Format time (seconds to minutes)
function formatTime(seconds) {
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) {
        return `${minutes} мин`;
    }
    const hours = Math.floor(minutes / 60);
    const remainingMinutes = minutes % 60;
    return `${hours}ч ${remainingMinutes}м`;
}

// Format difficulty level
function getDifficultyBadge(level) {
    const badges = {
        1: '<span class="px-2 py-1 bg-green-100 text-green-800 text-xs rounded">Легко</span>',
        2: '<span class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">Начальный</span>',
        3: '<span class="px-2 py-1 bg-yellow-100 text-yellow-800 text-xs rounded">Средний</span>',
        4: '<span class="px-2 py-1 bg-orange-100 text-orange-800 text-xs rounded">Сложный</span>',
        5: '<span class="px-2 py-1 bg-red-100 text-red-800 text-xs rounded">Эксперт</span>'
    };
    return badges[level] || badges[3];
}

// Get resource type icon
function getResourceTypeIcon(type) {
    const icons = {
        'video': '🎥',
        'reading': '📖',
        'exercise': '✍️',
        'quiz': '📝',
        'interactive': '🎮'
    };
    return icons[type] || '📚';
}

// Show loading spinner
function showLoading(elementId) {
    const element = document.getElementById(elementId);
    if (element) {
        element.innerHTML = `
            <div class="flex justify-center items-center py-12">
                <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
            </div>
        `;
    }
}

// Show error message
function showError(elementId, message) {
    const element = document.getElementById(elementId);
    if (element) {
        element.innerHTML = `
            <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
                <p>${message}</p>
            </div>
        `;
    }
}

// Show success toast
function showToast(message, type = 'success') {
    const bgColor = type === 'success' ? 'bg-green-500' : 'bg-red-500';
    
    const toast = document.createElement('div');
    toast.className = `fixed top-4 right-4 ${bgColor} text-white px-6 py-3 rounded-lg shadow-lg z-50 animate-fade-in`;
    toast.textContent = message;
    
    document.body.appendChild(toast);
    
    setTimeout(() => {
        toast.remove();
    }, 3000);
}

// Create rating stars
function createStars(rating, onChange) {
    let html = '<div class="flex space-x-1">';
    for (let i = 1; i <= 5; i++) {
        const filled = i <= rating;
        html += `
            <button 
                onclick="${onChange}(${i})" 
                class="text-2xl ${filled ? 'text-yellow-400' : 'text-gray-300'} hover:text-yellow-400 transition">
                ★
            </button>
        `;
    }
    html += '</div>';
    return html;
}

// Calculate level from points
function calculateLevel(points) {
    return Math.floor(points / 100) + 1;
}

// Calculate progress to next level
function calculateLevelProgress(points) {
    const currentLevelPoints = points % 100;
    return (currentLevelPoints / 100) * 100;
}