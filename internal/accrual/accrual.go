package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"gofemart/internal/model"
)

var (
	UrlAccrual string
)

func CheckOrder(order string) (*model.Accrual, error) {
	path_url, _ := url.JoinPath(UrlAccrual, "/api/orders/", order)
	fmt.Println(path_url)
	for i := 0; i < 3; i++ {
		resp, err := http.Get(path_url)
		if err != nil {
			return nil, err
		}
		fmt.Println(resp.StatusCode)
		switch resp.StatusCode {
		case 200:
			accrual := &model.Accrual{}
			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&accrual)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			return accrual, nil
		case 204:
			return nil, errors.New("the order is not registered in the payment system")
		case 429:
			fmt.Println("too many requests")
			time.Sleep(time.Second * 5 * time.Duration(i+1))
			continue
			return nil, errors.New("too many requests")
		case 500:
			return nil, errors.New("internal server error")
		default:
			return nil, errors.New("unknow error")
		}
	}
	return nil, errors.New("error timeout")
}
