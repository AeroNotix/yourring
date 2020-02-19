package main

import (
	"fmt"
	"github.com/AeroNotix/yourring"
	"os"
)

func main() {
	r := &yourring.Ring{}
	r.Init()
	f, err := os.Open("/dev/urandom")
	if err != nil {
		fmt.Println(err)
	}
	r.QueueRead(f.Fd(), 1024, 0)
}
