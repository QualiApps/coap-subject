#!/usr/bin/python

from coap.resource import CoapResource
from time import sleep


res = CoapResource("http://localhost:8888")

res.add_resource("r1", ("test",))
res.add_resource("r2", ("sensors", "temp"))

data = {"type": "python", "message": ""}
resources = res.get_resources().keys()

try:
    while True:
        for name in resources:
            data["message"] = name

            res.send_event(name, data)

        sleep(5)
except KeyboardInterrupt:
    for name in resources:
        res.remove_resource(name)
