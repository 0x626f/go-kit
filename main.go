package main

import (
	"github.com/0x626f/go-kit/console"
)

func main() {

	console.WithGlobalLogLevel(console.TRACE)
	console.WithGlobalColoring()
	console.WithGlobalTimestamp()
	c := console.NewConsole().WithName("Test")

	c.Errorf("Hello %v", 123)
	c.Warningf("Hello %v", 123)
	c.Debugf("Hello %v", 123)
	c.Infof("Hello %v", 123)
	c.Tracef("Hello %v", 123)
	c.Logf("Hello %v", 123)

}
