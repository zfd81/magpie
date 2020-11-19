magctl
========

`magctl` is a command line client for [magpie].



## Data operation commands

### LOAD  \<table name\> \<data file\>

Load data file

#### Examples

```bash
./magctl load userInfo data.csv
# [INFO] Start time: 2020-11-16 17:16:12.985 
# [INFO] End time: 2020-11-16 17:16:12.997 
# [INFO] Elapsed time: 12.179563ms 
# [INFO] Record Count: 1900 
# [INFO] Data loading complete 
```

### EXEC \<sql\>

Execute SQL

#### Examples

Add user information

```bash
./magctl exec insert into userInfo values ('id202011151212','zfd','pwd123',32)
# OK
./magctl exec insert into userInfo values ('id202011151212','zfd','pwd123',32),('id202011151215','lm','pwd456',31)
# OK
./magctl exec insert into userInfo (id,name,pwd,age) values ('id202011151212','zfd','pwd123',32)
# OK
./magctl exec insert into userInfo (id,name,pwd,age) values ('id202011151212','zfd','pwd123',32),('id202011151215','lm','pwd456',31)
# OK
```

Delete user information

```bash
./magctl exec delete from userInfo where id='id202011151212'
# OK
```

Modify user information

```bash
./magctl exec update userInfo set age=35 where id='id202011151212'
# OK
```


Query user information

```bash
./magctl exec select id,name,pwd,age from userInfo where id='id202011151212'
#[INFO] {"id":"id202011151212","name":"zfd","pwd":"pwd123","age":35} 

./magctl exec select id,name,pwd,if(age>20,'man','young man') age from userInfo where id='id202011151212'
#[INFO] {"id":"id202011151212","name":"zfd","pwd":"pwd123","age":"man"} 

```

## Table maintenance commands

### TABLE \<subcommand\>

TABLE provides commands for managing magpie cluster table information.

### TABLE ADD \<table definition file\>

TABLE ADD create a new table into the magpie cluster.

#### Table definition

```json
{
  "name": "userInfo",
  "cols": [
    {
      "name": "id",
      "dataType": "string"
    },
    {
      "name": "name",
      "dataType": "string"
    },
    {
      "name": "pwd",
      "dataType": "string"
    },
    {
      "name": "age",
      "dataType": "int"
    }
  ],
  "keys": [
    "id"
  ]
}
```

#### Example

```bash
./magctl table add define.json
# [INFO] Table userInfo created successfully
```

### TABLE DEL \<table name\>

TABLE DEL removes a table of a magpie cluster

#### Example

```bash
./magctl table del userInfo
# [INFO] Table userInfo deleted successfully
```

### TABLE DESC \<table name\> [options]

TABLE DESC print table details


#### Options

- file path -- File path of table details output

#### Example

```bash
./magctl table desc userInfo
# Table[userInfo] details:
  {
    "name": "userInfo",
    "cols": [
      {
        "name": "id",
        "dataType": "string",
      },
      {
        "name": "name",
        "dataType": "string",
      },
      {
        "name": "pwd",
        "dataType": "string",
      },
      {
        "name": "age",
        "dataType": "int",
      }
    ],
    "keys": [
      "id"
    ]
  }

./magctl table desc userInfo define.json
# [INFO] Export table structure succeeded
```

### TABLE LIST

TABLE LIST prints all tables.

#### Examples

```bash
./magctl table list
+----+--------------------------------+
| SN |                           Name |
+----+--------------------------------+
|  1 |                       userInfo |
+----+--------------------------------+
```

## Cluster maintenance commands

### MEMBER \<subcommand\>

MEMBER provides commands for managing magpie cluster membership.

### MEMBER LIST

MEMBER LIST prints the member details for all members associated with an magpie cluster.

#### Examples

```bash
./magctl member list
+----------------------+------------------+------------+-----------+---------------------+
|       ENDPOINT       |        ID        |    Team    | IS LEADER |    START-UP TIME    |
+----------------------+------------------+------------+-----------+---------------------+
|  192.168.31.158:8143 | 694d75db569fec75 |     magpie |      true | 2020-11-18 21:26:47 |
|  192.168.31.158:8888 | 694d75db569fec7e |     magpie |     false | 2020-11-18 21:27:33 |
|  192.168.31.158:7890 | 694d75db569fec85 |     magpie |     false | 2020-11-18 21:28:20 |
+----------------------+------------------+------------+-----------+---------------------+
```

### MEMBER STAT

MEMBER STAT prints the member status information.

#### Examples

```bash
./magctl member stat
+--------------------+--------------+-----------+-----------+
|     TABLE NAME     | COLUMN COUNT | ROW COUNT |  SIZE(KB) |
+--------------------+--------------+-----------+-----------+
|          tagDbInfo |         3000 |         0 |         0 |
|           userInfo |            4 |      9977 |       297 |
+--------------------+--------------+-----------+-----------+
```