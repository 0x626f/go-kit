package config

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestEnvConfig(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(file)
	envFile := filepath.Join(currentDir, ".test.env")

	type sample struct {
		First   int        `env:"FIRST"`
		Second  string     `env:"SECOND"`
		Third   float64    `env:"THIRD"`
		Fourth  []int8     `env:"FOURTH"`
		Fifth   bool       `env:"FIFTH"`
		Sixth   []string   `env:"SIXTH"`
		Seventh []bool     `env:"SEVENTH"`
		Eighth  [][]string `env:"EIGHTH"`
	}

	conf, err := FromFile[sample](envFile)

	if err != nil && err.Error() != "couldn't map dimensional arrays from .env" {
		t.Fatal(err)
	}

	if conf.First != 1 {
		t.Fatal("first mismatch")
	}

	if conf.Second != "test" {
		t.Fatal("second mismatch")
	}

	if conf.Third != 2.2 {
		t.Fatal("third mismatch")
	}

	fourthExpected := []int8{1, 2, 3, 4}
	if len(conf.Fourth) != len(fourthExpected) {
		t.Fatal("fourth len mismatch")
	}

	for index, expected := range fourthExpected {
		if conf.Fourth[index] != expected {
			t.Fatal("fourth mismatch")
		}
	}

	if !conf.Fifth {
		t.Fatal("fifth mismatch")
	}

	sixthExpected := []string{"one", "two", "three"}
	if len(conf.Sixth) != len(sixthExpected) {
		t.Fatal("sixth len mismatch")
	}

	for index, expected := range sixthExpected {
		if conf.Sixth[index] != expected {
			t.Fatal("sixth mismatch")
		}
	}

	seventhExpected := []bool{true, true, false, false}
	if len(conf.Seventh) != len(seventhExpected) {
		t.Fatal("seventh len mismatch")
	}

	for index, expected := range seventhExpected {
		if conf.Seventh[index] != expected {
			t.Fatal("seventh mismatch")
		}
	}

	if len(conf.Eighth) != 0 {
		t.Fatal("eight mismatch")
	}
}

func TestJsonConfig(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(file)
	envFile := filepath.Join(currentDir, ".test.json")

	type sample struct {
		First   int      `json:"first"`
		Second  string   `json:"second"`
		Third   float64  `json:"third"`
		Fourth  []int8   `json:"fourth"`
		Fifth   bool     `json:"fifth"`
		Sixth   []string `json:"sixth"`
		Seventh []bool   `json:"seventh"`
	}

	conf, err := FromFile[sample](envFile)

	if err != nil {
		t.Fatal(err)
	}

	if conf.First != 1 {
		t.Fatal("first mismatch")
	}

	if conf.Second != "test" {
		t.Fatal("second mismatch")
	}

	if conf.Third != 2.2 {
		t.Fatal("third mismatch")
	}

	fourthExpected := []int8{1, 2, 3, 4}
	if len(conf.Fourth) != len(fourthExpected) {
		t.Fatal("fourth len mismatch")
	}

	for index, expected := range fourthExpected {
		if conf.Fourth[index] != expected {
			t.Fatal("fourth mismatch")
		}
	}

	if !conf.Fifth {
		t.Fatal("fifth mismatch")
	}

	sixthExpected := []string{"one", "two", "three"}
	if len(conf.Sixth) != len(sixthExpected) {
		t.Fatal("sixth len mismatch")
	}

	for index, expected := range sixthExpected {
		if conf.Sixth[index] != expected {
			t.Fatal("sixth mismatch")
		}
	}

	seventhExpected := []bool{true, true, false, false}
	if len(conf.Seventh) != len(seventhExpected) {
		t.Fatal("seventh len mismatch")
	}

	for index, expected := range seventhExpected {
		if conf.Seventh[index] != expected {
			t.Fatal("seventh mismatch")
		}
	}
}

func TestEnvConfigWithLoad(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(file)
	envFile := filepath.Join(currentDir, ".test.env")

	type sample struct {
		First   int        `env:"FIRST"`
		Second  string     `env:"SECOND"`
		Third   float64    `env:"THIRD"`
		Fourth  []int8     `env:"FOURTH"`
		Fifth   bool       `env:"FIFTH"`
		Sixth   []string   `env:"SIXTH"`
		Seventh []bool     `env:"SEVENTH"`
		Eighth  [][]string `env:"EIGHTH"`
	}

	loadErr := LoadEnvs(envFile)

	if loadErr != nil {
		t.Fatal(loadErr)
	}

	conf, err := FromEnvs[sample]()

	if err != nil && err.Error() != "couldn't map dimensional arrays from .env" {
		t.Fatal(err)

	}

	if conf.First != 1 {
		t.Fatal("first mismatch")
	}

	if conf.Second != "test" {
		t.Fatal("second mismatch")
	}

	if conf.Third != 2.2 {
		t.Fatal("third mismatch")
	}

	fourthExpected := []int8{1, 2, 3, 4}
	if len(conf.Fourth) != len(fourthExpected) {
		t.Fatal("fourth len mismatch")
	}

	for index, expected := range fourthExpected {
		if conf.Fourth[index] != expected {
			t.Fatal("fourth mismatch")
		}
	}

	if !conf.Fifth {
		t.Fatal("fifth mismatch")
	}

	sixthExpected := []string{"one", "two", "three"}
	if len(conf.Sixth) != len(sixthExpected) {
		t.Fatal("sixth len mismatch")
	}

	for index, expected := range sixthExpected {
		if conf.Sixth[index] != expected {
			t.Fatal("sixth mismatch")
		}
	}

	seventhExpected := []bool{true, true, false, false}
	if len(conf.Seventh) != len(seventhExpected) {
		t.Fatal("seventh len mismatch")
	}

	for index, expected := range seventhExpected {
		if conf.Seventh[index] != expected {
			t.Fatal("seventh mismatch")
		}
	}

	if len(conf.Eighth) != 0 {
		t.Fatal("eight mismatch")
	}
}
