package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostView(t *testing.T) {

	tests := []struct {
		description   string
		route         string
		expectedError bool
		expectedCode  int
		expectedBody  string
	}{
		{
			description:   "posts index route",
			route:         "/posts",
			expectedError: false,
			expectedCode:  200,
			// expectedBody:  `{"message":"record not found"}`,
		},
		{
			description:   "non existing route",
			route:         "/posts/999",
			expectedError: false,
			expectedCode:  400,
			// expectedBody:  `{"message":"record not found"}`,
		},
	}

	// Setup the app as it is done in the main function
	app := Setup()

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route
		// from the test case
		req, _ := http.NewRequest(
			"GET",
			test.route,
			nil,
		)

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)
		fmt.Println(string(body))
		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		// Verify, that the reponse body equals the expected body
		// assert.Equalf(t, test.expectedBody, string(body), test.description)
	}
}
