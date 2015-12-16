# CoAP Server
-------------

### Available functions

    - resource discovery
    - observe (register, notify)
    - register resource
    - remove resource
    - send resource event


## Creating Statically Linked Executables

#### For arm5 (EV3)

```
$ CGO_ENABLED=0 GOARCH=arm GOARM=5 go get -a -ldflags '-s' github.com/qualiapps/subject
```
