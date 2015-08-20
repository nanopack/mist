# Mist

Mist is a simple pub/sub based on the idea that messages are tagged. To subscribe the client simply constructs a list of tags that it is interested in, and all messages that are tagged with ALL the tags are sent to the client. A client may have multiple subscriptions active at the same time.

## Protocol

The protocol to talk to mist is a simple line based tcp protocol. It was designed to be readable, debuggable and observable without specialized tools needed to decode framed packets.

### Client commands:

| command format | description | server response |
| --- | --- | --- |
| `ping` | ask for a pong response, mainly to ensure that the conenction is alive | `pong`
| `publish {tags} {data}` | publish a message `data` with a list of comma delimited tags | nil |
| `subscribe {tags}` | subscribe to messages that contain ALL tags in `tags` |  nil |
| `unsubscribe {tags}` | unsubscribe to a previous subscription to `tags`, order of the tags does not matter | nil |
| `list` | list all current subscriptions active with the current client, returns a space delimited set of subscriptions, where each tag in the subscription is delimited with a comma | `list {subscriptions}` |

### Published message format

Message that are published to clients as the result of a subscription are delivered in this format over the wire:

`publish {tags} {data}`

### Notes:

- Data flowing through mist is **NOT** touched in anyway. It is not verified in any way, but it **MUST NOT** contain a newline character as this will break the mist protocol.
- Messages are not guaranteed to be delivered, if the client is running behind on processing messages, newer messages could be dropped.
- Messages are not stored until they are delivered, if no client is available to receive the message, then it is dropped without being sent anywhere.