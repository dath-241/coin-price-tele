
### 1. **/help**
   - **Mô tả**: Hiển thị danh sách tất cả các lệnh hỗ trợ.
   - **Cú pháp**: `/help`
   - **Output**: Danh sách các lệnh và mô tả của từng lệnh.

### 2. **/start**
   - **Mô tả**: Đăng nhập người dùng vào hệ thống.
   - **Cú pháp**: `/start`
   - **Output**: Nếu đăng nhập thành công, trả về thông báo chào mừng; nếu không, trả về "Access denied."

### 3. **/login**
   - **Mô tả**: Đăng nhập vào hệ thống bằng tài khoản.
   - **Cú pháp**: `/login <username> <password>`
   - **Output**: Trả về thông báo đăng nhập thành công và lưu token người dùng, hoặc báo lỗi nếu có vấn đề.

### 4. **/getinfo**
   - **Mô tả**: Lấy thông tin người dùng.
   - **Cú pháp**: `/getinfo`
   - **Output**: Thông tin người dùng; hoặc lỗi nếu không lấy được token.

### 5. **/kline**
   - **Mô tả**: Lấy dữ liệu Kline cho một mã giao dịch.
   - **Cú pháp**: `/kline <symbol> <interval> [limit]`
   - **Output**: Gửi biểu đồ Kline với dữ liệu cho mã giao dịch và khung thời gian đã chọn.

### 6. **/menu**
   - **Mô tả**: Hiển thị menu các lựa chọn cho người dùng.
   - **Cú pháp**: `/menu`
   - **Output**: Gửi menu các tùy chọn.

### 7. **/price_spot**
   - **Mô tả**: Bắt đầu luồng cập nhật giá Spot cho một mã giao dịch.
   - **Cú pháp**: `/price_spot <symbol>`
   - **Output**: Cập nhật giá Spot theo thời gian thực cho mã giao dịch.

### 8. **/price_futures**
   - **Mô tả**: Bắt đầu luồng cập nhật giá Futures cho một mã giao dịch.
   - **Cú pháp**: `/price_futures <symbol>`
   - **Output**: Cập nhật giá Futures theo thời gian thực cho mã giao dịch.

### 9. **/funding_rate**
   - **Mô tả**: Bắt đầu luồng cập nhật funding rate cho một mã giao dịch.
   - **Cú pháp**: `/funding_rate <symbol>`
   - **Output**: Cập nhật funding rate cho mã giao dịch đã chọn.

### 10. **/kline_realtime**
   - **Mô tả**: Bắt đầu luồng cập nhật Kline thời gian thực cho một mã giao dịch.
   - **Cú pháp**: `/kline_realtime <symbol> <interval>`
   - **Output**: Cập nhật dữ liệu Kline thời gian thực cho mã giao dịch đã chọn.

### 11. **/stop**
   - **Mô tả**: Dừng cập nhật dữ liệu thời gian thực.
   - **Cú pháp**: `/stop`
   - **Output**: Thông báo dừng cập nhật thành công nếu có luồng đang hoạt động.

### 12. **/all_triggers**
   - **Mô tả**: Lấy danh sách tất cả các cảnh báo kích hoạt.
   - **Cú pháp**: `/all_triggers`
   - **Output**: Trả về danh sách các trigger hiện tại của người dùng.

### 13. **/delete_trigger**
   - **Mô tả**: Xóa một trigger cụ thể.
   - **Cú pháp**: `/delete_trigger <spot/future/price-difference/funding-rate> <symbol>`
   - **Output**: Xóa trigger đã chọn nếu có, hoặc trả về lỗi nếu không tìm thấy.

### 14. **/alert_price_with_threshold**
   - **Mô tả**: Tạo cảnh báo giá dựa trên ngưỡng đã đặt.
   - **Cú pháp**: `/alert_price_with_threshold <spot/future> <lower/above> <symbol> <threshold>`
   - **Output**: Cài đặt cảnh báo giá cho mã giao dịch với ngưỡng đã chọn.

### 15. **/price_difference**
   - **Mô tả**: Tạo cảnh báo khi sự khác biệt giá đạt ngưỡng đặt trước.
   - **Cú pháp**: `/price_difference <lower/above> <symbol> <threshold>`
   - **Output**: Cài đặt cảnh báo cho chênh lệch giá với ngưỡng đã chọn.

### 16. **/funding_rate_change**
   - **Mô tả**: Tạo cảnh báo khi funding rate đạt ngưỡng đặt trước.
   - **Cú pháp**: `/funding_rate_change <lower/above> <symbol> <threshold>`
   - **Output**: Cài đặt cảnh báo cho funding rate với ngưỡng đã chọn.
