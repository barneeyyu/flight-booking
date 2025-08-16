# CLAUDE.md — Best Practices

> 本文件由團隊資深工程師編寫，目的是確保程式碼品質、可測試性與長期維護性。
> 我們相信乾淨、可讀、可測試的程式碼比短期快速交付更重要。

---

## <critical_notes> 關鍵原則
- **以可測試性為首要考量**  
  - 所有業務邏輯必須可被單元測試覆蓋，避免硬編碼依賴（使用介面 / DI 注入）
  - 嚴禁將核心邏輯與 I/O、外部系統強耦合
- **遵循 Clean Code 原則**  
  - 單一職責（SRP）、短函數、語意清晰的命名
  - 移除重複程式碼（DRY）、避免不必要的複雜度（KISS）
  - 提前結束（guard clause）取代深層巢狀 if-else
- **安全優先**  
  - 所有外部輸入必須驗證
  - 憑證、金鑰與敏感資訊絕不進版本控制

---

## <paved_path> 推薦開發流程
1. **需求分析**  
   - 明確定義輸入、輸出與錯誤情境
2. **撰寫測試（TDD 優先）**  
   - 先定義測試案例與期望結果
3. **實作功能**  
   - 保持函數簡短、模組化、可單獨測試
4. **程式碼檢查**  
   - 必須通過 Linter、單元測試與整合測試
5. **Code Review**  
   - 對命名、可讀性、效能與可測試性給予回饋
6. **文件更新**  
   - 更新 README、API 規格與相關架構圖

---

## <patterns> 常見模式與實踐
- **Dependency Injection**：將依賴透過介面傳入以利測試  
- **Repository Pattern**：資料存取邏輯與業務邏輯分離  
- **Strategy Pattern**：針對可替換演算法的業務需求  
- **Factory Function**：建立具預設依賴的物件，避免直接使用 `new` 造成耦合

---

## <workflow> 新功能開發步驟
1. 建立 feature 分支（`feature/<name>`）
2. 撰寫或更新測試
3. 完成功能實作並確保本地測試全數通過
4. 發送 PR 並標記 reviewer
5. 經 Code Review 通過後合併至 main 分支

## <common_tasks> 常見任務指引
go test -race ./...
go vet ./...
go mod tidy