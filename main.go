package main

import "github.com/0x626f/go-kit/logger"

func main() {
	log := logger.NewLogger("test")

	log.Infof("Some message to test")
}
