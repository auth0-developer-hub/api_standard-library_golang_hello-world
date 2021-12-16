package cmd

import (
	"hello-golang-api/common"
	"hello-golang-api/server"
	"hello-golang-api/store/sqlite"
	"net"
	"net/http"
	_ "net/http/pprof" // Import for pprof
	"path"

	"go.uber.org/zap"
)

// Execute starts the program
func Execute(settings map[string]string, environment []byte) {
	common.InitLogger(settings)
	logger := zap.S().With("package", "cmd")
	cfg, err := common.InitConfig(environment)
	if err != nil {
		logger.Errorf("initializing config error: %v", err)
	}

	// Database
	dbstore, err := sqlite.New(settings, path.Join("store", "sqlite"))
	if err != nil {
		logger.Fatalw("Database error", "error", err)
	}
	defer dbstore.Close()

	if settings["profiler.enabled"] == "true" {
		hostPort := net.JoinHostPort(settings["profiler.host"], settings["profiler.port"])
		go func() {
			if err := http.ListenAndServe(hostPort, nil); err != nil {
				logger.Errorf("profiler server error: %v", err)
			}
		}()
		logger.Infof("profiler enabled on http://%s", hostPort)
	}

	server, err := server.New(cfg, settings, dbstore, logger)
	if err != nil {
		logger.Fatalw(err.Error())
	}

	httpServer := &http.Server{
		Addr:    settings["server.host"] + ":" + settings["server.port"],
		Handler: server.Handler(),
	}

	logger.Infof("API server listening on %s", httpServer.Addr)
	err = httpServer.ListenAndServe()
	if err != nil {
		logger.Fatalw(err.Error())
	}

}
