package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/CNFloss/my_go_webserver/internal/data"
	"github.com/CNFloss/my_go_webserver/internal/handlers"
	"github.com/CNFloss/my_go_webserver/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	envType := os.Getenv("GO_ENV")
	if envType == "" {
		envType = "development" // default to development
	}
	fmt.Println(strings.ToUpper(envType), "ENVIRONMENT\n ")

	envfile, err := os.ReadFile(".env."+envType); 
	if err != nil {
		log.Fatal("Failed to load environment: ", err)
	}

	kvMap := make(map[string]string)
	pairs := strings.Split(string(envfile), "\n")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			kvMap[kv[0]] = kv[1]
			os.Setenv(kv[0], kv[1])
		}
	}
	
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	port := ":80"
	if *debug {
		fmt.Println("   @@@@@@@@@@@@@@@@@@\n       Debug mode\n   @@@@@@@@@@@@@@@@@@ ")
		port = os.Getenv("PORT")
	}

	fmt.Println("\n>>>>>> Gopher Net <<<<<<")
	fmt.Println("\nlistening on port:", port)

	l := log.New(os.Stdout, "Chirpy ", log.LstdFlags)

	usersCache := data.NewCache()
	err = usersCache.Init("users.json", &data.User{})
	if err != nil {
		l.Println("cache intialization failed")
		l.Println(err)
	}

	config := handlers.NewApiConfig(l)

	hh := handlers.NewHealthz(l)
	rh := handlers.NewRoot(l)
	uh := handlers.NewUsers(usersCache,l)

	r := chi.NewRouter()
	apiR := chi.NewRouter()
	adminR := chi.NewRouter()

	r.Handle("/", rh)
	r.Handle("/user", uh)
	apiR.Method(http.MethodGet, "/healthz", hh)
	apiR.Method(http.MethodGet, "/reset", http.HandlerFunc(config.ResetHits))
	adminR.Method(http.MethodGet, "/metrics", config)
	const filepathRoot = "."
	fileServerHandler := http.StripPrefix("/app", config.MiddlewareMetricsInc(http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app/*", fileServerHandler)
	r.Handle("/app", fileServerHandler)

	corsSM := middleware.MiddlewareCors(r)
	r.Mount("/api", apiR)
	r.Mount("/admin", adminR)

	s := &http.Server{
		Addr: port,
		Handler: corsSM,
		IdleTimeout: 120*time.Second,
		ReadTimeout: 1*time.Second,
		WriteTimeout: 1*time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <- sigChan
	l.Println("Recieved terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(tc)
}