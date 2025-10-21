class SysPulseMonitor {
    constructor() {
        this.version = '1.0.0';
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.isConnected = false;
        
        // История метрик для графиков (60 значений)
        this.metricsHistory = {
            cpu: [],
            memory: [],
            disk: []
        };
        
        this.maxHistoryLength = 60;
        this.charts = {};
        
        // Настройки оповещений по умолчанию
        this.alertConfig = {
            cpu_threshold: 80,
            ram_threshold: 85,
            disk_threshold: 90,
            enabled: true
        };
        
        this.init();
    }

    async init() {
        console.log('🚀 Инициализация SysPulse Monitor...');
        console.log('📊 Подготовка графиков...');
        
        // Загружаем настройки оповещений
        await this.loadAlertConfig();
        
        // Инициализируем графики
        this.initCharts();
        
        // Загружаем версию приложения
        await this.loadVersion();
        
        // Подключаемся к WebSocket
        this.connectWebSocket();
        
        this.updateStatus('connecting', 'Подключение...');
    }

    initCharts() {
        const chartConfig = {
            type: 'line',
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: { display: false },
                    tooltip: { enabled: false }
                },
                scales: {
                    x: { display: false },
                    y: { 
                        display: false,
                        min: 0,
                        max: 100
                    }
                },
                elements: {
                    point: { radius: 0 },
                    line: { tension: 0.4 }
                },
                animation: { duration: 0 }
            }
        };

        // График CPU
        this.charts.cpu = new Chart(
            document.getElementById('cpu-chart'),
            {
                ...chartConfig,
                data: {
                    labels: Array(this.maxHistoryLength).fill(''),
                    datasets: [{
                        data: Array(this.maxHistoryLength).fill(0),
                        borderColor: '#4CAF50',
                        backgroundColor: 'rgba(76, 175, 80, 0.1)',
                        borderWidth: 2,
                        fill: true
                    }]
                }
            }
        );

        // График Memory
        this.charts.memory = new Chart(
            document.getElementById('memory-chart'),
            {
                ...chartConfig,
                data: {
                    labels: Array(this.maxHistoryLength).fill(''),
                    datasets: [{
                        data: Array(this.maxHistoryLength).fill(0),
                        borderColor: '#2196F3',
                        backgroundColor: 'rgba(33, 150, 243, 0.1)',
                        borderWidth: 2,
                        fill: true
                    }]
                }
            }
        );

        // График Disk
        this.charts.disk = new Chart(
            document.getElementById('disk-chart'),
            {
                ...chartConfig,
                data: {
                    labels: Array(this.maxHistoryLength).fill(''),
                    datasets: [{
                        data: Array(this.maxHistoryLength).fill(0),
                        borderColor: '#FF9800',
                        backgroundColor: 'rgba(255, 152, 0, 0.1)',
                        borderWidth: 2,
                        fill: true
                    }]
                }
            }
        );
    }

    connectWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;
        
        console.log(`🔌 Попытка подключения к: ${wsUrl}`);
        
        this.ws = new WebSocket(wsUrl);
        
        this.ws.onopen = () => {
            console.log('✅ WebSocket подключен');
            this.isConnected = true;
            this.reconnectAttempts = 0;
            this.updateStatus('connected', 'Подключено (WebSocket)');
        };
        
        this.ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                
                // Обновляем метрики
                this.updateMetricsDisplay(data);
                this.updateCharts(data);
                this.updateLastUpdateTime();
                
                // Обрабатываем оповещения если они есть
                if (data.alerts && data.alerts.length > 0) {
                    this.handleAlerts(data.alerts);
                }
                
            } catch (error) {
                console.error('❌ Ошибка парсинга метрик:', error);
            }
        };
        
        this.ws.onclose = (event) => {
            console.log(`🔌 WebSocket отключен: код=${event.code}, причина=${event.reason}`);
            this.isConnected = false;
            this.updateStatus('disconnected', `Отключено (код: ${event.code})`);
            this.handleReconnection();
        };
        
        this.ws.onerror = (error) => {
            console.error('❌ WebSocket ошибка:', error);
            this.updateStatus('error', 'Ошибка подключения');
        };
    }

    // Загружаем настройки оповещений с сервера
    async loadAlertConfig() {
        try {
            const response = await fetch('/api/alerts/config');
            this.alertConfig = await response.json();
            console.log('⚙️ Загружены настройки оповещений:', this.alertConfig);
        } catch (error) {
            console.error('Ошибка загрузки настроек оповещений:', error);
        }
    }

    updateMetricsDisplay(metrics) {
        // Обновляем CPU с цветовой индикацией
        const cpuElement = document.getElementById('cpu-usage');
        cpuElement.textContent = `${metrics.cpu.usage.toFixed(1)}%`;
        this.setMetricColor(cpuElement, metrics.cpu.usage, 80, 95);

        document.getElementById('cpu-details').textContent = 
            `${metrics.cpu.cores} ядер | Load: ${metrics.cpu.load1.toFixed(2)}`;

        // Обновляем Memory
        const memUsedGB = metrics.memory.used / 1024 / 1024 / 1024;
        const memTotalGB = metrics.memory.total / 1024 / 1024 / 1024;
        const memElement = document.getElementById('memory-usage');
        memElement.textContent = `${metrics.memory.usage.toFixed(1)}%`;
        this.setMetricColor(memElement, metrics.memory.usage, 85, 95);

        document.getElementById('memory-details').textContent = 
            `${memUsedGB.toFixed(2)}GB / ${memTotalGB.toFixed(2)}GB`;

        // Обновляем Disk
        const diskUsedGB = metrics.disk.used / 1024 / 1024 / 1024;
        const diskTotalGB = metrics.disk.total / 1024 / 1024 / 1024;
        const diskElement = document.getElementById('disk-usage');
        diskElement.textContent = `${metrics.disk.usage.toFixed(1)}%`;
        this.setMetricColor(diskElement, metrics.disk.usage, 90, 95);

        document.getElementById('disk-details').textContent = 
            `${diskUsedGB.toFixed(1)}GB / ${diskTotalGB.toFixed(1)}GB`;
    }

    // Обновляем графики новыми данными
    updateCharts(metrics) {
        // Добавляем новые данные в историю
        this.addToHistory('cpu', metrics.cpu.usage);
        this.addToHistory('memory', metrics.memory.usage);
        this.addToHistory('disk', metrics.disk.usage);

        // Обновляем все графики
        this.updateChart('cpu');
        this.updateChart('memory');
        this.updateChart('disk');
    }

    addToHistory(type, value) {
        this.metricsHistory[type].push(value);
        if (this.metricsHistory[type].length > this.maxHistoryLength) {
            this.metricsHistory[type].shift();
        }
    }

    updateChart(type) {
        if (this.charts[type] && this.metricsHistory[type].length > 0) {
            this.charts[type].data.datasets[0].data = [...this.metricsHistory[type]];
            this.charts[type].update('none');
        }
    }

    // Обрабатываем новые оповещения
    handleAlerts(alerts) {
        console.log('⚠️ Получены оповещения:', alerts);
        
        // Показываем оповещения в интерфейсе
        this.updateAlertsDisplay(alerts);
        
        // Воспроизводим звуковые оповещения
        this.playAlertSound(alerts);
    }

    // Обновляем отображение оповещений
    updateAlertsDisplay(alerts) {
        const container = document.getElementById('alerts-container');
        
        if (alerts.length === 0) {
            container.innerHTML = '<div class="alert-message">✅ Оповещений нет</div>';
            return;
        }

        container.innerHTML = alerts.map(alert => `
            <div class="alert-${alert.level}">
                <strong>${alert.type.toUpperCase()}</strong>: ${alert.message}
                <small>(${new Date(alert.timestamp).toLocaleTimeString()})</small>
            </div>
        `).join('');
    }

    // Воспроизводим звук оповещения
    playAlertSound(alerts) {
        if (!this.alertConfig.enabled) return;

        // Ищем критические оповещения
        const criticalAlerts = alerts.filter(alert => alert.level === 'critical');
        
        if (criticalAlerts.length > 0) {
            // Критическое оповещение - более срочный звук
            this.createBeep(800, 200);
        } else if (alerts.length > 0) {
            // Обычное предупреждение
            this.createBeep(400, 100);
        }
    }

    // Создаем звуковой сигнал
    createBeep(frequency, duration) {
        try {
            const audioContext = new (window.AudioContext || window.webkitAudioContext)();
            const oscillator = audioContext.createOscillator();
            const gainNode = audioContext.createGain();
            
            oscillator.connect(gainNode);
            gainNode.connect(audioContext.destination);
            
            oscillator.frequency.value = frequency;
            oscillator.type = 'sine';
            
            gainNode.gain.setValueAtTime(0.3, audioContext.currentTime);
            gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + duration / 1000);
            
            oscillator.start(audioContext.currentTime);
            oscillator.stop(audioContext.currentTime + duration / 1000);
        } catch (error) {
            console.log('🔇 Звуковые оповещения не поддерживаются');
        }
    }

    // Устанавливаем цвет в зависимости от значения
    setMetricColor(element, value, warningThreshold, criticalThreshold) {
        element.className = 'metric-value ';
        if (value >= criticalThreshold) {
            element.classList.add('metric-critical');
        } else if (value >= warningThreshold) {
            element.classList.add('metric-warning');
        } else {
            element.classList.add('metric-normal');
        }
    }

    handleReconnection() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            const delay = Math.min(1000 * this.reconnectAttempts, 10000);
            
            console.log(`🔄 Попытка переподключения ${this.reconnectAttempts}/${this.maxReconnectAttempts} через ${delay}мс`);
            this.updateStatus('reconnecting', `Переподключение... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
            
            setTimeout(() => {
                this.connectWebSocket();
            }, delay);
        } else {
            console.error('❌ Превышено максимальное количество попыток переподключения');
            this.updateStatus('error', 'Не удалось подключиться');
        }
    }

    async loadVersion() {
        try {
            const response = await fetch('/api/version');
            const data = await response.json();
            document.getElementById('version').textContent = data.version;
        } catch (error) {
            console.error('Ошибка загрузки версии:', error);
        }
    }

    updateLastUpdateTime() {
        const now = new Date();
        document.getElementById('lastUpdate').textContent = now.toLocaleTimeString();
    }

    updateStatus(status, message) {
        const statusElement = document.getElementById('status');
        statusElement.textContent = message;
        
        statusElement.className = 'status-value';
        if (status === 'connected') {
            statusElement.style.color = '#4CAF50';
        } else if (status === 'error' || status === 'disconnected') {
            statusElement.style.color = '#f44336';
        } else if (status === 'connecting' || status === 'reconnecting') {
            statusElement.style.color = '#ff9800';
        } else {
            statusElement.style.color = '#2196F3';
        }
    }

    // Показываем модальное окно с настройками оповещений
    showAlertSettings() {
        const modal = this.createAlertSettingsModal();
        document.body.appendChild(modal);
    }

    // Создаем модальное окно настроек
    createAlertSettingsModal() {
        const modal = document.createElement('div');
        modal.className = 'modal-overlay';
        modal.innerHTML = `
            <div class="modal-content">
                <div class="modal-header">
                    <h3>⚙️ Настройки оповещений</h3>
                    <button class="modal-close" onclick="this.closest('.modal-overlay').remove()">×</button>
                </div>
                <div class="modal-body">
                    <div class="setting-group">
                        <label>
                            <input type="checkbox" id="alerts-enabled" ${this.alertConfig.enabled ? 'checked' : ''}>
                            Включить оповещения
                        </label>
                    </div>
                    
                    <div class="setting-group">
                        <label>Порог CPU:</label>
                        <input type="range" id="cpu-threshold" min="50" max="95" value="${this.alertConfig.cpu_threshold}" 
                               oninput="document.getElementById('cpu-threshold-value').textContent = this.value + '%'">
                        <span id="cpu-threshold-value">${this.alertConfig.cpu_threshold}%</span>
                    </div>
                    
                    <div class="setting-group">
                        <label>Порог памяти:</label>
                        <input type="range" id="ram-threshold" min="50" max="95" value="${this.alertConfig.ram_threshold}"
                               oninput="document.getElementById('ram-threshold-value').textContent = this.value + '%'">
                        <span id="ram-threshold-value">${this.alertConfig.ram_threshold}%</span>
                    </div>
                    
                    <div class="setting-group">
                        <label>Порог диска:</label>
                        <input type="range" id="disk-threshold" min="50" max="95" value="${this.alertConfig.disk_threshold}"
                               oninput="document.getElementById('disk-threshold-value').textContent = this.value + '%'">
                        <span id="disk-threshold-value">${this.alertConfig.disk_threshold}%</span>
                    </div>
                </div>
                <div class="modal-footer">
                    <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Отмена</button>
                    <button class="btn btn-primary" onclick="window.monitor.saveAlertSettings()">Сохранить</button>
                </div>
            </div>
        `;
        
        return modal;
    }

    // Сохраняем настройки оповещений
    async saveAlertSettings() {
        const newConfig = {
            enabled: document.getElementById('alerts-enabled').checked,
            cpu_threshold: parseInt(document.getElementById('cpu-threshold').value),
            ram_threshold: parseInt(document.getElementById('ram-threshold').value),
            disk_threshold: parseInt(document.getElementById('disk-threshold').value)
        };

        try {
            const response = await fetch('/api/alerts/config', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(newConfig)
            });

            if (response.ok) {
                this.alertConfig = newConfig;
                console.log('✅ Настройки оповещений сохранены');
                
                // Закрываем модальное окно
                document.querySelector('.modal-overlay').remove();
                
                // Показываем уведомление
                this.showNotification('Настройки сохранены', 'success');
            }
        } catch (error) {
            console.error('❌ Ошибка сохранения настроек:', error);
            this.showNotification('Ошибка сохранения', 'error');
        }
    }

    // Загружаем историю оповещений
    async loadAlertHistory() {
        try {
            const response = await fetch('/api/alerts/history');
            const history = await response.json();
            this.showAlertHistoryModal(history);
        } catch (error) {
            console.error('Ошибка загрузки истории оповещений:', error);
        }
    }

    // Показываем историю оповещений
    showAlertHistoryModal(history) {
        const modal = document.createElement('div');
        modal.className = 'modal-overlay';
        
        const alertsHTML = history.alerts.map(alert => `
            <div class="alert-history-item alert-${alert.level}">
                <div class="alert-time">${new Date(alert.timestamp).toLocaleString()}</div>
                <div class="alert-message">${alert.message}</div>
                <div class="alert-value">${alert.value.toFixed(1)}% (порог: ${alert.threshold}%)</div>
            </div>
        `).join('');
        
        modal.innerHTML = `
            <div class="modal-content" style="max-width: 600px;">
                <div class="modal-header">
                    <h3>📋 История оповещений</h3>
                    <button class="modal-close" onclick="this.closest('.modal-overlay').remove()">×</button>
                </div>
                <div class="modal-body">
                    <div class="alert-stats">
                        <div class="stat-item">Всего: ${history.stats.total_alerts}</div>
                        <div class="stat-item">Активных: ${history.stats.active_alerts}</div>
                        <div class="stat-item">Сегодня: ${history.stats.today_alerts}</div>
                    </div>
                    <div class="alert-history-list">
                        ${alertsHTML.length > 0 ? alertsHTML : '<div class="no-alerts">Нет оповещений в истории</div>'}
                    </div>
                </div>
                <div class="modal-footer">
                    <button class="btn btn-secondary" onclick="window.monitor.clearAlertHistory()">Очистить историю</button>
                    <button class="btn btn-primary" onclick="this.closest('.modal-overlay').remove()">Закрыть</button>
                </div>
            </div>
        `;
        
        document.body.appendChild(modal);
    }

    // Очищаем историю оповещений
    async clearAlertHistory() {
        try {
            const response = await fetch('/api/alerts/clear', { method: 'POST' });
            if (response.ok) {
                this.showNotification('История оповещений очищена', 'success');
                document.querySelector('.modal-overlay').remove();
            }
        } catch (error) {
            console.error('Ошибка очистки истории:', error);
            this.showNotification('Ошибка очистки', 'error');
        }
    }

    // Показываем уведомление
    showNotification(message, type = 'info') {
        const notification = document.createElement('div');
        notification.className = `notification notification-${type}`;
        notification.textContent = message;
        notification.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 12px 20px;
            background: ${type === 'success' ? '#4CAF50' : type === 'error' ? '#f44336' : '#2196F3'};
            color: white;
            border-radius: 5px;
            z-index: 1000;
            animation: slideIn 0.3s ease;
        `;
        
        document.body.appendChild(notification);
        
        setTimeout(() => {
            notification.remove();
        }, 3000);
    }

    disconnect() {
        if (this.ws) {
            this.ws.close();
        }
    }
}

// Добавляем CSS анимацию
const style = document.createElement('style');
style.textContent = `
    @keyframes slideIn {
        from { transform: translateX(100%); opacity: 0; }
        to { transform: translateX(0); opacity: 1; }
    }
    
    .modal-overlay {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0,0,0,0.5);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
    }
    
    .modal-content {
        background: white;
        border-radius: 10px;
        padding: 0;
        max-width: 500px;
        width: 90%;
        max-height: 80vh;
        overflow-y: auto;
    }
    
    .modal-header {
        padding: 20px;
        border-bottom: 1px solid #eee;
        display: flex;
        justify-content: between;
        align-items: center;
    }
    
    .modal-header h3 {
        margin: 0;
        flex: 1;
    }
    
    .modal-close {
        background: none;
        border: none;
        font-size: 24px;
        cursor: pointer;
        color: #666;
    }
    
    .modal-body {
        padding: 20px;
    }
    
    .modal-footer {
        padding: 20px;
        border-top: 1px solid #eee;
        display: flex;
        gap: 10px;
        justify-content: flex-end;
    }
    
    .setting-group {
        margin-bottom: 20px;
    }
    
    .setting-group label {
        display: block;
        margin-bottom: 5px;
        font-weight: bold;
    }
    
    .setting-group input[type="range"] {
        width: 100%;
        margin: 10px 0;
    }
    
    .btn {
        padding: 10px 20px;
        border: none;
        border-radius: 5px;
        cursor: pointer;
        font-size: 14px;
    }
    
    .btn-primary {
        background: #667eea;
        color: white;
    }
    
    .btn-secondary {
        background: #6c757d;
        color: white;
    }
    
    .alert-stats {
        display: flex;
        gap: 20px;
        margin-bottom: 20px;
        padding: 15px;
        background: #f8f9fa;
        border-radius: 5px;
    }
    
    .stat-item {
        font-weight: bold;
    }
    
    .alert-history-item {
        padding: 10px;
        margin-bottom: 10px;
        border-radius: 5px;
        border-left: 4px solid;
    }
    
    .alert-history-item.alert-warning {
        background: #fff3cd;
        border-left-color: #ffc107;
    }
    
    .alert-history-item.alert-critical {
        background: #f8d7da;
        border-left-color: #dc3545;
    }
    
    .alert-time {
        font-size: 12px;
        color: #666;
        margin-bottom: 5px;
    }
    
    .alert-value {
        font-size: 12px;
        color: #666;
        margin-top: 5px;
    }
    
    .no-alerts {
        text-align: center;
        color: #666;
        padding: 20px;
    }
`;
document.head.appendChild(style);

// Запускаем приложение когда DOM загружен
document.addEventListener('DOMContentLoaded', () => {
    window.monitor = new SysPulseMonitor();
});

// Закрываем WebSocket при закрытии страницы
window.addEventListener('beforeunload', () => {
    if (window.monitor) {
        window.monitor.disconnect();
    }
});