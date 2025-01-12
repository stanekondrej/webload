# webload - a service to display system load through a web dashboard

My idea for this project is simple - you have a **client** on a machine that uses a
websocket connection to the **server** to send system load data in a specified
interval. The server then uses this information to display the system load
through a neat interface.

This interface can be accessed using a link returned to the user by the
client program after a negotiation with the server.

## WebSocket messages

The format spec for the messages sent over the WebSocket connections can be
found in [MESSAGES.md](./MESSAGES.md)
