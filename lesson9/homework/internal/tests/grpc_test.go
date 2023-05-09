package tests

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"homework9/internal/adapters/adfilters"
	"homework9/internal/adapters/userrepo"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/app"
	grpcPort "homework9/internal/ports/grpc"
)

func GetTestClient(t *testing.T) (grpcPort.AdServiceClient, context.Context) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(grpcPort.UnaryInterceptor, grpcPort.RecoveryInterceptor))
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New(), adfilters.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	return client, ctx
}

func TestGRPCCreateUser(t *testing.T) {
	client, ctx := GetTestClient(t)

	res, err := client.CreateUser(ctx, &grpcPort.UniversalUser{Nickname: "name", Email: "somemail@mail.com", UserId: 123})
	assert.NoError(t, err)
	assert.Equal(t, "name", res.Nickname)
	assert.Equal(t, "somemail@mail.com", res.Email)
	assert.Equal(t, int64(123), res.UserId)
}

func TestGRPCCreateAd(t *testing.T) {
	client, ctx := GetTestClient(t)

	a, err := client.CreateUser(ctx, &grpcPort.UniversalUser{Nickname: "name", Email: "somemail@mail.com", UserId: 123})
	assert.NoError(t, err)

	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "title", Text: "text", UserId: a.UserId})
	assert.NoError(t, err, "client.CreateAd")
	assert.Equal(t, "title", res.Title)
	assert.Equal(t, "text", res.Text)
	assert.Equal(t, a.UserId, res.AuthorId)
	assert.Equal(t, false, res.Published)

}

func TestGRPCChangeAdStatus(t *testing.T) {
	client, ctx := GetTestClient(t)

	a, _ := client.CreateUser(ctx, &grpcPort.UniversalUser{Nickname: "name", Email: "somemail@mail.com", UserId: 123})

	ad, _ := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "title", Text: "text", UserId: a.UserId})

	updatedAd, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: ad.Id, UserId: ad.AuthorId, Published: true})
	assert.NoError(t, err)
	assert.Equal(t, ad.Title, updatedAd.Title)
	assert.Equal(t, ad.Text, updatedAd.Text)
	assert.Equal(t, ad.AuthorId, updatedAd.AuthorId)
	assert.Equal(t, true, updatedAd.Published)
}

func TestGRPCUpdateAd(t *testing.T) {
	client, ctx := GetTestClient(t)

	a, _ := client.CreateUser(ctx, &grpcPort.UniversalUser{Nickname: "name", Email: "somemail@mail.com", UserId: 123})

	ad, _ := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "title", Text: "text", UserId: a.UserId})

	updatedAd, err := client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{AdId: ad.Id, Title: "title1", Text: "text1", UserId: ad.AuthorId})
	assert.NoError(t, err)
	assert.Equal(t, "title1", updatedAd.Title)
	assert.Equal(t, "text1", updatedAd.Text)
	assert.Equal(t, ad.AuthorId, updatedAd.AuthorId)
	assert.Equal(t, ad.Published, updatedAd.Published)
}

func TestGRPCDeleteAd(t *testing.T) {
	client, ctx := GetTestClient(t)

	a, _ := client.CreateUser(ctx, &grpcPort.UniversalUser{Nickname: "name", Email: "somemail@mail.com", UserId: 123})

	ad, _ := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "title", Text: "text", UserId: a.UserId})

	resp, err := client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AdId: ad.Id, AuthorId: ad.AuthorId})
	assert.NoError(t, err)
	assert.Equal(t, ad.Title, resp.Title)
	assert.Equal(t, ad.Text, resp.Text)
	assert.Equal(t, ad.AuthorId, resp.AuthorId)
	assert.Equal(t, ad.Published, resp.Published)
}
