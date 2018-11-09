// Copyright 2017 Sebastian Stauch.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Logger struct {
	LogFile *os.File
}

func (writer Logger) Write(bytes []byte) (int, error) {
	logline := fmt.Sprintf("%21s %s",
		time.Now().Local().Format("2006-01-02 15:04:05.999"),
		string(bytes))
	if writer.LogFile != nil {
		return writer.LogFile.WriteString(logline)
	} else {
		fmt.Printf(logline)
	}

	return len([]byte(logline)), nil

}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()

		// after request
		latency := time.Since(t)

		// access the status we are sending
		status := c.Writer.Status()
		logstring := fmt.Sprintf("remoteAddr=%-15s code=%-d method=%-6s duration=%-s path=%s",
			c.Request.RemoteAddr,
			status,
			c.Request.Method,
			latency,
			c.Request.RequestURI)

		log.Println(logstring)

	}
}
