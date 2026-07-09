package httpclient

type Request struct {
	Method          string
	Url             string
	Headers         map[string]string
	Body            string
	TimeoutMS       int
	FollowRedirects bool
}

type Response struct {
	Status     int
	StatusText string
	Headers    map[string][]string
	Body       string
	Json       any
}
