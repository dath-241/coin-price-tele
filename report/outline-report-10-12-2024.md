# Báo Cáo Tiến Độ Công Việc

**Timeline:** Từ ngày 2/11 đến 10/12  

## Tiến Độ Công Việc

Trong khoảng thời gian này, nhóm đã thực hiện các hoạt động chính sau:  

- Đọc tài liệu từ file notion do thầy cung cấp.  
- Nghiên cứu thiết kế từ Frontend (FE) để triển khai các chức năng trên Telegram theo yêu cầu.  

## Các Chức Năng Đã Hoàn Thành

1. **Giá Spot, Giá Future, Funding Rate, Funding Rate Countdown, Funding Rate Interval**  
   - Các dữ liệu này được cập nhật từ Binance với thời gian trả về gần như thời gian thực (~1 giây).  

2. **Thử Nghiệm Webhook**  
   - Chức năng webhook đã được thử nghiệm thành công, đảm bảo backend có thể gửi thông báo đến Telegram bot một cách ổn định.  

3. **API-Token cho Đăng Nhập**  
   - Chức năng đăng nhập sử dụng API-token đã được triển khai.  
   - Token được lưu trữ trong database để đảm bảo người dùng không bị logout khi server restart.  

## Các Chức Năng Mới và cập nhật

1. **Giá theo Kline và Biểu Đồ Nến**  
   - Tính năng lấy giá theo Kline và hiển thị biểu đồ nến đã được triển khai trên môi trường production.  
   - Trải nghiệm người dùng được tối ưu hóa: thay vì yêu cầu nhập đúng định dạng `/kline <symbol> <interval> <limit>`, bot hiện cung cấp giao diện đề xuất các tham số dưới dạng nút bấm (button) để người dùng dễ dàng lựa chọn.  
   - Quy trình sử dụng mới:  
     - Gõ lệnh `/kline`.  
     - Chọn loại dữ liệu (ondemand/realtime).  
     - Chọn `symbol` bằng nút bấm hoặc nhập nếu danh sách đề xuất không có.  
     - Chọn `interval` bằng nút bấm.  
     - Dữ liệu được gửi theo thời gian thực với 3 chế độ:  
       - **Resume:** Tiếp tục/tạm dừng xem luồng dữ liệu.   
       - **Chart:** Hiển thị dữ liệu dưới dạng biểu đồ nến (candlestick chart).  
       - **Stop:** Dừng nhận dữ liệu, kết thúc 1 quy trình lệnh. 

2. **Cập Nhật Domain Server Backend**  
   - Đã cập nhật domain của server backend, đảm bảo việc kết nối và giao tiếp giữa các hệ thống frontend, backend và Telegram bot diễn ra ổn định.  

3. **Tính Năng Group Chat**  
   - Đã triển khai tính năng hỗ trợ sử dụng bot trong group chat trên Telegram.  
   - Các thành viên trong nhóm chat có thể:  
     - Xem các cập nhật giá, funding rate và các thông tin khác trực tiếp trong group chat.  
     - Tính năng quản lý lệnh theo người dùng để tránh xung đột khi có nhiều yêu cầu cùng lúc.  

## Các Quy Tắc và Quy Trình Làm Việc

- **Commit Message Rule:** Thống nhất cách viết commit message để dễ theo dõi và quản lý mã nguồn.  
- **Pull Request Rule:** Thiết lập quy trình Pull Request nhằm đảm bảo mã nguồn được kiểm tra trước khi hợp nhất vào nhánh chính.  
- **Code Style:** Áp dụng quy tắc phong cách mã hóa với `golangci`.  
- **Unit Test:**  
  - Đã tạo file unit test, nhưng hiện gặp lỗi do các module phụ thuộc vào nhau và chưa mock được các function.  
- **Continuous Deployment (CD):** Thiết lập CD, sử dụng Docker để build và triển khai trên Heroku.  
- **Continuous Integration (CI):** Chưa triển khai CI vì unit test đang bị lỗi.  
- **Webhook:** Áp dụng webhook để giao tiếp với Telegram server, tiết kiệm tài nguyên server.  
- **Meeting:** Họp định kỳ 2 tuần/lần để báo cáo tiến độ và chia sẻ kiến thức trong dự án.  

## Kết Luận

Nhóm đã hoàn thiện nhiều tính năng cơ bản và đưa lên production. Các tính năng khác như CI, cải thiện unit test và tối ưu hệ thống vẫn đang được xử lý để đảm bảo hoạt động ổn định và dễ mở rộng trong tương lai.  

---
