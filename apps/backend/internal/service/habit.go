package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/reche13/habitum/internal/model/habit"
	"github.com/reche13/habitum/internal/repository"
)

type HabitService struct {
	*BaseService
	habitRepo *repository.HabitRepository
}

func NewHabitService(
	habitRepo *repository.HabitRepository,
) *HabitService {
	return &HabitService{
		BaseService: &BaseService{
			resourceName: "habit",
		},
		habitRepo: habitRepo,
	}
}

func (s *HabitService) CreateHabit(
	ctx context.Context,
	userID uuid.UUID,
	payload *habit.CreateHabitPayload,
) (*habit.Habit, error) {
	createdHabit, err := s.habitRepo.Create(ctx, userID, payload)
	if err != nil {
		return nil, s.wrapError(err)
	}

	return createdHabit, nil
}

func (s *HabitService) GetHabits(ctx context.Context, userID uuid.UUID) ([]habit.Habit, error) {
	habits, err := s.habitRepo.List(ctx, userID)
	if err != nil {
		return nil, s.wrapError(err)
	}

	return habits, nil
}