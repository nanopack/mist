[![mist logo](http://nano-assets.gopagoda.io/readme-headers/mist.png)](http://nanobox.io/open-source#mist)
 [![Build Status](https://travis-ci.org/nanopack/mist.svg)](https://travis-ci.org/nanopack/mist)

Mist is a simple pub/sub based on the idea that messages are tagged. To subscribe the client simply constructs a list of tags that it is interested in, and all messages that are tagged with ALL the tags are sent to the client. A client may have multiple subscriptions active at the same time.

## TCP Protocol

The protocol to talk to mist is a simple line-based TCP protocol. It was designed to be readable, debuggable and observable without specialized tools needed to decode framed packets.

### Client commands

| command format | description | server response |
| --- | --- | --- |
| `ping` | ask for a pong response, mainly to ensure that the connection is alive | `pong`
| `publish {tags} {data}` | publish a message `data` with a list of comma-delimited tags | nil |
| `subscribe {tags}` | subscribe to messages that contain ALL tags in `tags` |  nil |
| `unsubscribe {tags}` | unsubscribe to a previous subscription to `tags`; order of the tags does not matter | nil |
| `list` | list all current subscriptions active with the current client; returns a space-delimited set of subscriptions, where each tag in the subscription is delimited with a comma | `list {subscriptions}` |
| `register {tags} {token}` | register a token with a set of tags; this allows a WebSocket client to subscribe to tags | nil |
| `unregister {token}` | removes a token from mist completely | nil |
| `set {tags} {token}` | adds a set of tags to a token | nil |
| `unset {tags} {token}` | removes a set of tags from a token | nil |
| `tags {token}` | show tags that are associated with a token | `tags {token} {tags}` |

### Published message format

Message that are published to clients as the result of a subscription are delivered in this format over the wire:

`publish {tags} {data}`

### Notes

- Data flowing through mist is **NOT** touched in anyway. It is not verified in any way, but it **MUST NOT** contain a newline character as this will break the mist protocol.
- Messages are not guaranteed to be delivered, if the client is running behind on processing messages, newer messages could be dropped.
- Messages are not stored until they are delivered, if no client is available to receive the message, then it is dropped without being sent anywhere.

## WebSocket Endpoint

Mist also comes with an embeddable WebSocket api, that can be dropped into an already existing application. And by default has a layer of authentication.

To authenticate with the WebSocket endpoint, a valid token MUST be passed in as the `X-Auth-Token` header value. If the user is authenticated correctly, then the WebSocket is allowed to be established. The token is used to look up which tags can be used in a subscription, and all subscriptions MUST have at least 1 valid tag in them to be successful.

## Payloads

**note** - all frames are text frames

| Client Frame | Description | Server Frame |
| --- | --- | --- |
| `{"command": "subscribe", "tags": ["tag1", "Tag2"]}` | subscribe to events matching the tags field | `{"success": true, "command": "subscribe"}` |
| `{"command": "unsubscribe", "tags": ["tag1", "Tag2"]}` | unsubscribe to events matching the tags | `{"success": true, "command": "unsubscribe"}` |
| `{"command": "list"}` | list the subscriptions that are currently active | `{"success": true, "command": "list"}` |
| `{"command": "ping"}` | ping pong frame | `{"success": true, "command": "ping"}` |
| nil | Frame forwarded as a result of matching a subscription | `{"keys": ["tag1", "tag2"], "data": "Opaque Data encoded as a JSON string"}` |

### Notes
- publishing is not allowed over WebSockets.

## Subscription explanation

All events passing through Mist have a list of tags associated with them. For example, if Mist was forwarding IRC messages, then the format could look something like this:

```json
{
  "data": "How can I help?",
  "tags": ["user:nanobot", "type:public", "room:#nanobox"]
}
```

### Possible subscriptions

| Tags | Description |
| --- | --- |
| `["user:nanobot"]` | subscribe only to messages from the user nanobot, ignore everyone else |
| `["room:#nanobox"]` | subscribe only to messages from everyone in the #nanobox channel |
| `["room:#nanobox", "user:nanobot"]` | subscribe only to messages from nanobot in the #nanobox channel |

### Config
- Config files must contain space separated keys/values
```ini
tcp_listen_address 127.0.0.1:1445
http_listen_address 127.0.0.1:8080
log_level INFO
multicast_interface eth1
pg_user postgres
pg_database postgres
pg_address 127.0.0.1:5432
```
- Config passed as args must have a leading '-' and be separated by an '='
```
-tcp_listen_address=127.0.0.1:1445 \
-http_listen_address=127.0.0.1:8080 \
-log_level=INFO \
-multicast_interface=eth1 \
-pg_user=postgres \
-pg_database=postgres \
-pg_address=127.0.0.1:5432
```

### Running
- Run mist as follows `mist [CONFIG_FILE|-OPTION=VALUE]`

[![open source](http://nano-assets.gopagoda.io/open-src/nanobox-open-src.png)](http://nanobox.io/open-source)
