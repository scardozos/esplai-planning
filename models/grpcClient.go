package models

import (
	"context"
	"log"
	"time"

	wh "github.com/scardozos/ep-weekhandler/grpc/dates"
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
		log.Printf("me cago en todo: %v", err)
		return err
	}
	log.Printf("Successfully added date %v - Took %v", res.SetWeek, then)
	return nil
}

func (s *GrpcClient) GetNonWeeks() error {
	c := s.Context.DatesClient
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.GetStaticWeeks(ctx, &wh.GetStaticWeeksRequest{})
	if err != nil {
		log.Printf("Error when executing GetStaticWeeks: %v", err)
		return err
	}
	log.Printf("Got result from GetStaticWeeks:\n%v", res)
	return nil
}
