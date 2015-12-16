##How it works
-------------------------------

```
from coap.resource import CoapResource

# init
coap = CoapResource("http://localhost:8888")

# add resource
coap.add_resource("event", ("sensors", "event"))

# send event
coap.send_event("event", {"alert": "Message", ...})

# delete resource
coap.remove_resource("event")

```
