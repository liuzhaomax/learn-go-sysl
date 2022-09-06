package main

import (
	"context"
	"github.com/anz-bank/sysl-go-demo/src/handlers"

	"github.com/anz-bank/sysl-go/core"

	"github.com/anz-bank/sysl-go-demo/internal/gen/pkg/servers/Petdemo"
)

type AppConfig struct {
	// Define app-level config fields here.
}

func main() {
	petdemo.Serve(context.Background(),
		func(ctx context.Context, config AppConfig) (*petdemo.ServiceInterface, *core.Hooks, error) {
			// Perform one-time setup based on config here.
			return &petdemo.ServiceInterface{
				// Add handlers here.
				GetPetList: handlers.GetRandomPetPicListRead,
			}, nil, nil
		},
	)
}
