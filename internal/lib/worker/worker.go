package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"
)

var ErrRetryable = errors.New("retryable error")

type Worker struct {
	stream   Stream
	store    Store
	mutex    sync.RWMutex
	handlers map[string]H
	ctx      context.Context
	cancel   context.CancelFunc
	interval time.Duration
	wg       sync.WaitGroup
}

// New creates a new worker instance
func New(store Store, stream Stream) *Worker {
	ctx, cancel := context.WithCancel(context.Background())

	return &Worker{
		stream:   stream,
		store:    store,
		ctx:      ctx,
		cancel:   cancel,
		handlers: make(map[string]H),
		interval: 10 * time.Minute,
	}
}

func (w *Worker) Register(handler H) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.handlers[handler.Name] = handler
}

func (w *Worker) Start() error {
	w.wg.Add(1)
	go w.poll()

	w.wg.Add(1)
	go w.work()

	return nil
}

func (w *Worker) Wait() {
	w.wg.Wait()
}

func (w *Worker) Stop() {
	w.cancel()
	slog.Info("Worker stopping, waiting for ongoing jobs to finish...")
	w.wg.Wait()
}

func (w *Worker) poll() {
	defer w.wg.Done()
	for {
		if w.ctx.Err() != nil {
			return
		}

		jobs, err := w.store.Fetch(w.ctx)
		if err != nil {
			slog.Error("Failed to fetch jobs", "error", err)
		}

		for _, job := range jobs {

			if w.ctx.Err() != nil {
				return
			}

			if err := w.stream.Publish(w.ctx, job); err != nil {
				slog.Error("Failed to publish job", "id", job.ID, "error", err)
			}
		}

		select {
		case <-w.ctx.Done():
			slog.Info("Worker poller stopped")
			return
		case <-time.After(w.interval):
		}
	}

}

func (w *Worker) work() {
	defer w.wg.Done()
	for message, err := range w.stream.Read(w.ctx) {

		if w.ctx.Err() != nil {
			return
		}

		if err != nil {
			slog.Error("Failed to read from stream", "error", err)
			time.Sleep(5 * time.Second)
			continue
		}

		job := message.Job

		w.mutex.RLock()
		handler, ok := w.handlers[job.Name]
		if !ok {
			slog.Warn("No handler registered for job", "jobName", job.Name)
			w.mutex.RUnlock()
			continue
		}

		w.mutex.RUnlock()

		err := w.execute(handler, job)
		if err != nil {
			slog.Error("Failed to execute job", "job_id", job.ID, "error", err)
			continue
		}

		if err := message.Ack(w.ctx); err != nil {
			slog.Error("Failed to acknowledge message", "error", err)
			continue
		}

	}

	slog.Info("Worker stopped processing jobs")
}

func (w *Worker) execute(handler H, job Job) error {
	var execution Execution
	execution.StartedAt = time.Now()

	job.Status = StatusRunning
	job.UpdatedAt = time.Now()
	if err := w.store.Update(w.ctx, job); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	slog.Debug("Starting job execution", "id", job.ID, "name", job.Name)

	err := w.handle(handler, job)
	switch {

	case err == nil:
		job.Status = StatusCompleted

		slog.Debug("Job executed successfully",
			"id", job.ID,
			"name", job.Name,
			"duration", time.Since(execution.StartedAt))

	case errors.Is(err, ErrRetryable) && job.Retries < handler.MaxRetries:
		job.Retries++
		job.Status = StatusPending
		execution.Error = err.Error()

		slog.Warn("Job execution failed, will retry",
			"id", job.ID,
			"retries", job.Retries,
			"error", err,
			"duration", time.Since(execution.StartedAt))

	default:
		job.Status = StatusFailed
		execution.Error = err.Error()

		slog.Error("Job execution failed",
			"id", job.ID,
			"error", err,
			"duration", time.Since(execution.StartedAt))
	}

	execution.FinishedAt = time.Now()
	execution.FinishedIn = execution.FinishedAt.Sub(execution.StartedAt).String()

	job.Executions = append(job.Executions, execution)
	job.UpdatedAt = time.Now()

	if err := w.store.Update(w.ctx, job); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	return nil

}

func (w *Worker) handle(handler H, job Job) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			slog.Error("Panic recovered", "error", r, "stack", string(stack))
			err = fmt.Errorf("panic: %v: %s", r, string(stack))
		}
	}()

	return handler.Handle(w.ctx, job.Data)
}
