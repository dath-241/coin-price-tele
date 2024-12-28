# Coin-Price-Telegram-Bot

## Thông tin về nhóm

| STT | Tên thành viên        | Vai trò       | Mã số sinh viên | GitHub                                             |
| --- | --------------------- | ------------- | --------------- | -------------------------------------------------- |
| 1   | Trần Nguyễn Thanh Lâm | Product Owner | 2211822         | [Github](https://github.com/clgslsm)               |
| 2   | Nguyễn Trung Tín      | Developer     | 2213500         | [Github](https://github.com/TinnieTheCat198)       |
| 3   | Thái Thành Duy        | Developer     | 2210535         | [Github](https://github.com/ShaKk0722)             |
| 4   | Nguyễn Đăng Khoa      | Developer     | 2211621         | [Github](https://github.com/NguyenDangKhoaDepTrai) |
| 5   | Lê Thành Đạt          | Developer     | 2210683         | [Github](https://github.com/thnhdt)               |
| 6   | Nguyễn Hữu Đăng Khoa  | Developer     | 2211625         | [Github](https://github.com/thanhlam2000)          |
| 7   | Dương Hoàng Long      | Developer     | 2211873         | [Github](https://github.com/Long-noop)             |

## Repository

[Github](https://github.com/dath-241/coin-price-tele)

## Cấu trúc thư mục:

¦   docker-compose.yml
¦   Dockerfile
¦   package-lock.json
¦   README.MD
¦   tree.txt
¦   
+---.github
¦   +---workflows
¦           deploy-DO-vps-main.yml
¦           deploy-heroku-dev.yml
¦           
+---.vscode
¦       settings.json
¦       
+---docs
¦       .nojekyll
¦       command.md
¦       index.html
¦       install.md
¦       overview.md
¦       README.md
¦       _sidebar.md
¦       
+---report
¦       meeting-minutes-16-10-2024.md
¦       meeting-minutes-3-10-2024.md
¦       meeting-minutes-31-10-2024.md
¦       outline-report-1-11-2024.md
¦       outline-report-10-12-2024.md
¦       
+---src
    ¦   .env
    ¦   .env.template
    ¦   .gitignore
    ¦   .golangci.yml
    ¦   go.mod
    ¦   go.sum
    ¦   main.go
    ¦   
    +---.vscode
    ¦       settings.json
    ¦       
    +---bot
    ¦   ¦   bot.go
    ¦   ¦   
    ¦   +---handlers
    ¦           button.go
    ¦           candlestick.go
    ¦           chromedp-custom.go
    ¦           command.go
    ¦           constants.go
    ¦           getKline.go
    ¦           getPrices.go
    ¦           getUserID.go
    ¦           menuPrices.go
    ¦           RegisterCommand.go
    ¦           user.go
    ¦           websocket_stream.go
    ¦           
    +---cache
    ¦       volume_cache.go
    ¦       
    +---config
    ¦       config.go
    ¦       
    +---services
    ¦       auth.go
    ¦       auth_test.go
    ¦       constant.go
    ¦       database.go
    ¦       database_test.go
    ¦       fetchSymbol.go
    ¦       futures_symbols.txt
    ¦       futures_symbols_sorted.txt
    ¦       sort_symbol.go
    ¦       spot_symbols.txt
    ¦       spot_symbols_sorted.txt
    ¦       
    +---test
            .env
            cache.db
            main.go
            session.dat
            


## Liên hệ

---


