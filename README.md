# README

I made this to fix my problem when playing CS2—where the Steam client was extremely laggy or failed to connect entirely.
If I used 1.1.1.1 with WARP, Steam worked well, but CS2 suffered from high latency or disconnections.
So this was built with Gemini to fix that issue. After running this app, the Steam client connects smoothly without impacting in-game network stability.

<div align="center">

# 🚀 Steam Router Proxy

**Accelerate and seamlessly unblock Steam network services on Windows.**

[![Version](https://img.shields.io/github/v/release/toanbbpro/steam-router-proxy?style=flat-square&color=66c0f4)](https://github.com/toanbbpro/steam-router-proxy/releases)
[![License](https://img.shields.io/github/license/toanbbpro/steam-router-proxy?style=flat-square)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/toanbbpro/steam-router-proxy?style=flat-square)](https://github.com/toanbbpro/steam-router-proxy/stargazers)
[![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Ftoanbbpro%2Fsteam-router-proxy&count_bg=%2379C83D&title_bg=%23555555&icon=&icon_color=%23E7E7E7&title=views&edge_flat=true)](https://github.com/toanbbpro/steam-router-proxy)

[🇻🇳 Tiếng Việt](README_VN.md) | **🇬🇧 English**

</div>

---

**Steam Router Proxy** is a lightweight tool built with **Go + Wails** that automatically routes Steam traffic via DoH (DNS over HTTPS) & SNI Proxy to resolve connection issues with Steam Store, Community, Inventory, and backend services on Windows.

### ✨ Key Features

* ⚡ **Plug & Play:** Automatically runs the proxy and configures local `hosts` upon launch — zero manual setup required.
* 🛠️ **Windows Service Support:** Option to install as a background Windows Service (`Automatic`), running silently on boot without needing the GUI open.
* 🌐 **Bilingual & Modern UI:** Features a compact, Steam-themed dark UI with native window title integration and a single-click SVG flag language switcher (English / Vietnamese).
* 💬 **Custom Steam-Themed Modals:** Clean, multi-line notification modals replacing standard browser alerts.
* 📝 **Timestamped Logging:** Rotating log system with format `proxy-log-YYYYMMDD_HHMMSS.log` (Auto-rotates at 8MB).
* 🧹 **Clean Teardown:** Automatically restores original system network settings (`hosts` / DNS Cache) when closed or uninstalled.

---

### 📸 Screenshot

<div align="center">
  <img width="702" height="531" src="https://github.com/user-attachments/assets/571e3c4a-6787-4439-bda5-389ac451a06d" alt="Steam Router Proxy v0.3 UI" >
</div>

---

### 📥 Download

Latest Release: **v0.3**

👉 **[Download steam-router-proxy.exe (GitHub Releases)](https://github.com/toanbbpro/steam-router-proxy/releases/latest)**

1. Download the executable `steam-router-proxy.exe`.
2. Run the application (Accept Administrator privileges if prompted).
3. Use it immediately or click **"Install Auto-Start Service"** to keep it running silently in the background.

---

## 📊 Analytics & Stats

![Repository Traffic](https://repobeats.axiom.co/api/embed/toanbbpro/steam-router-proxy.svg "Repobeats analytics image")

---

<div align="center">

Developed by **[ToànBB](https://github.com/toanbbpro)**  
*If you find this project useful, please consider giving it a ⭐️!*

</div>
