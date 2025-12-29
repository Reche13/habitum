package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/reche13/habitum/internal/lib"
	"github.com/reche13/habitum/internal/model/habitlog"
	"github.com/reche13/habitum/internal/repository"
)

type HabitLogService struct {
	*BaseService
	habitLogRepo *repository.HabitLogRepository
}

func NewHabitLogService(
	habitLogRepo *repository.HabitLogRepository,
) *HabitLogService {
	return &HabitLogService{
		BaseService: &BaseService{
			resourceName: "habitlog",
		},
		habitLogRepo: habitLogRepo,
	}
}

func (s *HabitLogService) SetCompletion(
	ctx context.Context,
	userID uuid.UUID,
	payload *habitlog.HabitLogPayload,
) (*habitlog.HabitLog, error) {

	payload.LogDate = lib.NormalizeDate(payload.LogDate)

	return s.habitLogRepo.Create(ctx, userID, payload)
}


func (s *HabitLogService) GetByDate(
	ctx context.Context,
	userID uuid.UUID,
	date time.Time,
) ([]habitlog.HabitLog, error) {

	return s.habitLogRepo.GetByDate(
		ctx,
		userID,
		lib.NormalizeDate(date),
	)
}

func (s *HabitLogService) GetByDateRange(
	ctx context.Context,
	userID uuid.UUID,
	from time.Time,
	to time.Time,
) ([]habitlog.HabitLog, error) {

	return s.habitLogRepo.GetByDateRange(
		ctx,
		userID,
		lib.NormalizeDate(from),
		lib.NormalizeDate(to),
	)
}

func (s *HabitLogService) GetByHabit(
	ctx context.Context,
	userID uuid.UUID,
	habitID uuid.UUID,
	from time.Time,
	to time.Time,
) ([]habitlog.HabitLog, error) {

	return s.habitLogRepo.GetByHabit(
		ctx,
		userID,
		habitID,
		lib.NormalizeDate(from),
		lib.NormalizeDate(to),
	)
}


func (s *HabitLogService) GetToday(
	ctx context.Context,
	userID uuid.UUID,
) ([]habitlog.HabitLog, error) {

	today := lib.NormalizeDate(time.Now().UTC())

	return s.habitLogRepo.GetByDate(ctx, userID, today)
}


func (s *HabitLogService) GetThisWeek(
	ctx context.Context,
	userID uuid.UUID,
) ([]habitlog.HabitLog, error) {

	now := lib.NormalizeDate(time.Now().UTC())

	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	start := now.AddDate(0, 0, -(weekday - 1))
	end := start.AddDate(0, 0, 6)

	return s.habitLogRepo.GetByDateRange(ctx, userID, start, end)
}


func (s *HabitLogService) GetThisMonth(
	ctx context.Context,
	userID uuid.UUID,
	year int,
	month time.Month,
) ([]habitlog.HabitLog, error) {

	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)

	return s.habitLogRepo.GetByDateRange(ctx, userID, start, end)
}


func (s *HabitLogService) GetThisYear(
	ctx context.Context,
	userID uuid.UUID,
	year int,
) ([]habitlog.HabitLog, error) {

	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)

	return s.habitLogRepo.GetByDateRange(ctx, userID, start, end)
}


func (s *HabitLogService) GetHabitHistory(
	ctx context.Context,
	userID uuid.UUID,
	habitID uuid.UUID,
) ([]habitlog.HabitLog, error) {
	start := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	end := lib.NormalizeDate(time.Now().UTC())

	return s.habitLogRepo.GetByHabit(ctx, userID, habitID, start, end)
}


func (s *HabitLogService) GetCompletionRate(
	ctx context.Context,
	userID uuid.UUID,
	from time.Time,
	to time.Time,
) (float64, error) {

	logs, err := s.GetByDateRange(ctx, userID, from, to)
	if err != nil {
		return 0, err
	}

	if len(logs) == 0 {
		return 0, nil
	}

	completed := 0
	for _, l := range logs {
		if l.Completed {
			completed++
		}
	}

	return float64(completed) / float64(len(logs)), nil
}

func (s *HabitLogService) MarkComplete(
	ctx context.Context,
	userID uuid.UUID,
	habitID uuid.UUID,
	payload *habitlog.HabitLogPayload,
) (*habitlog.HabitLog, error) {
	payload.LogDate = lib.NormalizeDate(payload.LogDate)
	payload.HabitID = habitID
	payload.Completed = true

	log, err := s.habitLogRepo.Create(ctx, userID, payload)
	if err != nil {
		return nil, err
	}

	return log, nil
}

func (s *HabitLogService) UnmarkComplete(
	ctx context.Context,
	userID uuid.UUID,
	habitID uuid.UUID,
	logDate time.Time,
) error {
	normalizedDate := lib.NormalizeDate(logDate)
	return s.habitLogRepo.DeleteByHabitAndDate(ctx, userID, habitID, normalizedDate)
}

func (s *HabitLogService) GetCompletions(
	ctx context.Context,
	userID uuid.UUID,
	habitID uuid.UUID,
	startDate time.Time,
	endDate time.Time,
	limit int,
) ([]habitlog.HabitLog, int, error) {
	normalizedStart := lib.NormalizeDate(startDate)
	normalizedEnd := lib.NormalizeDate(endDate)

	return s.habitLogRepo.GetByHabitWithLimit(ctx, userID, habitID, normalizedStart, normalizedEnd, limit)
}

func (s *HabitLogService) GetCompletionHistory(
	ctx context.Context,
	userID uuid.UUID,
	habitID uuid.UUID,
	year *int,
	allTime bool,
) ([]time.Time, int, int, error) {
	var startDate, endDate time.Time

	if allTime {
		// Get all time history
		startDate = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate = lib.NormalizeDate(time.Now().UTC())
	} else if year != nil {
		// Get specific year
		startDate = time.Date(*year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(*year, 12, 31, 23, 59, 59, 0, time.UTC)
	} else {
		// Default to current year
		now := time.Now().UTC()
		currentYear := now.Year()
		startDate = time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(currentYear, 12, 31, 23, 59, 59, 0, time.UTC)
	}

	logs, err := s.habitLogRepo.GetByHabit(ctx, userID, habitID, startDate, endDate)
	if err != nil {
		return nil, 0, 0, err
	}

	// Extract only completed logs and get dates
	dates := make([]time.Time, 0)
	for _, log := range logs {
		if log.Completed {
			dates = append(dates, log.LogDate)
		}
	}

	// Calculate total days in period
	totalDays := int(endDate.Sub(startDate).Hours()/24) + 1
	completedDays := len(dates)

	return dates, totalDays, completedDays, nil
}
