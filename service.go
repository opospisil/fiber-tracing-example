package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type HeavyDutyService interface {
	DoSomethingExpansive(context.Context) error
}

type DefaultHeavyDutyService struct{}

func (hds *DefaultHeavyDutyService) DoSomethingExpansive(c context.Context) error {
	ctx, cancelFn := context.WithTimeout(c, time.Second*4)

	go func() {
		err := ActualFireForget(ctx)
		if err != nil {
			fmt.Printf("goroutine finished with an error: %v\n", err)
		}
		if err == nil {
			fmt.Println("This one was success")
		}
		cancelFn()
	}()
	return nil
}

// simulates random response time as well as random failures for different reasons
func ActualFireForget(c context.Context) error {
	n := rand.Intn(10)
	fmt.Printf("Starting fire forget for number %v\n", n)
	for range n {
		select {
		case <-c.Done():
			return c.Err()
		default:
			time.Sleep(time.Second)
		}
	}
	if n%2 == 0 {
		return nil
	}
	return fmt.Errorf("received %v, failing", n)
}

func NewDefaultHeavHeavyDutyService() *DefaultHeavyDutyService {
	return &DefaultHeavyDutyService{}
}
