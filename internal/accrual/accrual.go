package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gofemart/internal/model"
)

var (
	UrlAccrual string
)

func CheckOrder(order string) (*model.Accrual, error) {
	resp, err := http.Get(UrlAccrual + "/api/orders/" + order)
	if err != nil {
		return nil, err
	}
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
		fmt.Println("too many requests")
		time.Sleep(time.Second * 2)
		return nil, errors.New("too many requests")
	case 500:
		return nil, errors.New("internal server error")
	default:
		return nil, errors.New("unknow error")
	}

}
