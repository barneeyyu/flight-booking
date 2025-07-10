# 航班預訂系統 (Flight Booking System)

一個使用 Go 語言和 Gin 框架開發的簡單航班預訂 API 系統。

## 功能特色

- 🔍 航班搜尋功能
- 📊 分頁搜尋結果
- ✈️ 航班預訂功能
- 📋 查詢預訂詳情

## 開發說明

- 系統使用 SQLite 作為資料庫，資料庫檔案會自動創建為 `flights.db`
- 預訂系統使用資料庫事務和悲觀鎖 (SELECT ... FOR UPDATE) 來確保高併發下的座位數量一致性，避免超賣。
- 支援分頁搜尋以處理大量航班資料
- 搜尋航班 (GET /flights) 僅顯示部分欄位，若需完整航班資訊請使用查詢航班詳情 (GET /flights/:id) 端點。
- 專案程式碼中有多處以 `// TODO:` 註解標示，說明目前尚未實作但未來可擴充的功能（如複合索引、訂單狀態流轉、付款流程、通知用戶等），**可參考這些 TODO 註解**，作為後續優化與擴展的依據。


## 設計取捨（Trade-off Analysis）

本專案在設計時考量了多種情境，並做出以下取捨：

- **悲觀鎖（Pessimistic Lock）**  
  為確保高併發下的座位數一致性，採用資料庫悲觀鎖。雖然能避免超賣，但會降低吞吐量，若流量極大可能造成鎖競爭與延遲。
  > 如要因應高併發，可選用**樂觀鎖**，透過版本號（version）或 updated_at 欄位，先查詢再更新時比對版本，失敗就重試。

- **Repository/Service Interface 設計**  
  採用 interface + struct 的分層設計，提升可測試性與彈性，但初期開發會多一些樣板程式碼。

- **欄位設計的彈性**  
  欄位設計簡單，易於維護，但未來若需支援多狀態、付款、通知等，需再擴充，並且做migration。

- **Service 層直接依賴 ORM Transaction**  
  操作彈性大，但與 ORM 耦合較深，未來若要更換 ORM 或資料庫，重構成本較高。

- **Handler 層參數驗證**  
  目前直接在 handler 驗證，簡單直觀，但若驗證邏輯複雜，建議抽出 validator 層。

- **自動遞增 ID（Auto-increment ID）**  
  目前資料表主鍵 id 採用自動遞增（auto-increment）設計，雖然簡單易用，但有以下潛在缺點：
  - 在高併發或分散式系統下，容易產生衝突或熱點（hotspot），影響效能
  - id 可預測，可能被惡意用戶猜測、暴力查詢（安全性較低）
  - 若需改善，可考慮改用 UUID、雪花算法（Snowflake ID）等分散式唯一識別碼

- **高流量高併發所添加的技術選型**  
  - 先寫入快取（如 Redis），再異步寫入資料庫
    - 適合極高流量，減少 DB 壓力。
    - 但要考慮資料一致性與最終一致性。
  - 佇列排隊（如 Kafka、RabbitMQ）
    - 訂單請求先進佇列，後端消費者依序處理，避免同時搶票。
    - 適合搶票、秒殺等場景。
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
    > 這將會向資料庫中插入約 1000 筆航班資料。
    >
    > 每次執行 `make seed` 都會新增資料，請注意避免重複。
    > 如果您想清空資料庫並重新填充，可以先執行 `make clean` 再執行 `make seed`。

5.  運行應用程式
    ```bash
    make run
    ```
    應用程式將在 `http://localhost:8080` 啟動。


## API 端點

### 1. 搜尋航班
```
GET /flights?departure=TPE&arrival=HKG&date=2024-01-15&page=1&page_size=10
```

查詢參數：
- `departure`: 出發機場代碼
- `arrival`: 抵達機場代碼
- `airline`: 航空公司
- `date`: 出發日期 (YYYY-MM-DD)
- `page`: 頁碼 (預設: 1)
- `page_size`: 每頁筆數 (預設: 10)

> **注意：** 由於資料填充 (Seeding Data) 限制，目前可查詢的航班資料特性如下：
> -   **城市：** 僅限於 `Taipei`, `Tokyo`, `Seoul`, `Singapore`, `Hong Kong`。
> -   **航空公司：** 僅限於 `EVA Air`, `China Airlines`, `Japan Airlines`, `All Nippon Airways`, `Korean Air`, `Asiana Airlines`, `Singapore Airlines`, `Cathay Pacific`。
> -   **日期：** 僅限於 2025 年 8 月份。

### 2. 查詢航班詳情
```
GET /flights/:id
```

### 3. 建立預訂
```
POST /bookings
```

請求體範例：
```json
{
  "flight_id": 1,
  "passenger_name": "張三",
  "quantity": 1
}
```

### 4. 查詢預訂狀態
```
GET /bookings/:id
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
- `quantity`: 購買票數
- `total_price`: 總價

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
│       └── FlightBookingSystem.postman_collection.json        # Postman Collection 檔案
└── internal/
    ├── database/
    │   └── database.go    # 資料庫初始化和遷移邏輯
    ├── handler/
    │   ├── booking_handler.go # 預訂相關 API 處理函式
    │   └── flight_handler.go  # 航班相關 API 處理函式
    ├── models/
    │   └── models.go      # 資料模型定義 (Flight, Booking)
    ├── repository/
    │   ├── booking_repository.go # 預訂資料庫操作介面與實作
    │   └── flight_repository.go  # 航班資料庫操作介面與實作
    ├── router/
    │   └── router.go # Gin 路由設定
    └── service/
        └── booking_service.go # 預訂服務邏輯介面與實作（商業邏輯層）
```

## 授權

MIT License 