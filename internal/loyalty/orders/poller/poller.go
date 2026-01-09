package poller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/orders/model"
)

type orderService interface {
	FindOrdersToPoll(ctx context.Context) ([]model.AccrualOrder, error)
	UpdateOrderAccrual(ctx context.Context, order model.AccrualOrder) error
}

type balanceService interface {
	AddAccrual(ctx context.Context, userID int64, amount float64) error
}

// Poller polls accrual service for updates on the orders.
type Poller struct {
	cfg *Config
	os  orderService
	bs  balanceService
}

// New creates a new poller.
func New(os orderService, bs balanceService, address string) (*Poller, error) {
	cfg, err := NewConfig(address)
	if err != nil {
		return nil, fmt.Errorf("error creating new poller: %w", err)
	}
	return &Poller{
		cfg: cfg,
		os:  os,
		bs:  bs,
	}, nil
}

// accrualRequest makes a single accrual service request.
func (p *Poller) accrualRequest(number string) (*model.AccrualOrder, error) {
	req, err := http.NewRequest("GET", p.cfg.Address+number, nil)
	if err != nil {
		slog.Warn(err.Error())
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Warn(err.Error())
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	slog.Info("sent request to accrual")

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		slog.Warn(err.Error())
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var order model.AccrualOrder
	err = json.Unmarshal(buf.Bytes(), &order)
	if err != nil {
		slog.Warn(err.Error())
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &order, nil
}

func (p *Poller) pollWorker(
	ctx context.Context,
	ordersChan <-chan string,
	resChan chan<- *model.AccrualOrder,
	errChan chan<- error,
) {
	for order := range ordersChan {
		res, err := p.accrualRequest(order)
		if err != nil {
			errChan <- err
		}
		resChan <- res
	}
}

// startPolling creates new workers to go through polling channel
func (p *Poller) startPolling(
	ctx context.Context,
	pollable []model.AccrualOrder,
	ordersChan chan string,
	resChan chan *model.AccrualOrder,
) error {
	errChan := make(chan error)

	for i := 0; i < p.cfg.WorkerNum; i++ {
		go p.pollWorker(ctx, ordersChan, resChan, errChan)
	}

	for _, toPoll := range pollable {
		ordersChan <- toPoll.Order
	}

	close(ordersChan)

	for res := range resChan {
		p.os.UpdateOrderAccrual(ctx, *res)
	}

	return nil
}

// PollerLoop once in cfg.Interval gathers all pollable orders and polls all orders.
func (p *Poller) PollerLoop(ctx context.Context, errChan chan<- error) {
	ticker := time.NewTicker(p.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pollable, err := p.os.FindOrdersToPoll(ctx)
			if err != nil {
				errChan <- err
			}

			ordersChan := make(chan string, len(pollable))
			resChan := make(chan *model.AccrualOrder, len(pollable))

			err = p.startPolling(ctx, pollable, ordersChan, resChan)
			if err != nil {
				errChan <- err
			}
		case <-ctx.Done():
			return
		}
	}
}
