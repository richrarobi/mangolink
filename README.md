# mangolink
request / reply with nanomsg / mangos



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
