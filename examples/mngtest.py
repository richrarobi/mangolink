# mngtest.py
# simple test of python3 calling golang mangorep
# Now it is possible to write simple python scripts to use
# go programs on any system

import json, time

from nanomsg import Socket, REQ, REP, PUB, SUB, DONTWAIT, \
NanoMsgAPIError, EAGAIN, SUB_SUBSCRIBE

# first connect to the (remote) pi (b.local) running mangorep

s2 = Socket(REQ)
s2.connect('tcp://b.local:5555')

# set up the query in json
#s = json.dumps({"Cmnd":"sysType"})
s = json.dumps({"Cmnd":"cpuTemp"})

# run the query

while True:
    s2.send(s)
    msg = s2.recv()
#    print(msg)
    print(json.loads(msg.decode("utf-8")))
    
# note that the whole json response is printed here
# actually it showed up that I am returning a \n newline I don't want...
# Now corrected in mangorep (line 88 ish)

    time.sleep(2)

"""
Example output:-

{'Cmnd': 'cpuTemp', 'Args': None, 'rslt': "temp=42.9'C\n"}
{'Cmnd': 'cpuTemp', 'Args': None, 'rslt': "temp=41.9'C\n"}
{'Cmnd': 'cpuTemp', 'Args': None, 'rslt': "temp=41.9'C\n"}
{'Cmnd': 'cpuTemp', 'Args': None, 'rslt': "temp=41.9'C\n"}

"""
