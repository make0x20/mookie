package cron

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunner_Add(t *testing.T) {
	runner := NewRunner()
	
	t.Run("add single task", func(t *testing.T) {
		var executed bool
		task := func() error {
			executed = true
			return nil
		}
		
		runner.Add(task)
		
		// Run once manually to verify task was added
		runner.tasks[0]()
		
		if !executed {
			t.Error("task was not executed")
		}
	})
	
	t.Run("add multiple tasks", func(t *testing.T) {
		runner := NewRunner()
		var count int32
		
		// Add 3 tasks
		for i := 0; i < 3; i++ {
			runner.Add(func() error {
				atomic.AddInt32(&count, 1)
				return nil
			})
		}
		
		// Run all tasks manually
		for _, task := range runner.tasks {
			task()
		}
		
		if atomic.LoadInt32(&count) != 3 {
			t.Errorf("expected 3 task executions, got %d", count)
		}
	})
}

func TestRunner_Start(t *testing.T) {
	t.Run("tasks execute on schedule", func(t *testing.T) {
		runner := NewRunner()
		var count int32
		
		runner.Add(func() error {
			atomic.AddInt32(&count, 1)
			return nil
		})
		
		// Start runner with 100ms interval
		go runner.Start(100 * time.Millisecond)
		
		// Wait for ~3 executions
		time.Sleep(350 * time.Millisecond)
		runner.Stop()
		
		execCount := atomic.LoadInt32(&count)
		if execCount < 2 || execCount > 4 { // Allow for some timing flexibility
			t.Errorf("expected ~3 executions, got %d", execCount)
		}
	})
	
	t.Run("multiple tasks execute in order", func(t *testing.T) {
		runner := NewRunner()
		var sequence []int
		var mu sync.Mutex
		
		// Add tasks that record their execution order
		for i := 0; i < 3; i++ {
			taskNum := i
			runner.Add(func() error {
				mu.Lock()
				sequence = append(sequence, taskNum)
				mu.Unlock()
				return nil
			})
		}
		
		go runner.Start(100 * time.Millisecond)
		time.Sleep(150 * time.Millisecond) // Wait for one execution
		runner.Stop()
		
		if len(sequence) != 3 {
			t.Errorf("expected 3 task executions, got %d", len(sequence))
		}
		
		// Verify execution order
		for i := 0; i < len(sequence); i++ {
			if sequence[i] != i {
				t.Errorf("tasks executed out of order, got %v", sequence)
				break
			}
		}
	})
}

func TestRunner_Stop(t *testing.T) {
	t.Run("stops execution", func(t *testing.T) {
		runner := NewRunner()
		var count int32
		
		runner.Add(func() error {
			atomic.AddInt32(&count, 1)
			return nil
		})
		
		go runner.Start(100 * time.Millisecond)
		time.Sleep(250 * time.Millisecond) // Allow some executions
		runner.Stop()
		
		// Record the count
		countAfterStop := atomic.LoadInt32(&count)
		time.Sleep(200 * time.Millisecond) // Wait to verify no more executions
		
		if atomic.LoadInt32(&count) != countAfterStop {
			t.Error("tasks continued to execute after stop")
		}
	})
	
	t.Run("multiple stops are safe", func(t *testing.T) {
		runner := NewRunner()
		
		go runner.Start(100 * time.Millisecond)
		time.Sleep(50 * time.Millisecond)
		
		// Multiple stops should not panic
		runner.Stop()
		runner.Stop()
	})
}

func TestRunner_Concurrent(t *testing.T) {
	t.Run("concurrent task addition", func(t *testing.T) {
		runner := NewRunner()
		var wg sync.WaitGroup
		
		// Add tasks concurrently
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				runner.Add(func() error { return nil })
			}()
		}
		
		wg.Wait()
		
		if len(runner.tasks) != 10 {
			t.Errorf("expected 10 tasks, got %d", len(runner.tasks))
		}
	})
	
	t.Run("concurrent start/stop", func(t *testing.T) {
		runner := NewRunner()
		runner.Add(func() error { return nil })
		
		var wg sync.WaitGroup
		// Start and stop concurrently multiple times
		for i := 0; i < 5; i++ {
			wg.Add(2)
			go func() {
				defer wg.Done()
				go runner.Start(50 * time.Millisecond)
			}()
			go func() {
				defer wg.Done()
				time.Sleep(10 * time.Millisecond)
				runner.Stop()
			}()
		}
		
		wg.Wait() // Should not deadlock
	})
}

func BenchmarkRunner(b *testing.B) {
	b.Run("task addition", func(b *testing.B) {
		runner := NewRunner()
		task := func() error { return nil }
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			runner.Add(task)
		}
	})
	
	b.Run("task execution", func(b *testing.B) {
		runner := NewRunner()
		var count int32
		
		runner.Add(func() error {
			atomic.AddInt32(&count, 1)
			return nil
		})
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			runner.tasks[0]()
		}
	})
}
