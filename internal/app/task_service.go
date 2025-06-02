package app

import (
	"errors"
	"log"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type TaskService interface {
	Save(t domain.Task) (domain.Task, error)
	Find(id uint64) (interface{}, error)
	FindAll(uId uint64, status *string, date *time.Time) ([]domain.Task, error) // Оновлений інтерфейс для пошуку завдань за фільтрами
	Update(t domain.Task) (domain.Task, error)
	UpdateStatus(taskID uint64, userID uint64, status domain.TaskStatus) (domain.Task, error)
	// метод для оновлення статусу завдання
	Delete(id uint64) error
}

type taskService struct {
	taskRepo database.TaskRepository
}

func NewTaskService(tr database.TaskRepository) TaskService {
	return taskService{
		taskRepo: tr,
	}
}

func (s taskService) Save(t domain.Task) (domain.Task, error) {
	task, err := s.taskRepo.Save(t)
	if err != nil {
		log.Printf("taskService.Save(s.taskRepo.Save): %s", err)
		return domain.Task{}, err
	}

	return task, nil
}

func (s taskService) Find(id uint64) (interface{}, error) {
	task, err := s.taskRepo.Find(id)
	if err != nil {
		log.Printf("taskService.Find(s.taskRepo.Find): %s", err)
		return domain.Task{}, err
	}

	return task, nil
}

// повертаємо всі завдання користувача з можливістю фільтрації
func (s taskService) FindAll(uId uint64, status *string, date *time.Time) ([]domain.Task, error) {
	tasks, err := s.taskRepo.FindAllTasks(uId)
	if err != nil {
		log.Printf("taskService.FindAll(s.taskRepo.FindAllTasks): %s", err)
		return nil, err
	}

	var filtered []domain.Task
	for _, task := range tasks {
		if status != nil && string(task.Status) != *status {
			continue
		}
		if date != nil && (task.Date == nil || !sameDay(*task.Date, *date)) {
			continue
		}
		filtered = append(filtered, task)
	}

	return filtered, nil
}

// ігноруємо однаковий час
func sameDay(d1, d2 time.Time) bool {
	y1, m1, day1 := d1.Date()
	y2, m2, day2 := d2.Date()
	return y1 == y2 && m1 == m2 && day1 == day2
}

func (s taskService) Update(t domain.Task) (domain.Task, error) {
	task, err := s.taskRepo.Update(t)
	if err != nil {
		log.Printf("taskService.Update(s.taskRepo.Update): %s", err)
		return domain.Task{}, err
	}

	return task, nil
}

func (s taskService) Delete(id uint64) error {
	err := s.taskRepo.Delete(id)
	if err != nil {
		log.Printf("taskService.Delete(s.taskRepo.Delete): %s", err)
		return err
	}

	return nil
}

// оновлюємо статус завдання
func (s taskService) UpdateStatus(taskID uint64, userID uint64, status domain.TaskStatus) (domain.Task, error) {
	// знаходимо задачу
	task, err := s.taskRepo.Find(taskID)
	if err != nil {
		return domain.Task{}, err
	}

	// перевіряємо власника
	if task.UserId != userID {
		return domain.Task{}, errors.New("access denied")
	}

	task.Status = status

	return s.taskRepo.Update(task)
}
