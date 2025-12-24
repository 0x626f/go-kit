package env

import (
	"os"
	"path/filepath"
	"testing"
)

// TestConfig is a test configuration struct used across multiple tests
type TestConfig struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	Debug    bool   `env:"DEBUG"`
	Timeout  int    `env:"TIMEOUT"`
	MaxConns int    `env:"MAX_CONNS"`
}

// SimpleConfig is a minimal test configuration
type SimpleConfig struct {
	Value string `env:"VALUE"`
}

// TestNewManifest tests that NewManifest creates a valid instance
func TestNewManifest(t *testing.T) {
	manifest := NewManifest[TestConfig]()
	if manifest == nil {
		t.Fatal("NewManifest returned nil")
	}
}

// TestNewManifest_DifferentTypes tests that NewManifest works with different struct types
func TestNewManifest_DifferentTypes(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "TestConfig",
			test: func(t *testing.T) {
				m := NewManifest[TestConfig]()
				if m == nil {
					t.Fatal("NewManifest[TestConfig] returned nil")
				}
			},
		},
		{
			name: "SimpleConfig",
			test: func(t *testing.T) {
				m := NewManifest[SimpleConfig]()
				if m == nil {
					t.Fatal("NewManifest[SimpleConfig] returned nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// TestWithProperty tests setting individual properties
func TestWithProperty(t *testing.T) {
	// Clean up environment after test
	defer func() {
		os.Unsetenv("TEST_KEY")
		os.Unsetenv("ANOTHER_KEY")
	}()

	manifest := NewManifest[TestConfig]().
		WithProperty("TEST_KEY", "test_value").
		WithProperty("ANOTHER_KEY", "another_value")

	if manifest == nil {
		t.Fatal("WithProperty returned nil")
	}

	// Verify the environment variables were set
	if val := os.Getenv("TEST_KEY"); val != "test_value" {
		t.Errorf("Expected TEST_KEY=test_value, got %s", val)
	}
	if val := os.Getenv("ANOTHER_KEY"); val != "another_value" {
		t.Errorf("Expected ANOTHER_KEY=another_value, got %s", val)
	}
}

// TestWithProperty_Chaining tests that WithProperty supports method chaining
func TestWithProperty_Chaining(t *testing.T) {
	defer func() {
		os.Unsetenv("KEY1")
		os.Unsetenv("KEY2")
		os.Unsetenv("KEY3")
	}()

	manifest := NewManifest[TestConfig]().
		WithProperty("KEY1", "value1").
		WithProperty("KEY2", "value2").
		WithProperty("KEY3", "value3")

	if manifest == nil {
		t.Fatal("Method chaining failed")
	}

	vals := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
		"KEY3": "value3",
	}

	for key, expected := range vals {
		if val := os.Getenv(key); val != expected {
			t.Errorf("Expected %s=%s, got %s", key, expected, val)
		}
	}
}

// TestWithAppName tests setting an app name prefix
func TestWithAppName(t *testing.T) {
	// Save original prefix and restore after all tests
	originalPrefix := GetEnvPrefix()
	defer SetEnvPrefix(originalPrefix)

	// Clean up environment after all tests
	defer func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("TEST_APP_NAME")
	}()

	tests := []struct {
		name           string
		envVar         string
		envValue       string
		expectedPrefix string
	}{
		{
			name:           "Set prefix from environment variable",
			envVar:         "APP_NAME",
			envValue:       "myapp",
			expectedPrefix: "myapp",
		},
		{
			name:           "Set prefix with different variable",
			envVar:         "TEST_APP_NAME",
			envValue:       "testservice",
			expectedPrefix: "testservice",
		},
		{
			name:           "Empty app name",
			envVar:         "NONEXISTENT",
			envValue:       "",
			expectedPrefix: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment before each subtest
			os.Unsetenv("APP_NAME")
			os.Unsetenv("TEST_APP_NAME")
			SetEnvPrefix("") // Reset prefix before each subtest

			// Set up environment
			if tt.envValue != "" {
				os.Setenv(tt.envVar, tt.envValue)
			}

			manifest := NewManifest[TestConfig]().WithPrefix(tt.envVar)

			if manifest == nil {
				t.Fatal("WithPrefix returned nil")
			}

			if prefix := GetEnvPrefix(); prefix != tt.expectedPrefix {
				t.Errorf("Expected prefix=%s, got %s", tt.expectedPrefix, prefix)
			}
		})
	}
}

// TestWithSource tests loading environment variables from a file
func TestWithSource(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		verify      func(*testing.T)
		shouldPanic bool
	}{
		{
			name: "Valid env file",
			content: `HOST=localhost
PORT=8080
DEBUG=true`,
			verify: func(t *testing.T) {
				if val := os.Getenv("HOST"); val != "localhost" {
					t.Errorf("Expected HOST=localhost, got %s", val)
				}
				if val := os.Getenv("PORT"); val != "8080" {
					t.Errorf("Expected PORT=8080, got %s", val)
				}
				if val := os.Getenv("DEBUG"); val != "true" {
					t.Errorf("Expected DEBUG=true, got %s", val)
				}
			},
			shouldPanic: false,
		},
		{
			name:    "Empty file",
			content: "",
			verify: func(t *testing.T) {
				// No verification needed for empty file
			},
			shouldPanic: false,
		},
		{
			name: "File with comments",
			content: `# This is a comment
HOST=testhost
# Another comment
PORT=9000`,
			verify: func(t *testing.T) {
				if val := os.Getenv("HOST"); val != "testhost" {
					t.Errorf("Expected HOST=testhost, got %s", val)
				}
				if val := os.Getenv("PORT"); val != "9000" {
					t.Errorf("Expected PORT=9000, got %s", val)
				}
			},
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment variables
			defer func() {
				os.Unsetenv("HOST")
				os.Unsetenv("PORT")
				os.Unsetenv("DEBUG")
			}()

			// Create test file
			testFile := filepath.Join(tempDir, "test.env")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic but didn't get one")
					}
				}()
			}

			manifest := NewManifest[TestConfig]().WithSource(testFile)

			if !tt.shouldPanic {
				if manifest == nil {
					t.Fatal("WithSource returned nil")
				}
				tt.verify(t)
			}
		})
	}
}

// TestWithSource_NonexistentFile tests that WithSource panics on missing file
func TestWithSource_NonexistentFile(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nonexistent file but didn't get one")
		}
	}()

	NewManifest[TestConfig]().WithSource("/nonexistent/path/to/file.env")
}

// TestWithRelativeSource tests loading a file relative to an environment variable path
func TestWithRelativeSource(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Create a reference file and the env file in the same directory
	refFile := filepath.Join(configDir, "main.conf")
	envFile := filepath.Join(configDir, "secrets.env")

	if err := os.WriteFile(refFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create reference file: %v", err)
	}

	envContent := `SECRET_KEY=my-secret-key
API_TOKEN=abc123`
	if err := os.WriteFile(envFile, []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to create env file: %v", err)
	}

	// Clean up environment
	defer func() {
		os.Unsetenv("CONFIG_PATH")
		os.Unsetenv("SECRET_KEY")
		os.Unsetenv("API_TOKEN")
	}()

	// Set the reference path
	os.Setenv("CONFIG_PATH", refFile)

	manifest := NewManifest[TestConfig]().WithRelativeSource("CONFIG_PATH", "secrets.env")

	if manifest == nil {
		t.Fatal("WithRelativeSource returned nil")
	}

	// Verify the environment variables were loaded
	if val := os.Getenv("SECRET_KEY"); val != "my-secret-key" {
		t.Errorf("Expected SECRET_KEY=my-secret-key, got %s", val)
	}
	if val := os.Getenv("API_TOKEN"); val != "abc123" {
		t.Errorf("Expected API_TOKEN=abc123, got %s", val)
	}
}

// TestWithRelativeSource_NonexistentFile tests panic on missing relative file
func TestWithRelativeSource_NonexistentFile(t *testing.T) {
	defer func() {
		os.Unsetenv("CONFIG_PATH")
		if r := recover(); r == nil {
			t.Error("Expected panic for nonexistent relative file but didn't get one")
		}
	}()

	os.Setenv("CONFIG_PATH", "/some/path/config.conf")
	NewManifest[TestConfig]().WithRelativeSource("CONFIG_PATH", "nonexistent.env")
}

// TestLoad tests the Load method with various configurations
func TestLoad(t *testing.T) {
	// Clean up environment
	defer func() {
		os.Unsetenv("HOST")
		os.Unsetenv("PORT")
		os.Unsetenv("DEBUG")
		os.Unsetenv("TIMEOUT")
		os.Unsetenv("MAX_CONNS")
	}()

	tests := []struct {
		name   string
		setup  func()
		verify func(*testing.T, *TestConfig, error)
	}{
		{
			name: "Load with all fields set",
			setup: func() {
				os.Setenv("HOST", "example.com")
				os.Setenv("PORT", "3000")
				os.Setenv("DEBUG", "true")
				os.Setenv("TIMEOUT", "30")
				os.Setenv("MAX_CONNS", "100")
			},
			verify: func(t *testing.T, config *TestConfig, err error) {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if config.Host != "example.com" {
					t.Errorf("Expected Host=example.com, got %s", config.Host)
				}
				if config.Port != 3000 {
					t.Errorf("Expected Port=3000, got %d", config.Port)
				}
				if config.Debug != true {
					t.Errorf("Expected Debug=true, got %v", config.Debug)
				}
				if config.Timeout != 30 {
					t.Errorf("Expected Timeout=30, got %d", config.Timeout)
				}
				if config.MaxConns != 100 {
					t.Errorf("Expected MaxConns=100, got %d", config.MaxConns)
				}
			},
		},
		{
			name: "Load with partial fields",
			setup: func() {
				os.Setenv("HOST", "localhost")
				os.Setenv("PORT", "8080")
			},
			verify: func(t *testing.T, config *TestConfig, err error) {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if config.Host != "localhost" {
					t.Errorf("Expected Host=localhost, got %s", config.Host)
				}
				if config.Port != 8080 {
					t.Errorf("Expected Port=8080, got %d", config.Port)
				}
				// Other fields should be zero values
				if config.Debug != false {
					t.Errorf("Expected Debug=false, got %v", config.Debug)
				}
			},
		},
		{
			name: "Load with no environment variables",
			setup: func() {
				// Don't set any environment variables
			},
			verify: func(t *testing.T, config *TestConfig, err error) {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				// All fields should be zero values
				if config.Host != "" {
					t.Errorf("Expected empty Host, got %s", config.Host)
				}
				if config.Port != 0 {
					t.Errorf("Expected Port=0, got %d", config.Port)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean environment before each test
			os.Unsetenv("HOST")
			os.Unsetenv("PORT")
			os.Unsetenv("DEBUG")
			os.Unsetenv("TIMEOUT")
			os.Unsetenv("MAX_CONNS")

			tt.setup()

			config, err := NewManifest[TestConfig]().Load()
			tt.verify(t, config, err)
		})
	}
}

// TestLoad_NonStructType tests that Load returns an error for non-struct types
func TestLoad_NonStructType(t *testing.T) {
	_, err := NewManifest[string]().Load()
	if err == nil {
		t.Error("Expected error for non-struct type, got nil")
	}
}

// TestManifest_FullWorkflow tests a complete workflow with multiple methods
func TestManifest_FullWorkflow(t *testing.T) {
	// Save original prefix
	originalPrefix := GetEnvPrefix()
	defer SetEnvPrefix(originalPrefix)

	// Create temporary directory and files
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")
	envContent := `HOST=production.example.com
PORT=443
DEBUG=false`
	if err := os.WriteFile(envFile, []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to create env file: %v", err)
	}

	// Clean up environment
	defer func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("HOST")
		os.Unsetenv("PORT")
		os.Unsetenv("DEBUG")
		os.Unsetenv("TIMEOUT")
		os.Unsetenv("MAX_CONNS")
	}()

	// Set up app name
	os.Setenv("APP_NAME", "")

	// Execute full workflow
	config, err := NewManifest[TestConfig]().
		WithPrefix("APP_NAME").
		WithSource(envFile).
		WithProperty("TIMEOUT", "60").
		WithProperty("MAX_CONNS", "200").
		Load()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify configuration
	if config.Host != "production.example.com" {
		t.Errorf("Expected Host=production.example.com, got %s", config.Host)
	}
	if config.Port != 443 {
		t.Errorf("Expected Port=443, got %d", config.Port)
	}
	if config.Debug != false {
		t.Errorf("Expected Debug=false, got %v", config.Debug)
	}
	if config.Timeout != 60 {
		t.Errorf("Expected Timeout=60, got %d", config.Timeout)
	}
	if config.MaxConns != 200 {
		t.Errorf("Expected MaxConns=200, got %d", config.MaxConns)
	}
}

// TestManifest_WithPrefix tests that prefix is applied correctly
func TestManifest_WithPrefix(t *testing.T) {
	// Save original prefix
	originalPrefix := GetEnvPrefix()
	defer SetEnvPrefix(originalPrefix)

	// Clean up environment
	defer func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("myapp_HOST")
		os.Unsetenv("myapp_PORT")
	}()

	// Set environment variables with prefix
	os.Setenv("APP_NAME", "myapp")
	os.Setenv("myapp_HOST", "prefixed.example.com")
	os.Setenv("myapp_PORT", "9000")

	config, err := NewManifest[TestConfig]().
		WithPrefix("APP_NAME").
		Load()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if config.Host != "prefixed.example.com" {
		t.Errorf("Expected Host=prefixed.example.com, got %s", config.Host)
	}
	if config.Port != 9000 {
		t.Errorf("Expected Port=9000, got %d", config.Port)
	}
}

// TestManifest_MultipleSourceFiles tests loading from multiple source files
func TestManifest_MultipleSourceFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create first env file
	envFile1 := filepath.Join(tempDir, ".env.base")
	content1 := `HOST=base.example.com
PORT=8000`
	if err := os.WriteFile(envFile1, []byte(content1), 0644); err != nil {
		t.Fatalf("Failed to create first env file: %v", err)
	}

	// Create second env file (overrides)
	envFile2 := filepath.Join(tempDir, ".env.override")
	content2 := `PORT=9000
DEBUG=true`
	if err := os.WriteFile(envFile2, []byte(content2), 0644); err != nil {
		t.Fatalf("Failed to create second env file: %v", err)
	}

	// Clean up environment
	defer func() {
		os.Unsetenv("HOST")
		os.Unsetenv("PORT")
		os.Unsetenv("DEBUG")
	}()

	config, err := NewManifest[TestConfig]().
		WithSource(envFile1).
		WithSource(envFile2). // Should override PORT
		Load()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if config.Host != "base.example.com" {
		t.Errorf("Expected Host=base.example.com, got %s", config.Host)
	}
	// PORT should be overridden by second file
	if config.Port != 9000 {
		t.Errorf("Expected Port=9000 (overridden), got %d", config.Port)
	}
	if config.Debug != true {
		t.Errorf("Expected Debug=true, got %v", config.Debug)
	}
}

// BenchmarkNewManifest benchmarks the creation of a new manifest
func BenchmarkNewManifest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewManifest[TestConfig]()
	}
}

// BenchmarkWithProperty benchmarks setting a property
func BenchmarkWithProperty(b *testing.B) {
	defer os.Unsetenv("BENCH_KEY")

	for i := 0; i < b.N; i++ {
		_ = NewManifest[TestConfig]().WithProperty("BENCH_KEY", "value")
	}
}

// BenchmarkLoad benchmarks the Load operation
func BenchmarkLoad(b *testing.B) {
	defer func() {
		os.Unsetenv("HOST")
		os.Unsetenv("PORT")
		os.Unsetenv("DEBUG")
	}()

	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "8080")
	os.Setenv("DEBUG", "true")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewManifest[TestConfig]().Load()
	}
}

// BenchmarkFullWorkflow benchmarks a complete manifest workflow
func BenchmarkFullWorkflow(b *testing.B) {
	tempDir := b.TempDir()
	envFile := filepath.Join(tempDir, ".env")
	envContent := `HOST=localhost
PORT=8080`
	if err := os.WriteFile(envFile, []byte(envContent), 0644); err != nil {
		b.Fatalf("Failed to create env file: %v", err)
	}

	defer func() {
		os.Unsetenv("HOST")
		os.Unsetenv("PORT")
		os.Unsetenv("TIMEOUT")
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewManifest[TestConfig]().
			WithSource(envFile).
			WithProperty("TIMEOUT", "30").
			Load()
	}
}
