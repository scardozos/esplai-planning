package models

import (
	"context"
	"log"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	wh "github.com/scardozos/ep-weekhandler/grpc/dates"
	"google.golang.org/grpc"
)

type GrpcClientContext struct {
	DatesClient wh.DatesClient
}

type GrpcClient struct {
	Context *GrpcClientContext
}

func (s *GrpcClient) AddStaticDate(req *DateTime) error {
	c := s.Context.DatesClient
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	now := time.Now()
	defer cancel()
	log.Printf("Calling SetStaticWeek")
	res, err := c.SetStaticWeek(ctx, &wh.SetStaticWeekRequest{StaticWeek: &wh.Date{Year: req.Year, Month: req.Month, Day: req.Day}})
	then := time.Since(now)
	if err != nil {
		log.Printf("error calling SetStaticWeek: %v", err)
		return err
	}
	log.Printf("Successfully added date %v - Took %v", res.SetWeek, then)
	return nil
}

func (s *GrpcClient) GetNonWeeks(opts ...grpc.CallOption) ([]time.Time, error) {
	c := s.Context.DatesClient
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := c.GetStaticWeeks(ctx, &wh.GetStaticWeeksRequest{}, grpc_retry.WithMax(3))
	if err != nil {
		log.Printf("Error when executing GetStaticWeeks: %v", err)
		return nil, err
	}

	var retObj = make([]time.Time, len(res.StaticWeeks))
	for i, e := range res.StaticWeeks {
		retObj[i] = time.Date(int(e.Year), time.Month(e.Month), int(e.Day), 0, 0, 0, 0, time.UTC)
	}

	if len(retObj) == 0 {
		log.Print("GetNonWeeks() returned 0 static weeks")
	}
	return retObj, nil
}

func (s *GrpcClient) IsNonWeek(req *DateTime) error {
	c := s.Context.DatesClient
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.IsWeekStatic(ctx, &wh.IsWeekStaticRequest{
		Week: &wh.Date{Year: req.Year, Month: req.Month, Day: req.Day},
	})
	if err != nil {
		log.Printf("IsNonWeek error: %v", err)
		return err
	}
	switch res.IsStatic {
	case true:
		log.Printf("week %v is static", res.RequestedWeek)
	case false:
		log.Printf("week %v is not static", res.RequestedWeek)
	}
	return nil
}

func (s *GrpcClient) UnsetNonWeek(req *DateTime) error {
	c := s.Context.DatesClient
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.UnsetStaticWeek(ctx, &wh.UnsetStaticWeekRequest{
		StaticWeek: &wh.Date{Year: req.Year, Month: req.Month, Day: req.Day},
	})
	if err != nil {
		log.Printf("UnsetNonWeek error: %v", err)
		return err
	}
	w := res.UnsetWeek
	log.Printf("Successfully removed week %v-%v-%v.", w.Day, w.Month, w.Year)
	return nil
}
