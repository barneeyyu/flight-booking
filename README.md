# 航班預訂系統 (Flight Booking System)

一個使用 Go 語言和 Gin 框架開發的簡單航班預訂 API 系統。

## 功能特色

- 🔍 航班搜尋功能
- 📊 分頁搜尋結果
- ✈️ 航班預訂功能
- ➕ 支援超賣預訂
- 📋 查詢預訂詳情

## 開發說明

- 使用 SQLite 資料庫，支援自動遷移和索引優化
- 實現航班預訂的超賣和候補機制，確保高併發下的資料一致性
- 支援分頁搜尋以處理大量航班資料
- 詳細架構設計和技術決策請參考 [ARCHITECTURE.md](ARCHITECTURE.md)
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

6.  跑單元測試
    ```bash
    make test
    ```


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

## 資料庫

- **資料庫**: SQLite (flights.db)
- **模型**: Flight (航班), Booking (預訂)
- **特性**: 事務控制、索引優化、並發安全
- 詳細資料庫設計請參考 [ARCHITECTURE.md](ARCHITECTURE.md)

## 授權

MIT License 