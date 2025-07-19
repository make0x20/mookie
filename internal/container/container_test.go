// internal/container/container_test.go
package container

import (
	"testing"
)

func TestContainer_Register(t *testing.T) {
	c := New()
	service := "test service"

	c.Register("test", service)

	// Verify service was stored
	result, err := c.Get("test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != service {
		t.Errorf("got %v, want %v", result, service)
	}
}

func TestContainer_Get(t *testing.T) {
	tests := []struct {
		name        string
		service     any
		serviceName string
		wantErr     bool
	}{
		{
			name:        "existing service",
			service:     "test service",
			serviceName: "exists",
			wantErr:     false,
		},
		{
			name:        "non-existent service",
			service:     nil,
			serviceName: "doesnt-exist",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			if tt.service != nil {
				c.Register(tt.serviceName, tt.service)
			}

			result, err := c.Get(tt.serviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.service {
				t.Errorf("Get() got = %v, want %v", result, tt.service)
			}
		})
	}
}

func TestContainer_MustGet(t *testing.T) {
	t.Run("existing service", func(t *testing.T) {
		c := New()
		service := "test service"
		c.Register("test", service)

		// Should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustGet() panicked unexpectedly: %v", r)
			}
		}()

		result := c.MustGet("test")
		if result != service {
			t.Errorf("got %v, want %v", result, service)
		}
	})

	t.Run("non-existent service", func(t *testing.T) {
		c := New()

		// Should panic
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustGet() did not panic as expected")
			}
		}()

		c.MustGet("non-existent")
	})
}

func TestContainer_ConcurrentAccess(t *testing.T) {
	c := New()
	done := make(chan bool)

	// Test concurrent reads and writes
	go func() {
		for i := 0; i < 100; i++ {
			c.Register("test", i)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			c.Get("test")
		}
		done <- true
	}()

	// Wait for both goroutines to finish
	<-done
	<-done
}

func TestContainer_MultipleServices(t *testing.T) {
	c := New()

	// Register multiple services of different types
	c.Register("string", "string service")
	c.Register("int", 42)
	c.Register("bool", true)
	c.Register("struct", struct{ Name string }{"test"})

	// Verify each service
	if s, _ := c.Get("string"); s != "string service" {
		t.Errorf("got %v, want string service", s)
	}
	if i, _ := c.Get("int"); i != 42 {
		t.Errorf("got %v, want 42", i)
	}
	if b, _ := c.Get("bool"); b != true {
		t.Errorf("got %v, want true", b)
	}
	if s, _ := c.Get("struct"); s.(struct{ Name string }).Name != "test" {
		t.Errorf("got %v, want test", s)
	}
}
