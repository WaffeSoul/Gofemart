.PHONY: test
test: # Run all tests.
	./cmd/accrual/accrual_linux_amd64 &
	docker compose up -d
	sleep 5
	go test ./internal/app/service/
	docker compose down -v
	killall -9 accrual_linux_amd64