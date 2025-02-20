package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/eynopv/lac/pkg/client"
	"github.com/eynopv/lac/pkg/param"
	"github.com/eynopv/lac/pkg/printer"
	"github.com/eynopv/lac/pkg/request"
	"github.com/eynopv/lac/pkg/utils"
	"github.com/eynopv/lac/pkg/validators"
)

func runCommandFunction(
	args []string,
	variables map[string]string,
	headers map[string]string,
	clientConfig *client.ClientConfig,
) {
	req, err := request.LoadRequest(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	runRequest(req, variables, headers, client.NewClient(clientConfig))
}

func runRequest(
	req *request.Request,
	variables map[string]string,
	headers map[string]string,
	client *client.Client,
) {
	resolvedHeaders := map[string]string{}
	for key, value := range headers {
		resolvedHeaders[key] = param.Param(value).Resolve(variables)
	}

	req.Headers = utils.CombineMaps(resolvedHeaders, req.Headers)
	req.ResolveParameters(variables)

	result, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		os.Exit(1)
	}

	if client.PrinterConfig.PrintRequestHeaders {
		request, err := req.ToHttpRequest()
		if err != nil {
			fmt.Printf("failed to print request headers: %v\n", err)
		}

		fmt.Println(request.URL)
		printer.PrintHeaders(request.Header)
		fmt.Println()
	}

	if client.PrinterConfig.PrintRequestBody {
		if req.Body != nil {
			if validators.IsJson(req.Body) {
				var raw any
				if err := json.Unmarshal(req.Body, &raw); err == nil {
					printer.PrintPrettyJson(raw)
				} else {
					fmt.Println("Invalid body")
				}
			} else {
				fmt.Println(result.Text)
			}

			fmt.Println()
		}
	}

	if client.PrinterConfig.PrintResponseHeaders {
		protocol := result.Protocol
		status := result.Status
		statusCode := result.StatusCode
		elapsedTime := result.ElapsedTime

		if statusCode < 300 {
			fmt.Printf("%v %v [%v]\n", protocol, printer.Green(status), elapsedTime)
		} else if statusCode >= 300 && result.StatusCode < 400 {
			fmt.Printf("%v %v [%v]\n", result.Protocol, printer.Cyan(status), elapsedTime)
		} else {
			fmt.Printf("%v %v [%v]\n", result.Protocol, printer.Red(status), elapsedTime)
		}

		printer.PrintHeaders(result.Headers)
		fmt.Println()
	}

	if client.PrinterConfig.PrintResponseBody {
		if result.Body != nil {
			printer.PrintPrettyJson(result.Body)
		} else if result.Text != "" {
			fmt.Println(result.Text)
		}
	}
}
