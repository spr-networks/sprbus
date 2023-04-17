sprbus client

connect to a remote event bus using token or local unix socket
to listen for events.

run either on spr or using api
if remote will need to have notifications setup to \*

**remote**

```sh
./sprbus --addr 192.168.2.1
```

**local**

```sh
./sprbus
```

# TODO

- for now publish only works for local connections
- temp enable \* notifications:
  - send request to /notifications with `prefix:"", Notifications: False`
  - disable when we're done to not send excessive data over ws.
