package output

import (
	"fmt"

	"github.com/eynopv/lac/internal/httpclient"
)

func Print(r *httpclient.Response) {
	fmt.Println("Status: ", r.Status)
	fmt.Println(r.Body)
}
