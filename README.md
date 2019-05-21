# mangolink
request / reply with nanomsg / mangos

NEW EXAMPLE:- In examples, mreply_a.go and mreqgui.go work together.
On Raspi system a.local I have mreply_a.go running (from startup).
So... I run a web browser on another system, call up the picamera page on a.local. Then using the Gui, mreqgui.go on the
second system (or anywhere else!), I can move the pantilt hat down, up, left, right and to centre. See the picture...
p.s I have added a new version of the gui with 2 tabs (controlling 2 different pi pantilt heads)


Two examples that work together:- mreply_a.go and mreq.go. mreq line 90 ... calls mreply on system a, 
this controls the Pimoroni pantilt hat on the remote system......
Also in the example, the Pimoroni Blinkt led's are being driven remotely..
see the line:-  go request("tcp://c.local:5555", "setPix", p, r, g, b, l )
and don't forget to :- go request("tcp://c.local:5555", "clrAllPix")
on exit!!!



Next I plan to write a Go based gui, so I can control the pantilt heads with buttons.



mangorep has functions to be used by mangoreq.
I run mangorep on startup on Raspberry Pi. It will also work on a Linux PC.
Currently mangorep has some example functions. 

Added in the examples directory.. mangorep_pi.go and mangoreq.go

I have 2 Raspberry Pi's, each with a Pimoroni Blinkt (8 tricolour LED's). In this test, I can run the same mangorep_pi.go
on each of the pi's, and using the mangreq.go on another system (a linux Mint pc in this case) both the blinkt's 
will show the same random LED's at (virtually) the same time. 

see func clrAllPix and setPix in mangorep, and the calling code in func main in mangoreq....note the calls in the case statement.

Python example added in examples to call mangorep .

RichR

p.s. the request (mreq.go) timeout value of 3 seconds is a bit too short, especially when mechanical doo-hickies like servos are
involved. I have increased it to 10 in my system.

p.p.s I was probably remiss on not giving install info, this should help a bit ... :-
Linux first

   20  sudo apt install git

   26  go get nanomsg.org/go/mangos/v2
   27  go get github.com/richrarobi/mangolink

   31  go get github.com/andlabs/ui
   32  sudo apt install libgtk-3-dev



