### xLogger 

zap based logger 

### install 

```shell
  go get github.com/clearcodecn/xlogger@latest
```

### Usage

```shell
   import "github.com/clearcodecn/xlogger"
```

* add hooks to logger and inject log fields to logger. 

```go

    xlogger.AddHook(func(ctx context.Context) Field {
        reqid, ok := ctx.Value("reqid").(string)
        if !ok {
        return Field{}
        }
        return Any("reqid", reqid)
    })
    
    ctx := context.WithValue(context.Background(), "reqid", "123456")
    xlogger.Logger(ctx).Info("help me")
```

* gin middleware
```go

var conf xlogger.GinLogConfigure
conf.SkipPrefix("/static","/favico.ico")
conf.AddHeaderKeys("reqid")
gin.Use(xlogger.GinLog(conf))
```

* gorm middleware

```go
db := gorm.New()

db.Use(xlogger.NewLoggerPlugin())

```