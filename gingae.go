package gingae

import (
	"errors"
	"runtime"

	"github.com/gin-gonic/gin"

	"appengine"
)

// KeyGaeContext is a key of appengine.Context object in gin.Context
const KeyGaeContext = "appengine.Context"

// GaeHandlerFunc is type of a typical app route handler function
type GaeHandlerFunc func(c *gin.Context, gae appengine.Context)

// GaeToGinHandler initializes Google App Engine context.
func GaeToGinHandler(handler GaeHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if gae, ok := c.Get(KeyGaeContext); ok {
			handler(c, gae.(appengine.Context))
		} else {
			gae := appengine.NewContext(c.Request)
			c.Set(KeyGaeContext, gae)
			handler(c, gae)
		}
	}
}

// GaeErrorLogger adds GAE error/warning logger.
// If there were any errors, will:
// Calls `Warningf` if returned status code is < 500, otherwise `Errorf`
// Criticalf on panic
var GaeErrorLogger = GaeToGinHandler(func(c *gin.Context, gae appengine.Context) {
	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, 1<<16)
			runtime.Stack(stack, true)
			gae.Criticalf("%s\nStacktrace:%s\n", err, stack)
			c.AbortWithError(500, errors.New("Shouldn't happen. See GAE log."))
		}
	}()
	c.Next()
	for _, err := range c.Errors {
		msg := err.Error()
		if c.Writer.Status() < 500 {
			gae.Warningf(msg)
		} else {
			gae.Errorf(msg)
		}
	}
})
