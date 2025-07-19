package cron

import (
	"sync"
	"time"
)

/*
   Package cron provides a simple task scheduler that runs functions periodically.
   It supports multiple tasks, concurrent execution, and clean shutdown.

   How to use:
   1. Create a new Runner
   2. Add tasks to the Runner (functions that implement CronFunc)
   3. Start the Runner with a specified interval
   4. Stop the Runner when done

   Example basic usage:
       // Create and start runner
       runner := cron.NewRunner()
       runner.Add(func() error {
           fmt.Println("Task running...")
           return nil
       })
       go runner.Start(time.Minute)

   Example with dependencies:
       // Create task with database dependency
       func SaveMetrics(db *sql.DB) cron.CronFunc {
           return func() error {
               return db.Exec("INSERT INTO metrics...")
           }
       }

       // Use in application
       runner := cron.NewRunner()
       runner.Add(SaveMetrics(db))
       go runner.Start(time.Minute * 5)

       // Cleanup on shutdown
       defer runner.Stop()

   Example multiple tasks:
       runner := cron.NewRunner()

       // Add multiple tasks
       runner.Add(CleanupOldRecords(db))
       runner.Add(UpdateCache(cache))
       runner.Add(SendMetrics(metrics))

       // Run all tasks every 30 seconds
       go runner.Start(time.Second * 30)

   Notes:
   - Tasks run sequentially in the order they were added
   - All tasks share the same interval
   - Thread-safe
   - Supports graceful shutdown
   - Tasks should be idempotent
   - Error handling must be implemented in the task
   - Start() is blocking and should typically run in a goroutine
*/

// CronFunc is a function type that can be run on a schedule
type CronFunc func() error

// Runner runs tasks on a schedule
type Runner struct {
	tasks    []CronFunc
	stop     chan struct{}
	mu       sync.RWMutex
	stopOnce sync.Once
}

// NewRunner creates a new Runner
func NewRunner() *Runner {
	return &Runner{
		stop: make(chan struct{}),
	}
}

// Add adds a task to the Runner
func (r *Runner) Add(task CronFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks = append(r.tasks, task)
}

// Start starts the Runner and runs tasks on the specified interval
// Usually called in a goroutine for example: go runner.Start(time.Minute)
func (r *Runner) Start(runEvery time.Duration) {
	ticker := time.NewTicker(runEvery)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.mu.RLock()
			for _, task := range r.tasks {
				task()
			}
			r.mu.RUnlock()
		case <-r.stop:
			return
		}
	}
}

// Stop stops the Runner
func (r *Runner) Stop() {
	r.stopOnce.Do(func() {
		close(r.stop)
	})
}
