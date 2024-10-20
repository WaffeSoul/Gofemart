package accrual

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"gofemart/internal/config"
	"gofemart/internal/logger"
	"gofemart/internal/model"
	"gofemart/internal/storage"

	"go.uber.org/zap"
)

type Accrual struct {
	URL      string
	QueueCh  chan string
	DoneCh   chan struct{}
	ResultCh chan model.Accrual
	store    storage.Store
}

func NewAccrual(conf *config.Config, store *storage.Store) *Accrual {
	DoneCh := make(chan struct{})
	acc := Accrual{
		URL:    conf.Accrual,
		store:  *store,
		DoneCh: DoneCh,
	}
	acc.CheckQueue()
	acc.SaveToDB()
	return &acc
}

func (a *Accrual) Finish() {
	close(a.DoneCh)
}

func (a *Accrual) CheckOrder(order string) (*model.Accrual, error) {
	pathURL, _ := url.JoinPath(a.URL, "/api/orders/", order)
	resp, err := http.Get(pathURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		accrual := &model.Accrual{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&accrual)
		if err != nil {
			return nil, err
		}
		return accrual, nil
	case 204:
		return nil, errors.New("the order is not registered in the payment system")
	case 429:
		return nil, errors.New("too many requests")
	case 500:
		return nil, errors.New("internal server error")
	default:
		return nil, errors.New("unknow error")
	}
}

func (a *Accrual) CheckQueue() {
	a.QueueCh = make(chan string)

	go func() {
		defer close(a.QueueCh)
		for data := range a.QueueCh {
			logger.Info("check order", zap.String("order", data))
			res, err := a.CheckOrder(data)
			switch err {
			case nil:
				switch res.Status {
				case "REGISTERED", "PROCESSING":
					a.AddToQueue(data)
				case "INVALID":
					res.Accrual = 0
				}
				a.ResultCh <- *res
			case errors.New("the order is not registered in the payment system"):
				res = &model.Accrual{
					Order:   data,
					Status:  "INVALID",
					Accrual: 0,
				}
				a.ResultCh <- *res
			case errors.New("too many requests"):
				time.Sleep(5 * time.Second)
				a.AddToQueue(data)
			default:
				logger.Error("invalid get order status", zap.Error(err))
			}
		}
	}()

}

func (a *Accrual) AddToQueue(order string) {
	go func() {
		a.QueueCh <- order
	}()
}

func (a *Accrual) SaveToDB() {
	a.ResultCh = make(chan model.Accrual)
	go func() {
		defer close(a.ResultCh)
		for res := range a.ResultCh {
			data, err := a.store.Orders().FindByNumber(res.Order)
			if err != nil {
				logger.Error("invalid get number", zap.Error(err))
				continue
			}
			if !data.CheckStatus(res.Status) {
				data.AddAccrual(res.Accrual, res.Status)
				err = a.store.Orders().Update(data)
				if err != nil {
					logger.Error("invalid update order", zap.Error(err))
				}
			}
		}
	}()
}
