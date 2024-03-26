package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/eynopv/gorcli/internal"
)

func main() {
	showHeaders := flag.Bool("sh", false, "Show response headers")
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Usage: gorcli <request> [flags]")
		os.Exit(1)
	}

	internal.LoadDotEnv()

	config, err := internal.LoadConfig()
	if err != nil {
		fmt.Println("Unable to load config")
		os.Exit(1)
	}

	requestName := args[0]
	request, err := internal.NewRequest(requestName)
	if err != nil {
		fmt.Println("Unable to make request:", err)
		return
	}

	for key, value := range config.Headers {
		request.Header.Set(key, value)
	}

	start := time.Now()
	client := http.Client{}
	res, err := client.Do(request)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer res.Body.Close()

	fmt.Println("Status:", res.Status)
	fmt.Println("Elapsed Time:", elapsed)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Failed to read resonse body:", err)
		return
	}

	if *showHeaders {
		prettyHeaders, err := json.MarshalIndent(res.Header, "", " ")
		if err != nil {
			fmt.Println("Failed to parse headers", err)
			return
		}
		fmt.Println(string(prettyHeaders))
	}

	if len(body) == 0 {
		fmt.Println("Response body is empty")
		return
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("Failed to parse JSON:", err)
	}

	prettyJSON, err := json.MarshalIndent(responseData, "", "  ")
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return
	}

	fmt.Println(string(prettyJSON))
}
