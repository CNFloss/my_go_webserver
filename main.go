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

	"github.com/CNFloss/my_go_webserver/api/data"
	"github.com/CNFloss/my_go_webserver/api/handlers"
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


	hh := handlers.NewHealthz(l)
	rh := handlers.NewRoot(l)
	uh := handlers.NewUsers(usersCache,l)

	sm := http.NewServeMux()
	sm.Handle("/healthz", hh)
	sm.Handle("/", rh)
	sm.Handle("/user", uh)
	const filepathRoot = "."
	sm.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	corsSM := handlers.MiddlewareCors(sm)

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