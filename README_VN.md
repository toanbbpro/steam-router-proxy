# README

Ứng dụng này được phát triển nhằm giải quyết triệt để sự cố khi chơi game CS2: Client Steam bị giật lag hoặc không thể kết nối.
Nếu sử dụng 1.1.1.1 WARP, Steam hoạt động tốt nhưng CS2 lại bị gián đoạn hoặc giật lag.
Vì vậy, công cụ này ra đời (vibe coding cùng Gemini) để khắc phục tình trạng trên. Sau khi sử dụng, Steam Client hoạt động mượt mà và hoàn toàn không ảnh hưởng tới kết nối trong game.

<div align="center">

# 🚀 Steam Router Proxy

**Tăng tốc kết nối & mở khóa dịch vụ Steam tự động cho Windows**

[![Version](https://img.shields.io/github/v/release/toanbbpro/steam-router-proxy?style=flat-square&color=66c0f4)](https://github.com/toanbbpro/steam-router-proxy/releases)
[![License](https://img.shields.io/github/license/toanbbpro/steam-router-proxy?style=flat-square)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/toanbbpro/steam-router-proxy?style=flat-square)](https://github.com/toanbbpro/steam-router-proxy/stargazers)
[![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Ftoanbbpro%2Fsteam-router-proxy&count_bg=%2379C83D&title_bg=%23555555&icon=&icon_color=%23E7E7E7&title=views&edge_flat=true)](https://github.com/toanbbpro/steam-router-proxy)

**🇻🇳 Tiếng Việt** | [🇬🇧 English](README.md)

</div>

---

**Steam Router Proxy** là công cụ nhẹ, tối giản được xây dựng bằng **Go + Wails**, giúp tự động định tuyến traffic của Steam Client thông qua kỹ thuật DoH (DNS over HTTPS) & SNI Proxy để khắc phục tình trạng không truy cập được Store, Community, Inventory hoặc gián đoạn kết nối tới Steam trên Windows.

### ✨ Chức năng chính

* ⚡ **Zero Configuration:** Khởi chạy ứng dụng là tự động kích hoạt Proxy và điều hướng `hosts`, không cần thiết lập thủ công.
* 🛠️ **Chạy như Windows Service:** Tùy chọn cài đặt thành dịch vụ hệ thống tự khởi động cùng Windows (`Automatic`), chạy hoàn toàn ngầm kể cả khi tắt ứng dụng.
* 🌐 **Giao diện Song ngữ & Hiện đại:** Giao diện tối giản phong cách Steam, hiển thị phiên bản trên thanh tiêu đề Native Window, tích hợp nút chuyển đổi ngôn ngữ nhanh dạng cờ SVG (Anh / Việt).
* 💬 **Hộp thoại Thông báo Steam Style:** Popup Modal thông báo thiết kế riêng mượt mà, hỗ trợ ngắt dòng rõ ràng thay thế cho alert trình duyệt.
* 📝 **Hệ thống Log Timestamp:** Hỗ trợ ghi log chi tiết với cấu trúc tên dạng `proxy-log-YYYYMMDD_HHMMSS.log`, tự động xoay file khi đạt 8MB.
* 🧹 **Dọn dẹp tự động:** Tự động khôi phục cấu hình hệ thống (`hosts` / DNS Cache) ngay khi đóng ứng dụng hoặc gỡ Service.

---

### 📸 Ảnh màn hình (Screenshot)

<div align="center">
  <img width="702" height="531" src="https://github.com/user-attachments/assets/571e3c4a-6787-4439-bda5-389ac451a06d" alt="Steam Router Proxy v0.3 UI" >
</div>

---

### 📥 Link tải về (Download)

Phiên bản mới nhất: **v0.3**

👉 **[Tải file steam-router-proxy.exe tại đây (GitHub Releases)](https://github.com/toanbbpro/steam-router-proxy/releases/latest)**

1. Tải trực tiếp file `steam-router-proxy.exe`.
2. Mở file và khởi chạy (Chấp nhận quyền Admin nếu có).
3. Sử dụng ngay hoặc bấm **"Cài Service Tự Khởi Động"** để ứng dụng chạy ngầm cùng Windows.

---

## 📊 Thống kê dự án

![Repository Traffic](https://repobeats.axiom.co/api/embed/toanbbpro/steam-router-proxy.svg "Repobeats analytics image")

---

<div align="center">

Phát triển bởi **[ToànBB](https://github.com/toanbbpro)**  
*Nếu dự án này hữu ích với bạn, đừng quên tặng dự án một ngôi sao ⭐️ nhé!*

</div>
