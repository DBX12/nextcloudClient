package nextcloudClient

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

const simpleResponseOk = `<ocs><meta><status>ok</status><statuscode>100</statuscode><message/></meta><data/></ocs>`

func TestGetGroups(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://example.local/ocs/v1.php/cloud/groups",
		httpmock.NewStringResponder(200, "<?xml version=\"1.0\"?><ocs><meta><statuscode>100</statuscode><status>ok</status></meta><data><groups><element>admin</element><element>testGroup</element></groups></data></ocs>"),
	)
	c := NewClient("http://example.local", "the-user", "the-secret-password")
	groups, err := c.GetGroups()
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{"admin", "testGroup"}
	if !reflect.DeepEqual(groups, expected) {
		t.Fatal("Expectation not met")
	}
}

func TestClient_GetGroup(t *testing.T) {
	testOptions := DefaultTestOptions()
	type args struct {
		groupId string
	}
	tests := []struct {
		name         string
		clientData   clientData
		args         args
		statusCode   int
		responseBody string
		testOptions  RequestTestOptions
		want         string
		wantErr      bool
	}{
		{
			name:         "Successful request",
			clientData:   goodClient,
			args:         args{groupId: "admin"},
			statusCode:   200,
			responseBody: `<?xml version="1.0"?><ocs><meta><status>ok</status><statuscode>100</statuscode><message>OK</message><totalitems></totalitems><itemsperpage></itemsperpage></meta><data><groups><element>admin</element></groups></data></ocs>`,
			testOptions:  testOptions,
			want:         "admin",
			wantErr:      false,
		},
		{
			name:         "No group with that id",
			clientData:   goodClient,
			args:         args{groupId: "unknownGroup"},
			statusCode:   200,
			responseBody: `<?xml version="1.0"?><ocs><meta><status>ok</status><statuscode>100</statuscode><message>OK</message><totalitems></totalitems><itemsperpage></itemsperpage></meta><data><groups/></data></ocs>`,
			testOptions:  testOptions,
			want:         "",
			wantErr:      true,
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetResponder(fmt.Sprintf("%s/ocs/v1.php/cloud/groups?search=%s", HOST, tt.args.groupId), tt.statusCode, tt.responseBody, tt.testOptions)
			c := &Client{
				HostURL:    tt.clientData.HostURL,
				HTTPClient: tt.clientData.HTTPClient,
				username:   tt.clientData.username,
				password:   tt.clientData.password,
			}
			got, err := c.GetGroup(tt.args.groupId)
			CheckForResponderError(t, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGroup() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateGroup(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://example.local/ocs/v1.php/cloud/groups",
		func(req *http.Request) (*http.Response, error) {
			bodyContents, _ := ioutil.ReadAll(req.Body)
			expectedContent := []byte("groupid=testGroup01")
			if !reflect.DeepEqual(bodyContents, expectedContent) {
				t.Fatal("Request body mismatch")
			}
			response := httpmock.NewStringResponse(200, simpleResponseOk)
			return response, nil
		},
	)
	c := NewClient("http://example.local", "the-user", "the-secret-password")
	success, err := c.CreateGroup("testGroup01")
	if err != nil {
		t.Fatal(err)
	}
	if success != true {
		t.Fatal("Method returned false")
	}
}

func TestDeleteGroup(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", "http://example.local/ocs/v1.php/cloud/groups/testGroup01",
		httpmock.NewStringResponder(200, simpleResponseOk),
	)
	c := NewClient("http://example.local", "the-user", "the-secret-password")
	success, err := c.DeleteGroup("testGroup01")
	if err != nil {
		t.Fatal(err)
	}
	if success != true {
		t.Fatal("Method returned false")
	}
}

func TestGetGroupMembers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://example.local/ocs/v1.php/cloud/groups/testGroup",
		httpmock.NewStringResponder(200, "<?xml version=\"1.0\"?><ocs><meta><statuscode>100</statuscode><status>ok</status></meta><data><users><element>Frank</element><element>Jane</element></users></data></ocs>"),
	)
	c := NewClient("http://example.local", "the-user", "the-secret-password")
	groupMembers, err := c.GetGroupMembers("testGroup")
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{"Frank", "Jane"}
	if !reflect.DeepEqual(groupMembers, expected) {
		t.Fatal("Expectation not met")
	}
}

func TestGetGroupSubadmins(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://example.local/ocs/v1.php/cloud/groups/testGroup/subadmins",
		httpmock.NewStringResponder(200, "<?xml version=\"1.0\"?><ocs><meta><status>ok</status><statuscode>100</statuscode><message/></meta><data><element>Tom</element></data></ocs>"),
	)
	c := NewClient("http://example.local", "the-user", "the-secret-password")
	groupSubadmins, err := c.GetGroupSubadmins("testGroup")
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{"Tom"}
	if !reflect.DeepEqual(groupSubadmins, expected) {
		t.Fatal("Expectation not met")
	}
}
