# Architecture Documentation

本文檔詳細說明航班預訂系統的架構設計、技術決策和實作細節。

## Architecture Overview

這是一個基於 Clean Architecture 原則設計的航班預訂系統，採用分層架構模式，將業務邏輯與基礎設施分離，確保系統的可測試性、可維護性和可擴展性。

### High-Level System Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    Flight Booking System                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │
│  │   Client    │    │   Client    │    │   Client    │     │
│  │ (Postman)   │    │ (Frontend)  │    │   (Mobile)  │     │
│  └─────────────┘    └─────────────┘    └─────────────┘     │
│         │                   │                   │          │
│         └───────────────────┼───────────────────┘          │
│                             │                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                 HTTP API Layer                         │ │
│  │  ┌─────────────┐  ┌─────────────┐                     │ │
│  │  │   Flight    │  │   Booking   │                     │ │
│  │  │  Handler    │  │  Handler    │                     │ │
│  │  └─────────────┘  └─────────────┘                     │ │
│  └─────────────────────────────────────────────────────────┘ │
│                             │                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                Business Logic Layer                    │ │
│  │  ┌─────────────┐  ┌─────────────┐                     │ │
│  │  │   Flight    │  │   Booking   │                     │ │
│  │  │  Service    │  │  Service    │                     │ │
│  │  └─────────────┘  └─────────────┘                     │ │
│  └─────────────────────────────────────────────────────────┘ │
│                             │                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                Data Access Layer                       │ │
│  │  ┌─────────────┐  ┌─────────────┐                     │ │
│  │  │   Flight    │  │   Booking   │                     │ │
│  │  │ Repository  │  │ Repository  │                     │ │
│  │  └─────────────┘  └─────────────┘                     │ │
│  └─────────────────────────────────────────────────────────┘ │
│                             │                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                 Database Layer                         │ │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │ │
│  │  │   Flight    │  │   Booking   │  │   Indexes   │    │ │
│  │  │   Table     │  │   Table     │  │             │    │ │
│  │  └─────────────┘  └─────────────┘  └─────────────┘    │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Component Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      main.go                               │
│              (Dependency Injection)                        │
└─────────────────┬───────────────────────────────────────────┘
                  │
    ┌─────────────▼─────────────┐
    │         Router            │
    │    (Route Definition)     │
    └─────────────┬─────────────┘
                  │
    ┌─────────────▼─────────────┐
    │       Handlers            │
    │   (HTTP Request/Response) │
    │                           │
    │  ┌─────────┐ ┌─────────┐  │
    │  │ Flight  │ │ Booking │  │
    │  │Handler  │ │Handler  │  │
    │  └─────────┘ └─────────┘  │
    └─────────────┬─────────────┘
                  │
    ┌─────────────▼─────────────┐
    │       Services            │
    │   (Business Logic)        │
    │                           │
    │  ┌─────────┐ ┌─────────┐  │
    │  │ Flight  │ │ Booking │  │
    │  │Service  │ │Service  │  │
    │  └─────────┘ └─────────┘  │
    └─────────────┬─────────────┘
                  │
    ┌─────────────▼─────────────┐
    │     Repositories          │
    │   (Data Access)           │
    │                           │
    │  ┌─────────┐ ┌─────────┐  │
    │  │ Flight  │ │ Booking │  │
    │  │  Repo   │ │  Repo   │  │
    │  └─────────┘ └─────────┘  │
    └─────────────┬─────────────┘
                  │
    ┌─────────────▼─────────────┐
    │         Models            │
    │   (Data Structures)       │
    │                           │
    │  ┌─────────┐ ┌─────────┐  │
    │  │ Flight  │ │ Booking │  │
    │  │ Model   │ │ Model   │  │
    │  └─────────┘ └─────────┘  │
    └─────────────┬─────────────┘
                  │
    ┌─────────────▼─────────────┐
    │      Database             │
    │    (SQLite + GORM)        │
    └───────────────────────────┘
```

## Project Structure

```
flight-booking/
├── main.go                # 應用程式進入點，負責依賴注入和啟動
├── go.mod                 # Go 模組檔案
├── go.sum                 # 依賴版本鎖定
├── Makefile               # 專案自動化腳本
├── README.md              # 專案說明
├── ARCHITECTURE.md        # 架構設計文檔
├── CLAUDE.md              # Claude Code 指南
├── cmd/
│   └── seed/
│       └── main.go        # 獨立的資料填充程式
├── docs/
│   └── postman/
│       └── *.json         # Postman Collection 檔案
└── internal/              # 內部應用程式代碼
    ├── database/
    │   └── database.go    # 資料庫初始化和遷移邏輯
    ├── handler/           # HTTP 處理層 (Controller)
    │   ├── booking_handler.go     # 預訂相關 API 處理函式
    │   ├── booking_handler_test.go
    │   ├── flight_handler.go      # 航班相關 API 處理函式
    │   └── flight_handler_test.go
    ├── models/            # 資料模型定義
    │   └── models.go      # Flight, Booking 結構體定義
    ├── repository/        # 資料存取層 (Data Access)
    │   ├── booking_repository.go  # 預訂資料庫操作介面與實作
    │   └── flight_repository.go   # 航班資料庫操作介面與實作
    ├── router/            # 路由配置
    │   └── router.go      # Gin 路由設定和中介軟體
    └── service/           # 業務邏輯層 (Business Logic)
        └── booking_service.go     # 預訂服務邏輯介面與實作
```

### Layer Responsibilities

| Layer | Responsibility | Examples |
|-------|---------------|----------|
| **Handler** | HTTP 請求處理、參數驗證、回應格式化 | `CreateBooking()`, `SearchFlights()` |
| **Service** | 業務邏輯、事務管理、跨領域操作 | 超賣邏輯、價格計算、狀態管理 |
| **Repository** | 資料存取、查詢構建、ORM 封裝 | `FindAll()`, `Create()`, `Update()` |
| **Model** | 資料結構定義、驗證規則、關聯關係 | `Flight`, `Booking` 結構體 |

## 設計取捨（Trade-off Analysis）

本專案在設計時考量了多種情境，並做出以下取捨：

### 1. 並發控制策略

#### 悲觀鎖（Pessimistic Lock）
為確保高併發下的座位數一致性，採用資料庫悲觀鎖。

```go
// 使用 SELECT ... FOR UPDATE
if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
    Where("id = ?", booking.FlightID).
    First(&flight).Error; err != nil {
    return fmt.Errorf("failed to lock flight: %w", err)
}
```

**優點:**
- 強一致性保證
- 避免超賣問題
- 實作簡單直觀

**缺點:**
- 降低系統吞吐量
- 可能造成鎖競爭與延遲
- 不適合極高併發場景

**替代方案:**
樂觀鎖透過版本號或 updated_at 欄位，先查詢再更新時比對版本，失敗就重試。

### 2. 架構模式選擇

#### Repository/Service Interface 設計
採用 interface + struct 的分層設計。

```go
type BookingRepository interface {
    Create(booking *models.Booking) error
    FindByID(id uint) (*models.Booking, error)
    Update(booking *models.Booking) error
}
```

**優點:**
- 提升可測試性與彈性
- 依賴反轉，便於 mock
- 支援不同實作（如切換資料庫）

**缺點:**
- 初期開發需要更多樣板程式碼
- 增加程式碼複雜度

### 3. 資料設計考量

#### 自動遞增 ID
目前資料表主鍵 id 採用自動遞增設計。

**優點:**
- 簡單易用
- 效能佳（整數比較快）
- 天然排序

**潛在問題:**
- 高併發或分散式系統下容易產生衝突或熱點
- id 可預測，安全性較低
- 改善方案：UUID、雪花算法

### 4. ORM 依賴程度

#### Service 層直接依賴 ORM Transaction
操作彈性大，但與 ORM 耦合較深。

**優點:**
- 操作彈性大
- 事務控制精確
- 開發效率高

**缺點:**
- 與 ORM 耦合較深
- 未來更換 ORM 重構成本較高

## 資料庫設計

### 資料模型關係

```
┌─────────────────┐       ┌─────────────────┐
│     Flight      │   1   │     Booking     │
│                 │ ───── │                 │
│ - id (PK)       │   n   │ - id (PK)       │
│ - flight_number │       │ - flight_id (FK)│
│ - departure_*   │       │ - passenger_*   │
│ - arrival_*     │       │ - quantity      │
│ - price         │       │ - status        │
│ - available_*   │       │ - total_price   │
└─────────────────┘       └─────────────────┘
```

### 索引設計

為了提升航班搜尋效能，我們添加了以下索引：

#### Flight 表索引
- **複合索引** `idx_flight_search`: (departure_airport, arrival_airport, departure_time)
- **單欄位索引**: flight_number, airline, price

#### Booking 表索引  
- **複合索引** `idx_booking_search`: (flight_id, passenger_name)
- **單欄位索引**: booking_status

### 效能驗證

#### 1. 檢查索引是否建立成功
```bash
sqlite3 flights.db ".schema" | grep -E "(CREATE INDEX|idx_)"
```

#### 2. 驗證查詢計劃使用索引
```bash
sqlite3 flights.db "EXPLAIN QUERY PLAN SELECT * FROM flights WHERE departure_airport = 'Taipei' AND arrival_airport = 'Tokyo';"
```

**預期輸出:**
```
QUERY PLAN
`--SEARCH flights USING INDEX idx_flight_search
```

#### 3. 效能對比
- **有索引前**: 全表掃描 (SCAN TABLE flights)
- **有索引後**: 索引掃描 (SEARCH flights USING INDEX)
- **預期提升**: 查詢速度提升 10-100 倍

## 業務邏輯設計

### 超賣機制

```go
const (
    BookingStatusConfirmed  = "Confirmed"
    BookingStatusWaitlisted = "Waitlisted"
)

// 超賣邏輯
if flight.AvailableSeats >= booking.Quantity {
    booking.BookingStatus = BookingStatusConfirmed
} else if flight.AvailableSeats+oversellLimit >= booking.Quantity {
    booking.BookingStatus = BookingStatusWaitlisted
} else {
    return fmt.Errorf("not enough seats")
}
```

### 狀態流轉

```
[預訂請求] 
    ↓
[檢查座位]
    ↓
┌─────────────┬─────────────┬─────────────┐
│ 座位充足    │ 可超賣範圍  │ 超過限制    │
│ Confirmed   │ Waitlisted  │ Rejected    │
└─────────────┴─────────────┴─────────────┘
```

## 擴展性考量

### 未來優化方向

1. **付款流程整合**
   - 新增 PaymentStatus 欄位
   - 整合第三方支付 API
   - 實作付款超時取消機制

2. **通知系統**
   - Email/SMS 通知服務
   - 候補轉正通知
   - 異步訊息處理

3. **快取策略**
   - Redis 快取熱門航班
   - 查詢結果快取
   - 分散式快取設計

4. **監控與可觀測性**
   - 結構化日誌
   - Metrics 收集
   - 分散式追蹤

## 開發指引

### 測試策略
- **單元測試**: Handler, Service, Repository 層
- **整合測試**: 資料庫操作、API 端點
- **效能測試**: 並發預訂場景

### 程式碼品質
- 遵循 Go 語言慣例
- 介面設計優於實作
- 錯誤處理要明確
- 測試覆蓋率 > 80%

### 部署考量
- 容器化（Docker）
- 健康檢查端點
- 優雅關閉機制
- 配置外部化