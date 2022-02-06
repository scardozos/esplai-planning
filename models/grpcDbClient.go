package models

import (
	"context"
	"log"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	wh "github.com/scardozos/esplai-weeks-db/api/weeksdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO: IMPROVE DOCUMENTATION
type GrpcClientContext struct {
	DatesClient wh.WeeksDatabaseClient
}

type GrpcClient struct {
	Context *GrpcClientContext
}

func NewGrpcClientContext(endpoint string) (*GrpcClientContext, error) {
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(10 * time.Millisecond)),
		grpc_retry.WithCodes(codes.Internal, codes.Unavailable),
	}
	datesConn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
	)
	if err != nil {
		return nil, err
	}

	return &GrpcClientContext{DatesClient: wh.NewWeeksDatabaseClient(datesConn)}, nil

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

// Commented out as methods won't be used by esplai-planning, rather than esplain-planning-admin
/*
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
*/
