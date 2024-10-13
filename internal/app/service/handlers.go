package service

import (
	"encoding/json"
	"fmt"
	"gofemart/internal/accrual"
	"gofemart/internal/luhn"
	"gofemart/internal/model"
	"io"
	"net/http"
	"sort"
	"time"
)

func (s *Service) SetOrder() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		orderNumber, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		number := string(orderNumber)
		if !luhn.LuhnAlgorithm(number) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		userId := r.Context().Value("userId").(int)
		fmt.Println(number)
		check, err := s.store.Orders().FindByNumber(number)
		if err != nil {
			fmt.Println(err.Error())
			switch err.Error() {
			case "no number in db":
				break
			default:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if check != nil {
			if check.UserId == userId {
				w.WriteHeader(http.StatusOK)
				return
			} else {
				w.WriteHeader(http.StatusConflict)
				return
			}
		}
		accrualCheck, err := accrual.CheckOrder(string(orderNumber))
		if err != nil {
			switch err.Error() {
			default:
				fmt.Println("test")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		order := model.Order{
			Number:     number,
			UserId:     userId,
			UploadedAt: time.Now().Format("2006-01-02T15:04:05Z"),
		}
		order.AddAccrual(accrualCheck.Accrual, accrualCheck.Status)
		err = s.store.Orders().Create(&order)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	})
}

func (s *Service) GetOrders() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		userId := r.Context().Value("userId").(int)
		orders, err := s.store.Orders().FindByUserId(userId)
		fmt.Println(orders)
		if err != nil {
			switch err.Error() {
			case "no user_id in db":
				w.WriteHeader(http.StatusNoContent)
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		var res []model.Order
		for _, order := range *orders {
			switch order.Status {
			case "PROCESSING", "NEW":
				ac, err := accrual.CheckOrder(order.Number)
				if err != nil {
					fmt.Println(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if ac.Status == "PROCESSED" {
					order.AddAccrual(ac.Accrual, ac.Status)
					res = append(res, order)
					err = s.store.Orders().Update(&order)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				} else if ac.Status == "INVALID" {
					order.AddAccrual(0, ac.Status)
					err = s.store.Orders().Update(&order)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
			case "INVALID":
				continue
			case "PROCESSED":
				res = append(res, order)
			default:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		}
		sort.Slice(res, func(i, j int) bool {
			dateI, _ := time.Parse("2006-01-02T15:04:05Z", res[i].UploadedAt)
			dateJ, _ := time.Parse("2006-01-02T15:04:05Z", res[j].UploadedAt)
			return dateI.Before(dateJ)
		})
		jsonResp, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	})

}

func (s *Service) GetBalance() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		res := model.Balance{
			Current:  0,
			Withdraw: 0,
		}
		userId := r.Context().Value("userId").(int)
		orders, err := s.store.Orders().FindByUserId(userId)
		if err != nil {
			switch err.Error() {
			case "no user_id in db":
				jsonResp, err := json.Marshal(res)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write(jsonResp)
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		for _, order := range *orders {
			switch order.Status {
			case "PROCESSING", "NEW":
				ac, err := accrual.CheckOrder(order.Number)
				if err != nil {
					fmt.Println(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if ac.Status == "PROCESSED" {
					order.AddAccrual(ac.Accrual, ac.Status)
					res.Current += float64(ac.Accrual)
					err = s.store.Orders().Update(&order)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				} else if ac.Status == "INVALID" {
					order.AddAccrual(0, ac.Status)
					err = s.store.Orders().Update(&order)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
			case "INVALID":
				continue
			case "PROCESSED":
				res.Current += float64(order.Accrual)
			default:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		}
		withdraws, err := s.store.Withdrawals().FindByUserId(userId)
		if err != nil {
			switch err.Error() {
			case "no user_id in db":
				jsonResp, err := json.Marshal(res)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write(jsonResp)
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		for _, withdraw := range *withdraws {
			res.Withdraw += withdraw.Sum
		}
		res.Current -= res.Withdraw
		jsonResp, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	})
}

func (s *Service) Withdraw() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resJSON model.ReqWithdraw
		w.Header().Add("Content-Type", "application/json")
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&resJSON)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		userId := r.Context().Value("userId").(int)
		if !luhn.LuhnAlgorithm(resJSON.Order) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		orders, err := s.store.Orders().FindByUserId(userId)
		if err != nil {
			switch err.Error() {
			case "no user_id in db":
				w.WriteHeader(http.StatusPaymentRequired)
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		current := 0.0
		for _, order := range *orders {
			switch order.Status {
			case "PROCESSING", "NEW":
				ac, err := accrual.CheckOrder(order.Number)
				if err != nil {
					fmt.Println(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if ac.Status == "PROCESSED" {
					order.AddAccrual(ac.Accrual, ac.Status)
					current += float64(ac.Accrual)
					err = s.store.Orders().Update(&order)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				} else if ac.Status == "INVALID" {
					order.AddAccrual(0, ac.Status)
					err = s.store.Orders().Update(&order)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
			case "INVALID":
				continue
			case "PROCESSED":
				current += float64(order.Accrual)
			default:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		}
		draw := 0.0
		withdraws, err := s.store.Withdrawals().FindByUserId(userId)
		if err != nil {
			switch err.Error() {
			case "no user_id in db":

				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		for _, withdraw := range *withdraws {
			draw += withdraw.Sum
		}
		current -= draw
		if current < resJSON.Sum {
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}
		withdraw := model.Withdraw{
			Order:       resJSON.Order,
			UserId:      userId,
			Sum:         resJSON.Sum,
			ProcessedAt: time.Now().Format("2006-01-02T15:04:05Z"),
		}
		err = s.store.Withdrawals().Create(&withdraw)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	})
}

func (s *Service) Withdrawals() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		userId := r.Context().Value("userId").(int)

		res, err := s.store.Withdrawals().FindByUserId(userId)
		if err != nil {
			switch err.Error() {
			case "no user_id in db":
				w.WriteHeader(http.StatusNoContent)
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		resSort := *res
		sort.Slice(resSort, func(i, j int) bool {
			dateI, _ := time.Parse("2006-01-02T15:04:05Z", resSort[i].ProcessedAt)
			dateJ, _ := time.Parse("2006-01-02T15:04:05Z", resSort[j].ProcessedAt)
			return dateI.Before(dateJ)
		})
		jsonResp, err := json.Marshal(resSort)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	})
}
