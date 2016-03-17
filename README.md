[![mist logo](http://nano-assets.gopagoda.io/readme-headers/mist.png)](http://nanobox.io/open-source#mist)
[![Build Status](https://travis-ci.org/nanopack/mist.svg)](https://travis-ci.org/nanopack/mist)

Mist is a simple pub/sub based on the idea that messages are tagged. To subscribe, a client simply constructs a list of tags that it is interested in, and all messages that are tagged with *all* of those tags are sent to that client.

A client can not only be a subscriber (with multiple active subscriptions), but also a publisher. Clients will receive messages for any tags they are subscribed, *except message publish by themselves*.

## Available Commands
Mist comes with two sets of available commands out of the box. Basic commands and Admin commands. It also has the ability to accept custom commands and handlers.

You can connect to mist with something like netcat; once connected you can simply type commands:

```
>> nc 127.0.0.1 1445
{"command":"publish", "tags":["hello"], "data":"world!"}
```

#### Basic Commands
Basic command are what provide the core functionality of mist. They allow you subscribe to and publish messages, see all of your active subscriptions and also unsubscribe from any tags you no longer want to receive messages for.

| Command | Description | Example |
| --- | --- | --- |
| `ping` | ping the server to test for an active connection | `{"command":"ping"}` |
| `subscribe` | subscribe to messages for *all* `tags` in a group | `{"command":"subscribe", "tags":["hello"]}` |
| `unsubscribe` | unsubscribe `tags` (order does not matter) | `{"command":"unsubscribe", "tags":["hello"]}` |
| `publish` | publish `data` to the list of `tags` | `{"command":"publish", "tags":["hello"], "data":"world!"}` |
| `list` | list all active subscriptions for client | `{"command":"list"}` |

#### Admin Commands
If mist is started with an `authenticator` and a `token` then a client has the chance to validate that token on connect. Once validated mist adds some additional admin commands that allow the creation of `token`/`tag` combos that provide a layer of authentication when using basic commands.

| Command | Description | Example |
| --- | --- | --- |
| `register` | register a `token` with a set of `tags` | `{"command":"register", "tags":["hello"], "data":"TOKEN"}` |
| `unregister` | removes a `token` from mist completely | `{"command":"unregister", "data":"TOKEN"}` |
| `set` | adds a set of `tags` to a `token` | `{"command":"set", "tags":["hello"], "data":"TOKEN"}` |
| `unset` | removes a set of `tags` from a `token` | `{"command":"unset", "tags":["hello"], "data":"TOKEN"}` |
| `tags` | show `tags` that are associated with a `token` | `{"command":"tags", "data":"TOKEN"}` |

## Messages

All communications within mist are sent and received as JSON encoded/decoded messages:
```go
Message struct {
  Command string   `json:"command"`
  Tags    []string `json:"tags"`
  Data    string   `json:"data,omitemtpy"`
  Error   string   `json:"error,omitempty"`
}
```

Each Message has a set of `tags` and `data`. Tags can take any form you like, as they are just an array of strings.

``` json
{
  "tags": ["company:pagodabox", "product:mist", "repo:#nanopack"],
  "data": "Mist is awesome!"
}
```

### Subscribing / Publishing

Think of `tags` as a way to filter out messages you don't want to receive; the more tags that are added to a subscription the more direct a message has to be:

| Subscribed tags | Messages received from tags |
| --- | --- |
| `["onefish"]` | `["onefish"]`, `["onefish","twofish"]`, `["onefish","twofish","redfish"]` |
| `["onefish", "twofish"]` | `["onefish","twofish"]`, `["onefish","twofish","redfish"]` |
| `["onefish", "twofish", "redfish"]` | `["onefish","twofish","redfish"]` |

Message that are published to clients as the result of a subscription are delivered in this format:

`{"command":"<command>", "tags":["<tag>", "<tag>"], "data":"<data>"}`

A few things to not about how mist handles data:

* Data flowing through mist is *not touched or verified in anyway*, however, it **MUST NOT** contain a newline character as this will break the mist protocol.

* Messages are not guaranteed to be delivered, if the client is running behind on processing messages, newer messages could be dropped.

* Messages are not stored, if no client is available to receive the message, then it is dropped.

## Listeners

Out of the box mist supports three different types of servers (`TCP`, `HTTP`, and `Websocket`). By default, when mist starts, it will start one of each.

```
TCP server listening at '127.0.0.1:1445'...
HTTP server listening at '127.0.0.1:8080'...
WS server listening at '127.0.0.1:8888'...
```

When starting mist, you can specify any number and type of server you'd like as long as it follows the string URI protocol (If a listener is passed that mist doesn't support it will skip).

Also, if mist doesn't support a server you need it allows you to register custom servers that can be used on startup.

#### Available listeners:

`(scheme:[//[user:pass@]host[:port]][/]path[?query][#fragment])`

| Listener | URI scheme |
| --- | --- |
| tcp | `tcp://127.0.0.1:1445` |
| http | `http://127.0.0.1:8080` |
| websocket | `ws://127.0.0.1:8888` |

##### Example
```
./mist --server --listeners "tcp://127.0.0.1:1445", "http://127.0.0.1:8080", "ws://127.0.0.1:8888"
```

## Authenticators

Mist also provides support for authentication. This means that during startup you can provide mist with a `token` that you want to be used as authentication. Once enabled, any client that attempts to connect to mist *must* provide that token or be disconnected. By default mist does not use authentication.

Like listeners, mist allows for the registration of custom authenticators.

#### Available authenticators

`(scheme:[//[user:pass@]host[:port]][/]path[?query][#fragment])`

| Authenticator | URI scheme | description
| --- | --- | --- |
| memory | `memory://` | an in memory store |
| [scribble](https://github.com/nanobox-io/golang-scribble) | `scribble://?db=/tmp` | a tiny JSON database |
| postgres | `postgres://postgres@127.0.0.1:5432?db=postgres` | n/a |

##### Example

```
./mist --server --authenticator "memory://"
```

## Websockets

Since mist just uses a JSON message protocol internally, sending messages via websocket is easy.

NOTE: If authentication is enabled you'll need to provide a token when connecting the websocket:

* As a Header: `X-Auth-Token: token`
* As a query param: `x-auth-token=token`

##### Example

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

## Running mist

mist can be run as either a client or a server.

#### As a server:
To run mist as a server, using the following command will start mist as a daemon:

`mist --server`

If you need to override any default config options you can pass the path to a config file:

`mist --server --config /path/to/config`

##### example config

```yml
authenticator: memory://
listeners:
  - tcp://127.0.0.1:1445
log-level: INFO
token: TOKEN
```

Or you can just pass any configuration options as flags:

`mist --server --log-level DEBUG`

#### As a client (CLI):
You can also use mist as a client to any running mist.

```
Usage:
   [flags]
   [command]

Available Commands:
  ping        Ping a running mist server
  subscribe   Subscribe tags
  unsubscribe Unsubscribe tags
  publish     Publish a message
  list        List all subscriptions

Flags:
      --authenticator="": Setting this option enables authentication and uses the authenticator provided to store tokens
      --config="": /path/to/config.yml
  -h, --help[=false]: help for
      --listeners=[tcp://127.0.0.1:1445,http://127.0.0.1:8080,ws://127.0.0.1:8888]: A comma delimited list of servers to start
      --log-file="/var/log/mist.log": If log-type=file, the /path/to/logfile; ignored otherwise
      --log-level="INFO": Output level of logs (TRACE, DEBUG, INFO, WARN, ERROR, FATAL)
      --log-type="stdout": The type of logging (stdout, file)
      --replicator="": not yet implemented
      --server[=false]: Run mist as a server
      --token="": Auth token used when connecting to a Mist started with an authenticator
  -v, --version[=false]: Display the current version of this CLI

Use " [command] --help" for more information about a command.
```

## Contributing

Contributions to mist are welcome and encouraged. Mist is a [Nanobox](https://nanobox.io) project and contributions should follow the [Nanobox Contribution Process & Guidelines](https://docs.nanobox.io/contributing/).

[![open source](http://nano-assets.gopagoda.io/open-src/nanobox-open-src.png)](http://nanobox.io/open-source)
