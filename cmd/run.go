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
	variables map[string]interface{},
	headers map[string]string,
	clientConfig *client.ClientConfig,
) {
	requestTemplate, err := request.NewTemplate(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	requestTemplate = requestTemplate.Interpolate(variables, true)

	req, err := requestTemplate.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	auth, err := request.NewAuth(requestTemplate)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	runRequest(req, variables, headers, client.NewClient(clientConfig), auth)
}

func runRequest(
	req *request.Request,
	variables map[string]interface{},
	headers map[string]string,
	client *client.Client,
	auth request.Auth,
) {
	resolvedHeaders := map[string]request.StringOrStringList{}
	for key, value := range headers {
		resolvedHeaders[key] = []string{param.Param(value).Resolve(variables, true)}
	}

	req.Headers = utils.CombineMaps(resolvedHeaders, req.Headers)

	result, err := client.Do(req, auth)

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
				fmt.Println(string(req.Body))
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
