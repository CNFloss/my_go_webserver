package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/CNFloss/GopherNet/pkg/middleware"
)

func UnmarshalBytes(src []byte) (map[string]string, error) {
	out := make(map[string]string)
	err := parseBytes(src, out)

	return out, err
}

func Parse(r io.Reader) (map[string]string, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}

	return UnmarshalBytes(buf.Bytes())
}

func readFile(filename string) (envMap map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	return Parse(file)
}

func main() {
	envType := os.Getenv("GO_ENV")
	if envType == "" {
		envType = "development" // default to development
	}
	fmt.Println(strings.ToUpper(envType), " ENVIRONMENT\n ")

	var envfile *os.File
	envfile, err := os.Open(".env."+envType); 
	if err != nil {
		log.Fatal("Failed to load environment: ", err)
	}
	defer envfile.Close()

	envMap, err := Parse(envfile)	
	if err != nil {
		log.Fatal("Failed to load environment: ", err)
	}

	currentEnv := map[string]bool{}

	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range envMap {
		if !currentEnv[key] {
			err = os.Setenv(key, value)
			if err != nil {
				log.Fatal("Failed to load environment: ", err)
			}
		}
	}
	
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	port := ":80"
	if *debug {
		fmt.Println("   @@@@@@@@@@@@@@@@@@\n       Debug mode\n   @@@@@@@@@@@@@@@@@@\n ")
		port = os.Getenv("PORT")
	}

	fmt.Println(">>>>>> Gopher Net <<<<<<")
	fmt.Println("\nlistening on port:", port[1:])
	http.HandleFunc("/", middleware.Sayhelloname)
	http.HandleFunc("/login", middleware.Login)
	http.HandleFunc("/upload", middleware.Upload)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}