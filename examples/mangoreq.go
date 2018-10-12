package main

import (
    "os"
    "math/rand"
    "time"
    "encoding/json"
    "strconv"
    "os/signal"
    "syscall"
    "fmt"
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

func delay(ms int) {
    time.Sleep(time.Duration(ms) * time.Millisecond)
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
    
    for running {
        p := strconv.Itoa(rand.Intn(8))
        r := strconv.Itoa(rand.Intn(255))
        g := strconv.Itoa(rand.Intn(255))
        b := strconv.Itoa(rand.Intn(255))
        l := strconv.Itoa(rand.Intn(3))
        request("tcp://c.local:5555", "setPix", p, r, g, b, l )
        fmt.Println("Received from c: ",request("tcp://c.local:5555", "cpuTemp"))
        request("tcp://b.local:5555", "setPix", p, r, g, b, l )
        fmt.Println("Received from b: ",request("tcp://b.local:5555", "cpuTemp"))
        delay(200)
//        time.Sleep(1 * time.Second)
    }
    request("tcp://c.local:5555", "clrAllPix")
    request("tcp://b.local:5555", "clrAllPix")
    fmt.Println("Stopping")
}
