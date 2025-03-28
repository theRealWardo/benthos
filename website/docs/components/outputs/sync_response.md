---
title: sync_response
type: output
status: stable
categories: ["Utility"]
---

<!--
     THIS FILE IS AUTOGENERATED!

     To make changes please edit the contents of:
     lib/output/sync_response.go
-->

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


Returns the final message payload back to the input origin of the message, where
it is dealt with according to that specific input type.

```yml
# Config fields, showing default values
output:
  label: ""
  sync_response: {}
```

For most inputs this mechanism is ignored entirely, in which case the sync
response is dropped without penalty. It is therefore safe to use this output
even when combining input types that might not have support for sync responses.
An example of an input able to utilise this is the `http_server`.

It is safe to combine this output with others using broker types. For example,
with the `http_server` input we could send the payload to a Kafka
topic and also send a modified payload back with:

```yaml
input:
  http_server:
    path: /post
output:
  broker:
    pattern: fan_out
    outputs:
      - kafka:
          addresses: [ TODO:9092 ]
          topic: foo_topic
      - sync_response: {}
        processors:
          - mapping: 'root = content().uppercase()'
```

Using the above example and posting the message 'hello world' to the endpoint
`/post` Benthos would send it unchanged to the topic
`foo_topic` and also respond with 'HELLO WORLD'.

For more information please read [Synchronous Responses](/docs/guides/sync_responses).


