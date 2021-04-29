package nextcloudClient

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

const badLoginResponse = `<?xmlversion="1.0"?><ocs><meta><status>failure</status><statuscode>997</statuscode><message>Currentuserisnotloggedin</message><totalitems></totalitems><itemsperpage></itemsperpage></meta><data/></ocs>`
const responderErrorTag = "!Responder Error!"
const USER = "the-user"
const PASS = "the-secret-password"

type RequestTestOptions struct {
	ignoreAuthenticationTest bool
	ignoreHeaderTest         bool
	ignoreBodyTest           bool
	noResetBeforeRegister    bool
	username                 string
	password                 string
}

func DefaultTestOptions() RequestTestOptions {
	return RequestTestOptions{
		ignoreAuthenticationTest: false,
		ignoreHeaderTest:         false,
		ignoreBodyTest:           false,
		noResetBeforeRegister:    false,
		username:                 USER,
		password:                 PASS,
	}
}

func PostResponder(url string, expectedBody string, statusCode int, responseBody string, options RequestTestOptions) {
	GenericResponder("POST", url, expectedBody, statusCode, responseBody, options)
}

func PutResponder(url string, expectedBody string, statusCode int, responseBody string, options RequestTestOptions) {
	GenericResponder("PUT", url, expectedBody, statusCode, responseBody, options)
}

func GenericResponder(method string, url string, expectedBody string, statusCode int, responseBody string, options RequestTestOptions) {
	if options.noResetBeforeRegister == false {
		httpmock.Reset()
	}
	httpmock.RegisterResponder(method, url,
		func(request *http.Request) (*http.Response, error) {
			if options.ignoreHeaderTest == false {
				headerValue := request.Header.Get("OCS-APIRequest")
				headerValue = strings.ToLower(headerValue)
				if !reflect.DeepEqual(headerValue, "true") {
					return nil, errors.New(
						fmt.Sprintf(
							"%s, Bad OCS-APIRequest header value, want %s got %s",
							responderErrorTag,
							"true",
							headerValue,
						),
					)
				}
			}
			if options.ignoreAuthenticationTest == false {
				headerValue := request.Header.Get("Authorization")
				expectedHeaderValue := fmt.Sprintf("Basic %s",
					base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", options.username, options.password))),
				)
				if !reflect.DeepEqual(headerValue, expectedHeaderValue) {
					return httpmock.NewStringResponse(401, badLoginResponse), nil
				}
			}
			if options.ignoreBodyTest == false {
				if request.Body == nil {
					return nil, errors.New(fmt.Sprintf("%s, Body is empty, want %s got nil", responderErrorTag, expectedBody))
				}
				actualBody, _ := ioutil.ReadAll(request.Body)
				if !reflect.DeepEqual(actualBody, []byte(expectedBody)) {
					return nil, errors.New(fmt.Sprintf("%s, Body mismatched, want %s got %s", responderErrorTag, expectedBody, actualBody))
				}
			}
			return httpmock.NewStringResponse(statusCode, responseBody), nil
		})
}

func GetResponder(url string, statusCode int, responseBody string, options RequestTestOptions) {
	httpmock.RegisterResponder("GET", url,
		func(request *http.Request) (*http.Response, error) {
			if options.ignoreHeaderTest == false {
				headerValue := request.Header.Get("OCS-APIRequest")
				headerValue = strings.ToLower(headerValue)
				if !reflect.DeepEqual(headerValue, "true") {
					return nil, errors.New(
						fmt.Sprintf(
							"%s, Bad OCS-APIRequest header value, want %s got %s",
							responderErrorTag,
							"true",
							headerValue,
						),
					)
				}
			}
			if options.ignoreAuthenticationTest == false {
				headerValue := request.Header.Get("Authorization")
				expectedHeaderValue := fmt.Sprintf("Basic %s",
					base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", options.username, options.password))),
				)
				if !reflect.DeepEqual(headerValue, expectedHeaderValue) {
					return httpmock.NewStringResponse(401, badLoginResponse), nil
				}
			}
			return httpmock.NewStringResponse(statusCode, responseBody), nil
		})
}

func CheckForResponderError(t *testing.T, err error) {
	if err == nil {
		return
	}
	errMsg := err.Error()
	if strings.Contains(errMsg, responderErrorTag) {
		t.Fatal(errMsg)
	}
}
