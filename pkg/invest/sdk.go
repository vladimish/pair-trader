package invest

import (
	"context"
	"crypto/tls"
	sdk "github.com/tinkoff/invest-api-go-sdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	url = "invest-public-api.tinkoff.ru:443"
)

type SDK struct {
	ctx  context.Context
	conn *grpc.ClientConn
	md   metadata.MD

	Instruments sdk.InstrumentsServiceClient
	MarketData  sdk.MarketDataServiceClient
	Operations  sdk.OperationsServiceClient
	Accounts    sdk.UsersServiceClient
}

func NewSDK(token string) (*SDK, error) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	md := metadata.New(map[string]string{"Authorization": "Bearer " + token, "x-app-name": "vladimish"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	instruments := sdk.NewInstrumentsServiceClient(conn)
	marketData := sdk.NewMarketDataServiceClient(conn)
	operations := sdk.NewOperationsServiceClient(conn)
	users := sdk.NewUsersServiceClient(conn)

	return &SDK{
		ctx:  ctx,
		conn: conn,
		md:   md,

		Instruments: instruments,
		MarketData:  marketData,
		Operations:  operations,
		Accounts:    users,
	}, nil
}
