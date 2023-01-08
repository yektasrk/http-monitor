package main

import (
	"fmt"
	"github.com/yektasrk/http-monitor/inernal/httpserver"
)

func main() {

	fmt.Println("Http Monitor Init!")
	if err := httpserver.Serve(); err != nil {
		panic("error")
	}
}
