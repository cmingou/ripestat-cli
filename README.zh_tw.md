# ripestat-cli

一個基於 Go 語言的命令列工具，提供 RIPEstat API 的簡單包裝器。可直接從終端查詢 ASN（自治系統號碼）、IPv4 和 IPv6 位址的相關資訊。

[![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 功能特色

- **多重輸入支援**：可在單一指令中同時查詢 ASN、IPv4 和 IPv6 位址
- **自動辨識**：自動將輸入分類為 ASN、IPv4 或 IPv6
- **豐富資訊**：提供完整的資料，包括：
  - ASN 概覽和詳細資訊
  - 區域網際網路註冊機構 (RIR) 資訊
  - BGP 路由一致性資料
  - 地理位置資訊
- **表格輸出**：清晰的格式化表格顯示，便於閱讀

## 安裝方式

### 從原始碼安裝

```bash
git clone https://github.com/yourusername/ripestat-cli.git
cd ripestat-cli
make install
```

### 從原始碼建構

```bash
# 為目前平台建構
go build -o ripestat main.go

# 為所有平台建構
make all

# 為特定平台建構
make darwin  # macOS
make linux   # Linux
```

## 使用方法

工具接受多個參數並自動偵測其類型：

```bash
# 同時查詢多個資源
./ripestat 8.8.8.8 13335 2001:db8::1

# 查詢個別資源
./ripestat 8.8.8.8                # IPv4 位址
./ripestat 13335                   # ASN
./ripestat 2001:db8::1            # IPv6 位址
```

### 輸出範例

```
╭─────────────────────────────────────────────────────────────╮
│                        AS Information                        │
├─────────────────────────────────────────────────────────────┤
│ ASN: 13335                                                  │
│ Name: CLOUDFLARENET                                         │
│ Country: US                                                 │
│ Registry: ARIN                                              │
╰─────────────────────────────────────────────────────────────╯
```

## API 整合

此工具整合了多個 RIPEstat API 端點：

- **AS 概覽**：基本 ASN 資訊和元資料
- **RIR**：區域網際網路註冊機構資料
- **前綴路由一致性**：BGP 路由資訊
- **MaxMind GeoLite**：IP 位址的地理位置資料

## 開發

### 先決條件

- Go 1.21 或更新版本
- Make（選用，用於使用 Makefile 指令）

### 建構

```bash
# 清理先前的建構
make clean

# 執行測試
make test

# 程式碼檢查
make lint

# 開發建構
go build -o ripestat main.go
```

### 專案結構

```
├── main.go              # 進入點
├── cmd/
│   ├── root.go         # 主要指令邏輯 (Cobra CLI)
│   └── root_test.go    # 輸入驗證函式測試
├── pkg/
│   └── ripestat/       # API 客戶端套件
│       ├── api.go      # HTTP 客戶端函式
│       └── struts.go   # 回應結構體
└── internal/
    └── utils/          # 工具函式
        ├── asn.go      # ASN 查詢
        ├── ipv4.go     # IPv4 查詢
        ├── ipv6.go     # IPv6 查詢
        ├── util.go     # 通用工具
        └── util_test.go # 工具函式測試
```

## 相依性

- [github.com/spf13/cobra](https://github.com/spf13/cobra) - CLI 框架
- [github.com/olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) - 表格格式化
- Go 標準函式庫（網路和 HTTP）

## 貢獻

1. Fork 此儲存庫
2. 建立您的功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 開啟 Pull Request

## 授權

此專案採用 MIT 授權 - 詳見 [LICENSE](LICENSE) 檔案。

## 致謝

- [RIPE NCC](https://www.ripe.net/) 提供的 RIPEstat API
- Go 社群提供的優秀函式庫和工具

---

English version: [README.md](README.md)