package handlers

import (
	"context"
	"favorite_service/internal/models"
	"favorite_service/internal/repositories"
	"favorite_service/internal/services"
	"favorite_service/pkg/psql"
	"fmt"
	"os"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/gofiber/fiber/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

var app *fiber.App = fiber.New()

var ctx context.Context = context.Background()

type TestDB struct {
	DB        *gorm.DB
	Container testcontainers.Container
}

type HandlerSetup struct {
	DB                *gorm.DB
	App               *fiber.App
	MockUserClient    *MockUserClient
	MockProductClient *MockProductClient
}

func (h *HandlerSetup) SetupTestItemHandler() {
	itemRepository := repositories.NewFavoriteItemRepository(h.DB)
	listRepository := repositories.NewFavoriteListRepository(h.DB)
	itemService := services.NewFavoriItemService(itemRepository, listRepository, h.MockProductClient, h.MockUserClient)
	itemHandler := NewFavoriteItemHandler(itemService)
	itemHandler.SetRoutes(h.App)
}

func (h *HandlerSetup) SetupListHandler() {
	listRepository := repositories.NewFavoriteListRepository(h.DB)
	itemRepository := repositories.NewFavoriteItemRepository(h.DB)
	favoriteListService := services.NewFavoriteListService(listRepository, itemRepository, h.MockProductClient, h.MockUserClient)
	favoriteListHandler := NewFavoriteListHandler(favoriteListService)
	favoriteListHandler.SetRoutes(h.App)
}

func (t *TestDB) Setup() error {
	//ctx := context.Background()

	dbConfig := map[string]string{
		"POSTGRES_USER":     "user",
		"POSTGRES_PASSWORD": "password",
		"POSTGRES_DB":       "users",
	}

	defaultPort := nat.Port("5432/tcp")
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:latest",
			ExposedPorts: []string{defaultPort.Port()},
			Env:          dbConfig,
			WaitingFor: wait.ForAll(
				wait.ForLog("database system is ready to accept connections"),
				wait.ForListeningPort(defaultPort),
			),
		},
		Started: true,
	})
	if err != nil {
		return err
	}
	t.Container = container
	port, err := container.MappedPort(ctx, defaultPort)
	if err != nil {
		return err
	}
	fmt.Println("Veritabanı başladı port numarasi:", port)
	t.DB = psql.Connect("0.0.0.0", dbConfig["POSTGRES_USER"], dbConfig["POSTGRES_PASSWORD"], dbConfig["POSTGRES_DB"], port.Port())
	return t.loadSQLFiles()
}

func (t *TestDB) loadSQLFiles() error {
	fileCreate, err := os.ReadFile("../../psql/create_tables.sql")
	if err != nil {
		return err
	}
	if err := t.DB.Exec(string(fileCreate)).Error; err != nil {
		return err
	}

	fileFill, err := os.ReadFile("../../psql/fill_tables.sql")
	if err != nil {
		return err
	}
	return t.DB.Exec(string(fileFill)).Error
}

func (t *TestDB) CleanUp() {
	t.Container.Terminate(context.Background())
}

type MockProductClient struct{}

func (m *MockProductClient) VerifyProduct(ctx context.Context, productId int) (*models.Product, error) {

	if productId == 1 {
		return &models.Product{
			ID:    1,
			Name:  "Telefon",
			Price: 1599,
			Stock: 50}, nil
	} else if productId == 99 {
		return &models.Product{}, models.ErrRecordNotFound
	}
	return &models.Product{}, nil

}

type MockUserClient struct{}

func (m *MockUserClient) VerifyUser(token string, ctx context.Context) (*models.Users, error) {

	if token == "1" {
		return &models.Users{
			ID:       1,
			Email:    "ahmet.yilmaz@example.com",
			Password: "sifre123",
			Isim:     "Ahmet",
			Soyisim:  "Yılmaz",
			Resim:    "https://example.com/resim1.jpg"}, nil
	} else if token == "99" {
		return &models.Users{}, models.ErrUserUnauthorized
	}

	return &models.Users{}, models.ErrUserUnauthorized

}

func TestMain(m *testing.M) {
	testDB := &TestDB{}
	if err := testDB.Setup(); err != nil {
		fmt.Println("Veritabanı bağlantısı başarısız", err)
		os.Exit(1)
	}
	defer testDB.CleanUp()

	handlerSetup := &HandlerSetup{
		DB:             testDB.DB,
		App:            app,
		MockUserClient: &MockUserClient{},
	}

	handlerSetup.SetupTestItemHandler()
	handlerSetup.SetupListHandler()

	os.Exit(m.Run())
}
