class SysPulseMonitor {
    constructor() {
        this.version = '1.0.0';
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.isConnected = false;
        this.init();
    }

    async init() {
        console.log('🚀 Инициализация SysPulse Monitor...');
        console.log('🔌 Подключаемся к WebSocket...');
        
        // Загружаем версию приложения
        await this.loadVersion();
        
        // Подключаемся к WebSocket
        this.connectWebSocket();
        
        this.updateStatus('connecting', 'Подключение...');
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
                const metrics = JSON.parse(event.data);
                this.updateMetricsDisplay(metrics);
                this.updateLastUpdateTime();
            } catch (error) {
                console.error('❌ Ошибка парсинга метрик:', error);
            }
        };
        
        this.ws.onclose = (event) => {
            console.log('🔌 WebSocket отключен:', event.code, event.reason);
            this.isConnected = false;
            this.updateStatus('disconnected', 'Отключено');
            this.handleReconnection();
        };
        
        this.ws.onerror = (error) => {
            console.error('❌ WebSocket ошибка:', error);
            this.updateStatus('error', 'Ошибка подключения');
        };
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

    updateMetricsDisplay(metrics) {
        // Обновляем CPU
        document.getElementById('cpu-usage').textContent = `${metrics.cpu.usage.toFixed(1)}%`;
        document.getElementById('cpu-details').textContent = 
            `${metrics.cpu.cores} ядер | Load: ${metrics.cpu.load1.toFixed(2)}`;

        // Обновляем Memory с реальными данными
        const memUsedGB = metrics.memory.used / 1024 / 1024 / 1024;
        const memTotalGB = metrics.memory.total / 1024 / 1024 / 1024;
        document.getElementById('memory-usage').textContent = `${metrics.memory.usage.toFixed(1)}%`;
        document.getElementById('memory-details').textContent = 
            `${memUsedGB.toFixed(2)}GB / ${memTotalGB.toFixed(2)}GB`;

        // Обновляем Disk
        const diskUsedGB = metrics.disk.used / 1024 / 1024 / 1024;
        const diskTotalGB = metrics.disk.total / 1024 / 1024 / 1024;
        document.getElementById('disk-usage').textContent = `${metrics.disk.usage.toFixed(1)}%`;
        document.getElementById('disk-details').textContent = 
            `${diskUsedGB.toFixed(1)}GB / ${diskTotalGB.toFixed(1)}GB`;

        // Добавляем информацию о системе в консоль для отладки
        if (this.reconnectAttempts === 0) {
            console.log('📊 Получены метрики через WebSocket:', metrics.system);
        }
    }

    updateLastUpdateTime() {
        const now = new Date();
        document.getElementById('lastUpdate').textContent = now.toLocaleTimeString();
    }

    updateStatus(status, message) {
        const statusElement = document.getElementById('status');
        statusElement.textContent = message;
        
        // Сбрасываем классы
        statusElement.className = 'status-value';
        
        // Добавляем соответствующий класс
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

    // Метод для ручной отправки сообщений (для будущего использования)
    sendMessage(message) {
        if (this.ws && this.isConnected) {
            this.ws.send(JSON.stringify(message));
        }
    }

    // Метод для закрытия соединения
    disconnect() {
        if (this.ws) {
            this.ws.close();
        }
    }
}

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