package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hydralinkapp/hydralink/api/internal/controller"
	"github.com/hydralinkapp/hydralink/api/pkg/config"
	"github.com/hydralinkapp/hydralink/api/pkg/database"
	"github.com/hydralinkapp/hydralink/api/pkg/monitor"
	"github.com/hydralinkapp/hydralink/api/pkg/web"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// TODO : Refine logger configuration
	log.Print("Initializing logger...")
	logger := monitor.NewLogger()
	defer logger.Flush()

	// TODO : Refine database configuration
	logger.Infof("Initializing database...")
	db, err := database.New(os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer func() {
		logger.Infof("Stopping database...")
		if err := db.Close(); err != nil {
			logger.Errorf("%+v", err)
		}
	}()

	if err = router(&logger, &db).Run(os.Getenv("HOST_ADDRESS")); err != nil {
		logger.Errorf("%+v", err)
	}
	return nil
}

// router hooks all routes to controllers
func router(logger *monitor.Logger, db *database.DB) *gin.Engine {
	r := gin.Default()

	// TODO : Setup middleware

	// This allows frontend localhost:3000 to call backend localhost:4000 bypassing cors
	if config.IsDevEnvironment() {
		r.Use(cors.Default())
	}

	r.Use(web.RequestIDMiddleware())

	controller.Health(r, logger, db)
	controller.Session(r, logger, db)
	controller.User(r, logger, db)

	return r
}
