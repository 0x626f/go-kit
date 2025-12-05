package main

import (
	"github.com/0x626f/go-kit/logger"
	"time"
)

func main() {
	Logger, cancel := logger.NewLogger().WithTimestamp().WithName("Main Test").WithLogLevel(logger.TRACE).WithAsync(true, 100)

	type Event struct {
		Name   string `json:"eventName"`
		Source string `json:"eventSource"`
	}

	event1 := Event{Name: "test", Source: "main"}
	event2 := Event{Name: "test", Source: "dev"}

	Logger.ErrorJSONf(event1, "error test #%v", 1)
	Logger.TraceJSONf(event1, "trace test #%v", 1)
	Logger.DebugJSONf(event1, "debug test #%v", 1)
	Logger.InfoJSONf(event1, "info test #%v", 1)
	Logger.WarningJSONf(event1, "warning test #%v", 1)
	Logger.LogJSONf(event1, "warning test #%v", 1)

	Logger.ErrorJSONf(event2, "error test #%v", 2)
	Logger.TraceJSONf(event2, "trace test #%v", 2)
	Logger.DebugJSONf(event2, "debug test #%v", 2)
	Logger.InfoJSONf(event2, "info test #%v", 2)
	Logger.WarningJSONf(event2, "warning test #%v", 2)
	Logger.LogJSONf(event2, "warning test #%v", 2)

	Logger.InfoObjectf("User created: %v", "0x0012321").AssignString("name", "blabla").AssignInt("age", 50).AssignFloat32("weight", 81.663).Build()
	Logger.InfoObjectf("User created: %v", "0x0012321").AssignString("name", "blabla").AssignInt("age", 51).AssignByte("flag", '\t').Build()
	Logger.InfoObjectf("User created: %v", "0x0012321").AssignString("name", "blabla").AssignInt("age", 52).AssignStringArray("categories", []string{"science", "tech"}).AssignFloat32Array("some", []float32{1.2, 1.4, 0.000023}).Build()
	Logger.InfoObjectf("User created: %v", "0x0012321").AssignString("name", "blabla").AssignInt("age", 53).Build()

	Logger.DebugObjectf("User created: %v", "0x0012321").AssignString("name", "blabla").NestedStart("credential").AssignStringArray("options", []string{"password", "fingerprint"}).NestedEnd().Build()

	time.Sleep(5 * time.Second)
	cancel()
	time.Sleep(5 * time.Second)
}
