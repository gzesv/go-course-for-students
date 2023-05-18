package tests

import (
	"context"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"homework10/internal/adapters/adfilters"
	"homework10/internal/adapters/userrepo"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/app"
	grpcPort "homework10/internal/ports/grpc"
)

var (
	ErrorBadRequest = status.Error(codes.InvalidArgument, app.ErrWrongFormat.Error())
	ErrorForbidden  = status.Error(codes.PermissionDenied, app.ErrAccessDenied.Error())
)

type SuiteTest struct {
	suite.Suite
	client grpcPort.AdServiceClient
	ctx    context.Context
	lis    *bufconn.Listener
	srv    *grpc.Server
	cancel context.CancelFunc
	conn   *grpc.ClientConn
}

func (suite *SuiteTest) SetupTest() {
	suite.lis = bufconn.Listen(1024 * 1024)

	suite.srv = grpc.NewServer(grpc.ChainUnaryInterceptor(grpcPort.UnaryInterceptor, grpcPort.RecoveryInterceptor))

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New(), adfilters.New()))
	grpcPort.RegisterAdServiceServer(suite.srv, svc)

	go func() {
		srv := suite.srv
		lis := suite.lis
		_ = srv.Serve(lis)
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return suite.lis.Dial()
	}

	suite.ctx, suite.cancel = context.WithTimeout(context.Background(), 15*time.Second)

	var err error
	suite.conn, err = grpc.DialContext(suite.ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	suite.Assert().NoError(err)

	suite.client = grpcPort.NewAdServiceClient(suite.conn)
}

func (suite *SuiteTest) TestGRPCCreateUser() {
	res, err := suite.client.CreateUser(suite.ctx, &grpcPort.UniversalUser{Nickname: "name", Email: "somemail@mail.com", UserId: 123})
	suite.Assert().NoError(err)
	suite.Assert().Equal("name", res.Nickname)
	suite.Assert().Equal("somemail@mail.com", res.Email)
	suite.Assert().Equal(int64(123), res.UserId)

	_, err = suite.client.CreateUser(suite.ctx, &grpcPort.UniversalUser{Nickname: "name", Email: "somemail@mail.com", UserId: 123})
	suite.Assert().ErrorIs(err, ErrorBadRequest)
}

func (suite *SuiteTest) TestGRPCCreateAd() {
	_, _ = suite.client.CreateUser(suite.ctx, &grpcPort.UniversalUser{Nickname: "name", Email: "somemail@mail.com", UserId: 123})
	res, err := suite.client.CreateAd(suite.ctx, &grpcPort.CreateAdRequest{Title: "title", Text: "text", UserId: 123})
	suite.Assert().NoError(err)
	suite.Assert().Equal("title", res.Title)
	suite.Assert().Equal("text", res.Text)
	suite.Assert().Equal(int64(123), res.AuthorId)
	suite.Assert().Equal(false, res.Published)
}

func (suite *SuiteTest) TearDownTest() {
	_ = suite.lis.Close()
	suite.srv.Stop()
	suite.cancel()
	_ = suite.conn.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(SuiteTest))
}

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
	assert.NoError(t, err)

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

	_, err = client.CreateUser(ctx, &grpcPort.UniversalUser{Nickname: "abc", Email: "cat@mail.com", UserId: 123})
	assert.ErrorIs(t, err, ErrorBadRequest)
}

func TestGRPCCreateAd(t *testing.T) {
	client, ctx := GetTestClient(t)

	a, err := client.CreateUser(ctx, &grpcPort.UniversalUser{Nickname: "name", Email: "somemail@mail.com", UserId: 123})
	assert.NoError(t, err)

	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "title", Text: "text", UserId: a.UserId})
	assert.NoError(t, err)
	assert.Equal(t, "title", res.Title)
	assert.Equal(t, "text", res.Text)
	assert.Equal(t, a.UserId, res.AuthorId)
	assert.Equal(t, false, res.Published)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "cat", Text: "text", UserId: 5})
	assert.ErrorIs(t, err, ErrorBadRequest)
}

func TestGRPCChangeAdStatus(t *testing.T) {
	client, ctx := GetTestClient(t)

	a, _ := client.CreateUser(ctx, &grpcPort.UniversalUser{Nickname: "name", Email: "somemail@mail.com", UserId: 123})
	b, _ := client.CreateUser(ctx, &grpcPort.UniversalUser{Nickname: "name", Email: "somemail@mail.com", UserId: 124})
	ad, _ := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "title", Text: "text", UserId: a.UserId})

	updatedAd, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: ad.Id, UserId: ad.AuthorId, Published: true})
	assert.NoError(t, err)
	assert.Equal(t, ad.Title, updatedAd.Title)
	assert.Equal(t, ad.Text, updatedAd.Text)
	assert.Equal(t, ad.AuthorId, updatedAd.AuthorId)
	assert.Equal(t, true, updatedAd.Published)

	add, _ := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "ti", Text: "te", UserId: a.UserId})
	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: ad.Id, UserId: b.UserId, Published: add.Published})
	assert.ErrorIs(t, err, ErrorForbidden)
}

func TestGRPCUpdateAd(t *testing.T) {
	client, ctx := GetTestClient(t)

	a, _ := client.CreateUser(ctx, &grpcPort.UniversalUser{Nickname: "name", Email: "somemail@mail.com", UserId: 123})
	b, _ := client.CreateUser(ctx, &grpcPort.UniversalUser{Nickname: "name1", Email: "1@mail.com", UserId: 5})
	ad, _ := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "title", Text: "text", UserId: a.UserId})

	updatedAd, err := client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{AdId: ad.Id, Title: "title1", Text: "text1", UserId: ad.AuthorId})
	assert.NoError(t, err)
	assert.Equal(t, "title1", updatedAd.Title)
	assert.Equal(t, "text1", updatedAd.Text)
	assert.Equal(t, ad.AuthorId, updatedAd.AuthorId)
	assert.Equal(t, ad.Published, updatedAd.Published)

	_, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{AdId: ad.Id, Title: "new title", Text: "new text", UserId: b.UserId})
	assert.ErrorIs(t, err, ErrorForbidden)

}

func TestGRPCDeleteAd(t *testing.T) {
	client, ctx := GetTestClient(t)

	a, _ := client.CreateUser(ctx, &grpcPort.UniversalUser{Nickname: "name", Email: "somemail@mail.com", UserId: 123})
	b, _ := client.CreateUser(ctx, &grpcPort.UniversalUser{Nickname: "name1", Email: "somemail1@mail.com", UserId: 124})
	ad, _ := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "title", Text: "text", UserId: a.UserId})
	_, err := client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AdId: ad.Id, AuthorId: b.UserId})
	assert.ErrorIs(t, err, ErrorForbidden)
	_, err = client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AdId: ad.Id + 1, AuthorId: ad.AuthorId})
	assert.ErrorIs(t, err, ErrorBadRequest)

	resp, err := client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AdId: ad.Id, AuthorId: ad.AuthorId})
	assert.NoError(t, err)
	assert.Equal(t, ad.Title, resp.Title)
	assert.Equal(t, ad.Text, resp.Text)
	assert.Equal(t, ad.AuthorId, resp.AuthorId)
	assert.Equal(t, ad.Published, resp.Published)
}

func TestGRPCDeleteUserByID(t *testing.T) {
	client, ctx := GetTestClient(t)

	a, _ := client.CreateUser(ctx, &grpcPort.UniversalUser{Nickname: "name", Email: "somemail@mail.com", UserId: 123})

	_, err := client.DeleteUserByID(ctx, &grpcPort.DeleteUserRequest{Id: a.UserId + 1})
	assert.ErrorIs(t, err, ErrorBadRequest)

	resp, err := client.DeleteUserByID(ctx, &grpcPort.DeleteUserRequest{Id: a.UserId})
	assert.NoError(t, err)
	assert.Equal(t, a.Nickname, resp.Nickname)
	assert.Equal(t, a.Email, resp.Email)
	assert.Equal(t, a.UserId, resp.UserId)
}

func TestGRPCDefaultFilter(t *testing.T) {
	client, ctx := GetTestClient(t)

	_, _ = client.CreateUser(ctx, &grpcPort.UniversalUser{Nickname: "name", Email: "somemail@mail.com", UserId: 123})

	resp, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 123, Title: "hello", Text: "world"})
	assert.NoError(t, err)

	publishedAd, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: 123, AdId: resp.Id, Published: true})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 123, Title: "hello", Text: "world"})
	assert.NoError(t, err)

	ads, err := client.ListAds(ctx, &grpcPort.FilterRequest{})
	assert.NoError(t, err)
	assert.Len(t, ads.List, 1)
	assert.Equal(t, ads.List[0].Id, publishedAd.Id)
	assert.Equal(t, ads.List[0].Title, publishedAd.Title)
	assert.Equal(t, ads.List[0].Text, publishedAd.Text)
	assert.Equal(t, ads.List[0].AuthorId, publishedAd.AuthorId)
	assert.True(t, ads.List[0].Published)
}
