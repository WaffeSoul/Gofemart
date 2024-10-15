package accrual

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"gofemart/internal/config"
	"gofemart/internal/model"
)

type Accrual struct {
	URL string
}

func NewAccrual(conf *config.Config) *Accrual {
	return &Accrual{
		URL: conf.Accrual,
	}
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
