# README

I made this to fix my problem when play CS2 game, but Steam client verylag or can't connect.
If i use 1.1.1.1 with WARP, Steam good, but CS2 game lag or disconnected.
So this is vibe coding with Gemini to fix that problem. After use my app. Steam client ok and no affect game connect.

<div align="center">

# 🚀 Steam Network Router Proxy

**Tăng tốc kết nối & mở khóa dịch vụ Steam tự động cho Windows**  
*Accelerate and seamlessly unblock Steam network services on Windows.*

[![Version](https://img.shields.io/github/v/release/toanbbpro/steam-router-proxy?style=flat-square&color=66c0f4)](https://github.com/toanbbpro/steam-router-proxy/releases)
[![License](https://img.shields.io/github/license/toanbbpro/steam-router-proxy?style=flat-square)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/toanbbpro/steam-router-proxy?style=flat-square)](https://github.com/toanbbpro/steam-router-proxy/stargazers)
[![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Ftoanbbpro%2Fsteam-router-proxy&count_bg=%2379C83D&title_bg=%23555555&icon=&icon_color=%23E7E7E7&title=views&edge_flat=true)](https://github.com/toanbbpro/steam-router-proxy)

[🇻🇳 Tiếng Việt](#-tiếng-việt) | [🇬🇧 English](#-english)

</div>

---

## 🇻🇳 Tiếng Việt

**Steam Router Proxy** là công cụ nhẹ, tối giản được xây dựng bằng **Go + Wails**, giúp tự động định tuyến traffic của Steam Client thông qua kỹ thuật DoH (DNS over HTTPS) & SNI Proxy để khắc phục tình trạng không truy cập được Store, Community, Inventory hoặc gián đoạn kết nối tới Steam trên Windows.

### ✨ Chức năng chính

* ⚡ **Zero Configuration:** Khởi chạy ứng dụng là tự động kích hoạt Proxy và điều hướng `hosts`, không cần thiết lập thủ công.
* 🛠️ **Chạy như Windows Service:** Tùy chọn cài đặt thành dịch vụ hệ thống tự khởi động cùng Windows (`Automatic`), chạy hoàn toàn ngầm kể cả khi tắt ứng dụng.
* 📝 **Hệ thống Log Timestamp:** Hỗ trợ ghi log chi tiết với cấu trúc tên dạng `proxy-log-YYYYMMDD_HHMMSS.log`, tự động xoay file khi đạt 8MB.
* 🧹 **Dọn dẹp tự động:** Tự động khôi phục cấu hình hệ thống (`hosts` / DNS Cache) ngay khi đóng ứng dụng hoặc gỡ Service.
* 🎨 **Giao diện Steam Style:** Giao diện tối giản, trực quan, hỗ trợ quyền Administrator tự động khi khởi chạy.

---

### 📸 Ảnh màn hình (Screenshot)

<div align="center">
  <img width="746" height="510" src="https://github.com/user-attachments/assets/c6d87cdd-e792-4e23-bb34-fd7087893262" alt="Steam Router GUI" >
</div>

---

### 📥 Link tải về (Download)

Phiên bản mới nhất: **v0.2.3**

👉 **[Tải file SteamRouter-v0.2.3.zip tại đây (GitHub Releases)](https://github.com/toanbbpro/steam-router-proxy/releases/latest)**

1. Tải file `.zip` và giải nén.
2. Mở file `steam-router.exe` (Chấp nhận quyền Admin nếu có).
3. Sử dụng ngay hoặc bấm **"Cài Service Tự Khởi Động"** để app chạy ngầm cùng Windows.

---

## 🇬🇧 English

**Steam Router Proxy** is a lightweight tool built with **Go + Wails** that automatically routes Steam traffic via DoH (DNS over HTTPS) & SNI Proxy to resolve connection issues with Steam Store, Community, and backend services on Windows.

### ✨ Key Features

* ⚡ **Plug & Play:** Automatically runs proxy and configures local `hosts` upon launch — zero manual setup required.
* 🛠️ **Windows Service Support:** Option to install as a background Windows Service (`Automatic`), running silently on boot without needing the GUI open.
* 📝 **Timestamped Logging:** Rotating log system with format `proxy-log-YYYYMMDD_HHMMSS.log` (Auto-rotates at 8MB).
* 🧹 **Clean Teardown:** Restores original system network settings (`hosts` / DNS Cache) when closed or uninstalled.
* 🎨 **Clean & Lightweight:** Steam-inspired compact UI built with Go & Wails, featuring automatic Admin privilege elevation.

---

### 📥 Download

Latest Release: **v0.2.3**

👉 **[Download SteamRouter-v0.2.3.zip (GitHub Releases)](https://github.com/toanbbpro/steam-router-proxy/releases/latest)**

---

## 📊 Analytics & Stats

![Repository Traffic](https://repobeats.axiom.co/api/embed/0000000000000000000000000000000000000000.svg "Repobeats analytics image")

---

<div align="center">

Developed by **[ToànBB](https://github.com/toanbbpro)**  
*If you find this project useful, please consider giving it a ⭐️!*

</div>
