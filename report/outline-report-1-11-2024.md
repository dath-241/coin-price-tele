# Báo cáo tiến độ công việc

**Timeline:** Từ ngày 3/10 đến 1/11

## Tiến độ công việc

Trong khoảng thời gian này, tất cả các thành viên trong nhóm đã thực hiện các hoạt động sau:

- Đọc file notion của thầy.
- Nghiên cứu thiết kế từ Frontend (FE) để triển khai các chức năng trên Telegram theo yêu cầu.

## Các chức năng đã hoàn thành

1. **Giá Spot, Giá Future, Funding Rate, Funding Rate Countdown, Funding Rate Interval**
   - Các dữ liệu này được cập nhật từ Binance với thời gian trả về gần như thời gian thực (~ 1 giây).
   
2. **Giá theo kline và biểu đồ nến**
   - Đã thực hiện thành công việc lấy giá theo kline và vẽ được biểu đồ nến.

3. **Thử nghiệm webhook**
   - Chức năng webhook đã được thử nghiệm thành công, khi backend gửi webhook về Telegram thì tin nhắn đã được gửi đi. 

4. **API-token cho đăng nhập**
   - Đã triển khai chức năng đăng nhập sử dụng API-token và lưu token vào database nhằm tình trạng logout khi server restart.

## Tình trạng các tính năng

- Các tính năng cơ bản đã đạt được. Tuy nhiên, tính năng kline vẫn chưa được triển khai lên môi trường production.

## Các quy tắc và quy trình làm việc

- **Commit Message Rule:** Nhóm đã thống nhất quy tắc về cách viết commit message để dễ dàng theo dõi và quản lý code.
- **Pull Request Rule:** Đã thiết lập quy trình cho Pull Request để đảm bảo mã nguồn được xem xét và kiểm tra trước khi hợp nhất vào nhánh chính.
- **Code Style:** Đã áp dụng quy tắc phong cách mã hóa bằng cách sử dụng `golangci`.
- **Unit Test:** Đã có các file unit test. Tuy nhiên, do mã nguồn hiện tại các module đang bị phụ thuộc vào nhau nên các unit test đang bị lỗi vì chưa mock được các function.
- **Continuous Deployment (CD):** Đã thiết lập quy trình CD, build bằng Docker và triển khai lên Heroku.
- **Continuous Integration (CI):** Hiện tại chưa có CI vì các unit test đang bị lỗi.
- **Webhook:** Nhóm đã áp dụng kỹ thuật webhook để giao tiếp với Telegram server nhằm tiết kiệm chi phí server.
- **Meeting:** Nhóm đã tiến hành họp vào 2 tuần một lần để báo cáo tiến độ công việc và trình bày lại các kiến thức đã học được khi làm dự án. 
## Kết luận

Nhóm đang tiếp tục làm việc để hoàn thiện các tính năng còn lại và khắc phục các vấn đề còn tồn đọng.

---

