package gingae

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"

	"golang.org/x/net/context"

	"github.com/gin-gonic/gin"

	"google.golang.org/appengine"
)

// KeyGaeContext is a key of context.Context object in gin.Context
const KeyGaeContext = "context.Context"

// GaeHandlerFunc is type of a typical app route handler function
type GaeHandlerFunc func(c *gin.Context, gae context.Context)

// GaeToGinHandler initializes Google App Engine context.
func GaeToGinHandler(handler GaeHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if gae, ok := c.Get(KeyGaeContext); ok {
			handler(c, gae.(context.Context))
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
var GaeErrorLogger = GaeToGinHandler(func(c *gin.Context, gae context.Context) {
	errLogger := log.New(os.Stderr, "", 0)
	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, 1<<16)
			runtime.Stack(stack, true)
			fmt.Fprintf(os.Stderr, "%s\nStacktrace:%s\n", err, stack)
			c.AbortWithError(500, errors.New("Shouldn't happen - see GAE log"))
		}
	}()
	c.Next()
	for _, err := range c.Errors {
		msg := err.Error()
		if c.Writer.Status() < 500 {
			log.Printf(msg)
		} else {
			errLogger.Printf(msg)
		}
	}
})
