package handlers

import (
	"context"
	"net/http"

	petdemo "github.com/anz-bank/sysl-go-demo/internal/gen/pkg/servers/Petdemo"
	"github.com/anz-bank/sysl-go-demo/internal/gen/pkg/servers/Petdemo/petstore"
	"github.com/anz-bank/sysl-go/common"
)

// GetRandomPetPicListRead reads random pic from downstream
func GetRandomPetPicListRead(ctx context.Context,
	getRandomPetPicListRequest *petdemo.GetPetListRequest,
	client petdemo.GetPetListClient) (*petdemo.Pet, error) {

	// Retrieve the pets (using a fresh set of headers for the request)
	// This is required because by default Sysl-go reuses inbound headers for downstream requests
	ctx = common.RequestHeaderToContext(ctx, http.Header{})
	reqPetstore := petstore.GetPetListRequest{}
	pet, err := client.PetstoreGetPetList(ctx, &reqPetstore)
	if err != nil {
		return nil, err
	}

	// Return the result
	return &petdemo.Pet{
		Breed: string(*pet),
	}, nil
}
