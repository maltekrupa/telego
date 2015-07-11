# Telego

A - yet - not so powerfull library to work with Telegram bots in Golang.

*THIS PROJECT IS DISCONTINUED BECAUSE I STARTED USING [tucnak/telebot](https://github.com/tucnak/telebot).*

## Installation

`go get github.com/temal-/telego`

## Example

An echo bot.

```go
package main

import (
    "fmt"
    "time"

    "github.com/temal-/telego"
)

func main() {
    tgo := telego.NewTelego("APIKEY")
    nMsg, _ := tgo.GetLastMessage()
    oMsg, _ := tgo.GetLastMessage()
    fmt.Println("Waiting for messages ...")
    for {
        nMsg, _ = tgo.GetLastMessage()
        if nMsg != oMsg {
            fmt.Println("Got: " + nMsg.Text)
            _, _ = tgo.SendMessage(nMsg.Chat.Id, nMsg.Text)
            fmt.Println("Sending: " + nMsg.Text)
        }
        oMsg = nMsg
        time.Sleep(1000 * time.Millisecond)
    }
}
```
