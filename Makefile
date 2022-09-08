SYSLGO_SYSL=specs/petdemo.sysl
SYSLGO_PACKAGES=Petdemo
SYSLGO_APP.Petdemo = Petdemo

-include local.mk
include codegen.mk

.PHONY: clean
clean:
	rm -rf internal/gen

.PHONY: test
test: gen-all-servers
	go test -v ./...

run:
	go run cmd/Petdemo/main.go config/config.yaml
