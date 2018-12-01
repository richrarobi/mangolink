package main

import (
    "os"
    "math/rand"
    "time"
    "encoding/json"
//    "strconv"
    "os/signal"
    "syscall"
    "fmt"
    "nanomsg.org/go/mangos/v2"
    "nanomsg.org/go/mangos/v2/protocol/req"
// register transports
    _ "nanomsg.org/go/mangos/v2/transport/tcp"
    _ "nanomsg.org/go/mangos/v2/transport/ipc"
)

type Reqst struct {
    Cmnd  string  `json:cmnd`
    Args  []interface{} `json:args`
    Rslt  string  `json:"rslt,omitempty"`
}

func delay(ms int) {
    time.Sleep(time.Duration(ms) * time.Millisecond)
}

func request(url string, command string, args... interface{}) (res string) {
    var sock mangos.Socket
    var err error
    var msg []byte
    if sock, err = req.NewSocket(); err != nil {
        return "No Socket"
        }
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
    running := true
// initialise getout
    signalChannel := make(chan os.Signal, 2)
    signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
    go func() {
        sig := <-signalChannel
        switch sig {
        case os.Interrupt:
            fmt.Println("Stopping on Interrupt")
            running = false
            return
        case syscall.SIGTERM:
            fmt.Println("Stopping on Terminate")
            running = false
            return
        }
    }()

    fmt.Println("from a: ",request("tcp://a.local:5555", "ptDelta", "pan", 30))
    fmt.Println("from a: ",request("tcp://a.local:5555", "ptDelta", "tilt", 30))
    fmt.Println("from a: ",request("tcp://a.local:5555", "ptDelta", "pan", -60))
    fmt.Println("from a: ",request("tcp://a.local:5555", "ptDelta", "tilt", -60))
    fmt.Println("from a: ",request("tcp://a.local:5555", "ptDelta", "pan", 30))
    fmt.Println("from a: ",request("tcp://a.local:5555", "ptDelta", "tilt", 30))
    
    for running {
        p := rand.Intn(8)
        r := rand.Intn(255)
        g := rand.Intn(255)
        b := rand.Intn(255)
        l := rand.Intn(3)
        
//        fmt.Println("from pebble: ", request("tcp://pebble.local:5555", "sysType"))
//        go request("tcp://d.local:5555", "setPix", p, r, g, b, l )
        go request("tcp://c.local:5555", "setPix", p, r, g, b, l )
//        fmt.Println("from e: ",request("tcp://e.local:5555", "cpuTemp"))
//        go request("tcp://b.local:5555", "setPix", p, r, g, b, l )
//        fmt.Println("from pebble: ",request("tcp://pebble.local:5555", "adds", p, r, g, b, l))

//        fmt.Println("from a: ",request("tcp://a.local:5555", "ptGo", 30, 30))
//        fmt.Println("from a: ",request("tcp://a.local:5555", "ptGo", 0, 0))
        
        delay(50)
    }
    
//    fmt.Println("from f: ",request("tcp://f.local:5555", "pthome"))
//    go request("tcp://d.local:5555", "clrAllPix")
//    go request("tcp://c.local:5555", "clrAllPix")
    fmt.Println("Stopping")
    
// need to wait for completion of go routines
    time.Sleep(3 * time.Second)
}
