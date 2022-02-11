## 資料庫包
***
### MongoDB

```go
// ReplicaSet
host := "192.168.10.79:27017,192.168.10.80:27017,192.168.10.81:27017"
db, _ := mongox.New(host).
    //SetRegistry(bson.NewRegistryBuilder().
    //RegisterDecoder(reflect.TypeOf(decimal.Decimal{}), mongox.Decimal{}).
    //RegisterEncoder(reflect.TypeOf(decimal.Decimal{}), mongox.Decimal{}).
    //Build()).
    SetReplicaSet("sh-rs-3").
    SetPool(1, 5, 10).
    SetPoolMonitor().
    Connect()

// Direct
db, _ := mongox.New(host).
    SetDirect(true).
    SetPool(1, 5, 10).
    // SetPoolMonitor().
    Connect()
```

### PostgreSQL

```go
pgdb, err := postgrex.New("127.0.0.1", "6432", "user", "password", "db").
    // SetTimeZone("PRC").
    // SetLogger(logger.Default.LogMode(logger.Info)).
    Connect(postgrex.Pool(1, 10, 10))
```

### MySQL

```go

mydb, err := mysqlx.New("127.0.0.1", "3306", "user", "password", "db").
    // SetAppendParameter(mysqlx.NewParamsmeter()).
    // SetCharset("utf8").
    // SetLoc("UTC").
    // SetLogger(logger.Default.LogMode(logger.Info)).
    Connect(mysqlx.Pool(1, 2, 180))

```