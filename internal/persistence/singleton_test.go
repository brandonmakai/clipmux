package persistence

import (
	"path/filepath"
	"sync"
	"testing"

	"github.com/brandonmakai/clipmux/internal/logger"
)

// TestGetHistory_SingletonProperties tests the core behaviors of the GetHistory singleton.
func TestGetHistory_SingletonProperties(t *testing.T) {
	// This test relies on the global state of the singleton.
	// To test both chronological and recent-first, the test process would need to be run twice
	// with different parameters, as the first call determines the instance type for all subsequent tests.

	t.Run("returns the same instance on subsequent calls", func(t *testing.T) {
		// Create a temporary directory for the log file.
		tempDir := t.TempDir()
		logFile := filepath.Join(tempDir, "test.log")
		log := logger.GetLogger(logFile, false)

		instance1 := GetHistory(false, 10, 1024, log)
		instance2 := GetHistory(false, 10, 1024, log)

		// In Go, we can compare the pointers to see if they are the same instance.
		if instance1 != instance2 {
			t.Fatal("GetHistory returned different instances on subsequent calls")
		}
	})

	t.Run("is thread-safe and initializes only once", func(t *testing.T) {
		// To test concurrency, we need to reset the singleton state.
		// This is a common pattern in testing singletons.
		once = sync.Once{}
		instance = nil

		// Create a temporary directory for the log file.
		tempDir := t.TempDir()
		logFile := filepath.Join(tempDir, "test_concurrent.log")
		log := logger.GetLogger(logFile, false)

		var wg sync.WaitGroup
		const numGoroutines = 100

		instances := make(chan ClipboardHistory, numGoroutines)

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				// All goroutines will call GetHistory concurrently.
				// sync.Once should ensure only one of them creates the instance.
				instances <- GetHistory(true, 10, 1024, log)
			}()
		}

		wg.Wait()
		close(instances)

		// Check that all goroutines received the same instance.
		firstInstance := <-instances
		for inst := range instances {
			if inst != firstInstance {
				t.Fatal("Singleton initialization was not thread-safe; different instances were created")
			}
		}
	})
}
