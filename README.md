# gin-appengine

Easier integration between [Gin](https://github.com/gin-gonic/gin) and [App Engine](https://github.com/golang/appengine)

  * **GaeToGinHandler** is a helper method to save you from repetitive `appengine.Context` initialization code. It should be used to register gin handlers that have a more convenient function signature.
  * **GaeErrorLogger** is middleware to write App Engine logs. It will pick up erors from gin (if you call `AbortWithError`, `Error` etc.) and log either Errors or Warnings, depending on request status code set by your handler (anything <500 is a warning).

```go
router := gin.Default()
router.Use(gingae.GaeErrorLogger)
router.GET("/route", gingae.GaeToGinHandler(handler))

func handler(c *gin.Context, gae appengine.Context)  {
  // Your stuff
}
```
