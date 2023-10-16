package main

import (
	"os"

	"github.com/marsuc/SawitProTest/generated"
	"github.com/marsuc/SawitProTest/handler"
	customValidator "github.com/marsuc/SawitProTest/pkg/validator"
	"github.com/marsuc/SawitProTest/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Validator = customValidator.NewValidator()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})
	opts := handler.NewServerOptions{
		Repository: repo,
	}
	return handler.NewServer(opts)
}
