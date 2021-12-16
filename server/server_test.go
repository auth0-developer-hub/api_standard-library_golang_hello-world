package server

import (
	"context"
	"hello-golang-api/common"
	"hello-golang-api/entities"
	"hello-golang-api/store"
	"hello-golang-api/store/mock"
	"hello-golang-api/store/sqlite"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/google/uuid"
	"github.com/snowzach/queryp"

	"go.uber.org/zap"
)

var token *Token
var server *Server
var logger *zap.SugaredLogger

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	//test server setup
	cfg, err := common.InitConfig([]byte(`auth0-audience: https://api.example.com
auth0-domain: dev-qzkixsmr.eu.auth0.com
auth0-clientID: 8O50umOS1OiLfDJaY6ANhcOzQdTEIfS0
auth0-clientSecret: qZMCJfF8mztTSiKQy6lAEyCtecWX1zaJJHuvp_zDjxa9w_eXiSvLd0EB8oIbJ9Da
auth0-callbackURL: http://localhost:4040/callback`))
	if err != nil {
		log.Fatalf("initializing config error: %v", err)
	}

	dbstore := &mock.Client{}
	//common.InitLogger(settings)
	logger = zap.S().With("package", "test")
	server, err = New(cfg, settings, dbstore, logger)
	if err != nil {
		log.Fatal(err)
	}

	token, err = GetToken(cfg)
	if err != nil {
		log.Fatal(err)
	}

}

func shutdown() {

}

func updateServerStore(store store.MessageStore) {
	server.store = store
}
func TestRoutesWithoutToken(t *testing.T) {

	mockStore := &mock.Client{}
	updateServerStore(mockStore)
	srv := httptest.NewServer(server.Handler())
	defer srv.Close()

	// Make request and validate we get back proper response
	e := httpexpect.New(t, srv.URL)
	e.GET("/api/messages/public").Expect().Status(http.StatusOK)
	e.GET("/api/messages/protected").Expect().Status(http.StatusUnauthorized)
	e.GET("/api/messages/admin").Expect().Status(http.StatusUnauthorized)

	uuid, err := uuid.NewUUID()
	if err != nil {
		t.Error(err)
	}

	//mock and hit
	mockStore.MessageSavefn = func(ctx context.Context, user *entities.Message) error {
		return nil
	}
	e.POST("/api/messages").WithJSON(&entities.Message{
		Text: "hello world",
	}).Expect().Status(http.StatusUnauthorized)

	//mock and hit
	mockStore.MessagesListfn = func(ctx context.Context, qp *queryp.QueryParameters) ([]*entities.Message, int64, error) {
		return []*entities.Message{
			{
				Id:   uuid.String(),
				Text: "hello world",
				Date: time.Now(),
			},
		}, 1, nil
	}
	e.GET("/api/messages").Expect().Status(http.StatusOK).JSON().Object().Value("count").Equal(1)

	//mock and hit
	mockStore.MessageGetByIDfn = func(ctx context.Context, id string) (*entities.Message, error) {
		return &entities.Message{
			Id:   uuid.String(),
			Text: "hello world",
			Date: time.Now(),
		}, nil
	}
	e.GET("/api/messages/" + uuid.String()).Expect().Status(http.StatusOK).JSON().Object().Value("id").Equal(uuid.String())

	//mock and hit
	mockStore.MessageUpdatefn = func(ctx context.Context, user *entities.Message) error {
		return nil
	}
	e.PUT("/api/messages/" + uuid.String()).WithJSON(&entities.Message{
		Text: "hi world",
	}).Expect().Status(http.StatusUnauthorized)

	//mock and hit
	mockStore.MessageDeleteByIDfn = func(ctx context.Context, id string) error {
		return nil
	}
	e.DELETE("/api/messages/" + uuid.String()).Expect().Status(http.StatusUnauthorized)
}

func TestRoutesWithToken(t *testing.T) {
	mockStore := &mock.Client{}
	updateServerStore(mockStore)
	srv := httptest.NewServer(server.Handler())
	defer srv.Close()

	// Make request and validate we get back proper response
	e := httpexpect.New(t, srv.URL)
	e.GET("/").Expect().Status(http.StatusNotFound)
	e.GET("/ping").Expect().Status(http.StatusOK).Text().Equal("pong")
	e.GET("/version").Expect().Status(http.StatusOK).JSON().Object().Value("version").Equal(common.GitVersion)

	e.GET("/api/messages/public").WithHeader("Authorization", "Bearer "+token.AccessToken).Expect().Status(http.StatusOK)
	e.GET("/api/messages/protected").WithHeader("Authorization", "Bearer "+token.AccessToken).Expect().Status(http.StatusOK)
	e.GET("/api/messages/admin").WithHeader("Authorization", "Bearer "+token.AccessToken).Expect().Status(http.StatusForbidden)

	uuid, err := uuid.NewUUID()
	if err != nil {
		t.Error(err)
	}

	//mock and hit
	mockStore.MessageSavefn = func(ctx context.Context, user *entities.Message) error {
		return nil
	}
	e.POST("/api/messages").WithJSON(&entities.Message{
		Text: "hello world",
	}).WithHeader("Authorization", "Bearer "+token.AccessToken).Expect().Status(http.StatusOK)

	//mock and hit
	mockStore.MessagesListfn = func(ctx context.Context, qp *queryp.QueryParameters) ([]*entities.Message, int64, error) {
		return []*entities.Message{
			{
				Id:   uuid.String(),
				Text: "hello world",
				Date: time.Now(),
			},
		}, 1, nil
	}
	e.GET("/api/messages").Expect().Status(http.StatusOK).JSON().Object().Value("count").Equal(1)

	//mock and hit
	mockStore.MessageGetByIDfn = func(ctx context.Context, id string) (*entities.Message, error) {
		return &entities.Message{
			Id:   uuid.String(),
			Text: "hello world",
			Date: time.Now(),
		}, nil
	}
	e.GET("/api/messages/" + uuid.String()).Expect().Status(http.StatusOK).JSON().Object().Value("id").Equal(uuid.String())

	//mock and hit
	mockStore.MessageUpdatefn = func(ctx context.Context, user *entities.Message) error {
		return nil
	}
	e.PUT("/api/messages/"+uuid.String()).WithJSON(&entities.Message{
		Text: "hi world",
	}).WithHeader("Authorization", "Bearer "+token.AccessToken).Expect().Status(http.StatusOK)

	//mock and hit
	mockStore.MessageDeleteByIDfn = func(ctx context.Context, id string) error {
		return nil
	}
	e.DELETE("/api/messages/"+uuid.String()).WithHeader("Authorization", "Bearer "+token.AccessToken).Expect().Status(http.StatusForbidden)
}

func TestIntergation(t *testing.T) {
	// Sqlite Database
	dbstore, err := sqlite.New(settings, path.Join("../", "store", "sqlite"))
	if err != nil {
		logger.Fatalw("Database error", "error", err)
	}
	defer dbstore.Close()

	updateServerStore(dbstore)
	srv := httptest.NewServer(server.Handler())
	defer srv.Close()

	// Make request and validate we get back proper response
	e := httpexpect.New(t, srv.URL)
	e.GET("/api/messages/public").WithHeader("Authorization", "Bearer "+token.AccessToken).Expect().Status(http.StatusOK)
	e.GET("/api/messages/protected").WithHeader("Authorization", "Bearer "+token.AccessToken).Expect().Status(http.StatusOK)
	e.GET("/api/messages/admin").WithHeader("Authorization", "Bearer "+token.AccessToken).Expect().Status(http.StatusForbidden)

	//post msg
	createdMsg := e.POST("/api/messages").WithJSON(&entities.Message{
		Text: "hello world",
	}).WithHeader("Authorization", "Bearer "+token.AccessToken).Expect().Status(http.StatusOK).JSON().Object()

	createdID1 := createdMsg.Value("id").String().Raw()
	//list msgs
	e.GET("/api/messages").Expect().Status(http.StatusOK).JSON().Object().Value("count").Equal(1)

	//get msg by id
	rId := e.GET("/api/messages/" + createdID1).Expect().Status(http.StatusOK).JSON()
	rId.Object().Value("id").Equal(createdID1)

	//put msg by id
	e.PUT("/api/messages/"+createdID1).WithJSON(&entities.Message{
		Text: "hi world",
	}).WithHeader("Authorization", "Bearer "+token.AccessToken).Expect().Status(http.StatusOK).JSON().Object().Value("id").Equal(createdID1)

	//get msg by id to check updated text
	e.GET("/api/messages/" + createdID1).Expect().Status(http.StatusOK).JSON().Object().Value("text").Equal("hi world")

	//delete msg by id
	e.DELETE("/api/messages/"+createdID1).WithHeader("Authorization", "Bearer "+token.AccessToken).Expect().Status(http.StatusForbidden)

}
