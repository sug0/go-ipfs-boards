# WIP - The IPFS Textboard protocol

This is a textual protocol based on the [JSON](https://www.json.org/) standard.
Posts are objects comprised of metadata, with a field containing the CID of the
actual text content published by a user.

## Post format for version `0.1.2`

The post is a JSON object, with the following fields:

| Field | Type | Optional | Description |
| :--- | :---: | :---: | :--- |
| `Protocol` | `string` | `false` | Constant value `IPFS-TXT`. |
| `Version` | `string` | `false` | The protocol version. |
| `Topic` | `string` | `true` | Analogous to the board, in other software; formmated in UTF-8. Required if OP. |
| `Title` | `string` | `true` | The title of the OP, formatted in UTF-8. Required if OP. |
| `Thread` | `string` | `true` | The thread being replied to. Required if not OP. |
| `Content` | `string` | `false` | The content ID hash pointing to the post's text formatted in UTF-8. |
| `Posted` | `string` | `false` | The post date, as an RFC3339Nano time string. |
| `Extensions` | `object` | `true` | Application defined extensions. |

The character limits of the `Topic`, `Title` and `Content` should be application
dependent. The `Extensions` field can be used to implement
domain specific extensions, such as ~~tripfagging~~, and perhaps
the charater limits mentioned earlier.

### Example OP post

```json
{
  "Topic": "b",
  "Title": "I'm new here",
  "Protocol": "IPFS-TXT",
  "Version": "0.1.2",
  "Content": "QmWicCsiZuBdLPksfaj3qT6akFvHskwNDHSQSM3MLF1GRX",
  "Posted": "2020-05-21T00:04:24.3095002+01:00"
}
```

### Example reply post

```json
{
  "Thread": "QmPoZs6qaPVzCJTonnjPYifYbCM7Zn5nzsYwUyuQ8sLDY5",
  "Protocol": "IPFS-TXT",
  "Version": "0.1.2",
  "Content": "QmToUfTKxwXLGoXAGWeeJ7yoeaTkizE67iW8ybYVfdC9tR",
  "Posted": "2020-05-21T00:14:23.0384851+01:00"
}
```
