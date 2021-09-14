package ws

import (
    "github.com/labstack/echo/v4"
    . "github.com/jinuopti/janus/log"
)

var world *World

func GetWorld() *World {
    return world
}

func InitWebSocket() {
    if world == nil {
        world = NewWorld()
        go world.Run()
    }
    Logd("Initialize WebSocket! uri: /ws")
}

func WebSocketHandler(c echo.Context, readCallback func(*Client, []byte), disconnCallback func(*Client)) error {
    conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
    if err != nil {
        return err
    }
    defer func() { _ = conn.Close() }()

    ch := make(chan bool)
    client := NewClient(world, conn, ch, readCallback, disconnCallback)
    client.World.ChanEnter <- client

    <- ch

    Logd("Exit WebSocket Handler")

    return nil
}
