### xLogger 

zap based logger 

### install 

```shell
  go get github.com/clearcodecn/log@latest
```

### Usage

```shell
   import "github.com/clearcodecn/log"
```

* add hooks to logger and inject log fields to logger. 

```go

    log.AddHook(func(ctx context.Context) Field {
        reqid, ok := ctx.Value("reqid").(string)
        if !ok {
        return Field{}
        }
        return Any("reqid", reqid)
    })
    
    ctx := context.WithValue(context.Background(), "reqid", "123456")
	log.Logger(ctx).Info("help me")
```

* gin middleware
```go

var conf xlogger.GinLogConfigure
conf.SkipPrefix("/static","/favico.ico")
conf.AddHeaderKeys("reqid")
gin.Use(log.GinLog(conf))
```

* gorm middleware

```go
db := gorm.New()

db.Use(log.NewLoggerPlugin())

```