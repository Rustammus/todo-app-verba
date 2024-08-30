package app

import (
	"ToDoVerba/docs"
	"ToDoVerba/internal/config"
	"ToDoVerba/internal/crud"
	"ToDoVerba/internal/repos"
	"ToDoVerba/internal/route"
	"ToDoVerba/internal/service"
	"ToDoVerba/pkg/logging"
	"ToDoVerba/pkg/migrator"
	"encoding/json"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
)

func Run() {
	// Init config and logger
	logger := logging.GetLogger()
	logger.Info("Start application")
	conf := config.GetConfig(logger)

	loglvl, err := logrus.ParseLevel(conf.Server.LogLevel)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.SetLevel(loglvl)

	ex, _ := os.Executable()
	wd, _ := os.Getwd()

	logger.Debugf("Working directory: %s; Executable: %s", wd, ex)

	configJSON, _ := json.Marshal(conf)
	logger.Debugf("Config: %s", string(configJSON))

	// Run migrations
	RunMigration(conf, logger)

	// Init db connection
	pool := crud.GetPool(conf, logger)
	defer pool.Close()

	// Init repositories, service
	repositories := repos.NewRepositories(pool, logger)

	services := service.NewServices(service.Deps{
		Repos:  repositories,
		Logger: logger,
	})

	// Init router and handlers
	r := httprouter.New()

	if conf.Server.EnableSwag {
		r.Handler("GET", "/swagger/*any", httpSwagger.Handler(
			//todo refactor
			httpSwagger.URL("http://"+conf.Server.Host+":"+conf.Server.Port+"/swagger/doc.json"), //The url pointing to API definition
		))
		docs.SwaggerInfo.Host = conf.Server.Host + ":" + conf.Server.Port
		logger.Infof("Swagger enabled")
	}

	h := route.NewHandler(route.Deps{
		Services: services,
		Logger:   logger,
	})

	h.Init(r)
	logger.Fatal(http.ListenAndServe(":"+conf.Server.Port, r))
}

func RunMigration(conf *config.Config, logger logging.Logger) {
	if len(conf.Storage.Migration) == 0 {
		logger.Info("Migration file env not set in config. Skipping migration")
		return
	}

	m, err := migrator.NewMigrator(migrator.Deps{
		Username: conf.Storage.Username,
		Password: conf.Storage.Password,
		Host:     conf.Storage.Host,
		Port:     conf.Storage.Port,
		Database: conf.Storage.Database,
		Source:   conf.Storage.Migration,
	})
	defer func() {
		if m != nil {
			source, database := m.Close()
			if source != nil {
				logger.Errorf("Error while closing migrator source: %s", source.Error())
			}
			if database != nil {
				logger.Errorf("Error while closing migrator connection: %s", database.Error())
			}
		}
	}()

	if err != nil {
		logger.Fatalf("Error while initializing migrator: %s", err.Error())
	}

	err = m.Up()
	noChange := errors.Is(err, migrate.ErrNoChange)
	if err != nil && !noChange {
		logger.Fatalf("Error while up migrator: %s", err.Error())
	} else if noChange {
		logger.Infof("Database migration already up-to-date")
	} else {
		logger.Infof("Successfully migrated database to last version")
	}
}
