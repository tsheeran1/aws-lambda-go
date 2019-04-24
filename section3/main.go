package main

import "github.com/aws/aws-lambda-go/lambda"

type event struct {
	Question string
}

type response struct {
	Question string
	Answer   string
}

func handler(e event) (response, error) {
	return response{
		Question: e.Question,
		Answer:   "I don't know." + e.Question,
	}, nil
}

func main() {
	lambda.Start(handler)

}
