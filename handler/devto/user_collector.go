package devto

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"user-collector/pkg/forem"
)

type UserCollector struct {
	userRepo     UserRepo
	proxy        string
	foremService ForemService
}

func NewUserCollector(userRepo UserRepo, foremService ForemService) *UserCollector {
	return &UserCollector{
		userRepo:     userRepo,
		foremService: foremService,
	}
}

type ForemService interface {
	GetUser(ctx context.Context, userID int64) (*forem.User, error)
}

type UserRepo interface {
	Save(ctx context.Context, user forem.User) error
	GetLastUserID(ctx context.Context, max int) (int64, error)
}

func (u *UserCollector) Collect(ctx context.Context, rps, from, to int) {
	slog.Info(fmt.Sprintf("Conrurrent: %d - Proxy: %s - From: %d - To: %d", rps, u.proxy, from, to))

	latestID, err := u.userRepo.GetLastUserID(ctx, to)
	if err != nil {
		panic(err)
	}
	from = int(max(latestID+1, int64(from)))
	if from >= to {
		slog.Info("from is equals to to. DONE!")
		return
	}
	slog.Info(fmt.Sprintf("Start from %d", from))

	ch := make(chan int, 2)
	wg := sync.WaitGroup{}
	wg.Add(to - from)
	go func() {
		for i := from + 1; i <= to; i++ {
			ch <- i
		}
	}()
	for i := 0; i < rps; i++ {
		go func() {
			defer wg.Done()
			for id := range ch {
				// Create a separate devtoService to open new connection
				user, err := u.foremService.GetUser(ctx, int64(id))
				if err != nil {
					panic(err)
				}
				if user == nil {
					slog.Info(fmt.Sprintf("User not found: %d", id))
					continue
				}
				if err := u.userRepo.Save(ctx, *user); err != nil {
					panic(err)
				}
				slog.Info(fmt.Sprintf("Saved user: %d", id))
			}
		}()
	}
	wg.Wait()
}
