package cmd

import (
	"fmt"
	"os"

	"github.com/eynopv/lac/pkg/client"
	"github.com/eynopv/lac/pkg/param"
	"github.com/eynopv/lac/pkg/request"
	"github.com/eynopv/lac/pkg/utils"
)

func runCommandFunction(
	args []string,
	variables map[string]string,
	headers map[string]string,
	timeout int,
) {

	requestClient := client.NewClient(timeout)

		/*
	if strings.HasSuffix(args[0], ".flow.yaml") || strings.HasSuffix(args[0], ".flow.json") {
		flow, err := flow.LoadFlow(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, item := range flow.Items {
			req, err := request.LoadRequest(item.Uses)
		}
	}
		*/


	req, err := request.LoadRequest(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	runRequest(req, variables, headers, requestClient)
}

func runRequest(
	req *request.Request,
	variables map[string]string,
	headers map[string]string,
	client client.Client,
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

	result.Print()
}
