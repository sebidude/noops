package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/gin-gonic/gin"
	"github.com/sebidude/configparser"
)

var (
	builddate  string
	revision   string
	version    string
	appConfig  *AppConfig
	configfile string
	wg         sync.WaitGroup
	hostname   string
)

type AppConfig struct {
	ListenAddress string `json:"listenAddress"`
	LogFile       string `json:"logfile"`
}

func main() {

	kingpin.UsageTemplate(kingpin.CompactUsageTemplate)
	kingpin.Flag("config", "Full path to the configfile to be used.").Required().Short('c').StringVar(&configfile)
	kingpin.Parse()

	log.SetFlags(log.Lshortfile)

	// parse the configfile if it is set
	err := configparser.ParseYaml(configfile, &appConfig)
	if err != nil {
		log.Fatalf("Error reading configfile: %s", err.Error())
	}
	configparser.SetValuesFromEnvironment("NOOPS", appConfig)

	hostname, err = os.Hostname()
	if err != nil {
		log.Fatalln("Cannot get hostname!")
	}

	logger := new(Logger)
	if len(appConfig.LogFile) > 0 {
		logfile, err := os.OpenFile(appConfig.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
		if err != nil {
			log.Fatalf("error opening logfile: %v", err)
		}
		defer logfile.Close()
		logger.LogFile = logfile

	} else {
		appConfig.LogFile = "stdout"
	}
	log.SetOutput(logger)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	defaultRouter := router.Group("/")
	aliveRouter := router.Group("/alive")

	defaultRouter.Use(GinLogger())

	log.Println("=====  NoopsServer  =====")
	log.Printf("Builddate: %s", builddate)
	log.Printf("Version  : %s", version)
	log.Printf("Revision : %s", revision)
	log.Println(" ---")
	log.Printf("Hostname: %s", hostname)
	log.Printf("ListenAddress: %s", appConfig.ListenAddress)

	defaultRouter.GET("/", defaultGet)
	defaultRouter.GET("/say/:something", defaultGetSomething)

	aliveRouter.GET("", Alive)
	aliveRouter.GET("/", Alive)

	httpsrv := &http.Server{
		Addr:    appConfig.ListenAddress,
		Handler: router,
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Fatal(httpsrv.ListenAndServe())
	}()
	log.Println("HTTP server started.")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	log.Println("Shutdown signals registered.")
	<-signalChan
	log.Println("Shutdown signal received, exiting.")
	httpsrv.Shutdown(context.Background())
	wg.Wait()
	log.Println("Server exiting")
}

func defaultGet(c *gin.Context) {
	c.String(200, "Hello from %s", hostname)
}

func defaultGetSomething(c *gin.Context) {
	something := c.Param("something")
	c.String(200, "%s says %s", hostname, something)
}

func Alive(c *gin.Context) {
	c.String(200, "alive")
}
