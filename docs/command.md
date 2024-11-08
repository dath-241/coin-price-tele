# List command

### 1. **/help**
   - **Công dụng**: Hiển thị danh sách các lệnh hỗ trợ.
   - **Input**: Không yêu cầu tham số.
   - **Output**: Danh sách các lệnh khả dụng dưới dạng văn bản.

### 2. **/start**
   - **Công dụng**: Xác thực người dùng.
   - **Input**: Không yêu cầu tham số.
   - **Output**: "Access denied." nếu xác thực thất bại; nếu thành công, trả về thông báo xác nhận.

### 3. **/login <username> <password>**
   - **Công dụng**: Đăng nhập với tài khoản người dùng.
   - **Input**: `username` (tên người dùng), `password` (mật khẩu).
   - **Output**: "Đăng nhập thành công." hoặc lỗi nếu đăng nhập thất bại.

### 4. **/getinfo**
   - **Công dụng**: Lấy thông tin người dùng.
   - **Input**: Không yêu cầu tham số.
   - **Output**: Hiển thị thông tin người dùng hoặc thông báo lỗi nếu có sự cố.

### 5. **/kline <symbol> <interval> [limit]**
   - **Công dụng**: Lấy dữ liệu Kline (biểu đồ nến) của coin.
   - **Input**:
     - `symbol`: mã coin, ví dụ "BTCUSDT".
     - `interval`: khoảng thời gian (ví dụ: "1m" cho 1 phút).
     - `limit`: (tuỳ chọn) số lượng nến.
   - **Output**: Dữ liệu Kline hoặc lỗi nếu không thể lấy dữ liệu.

### 6. **/menu**
   - **Công dụng**: Hiển thị menu chính.
   - **Input**: Không yêu cầu tham số.
   - **Output**: Menu chính của bot.

### 7. **/p <symbol>**
   - **Công dụng**: Tìm kiếm và chọn mã tài sản gần nhất với `symbol` đã nhập.
   - **Input**: `symbol` - mã coin.
   - **Output**: Menu thông tin tài sản, hoặc thông báo nếu không tìm thấy.

### 8. **/price_spot <symbol>**
   - **Công dụng**: Lấy giá Spot cho coin.
   - **Input**: `symbol` - mã coin.
   - **Output**: Giá Spot của coin trong thời gian thực.

### 9. **/price_futures <symbol>**
   - **Công dụng**: Lấy giá Futures cho coin.
   - **Input**: `symbol` - mã coin.
   - **Output**: Giá Futures của coin trong thời gian thực.

### 10. **/funding_rate <symbol>**
   - **Công dụng**: Lấy Funding Rate cho một symbol.
   - **Input**: `symbol` - mã coin.
   - **Output**: Funding Rate cho một symbol trong thời gian thực.

### 11. **/kline_realtime <symbol> <interval>**
   - **Công dụng**: Theo dõi dữ liệu Kline theo thời gian thực.
   - **Input**:
     - `symbol`: mã coin.
     - `interval`: khoảng thời gian.
   - **Output**: Bắt đầu cập nhật dữ liệu Kline thời gian thực.

### 12. **/stop**
   - **Công dụng**: Dừng cập nhật thời gian thực.
   - **Input**: Không yêu cầu tham số.
   - **Output**: Thông báo đã dừng cập nhật dữ liệu Kline.

### 13. **/all_triggers**
   - **Công dụng**: Hiển thị tất cả các trigger (báo động).
   - **Input**: Không yêu cầu tham số.
   - **Output**: Danh sách tất cả các trigger đã cài đặt.

### 14. **/delete_trigger <spot/future/price-difference/funding-rate> <symbol>**
   - **Công dụng**: Xoá một trigger dựa trên loại giá và mã coin.
   - **Input**:
     - `spot/future/price-difference/funding-rate`: loại trigger.
     - `symbol`: mã tài sản.
   - **Output**: Thông báo xoá trigger thành công hoặc lỗi nếu không tìm thấy trigger.

### 15. **/alert_price_with_threshold <spot/future> <lower/above> <symbol> <threshold>**
   - **Công dụng**: Cài đặt cảnh báo khi giá vượt qua ngưỡng.
   - **Input**:
     - `spot/future`: loại giá.
     - `lower/above`: điều kiện cảnh báo.
     - `symbol`: mã coin.
     - `threshold`: ngưỡng giá.
   - **Output**: Thông báo cảnh báo đã được cài đặt.

### 16. **/price_difference <lower/above> <symbol> <threshold>**
   - **Công dụng**: Cài đặt cảnh báo khi có sự chênh lệch giá.
   - **Input**:
     - `lower/above`: điều kiện chênh lệch.
     - `symbol`: mã coin.
     - `threshold`: ngưỡng chênh lệch.
   - **Output**: Thông báo cảnh báo chênh lệch giá đã được cài đặt.

### 17. **/funding_rate_change <lower/above> <symbol> <threshold>**
   - **Công dụng**: Cài đặt cảnh báo khi Funding Rate thay đổi.
   - **Input**:
     - `lower/above`: điều kiện thay đổi.
     - `symbol`: mã coin.
     - `threshold`: ngưỡng thay đổi.
   - **Output**: Thông báo cảnh báo Funding Rate đã được cài đặt.
