# WIP - The IPFS Textboard protocol

This is a textual protocol based on the [JSON](https://www.json.org/) standard.
Posts are objects comprised of metadata, with a field containing the CID of the
actual text content published by a user. Optionally, individual textboards may
choose to validate posts out in the wild, using a cryptographic signature (e.g. RSA)
of the post content's CID.

## Post format

The post is a JSON object, with the following fields:

| Field | Type | Optional | Description |
| :--- | :---: | :---: | :--- |
| `Protocol` | `string` | `false` | Constant value `IPFS-TXT`. |
| `Version` | `string` | `false` | The protocol version. |
| `Topic` | `string` | `true` | Analogous to the board, in other software; formmated in UTF-8; should not exceed 64 characters. Only optional if not OP, otherwise ignored. |
| `Title` | `string` | `true` | The title of the OP, formatted in UTF-8; should not exceed 256 characters. Only optional if not OP, otherwise ignored. |
| `Thread` | `string` | `true` | The thread being replied to. Only optional if OP, otherwise ignored. |
| `Content` | `string` | `false` | The content ID hash pointing to the post's text formatted in UTF-8; should not exceed 1024 characters. |
| `Posted` | `string` | `false` | The post date, as an RFC3339Nano time string. |
| `Auth` | `object` | `true` | Authentication parameters. Explained below. |

The `Auth` field can be used to identify posts that have been validated by
certain textboard authorities, following a semi-decentralized architecture:

| Field | Type | Description |
| :--- | :---: | :--- |
| `Node` | `string` | Node identifier for the textboard in question. |
| `Signature` | `string` | The resulting signature. |
| `Format` | `string` | The format the signature is serialized in. |

The retrieval of the public key used to verify the signature should be handed to a higher
level protocol. The range of values accepted by `Format` is:

* `hex` - hexadecimal encoding
* `base64` - base64 encoding
* `base58` - base58 encoding
* `bin` - raw binary encoding
