package main

import (
    "fmt"
    "time"
    "encoding/json"
    "nanomsg.org/go-mangos"
    "nanomsg.org/go-mangos/protocol/req"
    "nanomsg.org/go-mangos/transport/ipc"
    "nanomsg.org/go-mangos/transport/tcp"
)

type Reqst struct {
    Cmnd  string  `json:cmnd`
    Args  []string `json:args`
    Rslt  string  `json:"rslt,omitempty"`
}

func request(url string, command string, args... string) (res string) {
    var sock mangos.Socket
    var err error
    var msg []byte
    if sock, err = req.NewSocket(); err != nil {
        return "No Socket"
        }
    sock.AddTransport(ipc.NewTransport())
    sock.AddTransport(tcp.NewTransport())
    if err = sock.Dial(url); err != nil {
        return "Cannot Dial"
        }
    reqst := Reqst{}
    reqst.Cmnd = command
    reqst.Args = args

    b, err := json.Marshal(reqst)
    if err != nil {
        b = []byte("error in json Marshal")
        }

    if err = sock.Send([]byte(b)); err != nil {
        return "Cannot Send"
        }
    var reply string
// time out after 3 seconds - gives lost connection
    err = sock.SetOption(mangos.OptionRecvDeadline, time.Second * 3)
    msg, err = sock.Recv()
    if err == nil {
// received data
        err = json.Unmarshal(msg, &reqst)
        if err != nil {
            reply = "Unmarshal Error"
        } else {
            reply = reqst.Rslt
        }
    } else { 
        reply = "Lost Connection"
        }
    sock.Close()
    return reply
}

func main() {
    for {
        fmt.Println("Received: ",request("tcp://slither.local:5555", "sysType"))
        fmt.Println("Received: ",request("tcp://slither.local:5555", "cpuTemp"))
        fmt.Println("Received: ",request("tcp://slither.local:5555", "adds", "1", "2"))
        fmt.Println("Received: ",request("tcp://c.local:5555", "adds", "1", "2","4"))
        fmt.Println("Received: ",request("tcp://c.local:5555", "sysType"))
        time.Sleep(1 * time.Second)
    }
}
