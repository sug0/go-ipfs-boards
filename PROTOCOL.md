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
| `Topic` | `string` | `true` | Analogous to the board, in other software; formmated in UTF-8. Only optional if not OP, otherwise ignored. |
| `Title` | `string` | `true` | The title of the OP, formatted in UTF-8. Only optional if not OP, otherwise ignored. |
| `Thread` | `string` | `true` | The thread being replied to. Only optional if OP, otherwise ignored. |
| `Content` | `string` | `false` | The content ID hash pointing to the post's text formatted in UTF-8. |
| `Posted` | `string` | `false` | The post date, as an RFC3339Nano time string. |
| `Auth` | `object` | `true` | Authentication parameters. Explained below. |

The `Auth` field can be used to identify posts that have been validated by
certain textboard authorities, following a semi-decentralized architecture:

| Field | Type | Description |
| :--- | :---: | :--- |
| `Signature` | `string` | The resulting signature, formatted in hexadecimal. |
| `PubKey` | `object` | The public key object. |

The `PubKey` field is described below:

| Field | Type | Description |
| :--- | :---: | :--- |
| `CID` | `string` | The content ID hash pointing to the public key. |
| `Scheme` | `string` | The cryptographic scheme used (e.g. RSA). |
| `Format` | `string` | The format the public key is serialized in. |

__TODO__ - specify the range of values accepted by `Format` and `Cipher`.
