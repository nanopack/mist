[![mist logo](http://nano-assets.gopagoda.io/readme-headers/mist.png)](http://nanobox.io/open-source#mist)

[![Build Status](https://travis-ci.org/nanopack/mist.svg)](https://travis-ci.org/nanopack/mist)

`mist` is a simple pub/sub based on the idea that messages are tagged. To subscribe, the client simply constructs a list of tags that it is interested in, and all messages that are tagged with **all** of those tags are sent to the client.

A client may have multiple subscriptions active at the same time.

## TCP

The protocol to talk to mist is a simple line-based TCP protocol. It was designed to be readable, debuggable, and observable without specialized tools needed to decode framed packets.

You can connect to a running `mist` with something like netcat:

```
nc 127.0.0.1 1445
```

Once connected you can simply type commands and the server will respond

#### Commands

| command format | description | server response |
| --- | --- | --- |
| `ping` | ping the server to test for an active connection | `pong`
| `publish tag,tag data` | publish `data` to the list of `tags` | nil |
| `subscribe tag,tag` | subscribe to messages for ***all*** `tags` in group | nil |
| `unsubscribe tag,tag` | unsubscribe `tags` (order of the `tags` does not matter) | nil |
| `list` | list all current active subscriptions for client | `list [[tag,tag] [tag]]` |
| `register {tags} {token}` | register a `token` with a set of `tags`; this allows a WebSocket client to subscribe to `tags` | nil |
| `unregister {token}` | removes a `token` from mist completely | nil |
| `set {tags} {token}` | adds a set of `tags` to a `token` | nil |
| `unset {tags} {token}` | removes a set of `tags` from a `token` | nil |
| `tags {token}` | show `tags` that are associated with a `token` | `tags {token} {tags}` |

### Published message format

Message that are published to clients as the result of a subscription are delivered in this format over the wire:

`publish tag,tag data`


* Data flowing through `mist` is **not touched or varified in anyway**, but it **MUST NOT** contain a newline character as this will break the mist protocol.

* Messages are not guaranteed to be delivered, if the client is running behind on processing messages, newer messages could be dropped.

* Messages are not stored until they are delivered, if no client is available to receive the message, then it is dropped without being sent anywhere.

## Websockets

`mist` also comes with an embeddable WebSocket api, that can be dropped into an already existing application. And by default has a layer of authentication. `mist` only accepts `text frames` as the form of communication across the socket.

To authenticate with the WebSocket endpoint, a valid token **must** be passed in one of the following methods:

* As a Header: `X-Auth-Token: token`
* As a query param: `x-auth-token=token`

If the user is authenticated correctly, then the WebSocket is allowed to connect. The `token` is used to look up which `tags` can be used in a subscription, and all subscriptions **must** have at least 1 valid `tag` in them to be successful.

#### Payloads

| Command | Description | Response
| --- | --- | --- |
| `{"command": "ping"}` | ping server | `{"success": true, "command": "ping"}` |
| `{"command": "subscribe", "tags": ["tag"]}` | subscribe to events matching `tags` | `{"success": true, "command": "subscribe"}` |
| `{"command": "unsubscribe", "tags": ["tag"]}` | unsubscribe from events matching `tags` | `{"success": true, "command": "unsubscribe"}` |
| `{"command": "list"}` | list active subscriptions | `{"success": true, "command": "list"}` |
| nil | Frame forwarded as a result of matching a subscription | `{"keys": ["tag"], "data": "Opaque Data encoded as a JSON string"}` |

#### Examples

``` javascript

  // connect the websocket
  var ws = new WebSocket("ws://localhost:8080/subscribe/websocket?x-auth-token=token")

  // handle responses from the server
  ws.onmessage = function(me){
    console.log("Response!", me.data)
  }

  // ping
  ws.send(JSON.stringify({"command": "ping"}))

  // subscribe
  ws.send(JSON.stringify({"command": "subscribe", "tags": ["hello", "world"]}))

  // unsubscribe
  ws.send(JSON.stringify({"command": "unsubscribe", "tags": ["hello", "world"]}))

  // list
  ws.send(JSON.stringify({"command": "list"}))
```

## Subscriptions

All events passing through `mist` have a list of `tags` associated with them. `Tags` can take any form you like, they are just an array of strings.

For example, if `mist` was forwarding IRC messages, then the format might look something like this:

``` json
{
  "data": "Mist is awesome!",
  "tags": ["user:nanobot", "type:admin", "room:#nanobox"]
}
```

#### Examples

| Tags | Description |
| --- | --- |
| `["user:nanobot"]` | subscribe only to messages from the user nanobot, ignore everyone else |
| `["room:#nanobox"]` | subscribe only to messages from everyone in the #nanobox channel |
| `["room:#nanobox", "user:nanobot"]` | subscribe only to messages from nanobot in the #nanobox channel |

## Config Options

`mist` will accept a config file on startup that can override any of the following defaults:

``` ini
tcp_listen_address 127.0.0.1:1445
http_listen_address 127.0.0.1:8080
log_level INFO
multicast_interface eth1
pg_user postgres
pg_database postgres
pg_address 127.0.0.1:5432
```

## Running `mist`

`mist` can be run as either a client or a server.

#### As a server:
To run `mist` as a server, using the following command will start `mist` as a daemon:

`mist -d`

If you need to override any default config options you can pass the path to a config file:

`mist -d --config path/to/config`

#### As a client:
You can also use `mist` as a client to connect to another running `mist`.

```
Usage:
   [flags]
   [command]

Available Commands:
  list        List all subscriptions
  ping        Ping a running mist server
  publish     Publish a message
  subscribe   Subscribe tags
  unsubscribe Unsubscribe tags

Flags:
  -c, --config="": Path to config options
  -d, --daemon[=false]: Run mist as a server
  -h, --help[=false]: help for
      --log-level="INFO": desc.
      --tcp-addr="127.0.0.1:1445": desc.
  -v, --version[=false]: Display the current version of this CLI

Use " [command] --help" for more information about a command.
```

[![open source](http://nano-assets.gopagoda.io/open-src/nanobox-open-src.png)](http://nanobox.io/open-source)
