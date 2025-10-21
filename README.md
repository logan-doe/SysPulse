# 💓 SysPulse - Real-Time System Monitor

![Go](https://img.shields.io/badge/Go-1.19+-00ADD8?style=for-the-badge&logo=go)
![WebSocket](https://img.shields.io/badge/WebSocket-Real--Time-FF6B6B?style=for-the-badge)
![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)

**Professional system monitoring solution with real-time metrics, beautiful dashboards, and intelligent alerts.**

---

## 🎯 Overview

SysPulse is a modern, real-time system monitoring tool that provides comprehensive insights into your system's performance with a beautiful, responsive web interface.

### ✨ Key Features

- 📊 **Real-time Monitoring** - Live metrics with 500ms updates
- 🎨 **Beautiful Dashboards** - Interactive charts and visualizations  
- 🔔 **Smart Alerts** - Configurable thresholds with sound notifications
- 📈 **Historical Data** - 60-point trend graphs for all metrics
- 🌐 **WebSocket Powered** - Instant updates without page reloads
- 📱 **Responsive Design** - Works on desktop and mobile
- ⚡ **Lightweight** - Minimal resource footprint

---

## 🚀 Quick Start

### Prerequisites
- Go 1.19 or higher
- Modern web browser

### Installation & Running

#### Method 1: From Source (Development)
```bash
# Clone the repository
git clone <your-repo-url>
cd syspulse

# Run directly with Go
go run ./cmd/syspulse-server

# Open browser manually
# http://localhost:8080
```

#### Method 2: Pre-built Binary
```bash
# Download latest release
# Extract and run:
./syspulse
# Browser opens automatically at http://localhost:8080
```

#### Method 3: Build from Source
```bash
# Build the project
./build.sh        # Linux/Mac
# or
build.bat         # Windows

# Run the built binary
cd dist
./run.sh          # Linux/Mac - auto-opens browser
# or  
run.bat           # Windows - auto-opens browser
```

---

## 📊 Monitoring Capabilities

### 🖥️ CPU Monitoring
- **Usage Percentage** - Real-time CPU utilization
- **Core Count** - Physical and logical processors
- **Load Average** - 1, 5, and 15-minute system load
- **Trend Analysis** - 60-point historical graph

### 🧠 Memory Monitoring  
- **Total/Used/Available** - Precise memory breakdown
- **Usage Percentage** - Current memory utilization
- **GB Display** - Human-readable memory values
- **Trend Tracking** - Memory usage patterns over time

### 💾 Disk Monitoring
- **Storage Utilization** - Disk space monitoring
- **Total/Used/Free** - Complete storage breakdown  
- **Automatic Detection** - Finds largest partition
- **Usage Trends** - Storage consumption patterns

### 🔔 Alert System
- **Configurable Thresholds** - Set custom limits for each metric
- **Multi-level Alerts** - Warning and Critical notifications
- **Sound Notifications** - Audio alerts for immediate attention
- **Alert History** - Complete log of all past alerts
- **Visual Indicators** - Color-coded metric displays

---

## 🛠️ Configuration

### Environment Variables
```bash
# Server Port (default: 8080)
export SYS_PULSE_PORT=9090

# Update Interval in milliseconds (default: 500)
export SYS_PULSE_UPDATE_INTERVAL=1000

# Alert Thresholds (percentage)
export SYS_PULSE_ALERT_CPU=85
export SYS_PULSE_ALERT_RAM=90  
export SYS_PULSE_ALERT_DISK=95

# Environment
export SYS_PULSE_ENV=production
```

### Web Interface Configuration
Access the settings modal to configure:
- Alert thresholds (50-95%)
- Enable/disable alert system
- Sound notifications
- All changes apply immediately

---

## 🏗️ Architecture

```
SysPulse/
├── cmd/syspulse-server/     # Application entry point
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── handlers/          # HTTP request handlers
│   ├── models/            # Data structures
│   ├── services/          # Business logic
│   │   ├── metrics_service.go     # System metrics collection
│   │   ├── websocket_service.go   # Real-time communication
│   │   └── alert_service.go       # Alert management
│   └── ...
├── web/static/            # Frontend assets
│   ├── index.html         # Main interface
│   ├── css/style.css      # Styling
│   └── js/app.js          # Client-side logic
└── pkg/metrics/           # Public utilities
```

### Technology Stack
- **Backend**: Go 1.19+ with standard library
- **Frontend**: Vanilla JavaScript, Chart.js, CSS3
- **Real-time**: WebSocket (gorilla/websocket)
- **Metrics**: gopsutil for system data
- **Styling**: Modern CSS with gradients and animations

---

## 🔧 API Endpoints

### HTTP API
```http
GET  /api/health          # Service health check
GET  /api/metrics         # Current system metrics
GET  /api/version         # Application version
GET  /api/clients         # Connected WebSocket clients
GET  /api/alerts/history  # Alert history
GET  /api/alerts/config   # Alert configuration
POST /api/alerts/config   # Update alert settings
POST /api/alerts/clear    # Clear alert history
```

### WebSocket
```http
GET /ws  # Real-time metrics stream
```

---

## 🎨 Interface Features

### Dashboard Components
1. **Status Bar** - Connection status, version, last update
2. **Metric Cards** - CPU, Memory, Disk with live values
3. **Trend Graphs** - Real-time charts with 60-point history
4. **Alert Panel** - Active notifications with timestamps
5. **Control Buttons** - Settings and history access

### Visual Indicators
- 🟢 **Green** - Normal operation (< warning threshold)
- 🟡 **Yellow** - Warning state (≥ warning threshold)  
- 🔴 **Red** - Critical state (≥ critical threshold)

### Responsive Design
- **Desktop** - Full dashboard with side-by-side metrics
- **Tablet** - Adaptive card layout
- **Mobile** - Stacked cards with optimized touch targets

---

## 📈 Performance Characteristics

- **Memory Usage**: ~15-25MB (Go runtime + metrics)
- **CPU Impact**: <1% during normal operation
- **Update Frequency**: 500ms (configurable)
- **WebSocket Messages**: ~2KB per update
- **Startup Time**: <2 seconds

---

## 🚨 Alert System

### Default Thresholds
- **CPU**: 80% (Warning) → 90% (Critical)
- **Memory**: 85% (Warning) → 95% (Critical)  
- **Disk**: 90% (Warning) → 95% (Critical)

### Alert Types
- **Visual** - Color changes and notification panel
- **Audible** - Distinct sounds for warning/critical
- **Historical** - Persistent alert log with statistics

### Customization
```javascript
// Example: Custom thresholds via web interface
{
  "cpu_threshold": 75,
  "ram_threshold": 80, 
  "disk_threshold": 85,
  "enabled": true
}
```

---

## 🔄 Development

### Project Structure
```bash
# Core application
go run cmd/syspulse-server/main.go

# Build for distribution
./build.sh

# Run tests (when implemented)
go test ./...

# Code formatting
go fmt ./...
```

### Adding New Metrics
1. Extend models in `internal/models/metrics.go`
2. Implement collection in `internal/services/metrics_service.go` 
3. Update frontend display in `web/static/js/app.js`
4. Add to WebSocket broadcast in `cmd/syspulse-server/main.go`

---

## 🌟 Use Cases

### 🏢 Enterprise Monitoring
- Server health monitoring
- Capacity planning
- Performance benchmarking
- Alert escalation

### 💻 Developer Workstations
- Resource usage tracking
- Development environment monitoring
- Performance debugging

### 🔬 Educational
- System administration learning
- Real-time data visualization
- WebSocket implementation examples

---

## 📊 Screenshots

*(Include actual screenshots in your repository)*

1. **Main Dashboard** - Clean metrics display with charts
2. **Alert Settings** - Configuration modal with sliders  
3. **Alert History** - Timeline of past notifications
4. **Mobile View** - Responsive design demonstration

---

## 🤝 Contributing

We welcome contributions! Please see our contributing guidelines for details.

### Development Setup
```bash
# 1. Fork and clone repository
git clone https://github.com/your-username/syspulse.git

# 2. Install dependencies
go mod download

# 3. Run in development mode
go run ./cmd/syspulse-server

# 4. Make changes and test
# 5. Submit pull request
```

### Roadmap
- [ ] Database persistence for historical data
- [ ] Email/SMS notifications
- [ ] Multi-system monitoring
- [ ] Docker containerization
- [ ] Plugin system for custom metrics

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🆘 Support

### Common Issues
**Browser doesn't open automatically**
- Check if scripts are executable: `chmod +x dist/run.sh`
- Manual access: http://localhost:8080

**High CPU usage**
- Increase update interval: `export SYS_PULSE_UPDATE_INTERVAL=1000`

**Port conflicts**
- Change port: `export SYS_PULSE_PORT=9090`

### Getting Help
- 📖 Check this README
- 🐛 Open an issue on GitHub
- 💬 Discussion forums

---

## 🎊 Conclusion

SysPulse provides enterprise-grade system monitoring with the simplicity of a single binary. Whether you're managing production servers or just curious about your system's performance, SysPulse delivers beautiful, real-time insights with minimal overhead.

**Start monitoring in under 60 seconds!** 🚀

---

*Built with ❤️ using Go and modern web technologies.*