package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

type Storage struct {
	filePath string
}

func NewStorage(filePath string) *Storage {
	if filePath == "" {
		currDir, err := os.Getwd()
		if err != nil {
			currDir = "."
		}
		filePath = filepath.Join(currDir, ".tasks.csv")
	}

	return &Storage{
		filePath: filePath,
	}
}

// utility methods to open and close file safely
func (s *Storage) loadFile() (*os.File, error) {
	// opens the file in rw mode and creates if doesn't exist
	f, err := os.OpenFile(s.filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading: %w", err)
	}

	// inline error checking - common go practice
	// request an exclusive lock to prevent any other process from accessing the file
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close() // if locking fails close the file
		return nil, fmt.Errorf("failed to lock file: %w", err)
	}

	return f, nil
}

func (s *Storage) closeFile(f *os.File) error {
	// release the lock
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		_ = f.Close()
		return fmt.Errorf("failed to unlock file: %w", err)
	}
	return f.Close()
}

// read tasks
func (s *Storage) LoadTasks() ([]Task, error) {
	file, err := s.loadFile()
	if err != nil {
		return nil, err
	}

	// the defer keyword ensures this function runs when the current
	// function exits, ensuring the file is always closed
	defer s.closeFile(file)

	// get file info
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stats: %w", err)
	}

	// if file is empty, write header and return empty tasks
	if info.Size() == 0 {
		// create csv writer
		writer := csv.NewWriter(file)
		err = writer.Write([]string{"ID", "Description", "CreatedAt", "IsComplete"})
		if err != nil {
			return nil, fmt.Errorf("failed to write CSV header: %w", err)
		}
		writer.Flush()
		return []Task{}, nil
	}

	// if file has content, read it
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	// no tasks present in csv
	if len(records) <= 1 {
		return []Task{}, nil
	}

	var tasks []Task
	for _, record := range records[1:] {
		// skip this row if it is malformed
		if len(record) != 4 {
			continue
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			continue
		}

		createdAt, err := time.Parse(time.RFC3339, record[2])
		if err != nil {
			continue
		}

		isComplete, err := strconv.ParseBool(record[3])
		if err != nil {
			continue
		}

		task := Task{
			ID:          id,
			Description: record[1],
			CreatedAt:   createdAt,
			IsComplete:  isComplete,
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// save tasks
func (s *Storage) SaveTasks(tasks []Task) error {
	file, err := s.loadFile()
	if err != nil {
		return err
	}
	defer s.closeFile(file)

	// truncate the file and move file ptr to beginning before writing
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek to beginning of file: %w", err)
	}

	writer := csv.NewWriter(file)

	// write header
	if err := writer.Write([]string{"ID", "Description", "CreatedAt", "IsComplete"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// write tasks
	for _, task := range tasks {
		record := []string{
			strconv.Itoa(task.ID),
			task.Description,
			task.CreatedAt.Format(time.RFC3339),
			strconv.FormatBool(task.IsComplete),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write task record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing CSV writer: %w", err)
	}

	return nil
}

// create a task
func (s *Storage) AddTask(description string) error {
	tasks, err := s.LoadTasks()
	if err != nil {
		return err
	}

	nextID := 1
	for _, task := range tasks {
		if task.ID >= nextID {
			nextID = task.ID + 1
		}
	}

	newTask := NewTask(nextID, description)
	tasks = append(tasks, newTask)

	return s.SaveTasks(tasks)
}

// read all tasks (based on filter)
func (s *Storage) ListTasks(showAll bool) ([]Task, error) {
	tasks, err := s.LoadTasks()
	if err != nil {
		return nil, err
	}

	if showAll {
		return tasks, nil
	}

	// only list incomplete task
	var incompleteTasks []Task
	for _, task := range tasks {
		if !task.IsComplete {
			incompleteTasks = append(incompleteTasks, task)
		}
	}

	return incompleteTasks, nil
}

// mark a task as complete
func (s *Storage) CompleteTask(id int) error {
	tasks, err := s.LoadTasks()
	if err != nil {
		return err
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].IsComplete = true
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}

	return s.SaveTasks(tasks)
}

// delete a task
func (s *Storage) DeleteTask(id int) error {
	tasks, err := s.LoadTasks()
	if err != nil {
		return err
	}

	found := false
	var updatedTasks []Task
	for _, task := range tasks {
		if task.ID == id {
			found = true
			continue // skip this task to remove it
		}
		updatedTasks = append(updatedTasks, task)
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}

	return s.SaveTasks(updatedTasks)
}
