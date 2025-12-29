package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/reche13/habitum/internal/model/habit"
	"github.com/reche13/habitum/internal/repository"
)

type HabitService struct {
	*BaseService
	habitRepo      *repository.HabitRepository
	habitLogService *HabitLogService
}

func NewHabitService(
	habitRepo *repository.HabitRepository,
	habitLogService *HabitLogService,
) *HabitService {
	return &HabitService{
		BaseService: &BaseService{
			resourceName: "habit",
		},
		habitRepo:       habitRepo,
		habitLogService: habitLogService,
	}
}

func (s *HabitService) CreateHabit(
	ctx context.Context,
	userID uuid.UUID,
	payload *habit.CreateHabitPayload,
) (*habit.HabitResponse, error) {
	createdHabit, err := s.habitRepo.Create(ctx, userID, payload)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Enrich with computed fields (will be zeros for new habit)
	enriched, err := s.enrichHabitWithStats(ctx, userID, *createdHabit)
	if err != nil {
		// If enrichment fails, return habit without computed fields
		return &habit.HabitResponse{
			Habit:            *createdHabit,
			CompletionRate:   0,
			CompletedToday:   false,
			CompletedThisWeek: 0,
			CompletionHistory: []string{},
		}, nil
	}

	return &enriched, nil
}

func (s *HabitService) GetHabits(ctx context.Context, userID uuid.UUID, filters *habit.ListFilters) ([]habit.HabitResponse, int, error) {
	habits, total, err := s.habitRepo.List(ctx, userID, filters)
	if err != nil {
		return nil, 0, s.wrapError(err)
	}

	// Enrich habits with computed fields
	enrichedHabits := make([]habit.HabitResponse, len(habits))
	for i := range habits {
		enriched, err := s.enrichHabitWithStats(ctx, userID, habits[i])
		if err != nil {
			// If enrichment fails, return habit without computed fields
			enrichedHabits[i] = habit.HabitResponse{
				Habit:            habits[i],
				CompletionRate:   0,
				CompletedToday:   false,
				CompletedThisWeek: 0,
				CompletionHistory: []string{},
			}
			continue
		}
		enrichedHabits[i] = enriched
	}

	return enrichedHabits, total, nil
}

func (s *HabitService) GetHabit(ctx context.Context, habitID uuid.UUID, userID uuid.UUID) (*habit.HabitResponse, error) {
	h, err := s.habitRepo.GetByID(ctx, habitID, userID)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Enrich with computed fields
	enriched, err := s.enrichHabitWithStats(ctx, userID, *h)
	if err != nil {
		// If enrichment fails, return habit without computed fields
		return &habit.HabitResponse{
			Habit:            *h,
			CompletionRate:   0,
			CompletedToday:   false,
			CompletedThisWeek: 0,
			CompletionHistory: []string{},
		}, nil
	}

	return &enriched, nil
}

// enrichHabitWithStats adds computed fields to a habit and returns HabitResponse
func (s *HabitService) enrichHabitWithStats(ctx context.Context, userID uuid.UUID, h habit.Habit) (habit.HabitResponse, error) {
	// Calculate current streak
	currentStreak, err := CalculateCurrentStreak(ctx, s.habitLogService.habitLogRepo, userID, h.ID, h.Frequency)
	if err != nil {
		currentStreak = h.CurrentStreak // Fallback to stored value
	}

	// Calculate longest streak
	longestStreak, err := CalculateLongestStreak(ctx, s.habitLogService.habitLogRepo, userID, h.ID, h.Frequency)
	if err != nil {
		longestStreak = h.LongestStreak // Fallback to stored value
	}

	// Update stored streaks if they differ
	if currentStreak != h.CurrentStreak || longestStreak != h.LongestStreak {
		// Update in database
		s.habitRepo.UpdateStreaks(ctx, h.ID, userID, currentStreak, longestStreak)
		h.CurrentStreak = currentStreak
		h.LongestStreak = longestStreak
	}

	// Calculate completion rate
	completionRate, err := CalculateCompletionRate(ctx, s.habitLogService.habitLogRepo, userID, h.ID, h.CreatedAt, h.Frequency)
	if err != nil {
		completionRate = 0
	}

	// Check if completed today
	completedToday, completedTodayAt, err := GetTodayStatus(ctx, s.habitLogService.habitLogRepo, userID, h.ID)
	if err != nil {
		completedToday = false
		completedTodayAt = nil
	}

	// Get completed this week count
	completedThisWeek, err := GetCompletedThisWeek(ctx, s.habitLogService.habitLogRepo, userID, h.ID)
	if err != nil {
		completedThisWeek = 0
	}

	// Get completion history (last year, limited to 365 dates)
	completionHistory, err := GetCompletionHistoryDates(ctx, s.habitLogService.habitLogRepo, userID, h.ID, 365)
	if err != nil {
		completionHistory = []string{}
	}

	return habit.HabitResponse{
		Habit:             h,
		CompletionRate:    completionRate,
		CompletedToday:    completedToday,
		CompletedTodayAt:  completedTodayAt,
		CompletedThisWeek: completedThisWeek,
		CompletionHistory: completionHistory,
	}, nil
}

func (s *HabitService) UpdateHabit(
	ctx context.Context,
	habitID uuid.UUID,
	userID uuid.UUID,
	payload *habit.UpdateHabitPayload,
) (*habit.HabitResponse, error) {
	// Verify habit exists and belongs to user
	_, err := s.habitRepo.GetByID(ctx, habitID, userID)
	if err != nil {
		return nil, s.wrapError(err)
	}

	updatedHabit, err := s.habitRepo.Update(ctx, habitID, userID, payload)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Enrich with computed fields
	enriched, err := s.enrichHabitWithStats(ctx, userID, *updatedHabit)
	if err != nil {
		// If enrichment fails, return habit without computed fields
		return &habit.HabitResponse{
			Habit:            *updatedHabit,
			CompletionRate:   0,
			CompletedToday:   false,
			CompletedThisWeek: 0,
			CompletionHistory: []string{},
		}, nil
	}

	return &enriched, nil
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