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

### 安裝步驟

1. clone project
```bash
git clone https://github.com/your-username/flight-booking.git
cd flight-booking
```

2. 安裝依賴
```bash
go mod tidy
```

3. 運行應用程式
```bash
go run main.go models.go
```

應用程式將在 `http://localhost:8080` 啟動

## API 端點

### 1. 健康檢查
```
GET /ping
```

### 2. 搜尋航班
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

### 3. 建立預訂
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

### 4. 查詢預訂狀態
```
GET /bookings/:id
```

### 5. 查詢航班詳情
```
GET /flights/:id
```

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
├── main.go          # 主程式和 API 路由
├── models.go        # 資料模型定義
├── go.mod           # Go 模組檔案
├── go.sum           # 依賴版本鎖定
└── README.md        # 專案說明
```

## 開發說明

- 系統使用 SQLite 作為資料庫，資料庫檔案會自動創建為 `flights.db`
- 預訂系統包含 20% 的超賣機率模擬真實情況
- 支援分頁搜尋以處理大量航班資料

## 授權

MIT License 