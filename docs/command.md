# List command

### 1. **/help** ✅

- **Công dụng**: Hiển thị danh sách các lệnh hỗ trợ.
- **Input**: Không yêu cầu tham số.
- **Output**: Danh sách các lệnh khả dụng dưới dạng văn bản.

### 2. **/start** ✅

- **Công dụng**: Xác thực người dùng.
- **Input**: Không yêu cầu tham số.
- **Output**: "Access denied." nếu xác thực thất bại; nếu thành công, trả về thông báo xác nhận.

### 3. **/login &lt;username&gt; &lt;password&gt;** ✅

- **Công dụng**: Đăng nhập với tài khoản người dùng.
- **Input**: `username` (tên người dùng), `password` (mật khẩu).
- **Output**: "Đăng nhập thành công." hoặc lỗi nếu đăng nhập thất bại.

### 4. **/getinfo** ✅

- **Công dụng**: Lấy thông tin người dùng.
- **Input**: Không yêu cầu tham số.
- **Output**: Hiển thị thông tin người dùng hoặc thông báo lỗi nếu có sự cố.

### 5. **/kline &lt;symbol&gt; &lt;interval&gt; &lt;limit&gt;** ✅

- **Công dụng**: Lấy dữ liệu Kline (biểu đồ nến) của coin.
- **Input**:
  - `symbol`: mã coin, ví dụ "BTCUSDT".
  - `interval`: khoảng thời gian (ví dụ: "1m" cho 1 phút).
  - `limit`: (tuỳ chọn) số lượng nến.
- **Output**: Dữ liệu Kline hoặc lỗi nếu không thể lấy dữ liệu.

### 6. **/menu** ✅

- **Công dụng**: Hiển thị menu chính.
- **Input**: Không yêu cầu tham số.
- **Output**: Menu chính của bot.

### 7. **&lt;symbol&gt;** ✅

- **Công dụng**: Tìm kiếm và chọn mã tài sản gần nhất với `symbol` đã nhập.
- **Input**: `symbol` - mã coin.
- **Output**: Giá spot và futures trong thời gian thực, hoặc thông báo không tìm thấy.

### 8. **/p &lt;symbol&gt;** ✅

- **Công dụng**: Tìm kiếm và chọn mã tài sản gần nhất với `symbol` đã nhập.
- **Input**: `symbol` - mã coin.
- **Output**: Menu thông tin tài sản, hoặc thông báo nếu không tìm thấy.

### 9. **/price_spot &lt;symbol&gt;** ✅

- **Công dụng**: Lấy giá Spot cho coin.
- **Input**: `symbol` - mã coin.
- **Output**: Giá Spot của coin trong thời gian thực.

### 10. **/price_futures &lt;symbol&gt;** ✅

- **Công dụng**: Lấy giá Futures cho coin.
- **Input**: `symbol` - mã coin.
- **Output**: Giá Futures của coin trong thời gian thực.

### 11. **/funding_rate &lt;symbol&gt;** ✅

- **Công dụng**: Lấy Funding Rate cho một symbol.
- **Input**: `symbol` - mã coin.
- **Output**: Funding Rate cho một symbol trong thời gian thực.

### 12. **/kline_realtime &lt;symbol&gt; &lt;interval&gt;** ✅

- **Công dụng**: Theo dõi dữ liệu Kline theo thời gian thực.
- **Input**:
  - `symbol`: mã coin.
  - `interval`: khoảng thời gian.
- **Output**: Bắt đầu cập nhật dữ liệu Kline thời gian thực.

### 13. **/stop** ✅

- **Công dụng**: Dừng cập nhật thời gian thực.
- **Input**: Không yêu cầu tham số.
- **Output**: Thông báo đã dừng cập nhật dữ liệu Kline.

### 14. **/all_triggers** ✅

- **Công dụng**: Hiển thị tất cả các trigger (báo động).
- **Input**: Không yêu cầu tham số.
- **Output**: Danh sách tất cả các trigger đã cài đặt.

### 15. **/delete_trigger &lt;spot/future/price-difference/funding-rate&gt; &lt;symbol&gt;** ✅

- **Công dụng**: Xoá một trigger dựa trên loại giá và mã coin.
- **Input**:
  - `spot/future/price-difference/funding-rate`: loại trigger.
  - `symbol`: mã tài sản.
- **Output**: Thông báo xoá trigger thành công hoặc lỗi nếu không tìm thấy trigger.

### 16. **/alert_price_with_threshold &lt;spot/future&gt; &lt;lower/above&gt; &lt;symbol&gt; &lt;threshold&gt;** ✅
- **Công dụng**: Cài đặt cảnh báo khi giá vượt qua ngưỡng.
- **Input**:
  - `spot/future`: loại giá.
  - `lower/above`: điều kiện cảnh báo.
  - `symbol`: mã coin.
  - `threshold`: ngưỡng giá.
- **Output**: Thông báo cảnh báo đã được cài đặt.

### 17. **/price_difference &lt;lower/above&gt; &lt;symbol&gt; &lt;threshold&gt;** ✅

- **Công dụng**: Cài đặt cảnh báo khi có sự chênh lệch giá.
- **Input**:
  - `lower/above`: điều kiện chênh lệch.
  - `symbol`: mã coin.
  - `threshold`: ngưỡng chênh lệch.
- **Output**: Thông báo cảnh báo chênh lệch giá đã được cài đặt.

### 18. **/funding_rate_change &lt;lower/above&gt; &lt;symbol&gt; &lt;threshold&gt;** ✅

- **Công dụng**: Cài đặt cảnh báo khi Funding Rate thay đổi.
- **Input**:
  - `lower/above`: điều kiện thay đổi.
  - `symbol`: mã coin.
  - `threshold`: ngưỡng thay đổi.
- **Output**: Thông báo cảnh báo Funding Rate đã được cài đặt.

### 19. **/register** ✅

- **Công dụng**: Đăng ký tài khoản người dùng.
- **Input**: /signup <email> <name> <username> <password>
- **Output**: Thông báo đăng ký thành công.

### 20. **/forgotpassword** ✅

- **Công dụng**: Quên mật khẩu.
- **Input**: /forgotpassword <username>
- **Output**: Thông báo quên mật khẩu.

### 21. **/marketcap &lt;symbol&gt;**

- **Công dụng**: Lấy thông tin về vốn hóa thị trường của coin.
- **Input**: `symbol` - mã coin.
- **Output**: Thông tin về vốn hóa và thứ hạng của coin.

### 22. **/volume &lt;symbol&gt;**

- **Công dụng**: Lấy thông tin về khối lượng giao dịch.
- **Input**: `symbol` - mã coin.
- **Output**: Thông tin về khối lượng giao dịch 24h.

### 23. **/indicator &lt;symbol&gt; &lt;indicator_type&gt; &lt;params...&gt;**

- **Công dụng**: Tính toán giá trị các chỉ báo kỹ thuật.
- **Input**:
  - `symbol`: mã coin.
  - `indicator_type`: loại chỉ báo (MA, EMA, BOLL,...).
  - `params`: các tham số cho chỉ báo.
- **Output**: Giá trị chỉ báo được tính toán.

### 24. **/load_indicator &lt;file_path&gt;**

- **Công dụng**: Tải plugin chỉ báo tùy chỉnh.
- **Input**: `file_path` - đường dẫn đến file plugin.
- **Output**: Thông báo tải plugin thành công hoặc thất bại.

### 25. **/alert_indicator &lt;symbol&gt; &lt;indicator&gt; &lt;condition&gt; &lt;value&gt;**

- **Công dụng**: Cài đặt cảnh báo dựa trên chỉ báo kỹ thuật.
- **Input**:
  - `symbol`: mã coin.
  - `indicator`: loại chỉ báo.
  - `condition`: điều kiện cảnh báo.
  - `value`: giá trị ngưỡng.
- **Output**: Thông báo cảnh báo đã được cài đặt.

### 26. **/snooze &lt;trigger_id&gt; &lt;duration&gt;**

- **Công dụng**: Tạm dừng cảnh báo trong một khoảng thời gian.
- **Input**:
  - `trigger_id`: ID của cảnh báo.
  - `duration`: thời gian tạm dừng (phút).
- **Output**: Xác nhận tạm dừng cảnh báo.

### 27. **/repeat &lt;trigger_id&gt; &lt;type&gt; &lt;value&gt;**

- **Công dụng**: Cài đặt lặp lại cho cảnh báo.
- **Input**:
  - `trigger_id`: ID của cảnh báo.
  - `type`: loại lặp lại (times/schedule/forever).
  - `value`: số lần/lịch lặp lại.
- **Output**: Xác nhận cài đặt lặp lại.

### 28. **/mute [on/off]**

- **Công dụng**: Bật/tắt chế độ im lặng của bot.
- **Input**: `on/off` (tùy chọn).
- **Output**: Thông báo trạng thái im lặng của bot.

### 29. **/chat_config &lt;setting&gt; &lt;value&gt;**

- **Công dụng**: Cấu hình chat cho nhóm/topic.
- **Input**:
  - `setting`: loại cài đặt.
  - `value`: giá trị cài đặt.
- **Output**: Xác nhận cài đặt chat.

### 30. **/forgotpassword &lt;username&gt;**
   - **Công dụng**: quên mật khẩu
   - **Input**:
     - `username`: tên người dùng.
   - **Output**: Gửi OTP qua mail.

### 31. **/changepassword &lt;old_password&gt; &lt;new_password&gt; &lt;confirm_newpassword&gt;**
   - **Công dụng**: đổi mật khẩu
   - **Input**:
     - `old_password`: mật khẩu cũ.
     - `new_password`: mật khẩu mới.
     - `confirm_newpassword`: xác nhận mật khẩu mới.
   - **Output**: Xác nhận thành công, yêu cầu đăng nhập lại.

### 32. **/changeinfo** 
   - **Công dụng**: đổi thông tin
   - **Input**:
   - **Output**: đổi thông tin.
