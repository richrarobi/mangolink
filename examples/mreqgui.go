package main

import (
    "time"
    "encoding/json"
//    "strconv"
    "github.com/andlabs/ui"
    _ "github.com/andlabs/ui/winmanifest"
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

var mainwin *ui.Window

func makeDataChoosersPage() ui.Control {
    hbox := ui.NewHorizontalBox()
    hbox.SetPadded(true)

    vbox := ui.NewVerticalBox()
    vbox.SetPadded(true)
    hbox.Append(vbox, true)

    grid := ui.NewGrid()
    grid.SetPadded(true)
    vbox.Append(grid, false)

    button := ui.NewButton("Up")
    button.OnClicked(func(*ui.Button) {
        request("tcp://a.local:5555", "ptDelta", "tilt", -5)
    })
    grid.Append(button, 1, 0, 1, 1, false, ui.AlignFill, false, ui.AlignFill)

    button = ui.NewButton("Left")
    button.OnClicked(func(*ui.Button) {
        request("tcp://a.local:5555", "ptDelta", "pan", 5)
    })
    grid.Append(button, 0, 1, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
        
    button = ui.NewButton("Centre")
    button.OnClicked(func(*ui.Button) {
        request("tcp://a.local:5555", "ptGo", 0, 0)
    })
    grid.Append(button, 1, 1, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
      
        
    button = ui.NewButton("Right")
    button.OnClicked(func(*ui.Button) {
        request("tcp://a.local:5555", "ptDelta", "pan", -5)
    })
    grid.Append(button, 2, 1, 1, 1, false, ui.AlignFill, false, ui.AlignFill)

    button = ui.NewButton("Down")
    button.OnClicked(func(*ui.Button) {
        request("tcp://a.local:5555", "ptDelta", "tilt", 5)
    })
    grid.Append(button, 1, 2, 1, 1, false, ui.AlignFill, false, ui.AlignFill)

    button = ui.NewButton("Exit")
    button.OnClicked(func(*ui.Button) { ui.Quit() })
    grid.Append(button, 1,4, 1, 1, false, ui.AlignFill, false, ui.AlignFill)

    return hbox
}

func setupUI() {
    mainwin = ui.NewWindow("PanTilt Control", 240, 200, true)
    mainwin.OnClosing(func(*ui.Window) bool {
        ui.Quit()
        return true
    })

    ui.OnShouldQuit(func() bool {
        mainwin.Destroy()
        return true
    })

    tab := ui.NewTab()
    mainwin.SetChild(tab)
    mainwin.SetMargined(true)

    tab.Append("Buttons", makeDataChoosersPage())
    tab.SetMargined(0, true)

    mainwin.Show()
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
// time out after 10 seconds - gives lost connection
    err = sock.SetOption(mangos.OptionRecvDeadline, time.Second * 10)
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
    ui.Main(setupUI)
}
