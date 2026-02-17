# ipbak

> **Professional Security-Grade IP Intelligence Utility**

`ipbak` is a high-performance, cross-platform CLI tool designed for developers, security researchers, and system administrators. It leverages MaxMind's industry-standard GeoLite2 databases to provide instant, detailed intelligence about any IP address.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20Windows%20%7C%20macOS-important?style=flat-square)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)

---

## ⚡ Features

*   **🌍 Precision Geolocation**: Instant access to Country, Region, City, Timezone, and Coordinates.
*   **📡 Network Intelligence**: Detailed ASN, Organization, and ISP data.
*   **🛡️ Smart Classification**: Automatically detects Hosting Providers (AWS, Google Cloud, etc.) vs Residential IPs.
*   **📍 My Location**: One-command lookup for your own public IP and location.
*   **📦 Batch Processing**: High-speed bulk analysis from file inputs.
*   **🤖 JSON Support**: Fully structured JSON output for easy integration with scripts and pipelines.
*   **🔄 Auto-Updates**: Smart database management with automatic `geoipupdate` installation and 24h cooldown enforcement.
*   **🔒 Privacy First**: Runs locally. No external API calls for lookups. No telemetry.

---

## 🚀 Installation

### From Source

Requires Go 1.21+

```bash
git clone https://github.com/yourusername/ipbak.git
cd ipbak
make build
```

This will create the `ipbak` binary in your current directory.

---

## 🛠️ Setup

`ipbak` uses MaxMind's GeoLite2 databases. Setup is a one-time process.

1.  **Register**: Create a free account at [MaxMind](https://www.maxmind.com).
2.  **Key**: Generate a License Key at [MaxMind License Key](https://www.maxmind.com/en/accounts/license-key).
3.  **Run Setup**:

    ```bash
    # Windows
    .\ipbak.exe setup

    # Linux / macOS / WSL
    ./ipbak setup
    ```

    Enter your credentials when prompted. `ipbak` will automatically install dependencies and download the latest databases.

---

## 📖 Usage

### 🔍 Single IP Lookup
Get detailed info about an IP address.

```bash
ipbak 8.8.8.8
```

### 📍 Where Am I?
Find your own public IP and location.

```bash
ipbak whereami
```

### 📂 Batch Processing
Analyze a list of IPs from a file (one per line).

```bash
ipbak batch ips.txt
```

### 📊 JSON Output
Perfect for piping into `jq` or other tools.

```bash
ipbak 8.8.8.8 --json
```

### 🧠 Advanced Mode
Show extra details like accuracy radius, postal code, and network classification.

```bash
ipbak 8.8.8.8 --advanced
```

---

## 🔧 Maintenance

### Updates
Update databases (respects 24h cooldown).

```bash
ipbak update
```

### Health Check
Diagnose installation and configuration issues.

```bash
ipbak doctor
```

---

## 📂 Configuration

| Platform | Config File | Database Location |
| :--- | :--- | :--- |
| **Windows** | `%APPDATA%\ipbak\config.yaml` | `%LOCALAPPDATA%\ipbak\` |
| **Linux/macOS** | `~/.config/ipbak/config.yaml` | `~/.local/share/ipbak/` |

---

## 📄 License

This project is licensed under the MIT License.

This product includes GeoLite2 data created by MaxMind, available from [https://www.maxmind.com](https://www.maxmind.com).
