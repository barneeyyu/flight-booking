# 航班預訂系統 (Flight Booking System)

一個使用 Go 語言和 Gin 框架開發的簡單航班預訂 API 系統。

## 功能特色

- 🔍 航班搜尋功能
- ✈️ 航班預訂功能
- 📋 預訂狀態查詢
- 🎫 座位可用性檢查
- 📊 分頁搜尋結果

## 技術棧

- **後端框架**: Gin (Go)
- **資料庫**: SQLite
- **ORM**: GORM
- **語言**: Go 1.23.2

## 安裝與運行

### 前置需求

- Go 1.23.2 或更高版本
- Git

#### 安裝 Go

如果您還沒有安裝 Go，請按照以下步驟：

##### macOS (使用 Homebrew)
```bash
brew install go
```
##### Windows
1. 前往 [Go 官方網站](https://golang.org/dl/) 下載 Windows 版本
2. 執行 `.msi` 安裝檔案


##### 驗證安裝
安裝完成後，在終端機執行：
```bash
go version
```
應該會顯示類似 `go version go1.23.2 darwin/amd64` 的訊息。

### 安裝步驟

1.  Clone project
    ```bash
    git clone https://github.com/your-username/flight-booking.git
    cd flight-booking
    ```

2.  安裝依賴
    ```bash
    go mod tidy
    ```

3.  建置應用程式
    ```bash
    make build
    ```

4.  資料填充 (Seeding Data)

    首次建置專案時需執行此語法將 mock 的機票 data insert 進 table（只需執行一次）
    ```bash
    make seed
    ```  
    > 這將會向資料庫中插入約 1000 筆航班資料。每次執行都會新增資料，請注意避免重複。
    >
    > 如果您想清空資料庫並重新填充，可以先執行 `make clean` 再執行 `make seed`。

5.  運行應用程式
    ```bash
    make run
    ```
    應用程式將在 `http://localhost:8080` 啟動。


## API 端點

### 1. 搜尋航班
```
GET /flights?departure_airport=TPE&arrival_airport=HKG&date=2024-01-15&page=1&page_size=10
```

查詢參數：
- `departure_airport`: 出發機場代碼
- `arrival_airport`: 抵達機場代碼
- `airline`: 航空公司
- `date`: 出發日期 (YYYY-MM-DD)
- `page`: 頁碼 (預設: 1)
- `page_size`: 每頁筆數 (預設: 10)

> **注意：** 由於資料填充 (Seeding Data) 限制，目前可查詢的航班資料特性如下：
> -   **城市：** 僅限於 `Taipei`, `Tokyo`, `Seoul`, `Singapore`, `Hong Kong`。
> -   **航空公司：** 僅限於 `EVA Air`, `China Airlines`, `Japan Airlines`, `All Nippon Airways`, `Korean Air`, `Asiana Airlines`, `Singapore Airlines`, `Cathay Pacific`。
> -   **日期：** 僅限於 2025 年 8 月份。

### 2. 建立預訂
```
POST /bookings
```

請求體範例：
```json
{
  "flight_id": 1,
  "passenger_name": "張三"
}
```

### 3. 查詢預訂狀態
```
GET /bookings/:id
```

### 4. 查詢航班詳情
```
GET /flights/:id
```

## Postman Collection

您可以匯入此 Postman Collection 檔案來測試所有 API 端點：

[下載 Postman Collection](docs/postman/FlightBookingSystem.postman_collection.json)

## 資料庫結構

### Flight 模型
- `flight_number`: 航班號碼
- `departure_airport`: 出發機場
- `arrival_airport`: 抵達機場
- `departure_time`: 出發時間
- `arrival_time`: 抵達時間
- `airline`: 航空公司
- `price`: 票價
- `available_seats`: 可用座位數

### Booking 模型
- `flight_id`: 航班 ID
- `passenger_name`: 乘客姓名
- `booking_status`: 預訂狀態 (Confirmed/Waitlisted)

## 專案結構

```
flight-booking/
├── main.go                # 應用程式進入點，負責依賴注入和啟動
├── go.mod                 # Go 模組檔案
├── go.sum                 # 依賴版本鎖定
├── Makefile               # 專案自動化腳本
├── README.md              # 專案說明
├── cmd/
│   └── seed/
│       └── main.go        # 獨立的資料填充程式
├── docs/
│   └── postman/
│       └── FlightBookingSystem.postman_collection.json        # 獨立的資料填充程式
└── internal/
    ├── database/
    │   └── database.go    # 資料庫初始化和遷移邏輯
    ├── handler/
    │   ├── booking_handler.go # 預訂相關 API 處理函式
    │   └── flight_handler.go  # 航班相關 API 處理函式
    ├── models/
    │   └── models.go      # 資料模型定義 (Flight, Booking)
    └── repository/
        ├── booking_repository.go # 預訂資料庫操作介面與實作
        └── flight_repository.go  # 航班資料庫操作介面與實作
```

## 開發說明

- 系統使用 SQLite 作為資料庫，資料庫檔案會自動創建為 `flights.db`
- 預訂系統包含 20% 的超賣機率模擬真實情況
- 支援分頁搜尋以處理大量航班資料
- 搜尋航班 (GET /flights) 僅顯示部分欄位，若需完整航班資訊請使用查詢航班詳情 (GET /flights/:id) 端點。

## 授權

MIT License 