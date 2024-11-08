## Hướng dẫn cài đặt

Sử dụng .env.template và thay thế bằng token tele của bạn

Tạo thêm file .env và để theo format .env.template

Tải xuống các package: `go mod tidy`

Để khởi động bot, cd vào thư mục src, chạy `go run main.go`

Để khởi động mock be server, cd vào thư mục mock-be

```
npm install
npm run start:dev
```

## Hướng dẫn local test webhook

### Cài đặt ngrok bằng choco
Tham khảo: https://ngrok.com/docs/getting-started/?os=windows
Chú ý tải choco, chạy Powershell/Cmp bằng quyền admin


### Local test
Tạo URL ngrok bằng lệnh: ```ngrok http 8443```

Copy đường dẫn tại Forwarding URL có dạng: https://ngrokurl.ngrok-free.app, thêm /webhook và copy vào WEBHOOK_URL của file .env

Chạy bot bằng lệnh ```go run main.go```

### Local test run bằng docker
Sử dụng lệnh ```docker build -t coin-price-tele .``` để build image

Sử dụng lệnh ```docker run -p 8443:8443 coin-price-tele``` để chạy container