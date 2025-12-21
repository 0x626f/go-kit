package env

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
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

	t.Run("ErrorHandling", func(t *testing.T) {
		if err != nil && err.Error() != "couldn't map dimensional arrays from .env" {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("IntField", func(t *testing.T) {
		if conf.First != 1 {
			t.Errorf("First = %v, want 1", conf.First)
		}
	})

	t.Run("StringField", func(t *testing.T) {
		if conf.Second != "test" {
			t.Errorf("Second = %v, want 'test'", conf.Second)
		}
	})

	t.Run("FloatField", func(t *testing.T) {
		if conf.Third != 2.2 {
			t.Errorf("Third = %v, want 2.2", conf.Third)
		}
	})

	t.Run("Int8SliceField", func(t *testing.T) {
		fourthExpected := []int8{1, 2, 3, 4}
		if len(conf.Fourth) != len(fourthExpected) {
			t.Errorf("Fourth length = %v, want %v", len(conf.Fourth), len(fourthExpected))
		}
		for index, expected := range fourthExpected {
			if conf.Fourth[index] != expected {
				t.Errorf("Fourth[%d] = %v, want %v", index, conf.Fourth[index], expected)
			}
		}
	})

	t.Run("BoolField", func(t *testing.T) {
		if !conf.Fifth {
			t.Errorf("Fifth = %v, want true", conf.Fifth)
		}
	})

	t.Run("StringSliceField", func(t *testing.T) {
		sixthExpected := []string{"one", "two", "three"}
		if len(conf.Sixth) != len(sixthExpected) {
			t.Errorf("Sixth length = %v, want %v", len(conf.Sixth), len(sixthExpected))
		}
		for index, expected := range sixthExpected {
			if conf.Sixth[index] != expected {
				t.Errorf("Sixth[%d] = %v, want %v", index, conf.Sixth[index], expected)
			}
		}
	})

	t.Run("BoolSliceField", func(t *testing.T) {
		seventhExpected := []bool{true, true, false, false}
		if len(conf.Seventh) != len(seventhExpected) {
			t.Errorf("Seventh length = %v, want %v", len(conf.Seventh), len(seventhExpected))
		}
		for index, expected := range seventhExpected {
			if conf.Seventh[index] != expected {
				t.Errorf("Seventh[%d] = %v, want %v", index, conf.Seventh[index], expected)
			}
		}
	})

	t.Run("DimensionalArrayField", func(t *testing.T) {
		if len(conf.Eighth) != 0 {
			t.Errorf("Eighth length = %v, want 0", len(conf.Eighth))
		}
	})
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

	t.Run("NoError", func(t *testing.T) {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("IntField", func(t *testing.T) {
		if conf.First != 1 {
			t.Errorf("First = %v, want 1", conf.First)
		}
	})

	t.Run("StringField", func(t *testing.T) {
		if conf.Second != "test" {
			t.Errorf("Second = %v, want 'test'", conf.Second)
		}
	})

	t.Run("FloatField", func(t *testing.T) {
		if conf.Third != 2.2 {
			t.Errorf("Third = %v, want 2.2", conf.Third)
		}
	})

	t.Run("Int8SliceField", func(t *testing.T) {
		fourthExpected := []int8{1, 2, 3, 4}
		if len(conf.Fourth) != len(fourthExpected) {
			t.Errorf("Fourth length = %v, want %v", len(conf.Fourth), len(fourthExpected))
		}
		for index, expected := range fourthExpected {
			if conf.Fourth[index] != expected {
				t.Errorf("Fourth[%d] = %v, want %v", index, conf.Fourth[index], expected)
			}
		}
	})

	t.Run("BoolField", func(t *testing.T) {
		if !conf.Fifth {
			t.Errorf("Fifth = %v, want true", conf.Fifth)
		}
	})

	t.Run("StringSliceField", func(t *testing.T) {
		sixthExpected := []string{"one", "two", "three"}
		if len(conf.Sixth) != len(sixthExpected) {
			t.Errorf("Sixth length = %v, want %v", len(conf.Sixth), len(sixthExpected))
		}
		for index, expected := range sixthExpected {
			if conf.Sixth[index] != expected {
				t.Errorf("Sixth[%d] = %v, want %v", index, conf.Sixth[index], expected)
			}
		}
	})

	t.Run("BoolSliceField", func(t *testing.T) {
		seventhExpected := []bool{true, true, false, false}
		if len(conf.Seventh) != len(seventhExpected) {
			t.Errorf("Seventh length = %v, want %v", len(conf.Seventh), len(seventhExpected))
		}
		for index, expected := range seventhExpected {
			if conf.Seventh[index] != expected {
				t.Errorf("Seventh[%d] = %v, want %v", index, conf.Seventh[index], expected)
			}
		}
	})
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

	t.Run("LoadEnvs", func(t *testing.T) {
		if loadErr != nil {
			t.Fatalf("LoadEnvs failed: %v", loadErr)
		}
	})

	conf, err := FromEnvs[sample]()

	t.Run("ErrorHandling", func(t *testing.T) {
		if err != nil && err.Error() != "couldn't map dimensional arrays from .env" {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("IntField", func(t *testing.T) {
		if conf.First != 1 {
			t.Errorf("First = %v, want 1", conf.First)
		}
	})

	t.Run("StringField", func(t *testing.T) {
		if conf.Second != "test" {
			t.Errorf("Second = %v, want 'test'", conf.Second)
		}
	})

	t.Run("FloatField", func(t *testing.T) {
		if conf.Third != 2.2 {
			t.Errorf("Third = %v, want 2.2", conf.Third)
		}
	})

	t.Run("Int8SliceField", func(t *testing.T) {
		fourthExpected := []int8{1, 2, 3, 4}
		if len(conf.Fourth) != len(fourthExpected) {
			t.Errorf("Fourth length = %v, want %v", len(conf.Fourth), len(fourthExpected))
		}
		for index, expected := range fourthExpected {
			if conf.Fourth[index] != expected {
				t.Errorf("Fourth[%d] = %v, want %v", index, conf.Fourth[index], expected)
			}
		}
	})

	t.Run("BoolField", func(t *testing.T) {
		if !conf.Fifth {
			t.Errorf("Fifth = %v, want true", conf.Fifth)
		}
	})

	t.Run("StringSliceField", func(t *testing.T) {
		sixthExpected := []string{"one", "two", "three"}
		if len(conf.Sixth) != len(sixthExpected) {
			t.Errorf("Sixth length = %v, want %v", len(conf.Sixth), len(sixthExpected))
		}
		for index, expected := range sixthExpected {
			if conf.Sixth[index] != expected {
				t.Errorf("Sixth[%d] = %v, want %v", index, conf.Sixth[index], expected)
			}
		}
	})

	t.Run("BoolSliceField", func(t *testing.T) {
		seventhExpected := []bool{true, true, false, false}
		if len(conf.Seventh) != len(seventhExpected) {
			t.Errorf("Seventh length = %v, want %v", len(conf.Seventh), len(seventhExpected))
		}
		for index, expected := range seventhExpected {
			if conf.Seventh[index] != expected {
				t.Errorf("Seventh[%d] = %v, want %v", index, conf.Seventh[index], expected)
			}
		}
	})

	t.Run("DimensionalArrayField", func(t *testing.T) {
		if len(conf.Eighth) != 0 {
			t.Errorf("Eighth length = %v, want 0", len(conf.Eighth))
		}
	})
}

// TestCanConvertFromEnv tests the canConvertFromEnv function with various reflect.Kind types
func TestCanConvertFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		kind     reflect.Kind
		expected bool
	}{
		// Supported primitive types
		{"String", reflect.String, true},
		{"Int", reflect.Int, true},
		{"Int8", reflect.Int8, true},
		{"Int16", reflect.Int16, true},
		{"Int32", reflect.Int32, true},
		{"Int64", reflect.Int64, true},
		{"Uint", reflect.Uint, true},
		{"Uint8", reflect.Uint8, true},
		{"Uint16", reflect.Uint16, true},
		{"Uint32", reflect.Uint32, true},
		{"Uint64", reflect.Uint64, true},
		{"Float32", reflect.Float32, true},
		{"Float64", reflect.Float64, true},
		{"Bool", reflect.Bool, true},
		{"Slice", reflect.Slice, true},

		// Unsupported types
		{"Map", reflect.Map, false},
		{"Struct", reflect.Struct, false},
		{"Ptr", reflect.Ptr, false},
		{"Interface", reflect.Interface, false},
		{"Chan", reflect.Chan, false},
		{"Func", reflect.Func, false},
		{"Array", reflect.Array, false},
		{"Complex64", reflect.Complex64, false},
		{"Complex128", reflect.Complex128, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := canConvertFromEnv(tt.kind)
			if result != tt.expected {
				t.Errorf("canConvertFromEnv(%v) = %v, want %v", tt.kind, result, tt.expected)
			}
		})
	}
}

// TestGetEnvAs tests the GetEnvAs generic function with various types
func TestGetEnvAs(t *testing.T) {
	// Clean up environment variables after test
	defer func() {
		os.Unsetenv("TEST_INT")
		os.Unsetenv("TEST_STRING")
		os.Unsetenv("TEST_BOOL")
		os.Unsetenv("TEST_FLOAT")
		os.Unsetenv("TEST_UINT")
		os.Unsetenv("TEST_SLICE_INT")
		os.Unsetenv("TEST_INVALID_INT")
	}()

	t.Run("GetEnvAs_Int", func(t *testing.T) {
		os.Setenv("TEST_INT", "42")
		result := GetEnvAs("TEST_INT", 0)
		if result != 42 {
			t.Errorf("GetEnvAs[int] = %v, want 42", result)
		}
	})

	t.Run("GetEnvAs_Int_Missing", func(t *testing.T) {
		os.Unsetenv("TEST_INT")
		result := GetEnvAs("TEST_INT", 100)
		if result != 100 {
			t.Errorf("GetEnvAs[int] with missing var = %v, want 100", result)
		}
	})

	t.Run("GetEnvAs_Int_InvalidValue", func(t *testing.T) {
		os.Setenv("TEST_INVALID_INT", "not_a_number")
		result := GetEnvAs("TEST_INVALID_INT", 99)
		if result != 99 {
			t.Errorf("GetEnvAs[int] with invalid value = %v, want 99", result)
		}
	})

	t.Run("GetEnvAs_String", func(t *testing.T) {
		os.Setenv("TEST_STRING", "hello")
		result := GetEnvAs("TEST_STRING", "default")
		if result != "hello" {
			t.Errorf("GetEnvAs[string] = %v, want 'hello'", result)
		}
	})

	t.Run("GetEnvAs_String_Missing", func(t *testing.T) {
		os.Unsetenv("TEST_STRING")
		result := GetEnvAs("TEST_STRING", "default")
		if result != "default" {
			t.Errorf("GetEnvAs[string] with missing var = %v, want 'default'", result)
		}
	})

	t.Run("GetEnvAs_Bool_True", func(t *testing.T) {
		os.Setenv("TEST_BOOL", "true")
		result := GetEnvAs("TEST_BOOL", false)
		if result != true {
			t.Errorf("GetEnvAs[bool] = %v, want true", result)
		}
	})

	t.Run("GetEnvAs_Bool_False", func(t *testing.T) {
		os.Setenv("TEST_BOOL", "false")
		result := GetEnvAs("TEST_BOOL", true)
		if result != false {
			t.Errorf("GetEnvAs[bool] = %v, want false", result)
		}
	})

	t.Run("GetEnvAs_Bool_NumericTrue", func(t *testing.T) {
		os.Setenv("TEST_BOOL", "1")
		result := GetEnvAs("TEST_BOOL", false)
		if result != true {
			t.Errorf("GetEnvAs[bool] with '1' = %v, want true", result)
		}
	})

	t.Run("GetEnvAs_Bool_NumericFalse", func(t *testing.T) {
		os.Setenv("TEST_BOOL", "0")
		result := GetEnvAs("TEST_BOOL", true)
		if result != false {
			t.Errorf("GetEnvAs[bool] with '0' = %v, want false", result)
		}
	})

	t.Run("GetEnvAs_Float64", func(t *testing.T) {
		os.Setenv("TEST_FLOAT", "3.14")
		result := GetEnvAs("TEST_FLOAT", 0.0)
		if result != 3.14 {
			t.Errorf("GetEnvAs[float64] = %v, want 3.14", result)
		}
	})

	t.Run("GetEnvAs_Uint", func(t *testing.T) {
		os.Setenv("TEST_UINT", "255")
		result := GetEnvAs("TEST_UINT", uint(0))
		if result != uint(255) {
			t.Errorf("GetEnvAs[uint] = %v, want 255", result)
		}
	})

	t.Run("GetEnvAs_SliceInt", func(t *testing.T) {
		os.Setenv("TEST_SLICE_INT", "1,2,3,4,5")
		result := GetEnvAs("TEST_SLICE_INT", []int{})
		expected := []int{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("GetEnvAs[[]int] = %v, want %v", result, expected)
		}
	})

	t.Run("GetEnvAs_SliceInt_Missing", func(t *testing.T) {
		os.Unsetenv("TEST_SLICE_INT")
		result := GetEnvAs("TEST_SLICE_INT", []int{99, 100})
		expected := []int{99, 100}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("GetEnvAs[[]int] with missing var = %v, want %v", result, expected)
		}
	})
}

// TestGetEnvDuration tests the GetEnvDuration function
func TestGetEnvDuration(t *testing.T) {
	// Clean up environment variables after test
	defer func() {
		os.Unsetenv("TEST_DURATION")
		os.Unsetenv("TEST_DURATION_MS")
		os.Unsetenv("TEST_DURATION_COMPLEX")
		os.Unsetenv("TEST_DURATION_INVALID")
	}()

	t.Run("GetEnvDuration_Seconds", func(t *testing.T) {
		os.Setenv("TEST_DURATION", "30s")
		result := GetEnvDuration("TEST_DURATION", 10*time.Second)
		if result != 30*time.Second {
			t.Errorf("GetEnvDuration = %v, want 30s", result)
		}
	})

	t.Run("GetEnvDuration_Milliseconds", func(t *testing.T) {
		os.Setenv("TEST_DURATION_MS", "500ms")
		result := GetEnvDuration("TEST_DURATION_MS", 100*time.Millisecond)
		if result != 500*time.Millisecond {
			t.Errorf("GetEnvDuration = %v, want 500ms", result)
		}
	})

	t.Run("GetEnvDuration_Complex", func(t *testing.T) {
		os.Setenv("TEST_DURATION_COMPLEX", "2h45m30s")
		expected := 2*time.Hour + 45*time.Minute + 30*time.Second
		result := GetEnvDuration("TEST_DURATION_COMPLEX", 1*time.Hour)
		if result != expected {
			t.Errorf("GetEnvDuration = %v, want %v", result, expected)
		}
	})

	t.Run("GetEnvDuration_Minutes", func(t *testing.T) {
		os.Setenv("TEST_DURATION", "15m")
		result := GetEnvDuration("TEST_DURATION", 5*time.Minute)
		if result != 15*time.Minute {
			t.Errorf("GetEnvDuration = %v, want 15m", result)
		}
	})

	t.Run("GetEnvDuration_Hours", func(t *testing.T) {
		os.Setenv("TEST_DURATION", "24h")
		result := GetEnvDuration("TEST_DURATION", 12*time.Hour)
		if result != 24*time.Hour {
			t.Errorf("GetEnvDuration = %v, want 24h", result)
		}
	})

	t.Run("GetEnvDuration_Missing", func(t *testing.T) {
		os.Unsetenv("TEST_DURATION")
		fallback := 42 * time.Second
		result := GetEnvDuration("TEST_DURATION", fallback)
		if result != fallback {
			t.Errorf("GetEnvDuration with missing var = %v, want %v", result, fallback)
		}
	})

	t.Run("GetEnvDuration_InvalidFormat", func(t *testing.T) {
		os.Setenv("TEST_DURATION_INVALID", "invalid_duration")
		fallback := 10 * time.Second
		result := GetEnvDuration("TEST_DURATION_INVALID", fallback)
		if result != fallback {
			t.Errorf("GetEnvDuration with invalid format = %v, want %v", result, fallback)
		}
	})

	t.Run("GetEnvDuration_EmptyString", func(t *testing.T) {
		os.Setenv("TEST_DURATION", "")
		fallback := 5 * time.Second
		result := GetEnvDuration("TEST_DURATION", fallback)
		if result != fallback {
			t.Errorf("GetEnvDuration with empty string = %v, want %v", result, fallback)
		}
	})

	t.Run("GetEnvDuration_Nanoseconds", func(t *testing.T) {
		os.Setenv("TEST_DURATION", "1000ns")
		result := GetEnvDuration("TEST_DURATION", 1*time.Nanosecond)
		if result != 1000*time.Nanosecond {
			t.Errorf("GetEnvDuration = %v, want 1000ns", result)
		}
	})

	t.Run("GetEnvDuration_Microseconds", func(t *testing.T) {
		os.Setenv("TEST_DURATION", "100us")
		result := GetEnvDuration("TEST_DURATION", 1*time.Microsecond)
		if result != 100*time.Microsecond {
			t.Errorf("GetEnvDuration = %v, want 100us", result)
		}
	})
}

// TestGetEnvAs_EdgeCases tests edge cases for GetEnvAs
func TestGetEnvAs_EdgeCases(t *testing.T) {
	defer func() {
		os.Unsetenv("TEST_EDGE")
	}()

	t.Run("GetEnvAs_Int8_Overflow", func(t *testing.T) {
		os.Setenv("TEST_EDGE", "500") // Exceeds int8 max (127)
		result := GetEnvAs("TEST_EDGE", int8(10))
		// Should return fallback due to overflow error
		if result != int8(10) {
			t.Errorf("GetEnvAs[int8] with overflow = %v, want 10", result)
		}
	})

	t.Run("GetEnvAs_EmptyString", func(t *testing.T) {
		os.Setenv("TEST_EDGE", "")
		result := GetEnvAs("TEST_EDGE", "default")
		if result != "" {
			t.Errorf("GetEnvAs[string] with empty value = %v, want empty string", result)
		}
	})

	t.Run("GetEnvAs_ZeroValue", func(t *testing.T) {
		os.Setenv("TEST_EDGE", "0")
		result := GetEnvAs("TEST_EDGE", 99)
		if result != 0 {
			t.Errorf("GetEnvAs[int] with zero = %v, want 0", result)
		}
	})

	t.Run("GetEnvAs_NegativeInt", func(t *testing.T) {
		os.Setenv("TEST_EDGE", "-42")
		result := GetEnvAs("TEST_EDGE", 0)
		if result != -42 {
			t.Errorf("GetEnvAs[int] with negative = %v, want -42", result)
		}
	})

	t.Run("GetEnvAs_NegativeFloat", func(t *testing.T) {
		os.Setenv("TEST_EDGE", "-3.14")
		result := GetEnvAs("TEST_EDGE", 0.0)
		if result != -3.14 {
			t.Errorf("GetEnvAs[float64] with negative = %v, want -3.14", result)
		}
	})
}

// TestSetEnvPrefix tests the SetEnvPrefix function
func TestSetEnvPrefix(t *testing.T) {
	// Save original prefix and restore after test
	originalPrefix := GetEnvPrefix()
	defer SetEnvPrefix(originalPrefix)

	t.Run("SetEnvPrefix_Basic", func(t *testing.T) {
		SetEnvPrefix("MYAPP")
		result := GetEnvPrefix()
		if result != "MYAPP" {
			t.Errorf("GetEnvPrefix() = %v, want 'MYAPP'", result)
		}
	})

	t.Run("SetEnvPrefix_Change", func(t *testing.T) {
		SetEnvPrefix("FIRST")
		if GetEnvPrefix() != "FIRST" {
			t.Errorf("GetEnvPrefix() = %v, want 'FIRST'", GetEnvPrefix())
		}

		SetEnvPrefix("SECOND")
		if GetEnvPrefix() != "SECOND" {
			t.Errorf("GetEnvPrefix() = %v, want 'SECOND'", GetEnvPrefix())
		}
	})

	t.Run("SetEnvPrefix_Empty", func(t *testing.T) {
		SetEnvPrefix("TEMP")
		SetEnvPrefix("")
		result := GetEnvPrefix()
		if result != "" {
			t.Errorf("GetEnvPrefix() after setting empty = %v, want ''", result)
		}
	})
}

// TestGetEnvPrefix tests the GetEnvPrefix function
func TestGetEnvPrefix(t *testing.T) {
	// Save original prefix and restore after test
	originalPrefix := GetEnvPrefix()
	defer SetEnvPrefix(originalPrefix)

	t.Run("GetEnvPrefix_Default", func(t *testing.T) {
		SetEnvPrefix("")
		result := GetEnvPrefix()
		if result != "" {
			t.Errorf("GetEnvPrefix() with no prefix = %v, want ''", result)
		}
	})

	t.Run("GetEnvPrefix_AfterSet", func(t *testing.T) {
		SetEnvPrefix("TESTPREFIX")
		result := GetEnvPrefix()
		if result != "TESTPREFIX" {
			t.Errorf("GetEnvPrefix() = %v, want 'TESTPREFIX'", result)
		}
	})
}

// TestGetPrefixedEnv tests the getPrefixedEnv internal function through public APIs
func TestGetPrefixedEnv(t *testing.T) {
	// Save original prefix and restore after test
	originalPrefix := GetEnvPrefix()
	defer func() {
		SetEnvPrefix(originalPrefix)
		os.Unsetenv("TEST_VAR")
		os.Unsetenv("MYAPP_TEST_VAR")
	}()

	t.Run("GetPrefixedEnv_NoPrefix", func(t *testing.T) {
		SetEnvPrefix("")
		os.Setenv("TEST_VAR", "value1")
		result := GetEnv("TEST_VAR", "default")
		if result != "value1" {
			t.Errorf("GetEnv without prefix = %v, want 'value1'", result)
		}
	})

	t.Run("GetPrefixedEnv_WithPrefix", func(t *testing.T) {
		SetEnvPrefix("MYAPP")
		os.Unsetenv("TEST_VAR")
		os.Setenv("MYAPP_TEST_VAR", "value2")
		result := GetEnv("TEST_VAR", "default")
		if result != "value2" {
			t.Errorf("GetEnv with prefix = %v, want 'value2'", result)
		}
	})

	t.Run("GetPrefixedEnv_PrefixNotFound", func(t *testing.T) {
		SetEnvPrefix("MYAPP")
		os.Unsetenv("TEST_VAR")
		os.Unsetenv("MYAPP_TEST_VAR")
		result := GetEnv("TEST_VAR", "fallback")
		if result != "fallback" {
			t.Errorf("GetEnv with prefix not found = %v, want 'fallback'", result)
		}
	})
}

// TestPrefixIntegration tests that all env functions respect the prefix
func TestPrefixIntegration(t *testing.T) {
	// Save original prefix and restore after test
	originalPrefix := GetEnvPrefix()
	defer func() {
		SetEnvPrefix(originalPrefix)
		os.Unsetenv("APP_PORT")
		os.Unsetenv("APP_DEBUG")
		os.Unsetenv("APP_TIMEOUT")
		os.Unsetenv("APP_MAX_CONN")
	}()

	SetEnvPrefix("APP")

	t.Run("Integration_GetEnv", func(t *testing.T) {
		os.Setenv("APP_PORT", "9000")
		result := GetEnv("PORT", "8080")
		if result != "9000" {
			t.Errorf("GetEnv with prefix = %v, want '9000'", result)
		}
	})

	t.Run("Integration_GetEnvAs", func(t *testing.T) {
		os.Setenv("APP_MAX_CONN", "100")
		result := GetEnvAs("MAX_CONN", 50)
		if result != 100 {
			t.Errorf("GetEnvAs with prefix = %v, want 100", result)
		}
	})

	t.Run("Integration_GetEnvAsBool", func(t *testing.T) {
		os.Setenv("APP_DEBUG", "true")
		result := GetEnvAs("DEBUG", false)
		if result != true {
			t.Errorf("GetEnvAs[bool] with prefix = %v, want true", result)
		}
	})

	t.Run("Integration_GetEnvDuration", func(t *testing.T) {
		os.Setenv("APP_TIMEOUT", "45s")
		result := GetEnvDuration("TIMEOUT", 30*time.Second)
		if result != 45*time.Second {
			t.Errorf("GetEnvDuration with prefix = %v, want 45s", result)
		}
	})
}

// TestPrefixWithFromEnvs tests that FromEnvs respects the prefix
func TestPrefixWithFromEnvs(t *testing.T) {
	// Save original prefix and restore after test
	originalPrefix := GetEnvPrefix()
	defer func() {
		SetEnvPrefix(originalPrefix)
		os.Unsetenv("MYSERVICE_HOST")
		os.Unsetenv("MYSERVICE_PORT")
		os.Unsetenv("MYSERVICE_ENABLED")
	}()

	type Config struct {
		Host    string `env:"HOST"`
		Port    int    `env:"PORT"`
		Enabled bool   `env:"ENABLED"`
	}

	SetEnvPrefix("MYSERVICE")
	os.Setenv("MYSERVICE_HOST", "localhost")
	os.Setenv("MYSERVICE_PORT", "3000")
	os.Setenv("MYSERVICE_ENABLED", "true")

	config, err := FromEnvs[Config]()
	if err != nil {
		t.Fatalf("FromEnvs with prefix failed: %v", err)
	}

	if config.Host != "localhost" {
		t.Errorf("Config.Host = %v, want 'localhost'", config.Host)
	}

	if config.Port != 3000 {
		t.Errorf("Config.Port = %v, want 3000", config.Port)
	}

	if !config.Enabled {
		t.Errorf("Config.Enabled = %v, want true", config.Enabled)
	}
}

// TestFromEnvs_NestedStructs tests mapping of nested structs from environment variables
func TestFromEnvs_NestedStructs(t *testing.T) {
	defer func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("APP_VERSION")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_SSL")
	}()

	type DatabaseConfig struct {
		Host     string `env:"HOST"`
		Port     int    `env:"PORT"`
		User     string `env:"USER"`
		Password string `env:"PASSWORD"`
		SSL      bool   `env:"SSL"`
	}

	type AppConfig struct {
		Name     string         `env:"APP_NAME"`
		Version  string         `env:"APP_VERSION"`
		Database DatabaseConfig `env:"DB"`
	}

	// Set environment variables
	os.Setenv("APP_NAME", "TestApp")
	os.Setenv("APP_VERSION", "1.0.0")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "admin")
	os.Setenv("DB_PASSWORD", "secret")
	os.Setenv("DB_SSL", "true")

	config, err := FromEnvs[AppConfig]()
	if err != nil {
		t.Fatalf("FromEnvs failed: %v", err)
	}

	// Test top-level fields
	if config.Name != "TestApp" {
		t.Errorf("AppConfig.Name = %v, want 'TestApp'", config.Name)
	}

	if config.Version != "1.0.0" {
		t.Errorf("AppConfig.Version = %v, want '1.0.0'", config.Version)
	}

	// Test nested struct fields
	if config.Database.Host != "localhost" {
		t.Errorf("DatabaseConfig.Host = %v, want 'localhost'", config.Database.Host)
	}

	if config.Database.Port != 5432 {
		t.Errorf("DatabaseConfig.Port = %v, want 5432", config.Database.Port)
	}

	if config.Database.User != "admin" {
		t.Errorf("DatabaseConfig.User = %v, want 'admin'", config.Database.User)
	}

	if config.Database.Password != "secret" {
		t.Errorf("DatabaseConfig.Password = %v, want 'secret'", config.Database.Password)
	}

	if !config.Database.SSL {
		t.Errorf("DatabaseConfig.SSL = %v, want true", config.Database.SSL)
	}
}

// TestFromEnvs_NestedStructPointers tests mapping of nested struct pointers from environment variables
func TestFromEnvs_NestedStructPointers(t *testing.T) {
	defer func() {
		os.Unsetenv("HOST")
		os.Unsetenv("PORT")
		os.Unsetenv("CACHE_HOST")
		os.Unsetenv("CACHE_PORT")
		os.Unsetenv("CACHE_TTL")
	}()

	type CacheConfig struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT"`
		TTL  int    `env:"TTL"`
	}

	type ServerConfig struct {
		Host  string       `env:"HOST"`
		Port  int          `env:"PORT"`
		Cache *CacheConfig `env:"CACHE"`
	}

	// Set environment variables
	os.Setenv("HOST", "0.0.0.0")
	os.Setenv("PORT", "8080")
	os.Setenv("CACHE_HOST", "redis.local")
	os.Setenv("CACHE_PORT", "6379")
	os.Setenv("CACHE_TTL", "3600")

	config, err := FromEnvs[ServerConfig]()
	if err != nil {
		t.Fatalf("FromEnvs failed: %v", err)
	}

	// Test top-level fields
	if config.Host != "0.0.0.0" {
		t.Errorf("ServerConfig.Host = %v, want '0.0.0.0'", config.Host)
	}

	if config.Port != 8080 {
		t.Errorf("ServerConfig.Port = %v, want 8080", config.Port)
	}

	// Test nested pointer to struct
	if config.Cache == nil {
		t.Fatal("ServerConfig.Cache is nil, expected initialized struct")
	}

	if config.Cache.Host != "redis.local" {
		t.Errorf("CacheConfig.Host = %v, want 'redis.local'", config.Cache.Host)
	}

	if config.Cache.Port != 6379 {
		t.Errorf("CacheConfig.Port = %v, want 6379", config.Cache.Port)
	}

	if config.Cache.TTL != 3600 {
		t.Errorf("CacheConfig.TTL = %v, want 3600", config.Cache.TTL)
	}
}

// TestFromEnvs_MultiLevelNesting tests multiple levels of nested structs
func TestFromEnvs_MultiLevelNesting(t *testing.T) {
	defer func() {
		os.Unsetenv("NAME")
		os.Unsetenv("DB_PRIMARY_HOST")
		os.Unsetenv("DB_PRIMARY_PORT")
		os.Unsetenv("DB_REPLICA_HOST")
		os.Unsetenv("DB_REPLICA_PORT")
	}()

	type Connection struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT"`
	}

	type DatabaseCluster struct {
		Primary Connection `env:"PRIMARY"`
		Replica Connection `env:"REPLICA"`
	}

	type Application struct {
		Name     string          `env:"NAME"`
		Database DatabaseCluster `env:"DB"`
	}

	// Set environment variables
	os.Setenv("NAME", "MultiLevelApp")
	os.Setenv("DB_PRIMARY_HOST", "primary.db.local")
	os.Setenv("DB_PRIMARY_PORT", "5432")
	os.Setenv("DB_REPLICA_HOST", "replica.db.local")
	os.Setenv("DB_REPLICA_PORT", "5433")

	config, err := FromEnvs[Application]()
	if err != nil {
		t.Fatalf("FromEnvs failed: %v", err)
	}

	if config.Name != "MultiLevelApp" {
		t.Errorf("Application.Name = %v, want 'MultiLevelApp'", config.Name)
	}

	if config.Database.Primary.Host != "primary.db.local" {
		t.Errorf("Primary.Host = %v, want 'primary.db.local'", config.Database.Primary.Host)
	}

	if config.Database.Primary.Port != 5432 {
		t.Errorf("Primary.Port = %v, want 5432", config.Database.Primary.Port)
	}

	if config.Database.Replica.Host != "replica.db.local" {
		t.Errorf("Replica.Host = %v, want 'replica.db.local'", config.Database.Replica.Host)
	}

	if config.Database.Replica.Port != 5433 {
		t.Errorf("Replica.Port = %v, want 5433", config.Database.Replica.Port)
	}
}

// TestFromEnvs_MixedNestedTypes tests a mix of nested structs, pointers, and primitive types
func TestFromEnvs_MixedNestedTypes(t *testing.T) {
	defer func() {
		os.Unsetenv("NAME")
		os.Unsetenv("PORT")
		os.Unsetenv("DEBUG")
		os.Unsetenv("AUTH_ENABLED")
		os.Unsetenv("AUTH_TOKEN")
		os.Unsetenv("LIMITS_MAX_CONN")
		os.Unsetenv("LIMITS_TIMEOUT")
	}()

	type AuthConfig struct {
		Enabled bool   `env:"ENABLED"`
		Token   string `env:"TOKEN"`
	}

	type LimitsConfig struct {
		MaxConnections int     `env:"MAX_CONN"`
		Timeout        float64 `env:"TIMEOUT"`
	}

	type ServiceConfig struct {
		Name   string       `env:"NAME"`
		Port   int          `env:"PORT"`
		Debug  bool         `env:"DEBUG"`
		Auth   *AuthConfig  `env:"AUTH"`
		Limits LimitsConfig `env:"LIMITS"`
	}

	// Set environment variables
	os.Setenv("NAME", "API")
	os.Setenv("PORT", "3000")
	os.Setenv("DEBUG", "true")
	os.Setenv("AUTH_ENABLED", "true")
	os.Setenv("AUTH_TOKEN", "secret-token-123")
	os.Setenv("LIMITS_MAX_CONN", "100")
	os.Setenv("LIMITS_TIMEOUT", "30.5")

	config, err := FromEnvs[ServiceConfig]()
	if err != nil {
		t.Fatalf("FromEnvs failed: %v", err)
	}

	if config.Name != "API" {
		t.Errorf("ServiceConfig.Name = %v, want 'API'", config.Name)
	}

	if config.Port != 3000 {
		t.Errorf("ServiceConfig.Port = %v, want 3000", config.Port)
	}

	if !config.Debug {
		t.Errorf("ServiceConfig.Debug = %v, want true", config.Debug)
	}

	if config.Auth == nil {
		t.Fatal("ServiceConfig.Auth is nil")
	}

	if !config.Auth.Enabled {
		t.Errorf("AuthConfig.Enabled = %v, want true", config.Auth.Enabled)
	}

	if config.Auth.Token != "secret-token-123" {
		t.Errorf("AuthConfig.Token = %v, want 'secret-token-123'", config.Auth.Token)
	}

	if config.Limits.MaxConnections != 100 {
		t.Errorf("LimitsConfig.MaxConnections = %v, want 100", config.Limits.MaxConnections)
	}

	if config.Limits.Timeout != 30.5 {
		t.Errorf("LimitsConfig.Timeout = %v, want 30.5", config.Limits.Timeout)
	}
}

// TestFromEnvs_NestedStructWithoutTags tests nested structs without env tags
func TestFromEnvs_NestedStructWithoutTags(t *testing.T) {
	defer func() {
		os.Unsetenv("ENABLED")
		os.Unsetenv("TIMEOUT")
	}()

	type Settings struct {
		Enabled bool `env:"ENABLED"`
		Timeout int  `env:"TIMEOUT"`
	}

	type Config struct {
		Settings Settings // No env tag - fields should map directly
	}

	os.Setenv("ENABLED", "true")
	os.Setenv("TIMEOUT", "60")

	config, err := FromEnvs[Config]()
	if err != nil {
		t.Fatalf("FromEnvs failed: %v", err)
	}

	if !config.Settings.Enabled {
		t.Errorf("Settings.Enabled = %v, want true", config.Settings.Enabled)
	}

	if config.Settings.Timeout != 60 {
		t.Errorf("Settings.Timeout = %v, want 60", config.Settings.Timeout)
	}
}

// TestFromEnvs_NestedStructWithSlices tests nested structs containing slices
func TestFromEnvs_NestedStructWithSlices(t *testing.T) {
	defer func() {
		os.Unsetenv("API_BASE_URL")
		os.Unsetenv("API_ENDPOINTS")
		os.Unsetenv("API_ALLOWED_ORIGINS")
	}()

	type APIConfig struct {
		BaseURL        string   `env:"BASE_URL"`
		Endpoints      []string `env:"ENDPOINTS"`
		AllowedOrigins []string `env:"ALLOWED_ORIGINS"`
	}

	type Config struct {
		API APIConfig `env:"API"`
	}

	os.Setenv("API_BASE_URL", "https://api.example.com")
	os.Setenv("API_ENDPOINTS", "/users,/posts,/comments")
	os.Setenv("API_ALLOWED_ORIGINS", "http://localhost:3000,https://example.com")

	config, err := FromEnvs[Config]()
	if err != nil {
		t.Fatalf("FromEnvs failed: %v", err)
	}

	if config.API.BaseURL != "https://api.example.com" {
		t.Errorf("APIConfig.BaseURL = %v, want 'https://api.example.com'", config.API.BaseURL)
	}

	expectedEndpoints := []string{"/users", "/posts", "/comments"}
	if !reflect.DeepEqual(config.API.Endpoints, expectedEndpoints) {
		t.Errorf("APIConfig.Endpoints = %v, want %v", config.API.Endpoints, expectedEndpoints)
	}

	expectedOrigins := []string{"http://localhost:3000", "https://example.com"}
	if !reflect.DeepEqual(config.API.AllowedOrigins, expectedOrigins) {
		t.Errorf("APIConfig.AllowedOrigins = %v, want %v", config.API.AllowedOrigins, expectedOrigins)
	}
}

// TestFromEnvs_NestedStructPartialMapping tests when only some nested fields have values
func TestFromEnvs_NestedStructPartialMapping(t *testing.T) {
	defer func() {
		os.Unsetenv("NAME")
		os.Unsetenv("DB_HOST")
	}()

	type Database struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT"` // This won't be set
	}

	type Config struct {
		Name     string   `env:"NAME"`
		Database Database `env:"DB"`
	}

	os.Setenv("NAME", "PartialApp")
	os.Setenv("DB_HOST", "localhost")
	// DB_PORT is not set

	config, err := FromEnvs[Config]()
	if err != nil {
		t.Fatalf("FromEnvs failed: %v", err)
	}

	if config.Name != "PartialApp" {
		t.Errorf("Config.Name = %v, want 'PartialApp'", config.Name)
	}

	if config.Database.Host != "localhost" {
		t.Errorf("Database.Host = %v, want 'localhost'", config.Database.Host)
	}

	if config.Database.Port != 0 {
		t.Errorf("Database.Port = %v, want 0 (zero value)", config.Database.Port)
	}
}
