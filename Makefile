gen:
	oapi-codegen --config ./internal/http/api/server.cfg.yaml ./api/oas-server.json
	oapi-codegen --config ./internal/http/api/types.cfg.yaml ./api/oas-server.json


run-test:
	go build -o ./bin/ ./...
	chmod -R +x ./bin/gophermart
	./cmd/gophermarttest/gophermarttest-darwin-arm64 \
		-test.v -test.run=^TestGophermart$ \
		-gophermart-binary-path=./bin/gophermart \
		-gophermart-host=localhost \
		-gophermart-port=8080 \
		-gophermart-database-uri="postgresql://postgres:postgres@localhost:5432/praktikum?sslmode=disable" \
		-accrual-binary-path=./cmd/accrual/accrual_darwin_arm64 \
		-accrual-host=localhost \
		-accrual-port=8081 \
		-accrual-database-uri="postgresql://postgres:postgres@localhost:5432/praktikum?sslmode=disable"
