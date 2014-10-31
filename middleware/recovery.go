package middleware

import (
	"fmt"
	"github.com/juju/errgo"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

type Recovery struct {
	Logger *log.Logger
}

func NewRecovery() *Recovery {
	return &Recovery{
		Logger: log.New(os.Stdout, "[klask] ", 0),
	}
}

func (self *Recovery) ServeHTTP(
	res http.ResponseWriter,
	req *http.Request,
	next http.HandlerFunc,
) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}

		res.WriteHeader(http.StatusInternalServerError)
		stack := debug.Stack()

		// send stack to server logger
		format := "PANIC: %s\n%s"
		self.Logger.Printf(format, err, stack)

		// respond with more informative stack
		for err != nil {
			switch e := err.(type) {
			case *errgo.Err:
				fmt.Fprintf(res, "in %s\n", e.Location())
				if msg := e.Message(); len(msg) > 0 {
					fmt.Fprintf(res, "  %s\n", e.Message())
				}
				err = e.Underlying()
			case error:
				fmt.Fprintf(res, "%s", e.Error())
				err = nil
			}
		}
		return
	}()

	next(res, req)
}
