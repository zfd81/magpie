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