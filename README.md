# zaipra

A library helps you improve the speed of getting responses when asking LLM.

**zaipra** 是一個協助您在與大型語言模型（LLM）互動時提高回應速度的 Go 語言函式庫。它通過智慧資訊路由來優化提示工程。

## 安裝

```bash
go get github.com/TimLai666/zaipra
```

## 功能

- **資訊路由 (Info Routing)**: 自動選擇與用戶問題相關的資訊，減少不必要的prompt內容。
<!-- - **提示擴散 (Prompt Diffusion)**: 優化提示結構以獲得更快更準確的回應 -->

## 與 LangChain Go 集成

`zaipra` 設計為與 [langchaingo](https://github.com/tmc/langchaingo) 協同工作。您需要提供一個實現了 `llms.Model` 介面的 LangChain Go 模型。

### 支持的 LangChain Go 模型

`zaipra` 可與任何實現了 `langchaingo/llms.Model` 接口的模型一起使用，包括：

- Ollama
- OpenAI
- Anthropic
- 以及其他 langchaingo 支援的模型

## 使用方法

以下是一個基本使用示例：

```go
package main

import (
    "fmt"

    "github.com/TimLai666/zaipra"
    "github.com/tmc/langchaingo/llms/ollama"
)

func main() {
    // 初始化 LangChain Go 的 LLM 模型
    llm, err := ollama.New(
        ollama.WithServerURL("http://localhost:11434"),
        ollama.WithModel("gemma3:12b"),
    )
    if err != nil {
        panic(err)
    }
    
    // 準備信息
    infos := []zaipra.Info{
        {
            Title:       "產品規格",
            Description: "產品的詳細規格和特性",
            Content:     "這是一個高性能的產品，具有多種功能...",
        },
        {
            Title:       "價格資訊",
            Description: "產品的價格和支付方式",
            Content:     "產品價格為 $299，支持信用卡和 PayPal 支付...",
        },
    }
    
    // 使用 zaipra 生成回答
    answer, err := zaipra.Answer(
        "客戶提問",         // 提問名稱
        "這個產品多少錢？", // 用戶問題
        "你是客服人員，請提供專業的回答", // 系統提示
        infos,
        llm,
    )
    
    if err != nil {
        panic(err)
    }
    
    fmt.Println("回答:", answer)
}
```

## 工作原理

1. **資訊路由**: zaipra 分析用戶問題，並從提供的資訊中選擇相關內容。
<!-- 2. **提示擴散**: 將選擇的信息與用戶問題和系統提示組合，生成優化的完整提示 -->
2. **LLM 回應**: 使用 langchaingo 將優化後的提示發送給 LLM 並取得回應。

## 許可證

MIT License