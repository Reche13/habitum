package service

import (
	"github.com/reche13/habitum/internal/repository"
	"github.com/reche13/habitum/internal/sqlerr"
)

type BaseService struct {
	resourceName string
}

func (b *BaseService) wrapError(err error) error {
	return sqlerr.WrapError(err, b.resourceName)
}

type Services struct{
	User *UserService
	Habit *HabitService
	HabitLog *HabitLogService
	Analytics *AnalyticsService
	Calendar *CalendarService
	Dashboard *DashboardService
}

func NewServices(repos *repository.Repositories) *Services {
	habitLogService := NewHabitLogService(repos.HabitLog)
	return &Services{
		User: NewUserService(repos.User),
		Habit: NewHabitService(repos.Habit, habitLogService),
		HabitLog: habitLogService,
		Analytics: NewAnalyticsService(repos.Habit, repos.HabitLog),
		Calendar: NewCalendarService(repos.Habit, repos.HabitLog),
		Dashboard: NewDashboardService(repos.Habit, repos.HabitLog),
	}
}