# WIP - The IPFS Textboard protocol

This is a textual protocol based on the [JSON](https://www.json.org/) standard.
Posts are objects comprised of metadata, with a field containing the CID of the
actual text content published by a user. Optionally, individual textboards may
choose to validate posts out in the wild, using a cryptographic signature (e.g. RSA)
of the post content's CID.

## Post format for version `0.1.2`

The post is a JSON object, with the following fields:

| Field | Type | Optional | Description |
| :--- | :---: | :---: | :--- |
| `Protocol` | `string` | `false` | Constant value `IPFS-TXT`. |
| `Version` | `string` | `false` | The protocol version. |
| `Topic` | `string` | `true` | Analogous to the board, in other software; formmated in UTF-8. Required if OP. |
| `Title` | `string` | `true` | The title of the OP, formatted in UTF-8. Required if OP. |
| `Thread` | `string` | `true` | The thread being replied to. Required if OP. |
| `Content` | `string` | `false` | The content ID hash pointing to the post's text formatted in UTF-8. |
| `Posted` | `string` | `false` | The post date, as an RFC3339Nano time string. |
| `Extensions` | `object` | `true` | Application defined extensions. |

The character limits of the `Topic`, `Title` and `Content` should be left to the
text board authority to decide. The `Extensions` field can be used to implement
domain specific extensions, such as ~~tripfagging~~.

### Example OP post

TODO

### Example reply post

TODO
