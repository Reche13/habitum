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

func (s *HabitService) GetHabits(ctx context.Context, userID uuid.UUID, filters *habit.ListFilters) ([]habit.Habit, int, error) {
	habits, total, err := s.habitRepo.List(ctx, userID, filters)
	if err != nil {
		return nil, 0, s.wrapError(err)
	}

	return habits, total, nil
}

func (s *HabitService) GetHabit(ctx context.Context, habitID uuid.UUID, userID uuid.UUID) (*habit.Habit, error) {
	h, err := s.habitRepo.GetByID(ctx, habitID, userID)
	if err != nil {
		return nil, s.wrapError(err)
	}

	return h, nil
}

func (s *HabitService) UpdateHabit(
	ctx context.Context,
	habitID uuid.UUID,
	userID uuid.UUID,
	payload *habit.UpdateHabitPayload,
) (*habit.Habit, error) {
	// Verify habit exists and belongs to user
	_, err := s.habitRepo.GetByID(ctx, habitID, userID)
	if err != nil {
		return nil, s.wrapError(err)
	}

	updatedHabit, err := s.habitRepo.Update(ctx, habitID, userID, payload)
	if err != nil {
		return nil, s.wrapError(err)
	}

	return updatedHabit, nil
}

func (s *HabitService) DeleteHabit(ctx context.Context, habitID uuid.UUID, userID uuid.UUID) error {
	// Verify habit exists and belongs to user
	_, err := s.habitRepo.GetByID(ctx, habitID, userID)
	if err != nil {
		return s.wrapError(err)
	}

	if err := s.habitRepo.Delete(ctx, habitID, userID); err != nil {
		return s.wrapError(err)
	}

	return nil
}