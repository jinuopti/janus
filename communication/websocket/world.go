package ws

import (
    . "github.com/jinuopti/janus/log"
    //"time"
)

type World struct {
    ClientMap map[*Client]bool
    ChanEnter chan *Client
    ChanLeave chan *Client
    Broadcast chan []byte
}

func NewWorld() *World {
    return &World{
        ClientMap: make(map[*Client]bool, 5),
        Broadcast: make(chan []byte),
    }
}

func (w *World) Run() {
    w.ChanEnter = make(chan *Client)
    w.ChanLeave = make(chan *Client)

    //ticker := time.NewTicker(10 * time.Second)

    for {
        select {
        case client := <- w.ChanEnter:
            w.ClientMap[client] = true
            Logd("New WebSocket Client, Len: %d", len(w.ClientMap))
        case client := <- w.ChanLeave:
            if _, ok := w.ClientMap[client]; ok {
                if client.DisconnectCallback != nil {
                    client.DisconnectCallback(client)
                }
                delete(w.ClientMap, client)
                close(client.Send)
            }
            Logd("Exit WebSocket Client, Len: %d", len(w.ClientMap))
        case message := <- w.Broadcast:
            for client := range w.ClientMap {
                client.Send <- message
            }
            // case tick := <- ticker.C:
            //     for client := range w.clientMap {
            //         client.send <- []byte(tick.String())
            //     }
        }
    }
}
