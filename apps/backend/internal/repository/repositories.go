package repository

import "github.com/jackc/pgx/v5/pgxpool"


type Repositories struct{
	User *UserRepository
	Habit *HabitRepository
	HabitLog *HabitLogRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		User: NewUserRepository(db),
		Habit: NewHabitRepository(db),
		HabitLog: NewHabitLogRepository(db),
	}
}