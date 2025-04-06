package cmd

import (
	"fmt"
	"os"

	"github.com/eynopv/lac/pkg/client"
	"github.com/eynopv/lac/pkg/param"
	"github.com/eynopv/lac/pkg/printer"
	"github.com/eynopv/lac/pkg/request"
	"github.com/eynopv/lac/pkg/request/authentication"
	"github.com/eynopv/lac/pkg/utils"
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

	auth, err := authentication.NewAuth(requestTemplate)
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
	auth authentication.Auth,
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

	prntr := printer.NewPrinter(client.PrinterConfig)
	prntr.Print(result)
}
