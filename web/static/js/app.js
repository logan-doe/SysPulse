class SysPulseMonitor {
    constructor() {
        this.version = '1.0.0';
        this.updateInterval = 2000;
        this.isConnected = false;
        this.init();
    }

    async init() {
        console.log('🚀 Инициализация SysPulse Monitor...');
        
        // Загружаем версию приложения
        await this.loadVersion();
        
        // Загружаем начальные метрики
        await this.loadMetrics();
        
        // Запускаем периодическое обновление
        this.startAutoUpdate();
        
        this.updateStatus('connected', 'Подключено');
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

    async loadMetrics() {
        try {
            const response = await fetch('/api/metrics');
            const metrics = await response.json();
            
            this.updateMetricsDisplay(metrics);
            this.updateLastUpdateTime();
            
        } catch (error) {
            console.error('Ошибка загрузки метрик:', error);
            this.updateStatus('error', 'Ошибка подключения');
        }
    }

    updateMetricsDisplay(metrics) {
        // Обновляем CPU
        document.getElementById('cpu-usage').textContent = `${metrics.cpu.usage.toFixed(1)}%`;
        document.getElementById('cpu-details').textContent = 
            `${metrics.cpu.cores} ядер | Load: ${metrics.cpu.load1.toFixed(2)}`;

        // Обновляем Memory
        const memUsedGB = metrics.memory.used / 1024 / 1024 / 1024;
        const memTotalGB = metrics.memory.total / 1024 / 1024 / 1024;
        document.getElementById('memory-usage').textContent = `${metrics.memory.usage.toFixed(1)}%`;
        document.getElementById('memory-details').textContent = 
            `${memGB.toFixed(2)}GB / ${memTotalGB.toFixed(2)}GB`;

        // Обновляем Disk
        const diskGB = metrics.disk.used / 1024 / 1024 / 1024;
        const diskTotalGB = metrics.disk.total / 1024 / 1024 / 1024;
        document.getElementById('disk-usage').textContent = `${metrics.disk.usage.toFixed(1)}%`;
        document.getElementById('disk-details').textContent = 
            `${diskGB.toFixed(1)}GB / ${diskTotalGB.toFixed(1)}GB`;

        console.log('System:', metrics.system);
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
        } else if (status === 'error') {
            statusElement.style.color = '#f44336';
        } else {
            statusElement.style.color = '#ff9800';
        }
    }

    startAutoUpdate() {
        setInterval(() => {
            this.loadMetrics();
        }, this.updateInterval);
    }
}

// Запускаем приложение когда DOM загружен
document.addEventListener('DOMContentLoaded', () => {
    new SysPulseMonitor();
});