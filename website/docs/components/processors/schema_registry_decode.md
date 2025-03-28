---
title: schema_registry_decode
type: processor
status: experimental
categories: ["Parsing","Integration"]
---

<!--
     THIS FILE IS AUTOGENERATED!

     To make changes please edit the contents of:
     lib/processor/schema_registry_decode.go
-->

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

:::caution EXPERIMENTAL
This component is experimental and therefore subject to change or removal outside of major version releases.
:::
Automatically decodes and validates messages with schemas from a Confluent Schema Registry service.


<Tabs defaultValue="common" values={[
  { label: 'Common', value: 'common', },
  { label: 'Advanced', value: 'advanced', },
]}>

<TabItem value="common">

```yml
# Common config fields, showing default values
label: ""
schema_registry_decode:
  url: ""
```

</TabItem>
<TabItem value="advanced">

```yml
# All config fields, showing default values
label: ""
schema_registry_decode:
  avro_raw_json: false
  url: ""
  oauth:
    enabled: false
    consumer_key: ""
    consumer_secret: ""
    access_token: ""
    access_token_secret: ""
  basic_auth:
    enabled: false
    username: ""
    password: ""
  jwt:
    enabled: false
    private_key_file: ""
    signing_method: ""
    claims: {}
    headers: {}
  tls:
    skip_cert_verify: false
    enable_renegotiation: false
    root_cas: ""
    root_cas_file: ""
    client_certs: []
```

</TabItem>
</Tabs>

Decodes messages automatically from a schema stored within a [Confluent Schema Registry service](https://docs.confluent.io/platform/current/schema-registry/index.html) by extracting a schema ID from the message and obtaining the associated schema from the registry. If a message fails to match against the schema then it will remain unchanged and the error can be caught using error handling methods outlined [here](/docs/configuration/error_handling).

Currently only Avro schemas are supported.

### Avro JSON Format

This processor creates documents formatted as [Avro JSON](https://avro.apache.org/docs/current/specification/_print/#json-encoding) when decoding with Avro schemas. In this format the value of a union is encoded in JSON as follows:

- if its type is `null`, then it is encoded as a JSON `null`;
- otherwise it is encoded as a JSON object with one name/value pair whose name is the type's name and whose value is the recursively encoded value. For Avro's named types (record, fixed or enum) the user-specified name is used, for other types the type name is used.

For example, the union schema `["null","string","Foo"]`, where `Foo` is a record name, would encode:

- `null` as `null`;
- the string `"a"` as `{"string": "a"}`; and
- a `Foo` instance as `{"Foo": {...}}`, where `{...}` indicates the JSON encoding of a `Foo` instance.

However, it is possible to instead create documents in [standard/raw JSON format](https://pkg.go.dev/github.com/linkedin/goavro/v2#NewCodecForStandardJSONFull) by setting the field [`avro_raw_json`](#avro_raw_json) to `true`.

## Fields

### `avro_raw_json`

Whether Avro messages should be decoded into normal JSON ("json that meets the expectations of regular internet json") rather than [Avro JSON](https://avro.apache.org/docs/current/specification/_print/#json-encoding). If `true` the schema returned from the subject should be decoded as [standard json](https://pkg.go.dev/github.com/linkedin/goavro/v2#NewCodecForStandardJSONFull) instead of as [avro json](https://pkg.go.dev/github.com/linkedin/goavro/v2#NewCodec). There is a [comment in goavro](https://github.com/linkedin/goavro/blob/5ec5a5ee7ec82e16e6e2b438d610e1cab2588393/union.go#L224-L249), the [underlining library used for avro serialization](https://github.com/linkedin/goavro), that explains in more detail the difference between the standard json and avro json.


Type: `bool`  
Default: `false`  

### `url`

The base URL of the schema registry service.


Type: `string`  

### `oauth`

Allows you to specify open authentication via OAuth version 1.


Type: `object`  
Requires version 4.7.0 or newer  

### `oauth.enabled`

Whether to use OAuth version 1 in requests.


Type: `bool`  
Default: `false`  

### `oauth.consumer_key`

A value used to identify the client to the service provider.


Type: `string`  
Default: `""`  

### `oauth.consumer_secret`

A secret used to establish ownership of the consumer key.
:::warning Secret
This field contains sensitive information that usually shouldn't be added to a config directly, read our [secrets page for more info](/docs/configuration/secrets).
:::


Type: `string`  
Default: `""`  

### `oauth.access_token`

A value used to gain access to the protected resources on behalf of the user.


Type: `string`  
Default: `""`  

### `oauth.access_token_secret`

A secret provided in order to establish ownership of a given access token.
:::warning Secret
This field contains sensitive information that usually shouldn't be added to a config directly, read our [secrets page for more info](/docs/configuration/secrets).
:::


Type: `string`  
Default: `""`  

### `basic_auth`

Allows you to specify basic authentication.


Type: `object`  
Requires version 4.7.0 or newer  

### `basic_auth.enabled`

Whether to use basic authentication in requests.


Type: `bool`  
Default: `false`  

### `basic_auth.username`

A username to authenticate as.


Type: `string`  
Default: `""`  

### `basic_auth.password`

A password to authenticate with.
:::warning Secret
This field contains sensitive information that usually shouldn't be added to a config directly, read our [secrets page for more info](/docs/configuration/secrets).
:::


Type: `string`  
Default: `""`  

### `jwt`

BETA: Allows you to specify JWT authentication.


Type: `object`  
Requires version 4.7.0 or newer  

### `jwt.enabled`

Whether to use JWT authentication in requests.


Type: `bool`  
Default: `false`  

### `jwt.private_key_file`

A file with the PEM encoded via PKCS1 or PKCS8 as private key.


Type: `string`  
Default: `""`  

### `jwt.signing_method`

A method used to sign the token such as RS256, RS384 or RS512.


Type: `string`  
Default: `""`  

### `jwt.claims`

A value used to identify the claims that issued the JWT.


Type: `object`  

### `jwt.headers`

Add optional key/value headers to the JWT.


Type: `object`  

### `tls`

Custom TLS settings can be used to override system defaults.


Type: `object`  

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


