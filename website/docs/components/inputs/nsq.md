---
title: nsq
type: input
status: stable
categories: ["Services"]
---

<!--
     THIS FILE IS AUTOGENERATED!

     To make changes please edit the contents of:
     lib/input/nsq.go
-->

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

Subscribe to an NSQ instance topic and channel.


<Tabs defaultValue="common" values={[
  { label: 'Common', value: 'common', },
  { label: 'Advanced', value: 'advanced', },
]}>

<TabItem value="common">

```yml
# Common config fields, showing default values
input:
  label: ""
  nsq:
    nsqd_tcp_addresses: []
    lookupd_http_addresses: []
    topic: ""
    channel: ""
    user_agent: ""
    max_in_flight: 100
```

</TabItem>
<TabItem value="advanced">

```yml
# All config fields, showing default values
input:
  label: ""
  nsq:
    nsqd_tcp_addresses: []
    lookupd_http_addresses: []
    tls:
      enabled: false
      skip_cert_verify: false
      enable_renegotiation: false
      root_cas: ""
      root_cas_file: ""
      client_certs: []
    topic: ""
    channel: ""
    user_agent: ""
    max_in_flight: 100
```

</TabItem>
</Tabs>

## Fields

### `nsqd_tcp_addresses`

A list of nsqd addresses to connect to.


Type: `array`  
Default: `[]`  

### `lookupd_http_addresses`

A list of nsqlookupd addresses to connect to.


Type: `array`  
Default: `[]`  

### `tls`

Custom TLS settings can be used to override system defaults.


Type: `object`  

### `tls.enabled`

Whether custom TLS settings are enabled.


Type: `bool`  
Default: `false`  

### `tls.skip_cert_verify`

Whether to skip server side certificate verification.


Type: `bool`  
Default: `false`  

### `tls.enable_renegotiation`

Whether to allow the remote server to repeatedly request renegotiation. Enable this option if you're seeing the error message `local error: tls: no renegotiation`.


Type: `bool`  
Default: `false`  
Requires version 3.45.0 or newer  

### `tls.root_cas`

An optional root certificate authority to use. This is a string, representing a certificate chain from the parent trusted root certificate, to possible intermediate signing certificates, to the host certificate.
:::warning Secret
This field contains sensitive information that usually shouldn't be added to a config directly, read our [secrets page for more info](/docs/configuration/secrets).
:::


Type: `string`  
Default: `""`  

```yml
# Examples

root_cas: |-
  -----BEGIN CERTIFICATE-----
  ...
  -----END CERTIFICATE-----
```

### `tls.root_cas_file`

An optional path of a root certificate authority file to use. This is a file, often with a .pem extension, containing a certificate chain from the parent trusted root certificate, to possible intermediate signing certificates, to the host certificate.


Type: `string`  
Default: `""`  

```yml
# Examples

root_cas_file: ./root_cas.pem
```

### `tls.client_certs`

A list of client certificates to use. For each certificate either the fields `cert` and `key`, or `cert_file` and `key_file` should be specified, but not both.


Type: `array`  
Default: `[]`  

```yml
# Examples

client_certs:
  - cert: foo
    key: bar

client_certs:
  - cert_file: ./example.pem
    key_file: ./example.key
```

### `tls.client_certs[].cert`

A plain text certificate to use.


Type: `string`  
Default: `""`  

### `tls.client_certs[].key`

A plain text certificate key to use.
:::warning Secret
This field contains sensitive information that usually shouldn't be added to a config directly, read our [secrets page for more info](/docs/configuration/secrets).
:::


Type: `string`  
Default: `""`  

### `tls.client_certs[].cert_file`

The path of a certificate to use.


Type: `string`  
Default: `""`  

### `tls.client_certs[].key_file`

The path of a certificate key to use.


Type: `string`  
Default: `""`  

### `tls.client_certs[].password`

A plain text password for when the private key is password encrypted in PKCS#1 or PKCS#8 format. The obsolete `pbeWithMD5AndDES-CBC` algorithm is not supported for the PKCS#8 format. Warning: Since it does not authenticate the ciphertext, it is vulnerable to padding oracle attacks that can let an attacker recover the plaintext.
:::warning Secret
This field contains sensitive information that usually shouldn't be added to a config directly, read our [secrets page for more info](/docs/configuration/secrets).
:::


Type: `string`  
Default: `""`  

```yml
# Examples

password: foo

password: ${KEY_PASSWORD}
```

### `topic`

The topic to consume from.


Type: `string`  
Default: `""`  

### `channel`

The channel to consume from.


Type: `string`  
Default: `""`  

### `user_agent`

A user agent to assume when connecting.


Type: `string`  
Default: `""`  

### `max_in_flight`

The maximum number of pending messages to consume at any given time.


Type: `int`  
Default: `100`  


