package main

import (
        "fmt"
        "time"
        "encoding/json"
        "strings"
        "os/exec"
        "context"
        "bufio"
        "strconv"
        "nanomsg.org/go/mangos/v2"
        "nanomsg.org/go/mangos/v2/protocol/rep"
        // register transports
        _ "nanomsg.org/go/mangos/v2/transport/tcp"

        "github.com/richrarobi/periBlink"
)

type Reqst struct {
    Cmnd  string  `json:cmnd`
    Args  []interface{} `json:args`
    Rslt  string  `json:"rslt,omitempty"`
}

const url = "tcp://*:5555"

func main() {
        var sock mangos.Socket
        var err error
        var msg []byte
        if sock, err = rep.NewSocket(); err != nil {
                fmt.Println( err )
        }

        if err = sock.Listen(url); err != nil {
                fmt.Println( err )
        }
        
        periBlink.Setup()
        periBlink.SetLuminance(1)
        periBlink.Clear()
        periBlink.Show()
        
        for {
                reqst := Reqst{}
                msg, err = sock.Recv()
                err := json.Unmarshal(msg, &reqst)
                if err != nil {
                        fmt.Println("error:", err)
                }
                
        // call the function and loadup the result
        if reqst.Cmnd != "" {
            switch reqst.Cmnd {
            case "cpuTemp":
                reqst.Rslt = cpuTemp()
            case "sysType":
                reqst.Rslt = sysType()
            case "clrAllPix":
                reqst.Rslt = clrAllPix()
            case "setPix":
                args := reqst.Args
                reqst.Rslt = setPix(args ...)
            case "adds":
                args := reqst.Args
                reqst.Rslt = adds(args ...)
            default:
                reqst.Rslt = none()
            }
        }

//  Send reply back to client
        b, err := json.Marshal(&reqst)
        if err != nil {
            b = []byte("error in json Marshal")
            }
                err = sock.Send([]byte(b))
                if err != nil {
                        fmt.Println( err )
                }
        }
}

func setPix(args... interface{}) ( res string) {
    p := int(args[0].(float64))
    r := int(args[1].(float64))
    g := int(args[2].(float64))
    b := int(args[3].(float64))
    l := int(args[4].(float64))
    periBlink.SetPixel( p, r, g, b, l )
    periBlink.Show()
    return "done"
}

func clrAllPix() ( res string) {
    periBlink.Clear()
    periBlink.Show()
    return "done"
}


func adds ( args... interface{} ) ( res string) {
    var z int
    for i := 0 ; i < len(args); i++ {
            x := int(args[i].(float64))
            z = x + z
        }
    res = strconv.Itoa(z)
    return res
}

func cpuTemp()(res string) {
    var x string
    if strings.Contains( sysType(), "ARM" ) {
        x = strings.TrimSpace(exeCmd( "/opt/vc/bin/vcgencmd", "measure_temp"))
    } else {
        x = "cpuTemp: SysType Not ARM"
    }
    return x
}

func sysType()(res string) {
    var x string
    scanner := bufio.NewScanner(strings.NewReader(exeCmd( "cat", "/proc/cpuinfo")))
    for scanner.Scan() {
        if strings.Contains( scanner.Text(), "model name") {
            x = scanner.Text() [ 13:len( scanner.Text() ) ]
            return x
            break
        }
    }
     return "SysType Not Found"
}

func exeCmd(command string, args... string) (res string) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctx, command, args... )
    out, err := cmd.Output()

    if ctx.Err() == context.DeadlineExceeded {
        fmt.Println("exeCmd: Command timed out")
        return
    }
    if err != nil {
        fmt.Println("Non-zero exit code:", err)
        return "exeCmd: Error in external command"
        }
    return string(out)
}

func none () (res string) {
    return ("Function Not Found")
}
