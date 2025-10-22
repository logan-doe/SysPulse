class SysPulseMonitor {
    constructor() {
        this.ws = null;
        this.isConnected = false;
        this.currentTheme = localStorage.getItem('theme') || 'light';
        this.charts = {};
        this.init();
    }

    init() {
        console.log('🚀 SysPulse Monitor запускается...');
        this.applyTheme(this.currentTheme);
        this.initCharts();
        this.connectWebSocket();
        this.setupEventListeners();
    }

    applyTheme(theme) {
        this.currentTheme = theme;
        document.documentElement.setAttribute('data-theme', theme);
        localStorage.setItem('theme', theme);
        
        const themeButton = document.getElementById('theme-toggle');
        if (themeButton) {
            themeButton.textContent = theme === 'light' ? '🌙 Тёмная' : '☀️ Светлая';
        }
    }

    toggleTheme() {
        const newTheme = this.currentTheme === 'light' ? 'dark' : 'light';
        this.applyTheme(newTheme);
        this.updateChartColors();
    }

    setupEventListeners() {
        const themeButton = document.getElementById('theme-toggle');
        if (themeButton) {
            themeButton.addEventListener('click', () => this.toggleTheme());
        }
    }

    initCharts() {
        const chartsConfig = {
            cpu: { color: this.getChartColor('cpu') },
            memory: { color: this.getChartColor('memory') },
            disk: { color: this.getChartColor('disk') }
        };

        this.charts = {};
        
        Object.keys(chartsConfig).forEach(type => {
            const ctx = document.getElementById(`${type}-chart`);
            if (!ctx) return;

            this.charts[type] = new Chart(ctx, {
                type: 'line',
                data: {
                    labels: [],
                    datasets: [{
                        data: [],
                        borderColor: chartsConfig[type].color,
                        borderWidth: 3,
                        fill: false,
                        tension: 0.4,
                        pointRadius: 0
                    }]
                },
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
                    animation: { duration: 0 }
                }
            });
        });
    }

    getChartColor(type) {
        const styles = getComputedStyle(document.documentElement);
        return styles.getPropertyValue(`--chart-${type}`).trim() || '#475569';
    }

    updateChartColors() {
        Object.keys(this.charts).forEach(type => {
            if (this.charts[type]) {
                this.charts[type].data.datasets[0].borderColor = this.getChartColor(type);
                this.charts[type].update('none');
            }
        });
    }

    connectWebSocket() {
        this.ws = new WebSocket(`ws://${window.location.host}/ws`);
        
        this.ws.onopen = () => {
            this.isConnected = true;
            this.updateStatus('connected', '✅ Подключено');
        };
        
        this.ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                this.updateAllMetrics(data);
            } catch (error) {
                console.error('❌ Ошибка парсинга:', error);
            }
        };
        
        this.ws.onclose = () => {
            this.isConnected = false;
            this.updateStatus('disconnected', '❌ Отключено');
        };
        
        this.ws.onerror = (error) => {
            console.error('❌ WebSocket ошибка:', error);
            this.updateStatus('error', '⚠️ Ошибка подключения');
        };
    }

    updateAllMetrics(data) {
        this.updateSystemMetrics(data);
        this.updateNetworkMetrics(data);
        this.updateProcesses(data);
        this.updateCharts(data);
        this.updateTime();
    }

    updateSystemMetrics(data) {
        if (data.cpu) {
            this.updateElement('cpu-usage', `${data.cpu.usage.toFixed(1)}%`);
            this.updateElement('cpu-details', `${data.cpu.cores} ядер | Load: ${data.cpu.load1.toFixed(2)}`);
        }
        
        if (data.memory) {
            this.updateElement('memory-usage', `${data.memory.usage.toFixed(1)}%`);
            if (data.memory.used && data.memory.total) {
                const usedGB = (data.memory.used / 1024 / 1024 / 1024).toFixed(1);
                const totalGB = (data.memory.total / 1024 / 1024 / 1024).toFixed(1);
                this.updateElement('memory-details', `${usedGB}GB / ${totalGB}GB`);
            }
        }
        
        if (data.disk) {
            this.updateElement('disk-usage', `${data.disk.usage.toFixed(1)}%`);
            if (data.disk.used && data.disk.total) {
                const usedGB = (data.disk.used / 1024 / 1024 / 1024).toFixed(1);
                const totalGB = (data.disk.total / 1024 / 1024 / 1024).toFixed(1);
                this.updateElement('disk-details', `${usedGB}GB / ${totalGB}GB`);
            }
        }
    }

    updateNetworkMetrics(data) {
        if (!data.network) return;
        
        const net = data.network;
        this.updateElement('network-status', net.is_online ? '🟢 ONLINE' : '🔴 OFFLINE');
        this.updateElement('network-upload', `${net.current_upload.toFixed(2)} Mb/s`);
        this.updateElement('network-download', `${net.current_download.toFixed(2)} Mb/s`);
        this.updateElement('network-ping', `${net.ping.toFixed(1)} ms`);
        this.updateElement('network-ip', net.local_ip || '-');
    }

    updateProcesses(data) {
        if (!data.processes || !Array.isArray(data.processes)) return;
        
        const topCPU = data.processes
            .filter(p => p.cpu_percent !== undefined)
            .sort((a, b) => b.cpu_percent - a.cpu_percent)
            .slice(0, 10);
        
        const topMemory = data.processes
            .filter(p => p.memory_percent !== undefined)
            .sort((a, b) => b.memory_percent - a.memory_percent)
            .slice(0, 10);
        
        this.updateProcessTable('cpu-processes', topCPU);
        this.updateProcessTable('memory-processes', topMemory);
    }

    updateProcessTable(containerId, processes) {
        const container = document.getElementById(containerId);
        if (!container) return;
        
        if (processes.length === 0) {
            container.innerHTML = '<div class="no-processes">Нет данных</div>';
            return;
        }
        
        const rows = processes.map(proc => `
            <tr onclick="window.monitor.showProcessDetails(${JSON.stringify(proc).replace(/"/g, '&quot;')})">
                <td class="process-name">${this.truncateText(proc.process || 'Unknown', 25)}</td>
                <td class="process-cpu">${proc.cpu_percent.toFixed(1)}%</td>
                <td class="process-memory">${proc.memory_percent.toFixed(1)}%</td>
                <td class="process-status">${proc.status || 'N/A'}</td>
            </tr>
        `).join('');
        
        container.innerHTML = `
            <table class="processes-table">
                <thead>
                    <tr>
                        <th>Процесс</th>
                        <th>CPU</th>
                        <th>Память</th>
                        <th>Статус</th>
                    </tr>
                </thead>
                <tbody>${rows}</tbody>
            </table>
        `;
    }

    showProcessDetails(process) {
        const modal = document.createElement('div');
        modal.className = 'modal-overlay';
        modal.innerHTML = `
            <div class="modal-content">
                <div class="modal-header">
                    <h3>🔍 Детали процесса</h3>
                    <button class="modal-close" onclick="this.closest('.modal-overlay').remove()">×</button>
                </div>
                <div class="modal-body">
                    <div class="process-detail">
                        <label>Процесс:</label>
                        <span>${process.process || 'Unknown'}</span>
                    </div>
                    <div class="process-detail">
                        <label>PID:</label>
                        <span>${process.pid}</span>
                    </div>
                    <div class="process-detail">
                        <label>Пользователь:</label>
                        <span>${process.user || 'N/A'}</span>
                    </div>
                    <div class="process-detail">
                        <label>Потоков:</label>
                        <span>${process.threads || 0}</span>
                    </div>
                    <div class="process-detail">
                        <label>Запущен:</label>
                        <span>${process.createtime ? new Date(process.createtime * 1000).toLocaleString() : 'N/A'}</span>
                    </div>
                    <div class="process-detail full-width">
                        <label>Командная строка:</label>
                        <div class="command-line">${process.commandline || 'N/A'}</div>
                    </div>
                </div>
            </div>
        `;
        document.body.appendChild(modal);
    }

    updateCharts(data) {
        if (data.cpu && this.charts.cpu) {
            this.updateChart('cpu', data.cpu.usage);
        }
        if (data.memory && this.charts.memory) {
            this.updateChart('memory', data.memory.usage);
        }
        if (data.disk && this.charts.disk) {
            this.updateChart('disk', data.disk.usage);
        }
    }

    updateChart(type, value) {
        const chart = this.charts[type];
        if (!chart) return;
        
        chart.data.labels.push('');
        chart.data.datasets[0].data.push(value);
        
        if (chart.data.labels.length > 30) {
            chart.data.labels.shift();
            chart.data.datasets[0].data.shift();
        }
        
        chart.update('none');
    }

    updateElement(id, value) {
        const element = document.getElementById(id);
        if (element) element.textContent = value;
    }

    updateStatus(status, message) {
        this.updateElement('status', message);
    }

    updateTime() {
        this.updateElement('lastUpdate', new Date().toLocaleTimeString());
    }

    truncateText(text, length) {
        return text.length > length ? text.substring(0, length) + '...' : text;
    }
}

// Запуск приложения
document.addEventListener('DOMContentLoaded', () => {
    window.monitor = new SysPulseMonitor();
});

window.addEventListener('beforeunload', () => {
    if (window.monitor && window.monitor.ws) {
        window.monitor.ws.close();
    }
});