# Message spec

To serialize messages, I chose [protocol buffers](https://protobuf.dev). See
[messages.proto](./messages.proto) for more info.

Language-native code needs to get generated from the `.proto` spec file. To do
so, run the `./gen_proto.sh` shell script.

## Connection flow

There are two endpoints on the backend server:

- `/query` - **query** status data of machines by ID
- `/provide` - **provide** status data for a machine

### `/query`

1. Establish WS connection
2. Client sends an ID to query over the WS connection
3. 
    1.  Server responds and keeps responding with status data for the requested
        machine in intervals specified by the configuration
    2.  Server responds with an error because the requested ID doesn't exist,
        connection is terminated
4. Client goes offline, connection is terminated

### `/provide`

1. Establish WS connection
2. Server sends a **randomly generated** ID over the WS connection
3. Client starts sending status data in intervals
    1. If the client sends data too fast, the server may terminate the
       connection
    2. The server periodically sends keep-alive pings to let the client know
       that the connection didn't drop
4. Client goes offline, connection is terminated
