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

// TestDefaultTag_BasicTypes tests the default tag with various basic types
func TestDefaultTag_BasicTypes(t *testing.T) {
	defer func() {
		os.Unsetenv("STRING_VAR")
		os.Unsetenv("INT_VAR")
		os.Unsetenv("BOOL_VAR")
		os.Unsetenv("FLOAT_VAR")
	}()

	type Config struct {
		StringWithDefault string  `env:"STRING_VAR" default:"default_string"`
		IntWithDefault    int     `env:"INT_VAR" default:"42"`
		BoolWithDefault   bool    `env:"BOOL_VAR" default:"true"`
		FloatWithDefault  float64 `env:"FLOAT_VAR" default:"3.14"`
	}

	t.Run("UseDefaultWhenEnvNotSet", func(t *testing.T) {
		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.StringWithDefault != "default_string" {
			t.Errorf("StringWithDefault = %v, want 'default_string'", config.StringWithDefault)
		}

		if config.IntWithDefault != 42 {
			t.Errorf("IntWithDefault = %v, want 42", config.IntWithDefault)
		}

		if !config.BoolWithDefault {
			t.Errorf("BoolWithDefault = %v, want true", config.BoolWithDefault)
		}

		if config.FloatWithDefault != 3.14 {
			t.Errorf("FloatWithDefault = %v, want 3.14", config.FloatWithDefault)
		}
	})

	t.Run("EnvOverridesDefault", func(t *testing.T) {
		os.Setenv("STRING_VAR", "env_string")
		os.Setenv("INT_VAR", "100")
		os.Setenv("BOOL_VAR", "false")
		os.Setenv("FLOAT_VAR", "2.71")

		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.StringWithDefault != "env_string" {
			t.Errorf("StringWithDefault = %v, want 'env_string'", config.StringWithDefault)
		}

		if config.IntWithDefault != 100 {
			t.Errorf("IntWithDefault = %v, want 100", config.IntWithDefault)
		}

		if config.BoolWithDefault {
			t.Errorf("BoolWithDefault = %v, want false", config.BoolWithDefault)
		}

		if config.FloatWithDefault != 2.71 {
			t.Errorf("FloatWithDefault = %v, want 2.71", config.FloatWithDefault)
		}
	})
}

// TestDefaultTag_NumericTypes tests default tag with various numeric types
func TestDefaultTag_NumericTypes(t *testing.T) {
	defer func() {
		os.Unsetenv("INT8_VAR")
		os.Unsetenv("INT16_VAR")
		os.Unsetenv("INT32_VAR")
		os.Unsetenv("INT64_VAR")
		os.Unsetenv("UINT_VAR")
		os.Unsetenv("UINT8_VAR")
		os.Unsetenv("UINT16_VAR")
		os.Unsetenv("UINT32_VAR")
		os.Unsetenv("UINT64_VAR")
		os.Unsetenv("FLOAT32_VAR")
	}()

	type Config struct {
		Int8Val    int8    `env:"INT8_VAR" default:"127"`
		Int16Val   int16   `env:"INT16_VAR" default:"32767"`
		Int32Val   int32   `env:"INT32_VAR" default:"2147483647"`
		Int64Val   int64   `env:"INT64_VAR" default:"9223372036854775807"`
		UintVal    uint    `env:"UINT_VAR" default:"42"`
		Uint8Val   uint8   `env:"UINT8_VAR" default:"255"`
		Uint16Val  uint16  `env:"UINT16_VAR" default:"65535"`
		Uint32Val  uint32  `env:"UINT32_VAR" default:"4294967295"`
		Uint64Val  uint64  `env:"UINT64_VAR" default:"18446744073709551615"`
		Float32Val float32 `env:"FLOAT32_VAR" default:"3.14159"`
	}

	config, err := FromEnvs[Config]()
	if err != nil {
		t.Fatalf("FromEnvs failed: %v", err)
	}

	t.Run("Int8Default", func(t *testing.T) {
		if config.Int8Val != 127 {
			t.Errorf("Int8Val = %v, want 127", config.Int8Val)
		}
	})

	t.Run("Int16Default", func(t *testing.T) {
		if config.Int16Val != 32767 {
			t.Errorf("Int16Val = %v, want 32767", config.Int16Val)
		}
	})

	t.Run("Int32Default", func(t *testing.T) {
		if config.Int32Val != 2147483647 {
			t.Errorf("Int32Val = %v, want 2147483647", config.Int32Val)
		}
	})

	t.Run("Int64Default", func(t *testing.T) {
		if config.Int64Val != 9223372036854775807 {
			t.Errorf("Int64Val = %v, want 9223372036854775807", config.Int64Val)
		}
	})

	t.Run("UintDefault", func(t *testing.T) {
		if config.UintVal != 42 {
			t.Errorf("UintVal = %v, want 42", config.UintVal)
		}
	})

	t.Run("Uint8Default", func(t *testing.T) {
		if config.Uint8Val != 255 {
			t.Errorf("Uint8Val = %v, want 255", config.Uint8Val)
		}
	})

	t.Run("Uint16Default", func(t *testing.T) {
		if config.Uint16Val != 65535 {
			t.Errorf("Uint16Val = %v, want 65535", config.Uint16Val)
		}
	})

	t.Run("Uint32Default", func(t *testing.T) {
		if config.Uint32Val != 4294967295 {
			t.Errorf("Uint32Val = %v, want 4294967295", config.Uint32Val)
		}
	})

	t.Run("Uint64Default", func(t *testing.T) {
		if config.Uint64Val != 18446744073709551615 {
			t.Errorf("Uint64Val = %v, want 18446744073709551615", config.Uint64Val)
		}
	})

	t.Run("Float32Default", func(t *testing.T) {
		if config.Float32Val != 3.14159 {
			t.Errorf("Float32Val = %v, want 3.14159", config.Float32Val)
		}
	})
}

// TestDefaultTag_Slices tests default tag with slice types
func TestDefaultTag_Slices(t *testing.T) {
	defer func() {
		os.Unsetenv("INT_SLICE")
		os.Unsetenv("STRING_SLICE")
		os.Unsetenv("BOOL_SLICE")
	}()

	type Config struct {
		IntSlice    []int    `env:"INT_SLICE" default:"1,2,3,4,5"`
		StringSlice []string `env:"STRING_SLICE" default:"one,two,three"`
		BoolSlice   []bool   `env:"BOOL_SLICE" default:"true,false,true"`
	}

	t.Run("UseDefaultSlices", func(t *testing.T) {
		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		expectedIntSlice := []int{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(config.IntSlice, expectedIntSlice) {
			t.Errorf("IntSlice = %v, want %v", config.IntSlice, expectedIntSlice)
		}

		expectedStringSlice := []string{"one", "two", "three"}
		if !reflect.DeepEqual(config.StringSlice, expectedStringSlice) {
			t.Errorf("StringSlice = %v, want %v", config.StringSlice, expectedStringSlice)
		}

		expectedBoolSlice := []bool{true, false, true}
		if !reflect.DeepEqual(config.BoolSlice, expectedBoolSlice) {
			t.Errorf("BoolSlice = %v, want %v", config.BoolSlice, expectedBoolSlice)
		}
	})

	t.Run("EnvOverridesDefaultSlices", func(t *testing.T) {
		os.Setenv("INT_SLICE", "10,20,30")
		os.Setenv("STRING_SLICE", "alpha,beta")
		os.Setenv("BOOL_SLICE", "false,false")

		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		expectedIntSlice := []int{10, 20, 30}
		if !reflect.DeepEqual(config.IntSlice, expectedIntSlice) {
			t.Errorf("IntSlice = %v, want %v", config.IntSlice, expectedIntSlice)
		}

		expectedStringSlice := []string{"alpha", "beta"}
		if !reflect.DeepEqual(config.StringSlice, expectedStringSlice) {
			t.Errorf("StringSlice = %v, want %v", config.StringSlice, expectedStringSlice)
		}

		expectedBoolSlice := []bool{false, false}
		if !reflect.DeepEqual(config.BoolSlice, expectedBoolSlice) {
			t.Errorf("BoolSlice = %v, want %v", config.BoolSlice, expectedBoolSlice)
		}
	})
}

// TestDefaultTag_EmptyString tests default tag with empty string values
func TestDefaultTag_EmptyString(t *testing.T) {
	defer func() {
		os.Unsetenv("EMPTY_VAR")
	}()

	type Config struct {
		EmptyDefault string `env:"EMPTY_VAR" default:""`
		NoDefault    string `env:"NO_DEFAULT_VAR"`
	}

	config, err := FromEnvs[Config]()
	if err != nil {
		t.Fatalf("FromEnvs failed: %v", err)
	}

	t.Run("EmptyStringDefault", func(t *testing.T) {
		if config.EmptyDefault != "" {
			t.Errorf("EmptyDefault = %v, want empty string", config.EmptyDefault)
		}
	})

	t.Run("NoDefaultStaysZero", func(t *testing.T) {
		if config.NoDefault != "" {
			t.Errorf("NoDefault = %v, want empty string (zero value)", config.NoDefault)
		}
	})
}

// TestDefaultTag_NestedStructs tests default tag with nested structs
func TestDefaultTag_NestedStructs(t *testing.T) {
	// Clean up before and after test
	os.Unsetenv("NAME")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")

	defer func() {
		os.Unsetenv("NAME")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
	}()

	type DatabaseConfig struct {
		Host string `env:"HOST" default:"localhost"`
		Port int    `env:"PORT" default:"5432"`
		User string `env:"USER" default:"postgres"`
	}

	type AppConfig struct {
		Name     string         `env:"NAME" default:"MyApp"`
		Database DatabaseConfig `env:"DB"`
	}

	t.Run("UseNestedDefaults", func(t *testing.T) {
		config, err := FromEnvs[AppConfig]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.Name != "MyApp" {
			t.Errorf("Name = %v, want 'MyApp'", config.Name)
		}

		if config.Database.Host != "localhost" {
			t.Errorf("Database.Host = %v, want 'localhost'", config.Database.Host)
		}

		if config.Database.Port != 5432 {
			t.Errorf("Database.Port = %v, want 5432", config.Database.Port)
		}

		if config.Database.User != "postgres" {
			t.Errorf("Database.User = %v, want 'postgres'", config.Database.User)
		}
	})

	t.Run("EnvOverridesNestedDefaults", func(t *testing.T) {
		os.Setenv("NAME", "CustomApp")
		os.Setenv("DB_HOST", "db.example.com")
		os.Setenv("DB_PORT", "3306")
		os.Setenv("DB_USER", "admin")

		config, err := FromEnvs[AppConfig]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.Name != "CustomApp" {
			t.Errorf("Name = %v, want 'CustomApp'", config.Name)
		}

		if config.Database.Host != "db.example.com" {
			t.Errorf("Database.Host = %v, want 'db.example.com'", config.Database.Host)
		}

		if config.Database.Port != 3306 {
			t.Errorf("Database.Port = %v, want 3306", config.Database.Port)
		}

		if config.Database.User != "admin" {
			t.Errorf("Database.User = %v, want 'admin'", config.Database.User)
		}
	})
}

// TestDefaultTag_WithPrefix tests default tag with environment variable prefixing
func TestDefaultTag_WithPrefix(t *testing.T) {
	originalPrefix := GetEnvPrefix()
	defer func() {
		SetEnvPrefix(originalPrefix)
		os.Unsetenv("MYAPP_PORT")
		os.Unsetenv("MYAPP_DEBUG")
	}()

	type Config struct {
		Port  int  `env:"PORT" default:"8080"`
		Debug bool `env:"DEBUG" default:"false"`
	}

	t.Run("DefaultWithPrefix", func(t *testing.T) {
		SetEnvPrefix("MYAPP")

		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.Port != 8080 {
			t.Errorf("Port = %v, want 8080", config.Port)
		}

		if config.Debug {
			t.Errorf("Debug = %v, want false", config.Debug)
		}
	})

	t.Run("PrefixedEnvOverridesDefault", func(t *testing.T) {
		SetEnvPrefix("MYAPP")
		os.Setenv("MYAPP_PORT", "9000")
		os.Setenv("MYAPP_DEBUG", "true")

		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.Port != 9000 {
			t.Errorf("Port = %v, want 9000", config.Port)
		}

		if !config.Debug {
			t.Errorf("Debug = %v, want true", config.Debug)
		}
	})
}

// TestDefaultTag_MixedConfiguration tests a mix of fields with and without defaults
func TestDefaultTag_MixedConfiguration(t *testing.T) {
	defer func() {
		os.Unsetenv("REQUIRED_VAR")
		os.Unsetenv("OPTIONAL_VAR")
	}()

	type Config struct {
		RequiredField string `env:"REQUIRED_VAR"`
		OptionalField string `env:"OPTIONAL_VAR" default:"optional_default"`
		DefaultField  int    `env:"DEFAULT_VAR" default:"100"`
	}

	t.Run("MixedWithoutEnv", func(t *testing.T) {
		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		// Required field not set, should be zero value
		if config.RequiredField != "" {
			t.Errorf("RequiredField = %v, want empty string (zero value)", config.RequiredField)
		}

		// Optional field should use default
		if config.OptionalField != "optional_default" {
			t.Errorf("OptionalField = %v, want 'optional_default'", config.OptionalField)
		}

		// Default field should use default
		if config.DefaultField != 100 {
			t.Errorf("DefaultField = %v, want 100", config.DefaultField)
		}
	})

	t.Run("MixedWithEnv", func(t *testing.T) {
		os.Setenv("REQUIRED_VAR", "required_value")
		os.Setenv("OPTIONAL_VAR", "env_value")

		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.RequiredField != "required_value" {
			t.Errorf("RequiredField = %v, want 'required_value'", config.RequiredField)
		}

		if config.OptionalField != "env_value" {
			t.Errorf("OptionalField = %v, want 'env_value'", config.OptionalField)
		}

		// Default field still uses default since env not set
		if config.DefaultField != 100 {
			t.Errorf("DefaultField = %v, want 100", config.DefaultField)
		}
	})
}

// TestDefaultTag_InvalidDefaultValue tests handling of invalid default values
func TestDefaultTag_InvalidDefaultValue(t *testing.T) {
	type Config struct {
		InvalidInt   int     `env:"INVALID_INT" default:"not_a_number"`
		InvalidBool  bool    `env:"INVALID_BOOL" default:"not_a_bool"`
		InvalidFloat float64 `env:"INVALID_FLOAT" default:"not_a_float"`
	}

	t.Run("InvalidDefaultsReturnError", func(t *testing.T) {
		_, err := FromEnvs[Config]()
		if err == nil {
			t.Error("FromEnvs should return error for invalid default values")
		}
	})
}

// TestDefaultTag_NestedStructPointers tests default tag with nested struct pointers
func TestDefaultTag_NestedStructPointers(t *testing.T) {
	defer func() {
		os.Unsetenv("CACHE_HOST")
		os.Unsetenv("CACHE_PORT")
	}()

	type CacheConfig struct {
		Host string `env:"HOST" default:"localhost"`
		Port int    `env:"PORT" default:"6379"`
	}

	type Config struct {
		Cache *CacheConfig `env:"CACHE"`
	}

	t.Run("DefaultsInPointerStruct", func(t *testing.T) {
		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.Cache == nil {
			t.Fatal("Cache should not be nil")
		}

		if config.Cache.Host != "localhost" {
			t.Errorf("Cache.Host = %v, want 'localhost'", config.Cache.Host)
		}

		if config.Cache.Port != 6379 {
			t.Errorf("Cache.Port = %v, want 6379", config.Cache.Port)
		}
	})

	t.Run("EnvOverridesDefaultsInPointerStruct", func(t *testing.T) {
		os.Setenv("CACHE_HOST", "redis.example.com")
		os.Setenv("CACHE_PORT", "6380")

		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.Cache == nil {
			t.Fatal("Cache should not be nil")
		}

		if config.Cache.Host != "redis.example.com" {
			t.Errorf("Cache.Host = %v, want 'redis.example.com'", config.Cache.Host)
		}

		if config.Cache.Port != 6380 {
			t.Errorf("Cache.Port = %v, want 6380", config.Cache.Port)
		}
	})
}

// TestDefaultTag_MultiLevelNesting tests default tag with multiple levels of nesting
func TestDefaultTag_MultiLevelNesting(t *testing.T) {
	defer func() {
		os.Unsetenv("NAME")
		os.Unsetenv("DB_PRIMARY_HOST")
		os.Unsetenv("DB_PRIMARY_PORT")
		os.Unsetenv("DB_REPLICA_HOST")
		os.Unsetenv("DB_REPLICA_PORT")
	}()

	type Connection struct {
		Host string `env:"HOST" default:"localhost"`
		Port int    `env:"PORT" default:"5432"`
	}

	type DatabaseCluster struct {
		Primary Connection `env:"PRIMARY"`
		Replica Connection `env:"REPLICA"`
	}

	type Application struct {
		Name     string          `env:"NAME" default:"DefaultApp"`
		Database DatabaseCluster `env:"DB"`
	}

	t.Run("MultiLevelDefaults", func(t *testing.T) {
		config, err := FromEnvs[Application]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.Name != "DefaultApp" {
			t.Errorf("Name = %v, want 'DefaultApp'", config.Name)
		}

		if config.Database.Primary.Host != "localhost" {
			t.Errorf("Database.Primary.Host = %v, want 'localhost'", config.Database.Primary.Host)
		}

		if config.Database.Primary.Port != 5432 {
			t.Errorf("Database.Primary.Port = %v, want 5432", config.Database.Primary.Port)
		}

		if config.Database.Replica.Host != "localhost" {
			t.Errorf("Database.Replica.Host = %v, want 'localhost'", config.Database.Replica.Host)
		}

		if config.Database.Replica.Port != 5432 {
			t.Errorf("Database.Replica.Port = %v, want 5432", config.Database.Replica.Port)
		}
	})

	t.Run("PartialEnvWithDefaults", func(t *testing.T) {
		os.Setenv("NAME", "ProductionApp")
		os.Setenv("DB_PRIMARY_HOST", "primary.db.local")
		// PRIMARY_PORT uses default
		// REPLICA_HOST and REPLICA_PORT use defaults

		config, err := FromEnvs[Application]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.Name != "ProductionApp" {
			t.Errorf("Name = %v, want 'ProductionApp'", config.Name)
		}

		if config.Database.Primary.Host != "primary.db.local" {
			t.Errorf("Database.Primary.Host = %v, want 'primary.db.local'", config.Database.Primary.Host)
		}

		if config.Database.Primary.Port != 5432 {
			t.Errorf("Database.Primary.Port = %v, want 5432 (default)", config.Database.Primary.Port)
		}

		if config.Database.Replica.Host != "localhost" {
			t.Errorf("Database.Replica.Host = %v, want 'localhost' (default)", config.Database.Replica.Host)
		}

		if config.Database.Replica.Port != 5432 {
			t.Errorf("Database.Replica.Port = %v, want 5432 (default)", config.Database.Replica.Port)
		}
	})
}

// TestDefaultTag_ZeroValues tests that zero values can be set as defaults
func TestDefaultTag_ZeroValues(t *testing.T) {
	defer func() {
		os.Unsetenv("ZERO_INT")
		os.Unsetenv("ZERO_BOOL")
		os.Unsetenv("ZERO_FLOAT")
	}()

	type Config struct {
		ZeroInt   int     `env:"ZERO_INT" default:"0"`
		ZeroBool  bool    `env:"ZERO_BOOL" default:"false"`
		ZeroFloat float64 `env:"ZERO_FLOAT" default:"0.0"`
	}

	config, err := FromEnvs[Config]()
	if err != nil {
		t.Fatalf("FromEnvs failed: %v", err)
	}

	t.Run("ZeroIntDefault", func(t *testing.T) {
		if config.ZeroInt != 0 {
			t.Errorf("ZeroInt = %v, want 0", config.ZeroInt)
		}
	})

	t.Run("ZeroBoolDefault", func(t *testing.T) {
		if config.ZeroBool {
			t.Errorf("ZeroBool = %v, want false", config.ZeroBool)
		}
	})

	t.Run("ZeroFloatDefault", func(t *testing.T) {
		if config.ZeroFloat != 0.0 {
			t.Errorf("ZeroFloat = %v, want 0.0", config.ZeroFloat)
		}
	})
}

// TestTimeDuration tests mapping of time.Duration type from environment variables
func TestTimeDuration(t *testing.T) {
	defer func() {
		os.Unsetenv("TIMEOUT")
		os.Unsetenv("RETRY_DELAY")
		os.Unsetenv("MAX_DURATION")
		os.Unsetenv("MIN_DURATION")
	}()

	type Config struct {
		Timeout     time.Duration `env:"TIMEOUT" default:"30s"`
		RetryDelay  time.Duration `env:"RETRY_DELAY" default:"5s"`
		MaxDuration time.Duration `env:"MAX_DURATION" default:"1h"`
		MinDuration time.Duration `env:"MIN_DURATION" default:"100ms"`
	}

	t.Run("UseDefaultDurations", func(t *testing.T) {
		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.Timeout != 30*time.Second {
			t.Errorf("Timeout = %v, want 30s", config.Timeout)
		}

		if config.RetryDelay != 5*time.Second {
			t.Errorf("RetryDelay = %v, want 5s", config.RetryDelay)
		}

		if config.MaxDuration != 1*time.Hour {
			t.Errorf("MaxDuration = %v, want 1h", config.MaxDuration)
		}

		if config.MinDuration != 100*time.Millisecond {
			t.Errorf("MinDuration = %v, want 100ms", config.MinDuration)
		}
	})

	t.Run("EnvOverridesDefaultDurations", func(t *testing.T) {
		os.Setenv("TIMEOUT", "45s")
		os.Setenv("RETRY_DELAY", "10s")
		os.Setenv("MAX_DURATION", "2h30m")
		os.Setenv("MIN_DURATION", "500ms")

		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.Timeout != 45*time.Second {
			t.Errorf("Timeout = %v, want 45s", config.Timeout)
		}

		if config.RetryDelay != 10*time.Second {
			t.Errorf("RetryDelay = %v, want 10s", config.RetryDelay)
		}

		if config.MaxDuration != 2*time.Hour+30*time.Minute {
			t.Errorf("MaxDuration = %v, want 2h30m", config.MaxDuration)
		}

		if config.MinDuration != 500*time.Millisecond {
			t.Errorf("MinDuration = %v, want 500ms", config.MinDuration)
		}
	})

	t.Run("VariousDurationFormats", func(t *testing.T) {
		testCases := []struct {
			name     string
			envValue string
			expected time.Duration
		}{
			{"Nanoseconds", "1000ns", 1000 * time.Nanosecond},
			{"Microseconds", "500us", 500 * time.Microsecond},
			{"Milliseconds", "250ms", 250 * time.Millisecond},
			{"Seconds", "60s", 60 * time.Second},
			{"Minutes", "5m", 5 * time.Minute},
			{"Hours", "24h", 24 * time.Hour},
			{"Complex", "1h30m45s", 1*time.Hour + 30*time.Minute + 45*time.Second},
			{"MultipleUnits", "2h15m30s500ms", 2*time.Hour + 15*time.Minute + 30*time.Second + 500*time.Millisecond},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				os.Setenv("TIMEOUT", tc.envValue)

				config, err := FromEnvs[Config]()
				if err != nil {
					t.Fatalf("FromEnvs failed: %v", err)
				}

				if config.Timeout != tc.expected {
					t.Errorf("Timeout = %v, want %v", config.Timeout, tc.expected)
				}
			})
		}
	})

	t.Run("InvalidDurationFormat", func(t *testing.T) {
		os.Setenv("TIMEOUT", "invalid-duration")

		_, err := FromEnvs[Config]()
		if err == nil {
			t.Error("Expected error for invalid duration format, got nil")
		}
	})

	t.Run("ZeroDuration", func(t *testing.T) {
		os.Setenv("TIMEOUT", "0s")

		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.Timeout != 0 {
			t.Errorf("Timeout = %v, want 0s", config.Timeout)
		}
	})

	t.Run("NegativeDuration", func(t *testing.T) {
		os.Setenv("TIMEOUT", "-10s")

		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.Timeout != -10*time.Second {
			t.Errorf("Timeout = %v, want -10s", config.Timeout)
		}
	})
}

// TestTimeDuration_WithPrefix tests time.Duration with environment variable prefix
func TestTimeDuration_WithPrefix(t *testing.T) {
	originalPrefix := GetEnvPrefix()
	defer func() {
		SetEnvPrefix(originalPrefix)
		os.Unsetenv("APP_TIMEOUT")
		os.Unsetenv("APP_RETRY")
	}()

	type Config struct {
		Timeout time.Duration `env:"TIMEOUT" default:"30s"`
		Retry   time.Duration `env:"RETRY" default:"5s"`
	}

	t.Run("DefaultWithPrefix", func(t *testing.T) {
		SetEnvPrefix("APP")

		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.Timeout != 30*time.Second {
			t.Errorf("Timeout = %v, want 30s", config.Timeout)
		}

		if config.Retry != 5*time.Second {
			t.Errorf("Retry = %v, want 5s", config.Retry)
		}
	})

	t.Run("PrefixedEnvOverrides", func(t *testing.T) {
		SetEnvPrefix("APP")
		os.Setenv("APP_TIMEOUT", "1m")
		os.Setenv("APP_RETRY", "10s")

		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.Timeout != 1*time.Minute {
			t.Errorf("Timeout = %v, want 1m", config.Timeout)
		}

		if config.Retry != 10*time.Second {
			t.Errorf("Retry = %v, want 10s", config.Retry)
		}
	})
}

// TestTimeDuration_NestedStructs tests time.Duration in nested structures
func TestTimeDuration_NestedStructs(t *testing.T) {
	defer func() {
		os.Unsetenv("SERVER_TIMEOUT")
		os.Unsetenv("SERVER_KEEPALIVE")
		os.Unsetenv("DB_TIMEOUT")
		os.Unsetenv("DB_RETRY")
	}()

	type ServerConfig struct {
		Timeout   time.Duration `env:"TIMEOUT" default:"30s"`
		KeepAlive time.Duration `env:"KEEPALIVE" default:"60s"`
	}

	type DatabaseConfig struct {
		Timeout time.Duration `env:"TIMEOUT" default:"10s"`
		Retry   time.Duration `env:"RETRY" default:"5s"`
	}

	type Config struct {
		Server   ServerConfig   `env:"SERVER"`
		Database DatabaseConfig `env:"DB"`
	}

	t.Run("NestedDefaults", func(t *testing.T) {
		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.Server.Timeout != 30*time.Second {
			t.Errorf("Server.Timeout = %v, want 30s", config.Server.Timeout)
		}

		if config.Server.KeepAlive != 60*time.Second {
			t.Errorf("Server.KeepAlive = %v, want 60s", config.Server.KeepAlive)
		}

		if config.Database.Timeout != 10*time.Second {
			t.Errorf("Database.Timeout = %v, want 10s", config.Database.Timeout)
		}

		if config.Database.Retry != 5*time.Second {
			t.Errorf("Database.Retry = %v, want 5s", config.Database.Retry)
		}
	})

	t.Run("NestedEnvOverrides", func(t *testing.T) {
		os.Setenv("SERVER_TIMEOUT", "45s")
		os.Setenv("SERVER_KEEPALIVE", "90s")
		os.Setenv("DB_TIMEOUT", "20s")
		os.Setenv("DB_RETRY", "10s")

		config, err := FromEnvs[Config]()
		if err != nil {
			t.Fatalf("FromEnvs failed: %v", err)
		}

		if config.Server.Timeout != 45*time.Second {
			t.Errorf("Server.Timeout = %v, want 45s", config.Server.Timeout)
		}

		if config.Server.KeepAlive != 90*time.Second {
			t.Errorf("Server.KeepAlive = %v, want 90s", config.Server.KeepAlive)
		}

		if config.Database.Timeout != 20*time.Second {
			t.Errorf("Database.Timeout = %v, want 20s", config.Database.Timeout)
		}

		if config.Database.Retry != 10*time.Second {
			t.Errorf("Database.Retry = %v, want 10s", config.Database.Retry)
		}
	})
}
