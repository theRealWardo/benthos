---
title: sql_insert
type: processor
status: stable
categories: ["Integration"]
---

<!--
     THIS FILE IS AUTOGENERATED!

     To make changes please edit the contents of:
     lib/processor/sql_insert.go
-->

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

Inserts rows into an SQL database for each message, and leaves the message unchanged.

Introduced in version 3.59.0.


<Tabs defaultValue="common" values={[
  { label: 'Common', value: 'common', },
  { label: 'Advanced', value: 'advanced', },
]}>

<TabItem value="common">

```yml
# Common config fields, showing default values
label: ""
sql_insert:
  driver: ""
  dsn: ""
  table: ""
  columns: []
  args_mapping: ""
```

</TabItem>
<TabItem value="advanced">

```yml
# All config fields, showing default values
label: ""
sql_insert:
  driver: ""
  dsn: ""
  table: ""
  columns: []
  args_mapping: ""
  prefix: ""
  suffix: ""
  init_files: []
  init_statement: ""
  conn_max_idle_time: ""
  conn_max_life_time: ""
  conn_max_idle: 0
  conn_max_open: 0
```

</TabItem>
</Tabs>

If the insert fails to execute then the message will still remain unchanged and the error can be caught using error handling methods outlined [here](/docs/configuration/error_handling).

## Examples

<Tabs defaultValue="Table Insert (MySQL)" values={[
{ label: 'Table Insert (MySQL)', value: 'Table Insert (MySQL)', },
]}>

<TabItem value="Table Insert (MySQL)">


Here we insert rows into a database by populating the columns id, name and topic with values extracted from messages and metadata:

```yaml
pipeline:
  processors:
    - sql_insert:
        driver: mysql
        dsn: foouser:foopassword@tcp(localhost:3306)/foodb
        table: footable
        columns: [ id, name, topic ]
        args_mapping: |
          root = [
            this.user.id,
            this.user.name,
            meta("kafka_topic"),
          ]
```

</TabItem>
</Tabs>

## Fields

### `driver`

A database [driver](#drivers) to use.


Type: `string`  
Options: `mysql`, `postgres`, `clickhouse`, `mssql`, `sqlite`, `oracle`, `snowflake`.

### `dsn`

A Data Source Name to identify the target database.

#### Drivers

The following is a list of supported drivers, their placeholder style, and their respective DSN formats:

| Driver | Data Source Name Format |
|---|---|
| `clickhouse` | [`clickhouse://[username[:password]@][netloc][:port]/dbname[?param1=value1&...&paramN=valueN]`](https://github.com/ClickHouse/clickhouse-go#dsn) |
| `mysql` | `[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]` |
| `postgres` | `postgres://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]` |
| `mssql` | `sqlserver://[user[:password]@][netloc][:port][?database=dbname&param1=value1&...]` |
| `sqlite` | `file:/path/to/filename.db[?param&=value1&...]` |
| `oracle` | `oracle://[username[:password]@][netloc][:port]/service_name?server=server2&server=server3` |
| `snowflake` | `username[:password]@account_identifier/dbname/schemaname[?param1=value&...&paramN=valueN]` |

Please note that the `postgres` driver enforces SSL by default, you can override this with the parameter `sslmode=disable` if required.

The `snowflake` driver supports multiple DSN formats. Please consult [the docs](https://pkg.go.dev/github.com/snowflakedb/gosnowflake#hdr-Connection_String) for more details. For [key pair authentication](https://docs.snowflake.com/en/user-guide/key-pair-auth.html#configuring-key-pair-authentication), the DSN has the following format: `<snowflake_user>@<snowflake_account>/<db_name>/<schema_name>?warehouse=<warehouse>&role=<role>&authenticator=snowflake_jwt&privateKey=<base64_url_encoded_private_key>`, where the value for the `privateKey` parameter can be constructed from an unencrypted RSA private key file `rsa_key.p8` using `openssl enc -d -base64 -in rsa_key.p8 | basenc --base64url -w0` (you can use `gbasenc` insted of `basenc` on OSX if you install `coreutils` via Homebrew). If you have a password-encrypted private key, you can decrypt it using `openssl pkcs8 -in rsa_key_encrypted.p8 -out rsa_key.p8`. Also, make sure fields such as the username are URL-encoded.


Type: `string`  

```yml
# Examples

dsn: clickhouse://username:password@host1:9000,host2:9000/database?dial_timeout=200ms&max_execution_time=60

dsn: foouser:foopassword@tcp(localhost:3306)/foodb

dsn: postgres://foouser:foopass@localhost:5432/foodb?sslmode=disable

dsn: oracle://foouser:foopass@localhost:1521/service_name
```

### `table`

The table to insert to.


Type: `string`  

```yml
# Examples

table: foo
```

### `columns`

A list of columns to insert.


Type: `array`  

```yml
# Examples

columns:
  - foo
  - bar
  - baz
```

### `args_mapping`

A [Bloblang mapping](/docs/guides/bloblang/about) which should evaluate to an array of values matching in size to the number of columns specified.


Type: `string`  

```yml
# Examples

args_mapping: root = [ this.cat.meow, this.doc.woofs[0] ]

args_mapping: root = [ meta("user.id") ]
```

### `prefix`

An optional prefix to prepend to the insert query (before INSERT).


Type: `string`  

### `suffix`

An optional suffix to append to the insert query.


Type: `string`  

```yml
# Examples

suffix: ON CONFLICT (name) DO NOTHING
```

### `init_files`

An optional list of file paths containing SQL statements to execute immediately upon the first connection to the target database. This is a useful way to initialise tables before processing data. Glob patterns are supported, including super globs (double star).

Care should be taken to ensure that the statements are idempotent, and therefore would not cause issues when run multiple times after service restarts. If both `init_statement` and `init_files` are specified the `init_statement` is executed _after_ the `init_files`.

If a statement fails for any reason a warning log will be emitted but the operation of this component will not be stopped.


Type: `array`  
Requires version 4.10.0 or newer  

```yml
# Examples

init_files:
  - ./init/*.sql

init_files:
  - ./foo.sql
  - ./bar.sql
```

### `init_statement`

An optional SQL statement to execute immediately upon the first connection to the target database. This is a useful way to initialise tables before processing data. Care should be taken to ensure that the statement is idempotent, and therefore would not cause issues when run multiple times after service restarts.

If both `init_statement` and `init_files` are specified the `init_statement` is executed _after_ the `init_files`.

If the statement fails for any reason a warning log will be emitted but the operation of this component will not be stopped.


Type: `string`  
Requires version 4.10.0 or newer  

```yml
# Examples

init_statement: |2
  CREATE TABLE IF NOT EXISTS some_table (
    foo varchar(50) not null,
    bar integer,
    baz varchar(50),
    primary key (foo)
  ) WITHOUT ROWID;
```

### `conn_max_idle_time`

An optional maximum amount of time a connection may be idle. Expired connections may be closed lazily before reuse. If value <= 0, connections are not closed due to a connection's idle time.


Type: `string`  

### `conn_max_life_time`

An optional maximum amount of time a connection may be reused. Expired connections may be closed lazily before reuse. If value <= 0, connections are not closed due to a connection's age.


Type: `string`  

### `conn_max_idle`

An optional maximum number of connections in the idle connection pool. If conn_max_open is greater than 0 but less than the new conn_max_idle, then the new conn_max_idle will be reduced to match the conn_max_open limit. If value <= 0, no idle connections are retained. The default max idle connections is currently 2. This may change in a future release.


Type: `int`  

### `conn_max_open`

An optional maximum number of open connections to the database. If conn_max_idle is greater than 0 and the new conn_max_open is less than conn_max_idle, then conn_max_idle will be reduced to match the new conn_max_open limit. If value <= 0, then there is no limit on the number of open connections. The default is 0 (unlimited).


Type: `int`  


