// mreply.go
// Oct 2018
// Rich Robinson

package main

import (
    "bufio"
    "context"
    "encoding/json"
    "fmt"
    "os/exec"
    "strconv"
    "strings"
    "time"
    pt "github.com/richrarobi/pantilt"

    "nanomsg.org/go/mangos/v2"
    "nanomsg.org/go/mangos/v2/protocol/rep"

    // register transports
    _ "nanomsg.org/go/mangos/v2/transport/tcp"
)

type Reqst struct {
    Cmnd string        `json:cmnd`
    Args []interface{} `json:args`
    Rslt string        `json:"rslt,omitempty"`
}

const url = "tcp://*:5555"

func main() {
    var sock mangos.Socket
    var err error
    var msg []byte
    if sock, err = rep.NewSocket(); err != nil {
            fmt.Println(err)
    }
    if err = sock.Listen(url); err != nil {
            fmt.Println(err)
    }

    pt.Open()

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
                    case "adds":
                            args := reqst.Args
                            reqst.Rslt = adds(args...)
                    case "ptDelta":
                            args := reqst.Args
                            reqst.Rslt = ptDelta(args...)
                    case "ptGo":
                            args := reqst.Args
                            reqst.Rslt = ptGo(args...)
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
                fmt.Println(err)
        }
    }
}

func ptDelta(args ...interface{}) (res string){
    pt.ServoEnable(args[0].( string ), true)
    pt.Delta(args[0].( string ), int(args[1].(float64)))
    pt.ServoEnable(args[0].( string ), false)
    return "done"
}

func ptGo(args ...interface{}) (res string){
    pt.ServoEnable("pan", true)
    pt.ServoEnable("tilt", true)
    pt.Go(int(args[0].(float64)), int(args[1].(float64)))
    pt.ServoEnable("pan", false)
    pt.ServoEnable("tilt", false)
    return "done"
}

func adds(args ...interface{}) (res string) {
    var z int
    /*
   // for testing
       for i := 0 ; i < len(args); i++ {
           switch args[i].(type) {
               case string:
                   fmt.Println("String ",  args[i].( string ) )
               case float64:
                   fmt.Println( "Float", args[i].(float64))
               case bool:
                   fmt.Println( "Bool ", args[i].( bool ))
               case nil:
                   fmt.Println( " Nil ")
                default:
                    fmt.Println( "Unknown Type")
           }
       }
    */
    for i := 0; i < len(args); i++ {
        x := int(args[i].(float64))
        z = x + z
    }
    res = strconv.Itoa(z)
    return res
}

func cpuTemp() (res string) {
    var x string
    if strings.Contains(sysType(), "ARM") {
            x = strings.TrimSpace(exeCmd("/opt/vc/bin/vcgencmd", "measure_temp"))
    } else {
        x = "cpuTemp: SysType Not ARM"
    }
    return x
}

func sysType() (res string) {
    var x string
    scanner := bufio.NewScanner(strings.NewReader(exeCmd("cat", "/proc/cpuinfo")))
    for scanner.Scan() {
        if strings.Contains(scanner.Text(), "model name") {
            x = scanner.Text()[13:len(scanner.Text())]
            return x
            break
        }
    }
    return "SysType Not Found"
}

func exeCmd(command string, args ...string) (res string) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctx, command, args...)
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

func none() (res string) {
    return ("Function Not Found")
}
