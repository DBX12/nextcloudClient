package nextcloudClient

import (
	"encoding/xml"
	"fmt"
	"github.com/jarcoal/httpmock"
	"net/http"
	"reflect"
	"testing"
	"time"
)

const HOST = "http://example.local"

var hostUrl = fmt.Sprintf("%s/ocs/v1.php", HOST)

type clientData struct {
	HostURL    string
	HTTPClient *http.Client
	username   string
	password   string
}

// creates a client with valid credentials
var goodClient = clientData{
	HostURL:    hostUrl,
	HTTPClient: &http.Client{Timeout: 10 * time.Second},
	username:   USER,
	password:   PASS,
}

// client with bad credentials
var badClient = clientData{
	HostURL:    hostUrl,
	HTTPClient: &http.Client{Timeout: 10 * time.Second},
	username:   "bad-user",
	password:   "bad-pass",
}

func TestValidateUserData(t *testing.T) {
	type args struct {
		userData *UserData
	}
	tests := []struct {
		name     string
		args     args
		result   bool
		problems []string
	}{
		{
			name: "Missing userId",
			args: args{userData: &UserData{
				Password: "secret",
				Email:    "test@example.local",
			}},
			result:   false,
			problems: []string{"UserId must not be empty"},
		},
		{
			name: "Missing password and email",
			args: args{userData: &UserData{
				UserId: "john.doe",
			}},
			result:   false,
			problems: []string{"Either Password or Email must be set"},
		},
		{
			name: "Correct user data",
			args: args{userData: &UserData{
				UserId: "john.doe",
				Email:  "john.doe@example.local",
			}},
			result:   true,
			problems: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.args.userData.Validate()
			if got != tt.result {
				t.Errorf("ValidateUserData() got = %v, result %v", got, tt.result)
			}
			if !reflect.DeepEqual(got1, tt.problems) {
				t.Errorf("ValidateUserData() got1 = %v, result %v", got1, tt.problems)
			}
		})
	}
}

func TestClient_CreateUser(t *testing.T) {
	type args struct {
		userData *UserData
	}
	var tests = []struct {
		name         string
		clientData   clientData
		args         args
		expectedBody string
		statusCode   int
		response     string
		testOptions  RequestTestOptions
		want         bool
		wantErr      bool
	}{
		{
			name:       "Successful creation",
			clientData: goodClient,
			args: args{userData: &UserData{
				UserId:           "john.doe",
				Password:         "johnsPassword",
				DisplayName:      "John Doe",
				Email:            "john.doe@example.local",
				GroupIds:         []string{"employees", "development"},
				SubadminGroupIds: []string{"employees", "development", "accounting"},
				Quota:            QuotaUnlimited,
				Language:         "en",
			}},
			expectedBody: "displayName=John+Doe&email=john.doe%40example.local&groups%5B%5D=employees&groups%5B%5D=development&language=en&password=johnsPassword&quota=none&subadmin%5B%5D=employees&subadmin%5B%5D=development&subadmin%5B%5D=accounting&userid=john.doe",
			statusCode:   200,
			response:     simpleResponseOk,
			testOptions:  DefaultTestOptions(),
			want:         true,
			wantErr:      false,
		},
		{
			name:       "Missing userId",
			clientData: goodClient,
			args: args{userData: &UserData{
				Email: "john.doe@example.local",
			}},
			expectedBody: "email=john.doe%40example.local",
			statusCode:   0,
			response:     "",
			testOptions:  DefaultTestOptions(),
			want:         false,
			wantErr:      true,
		},
		{
			name:       "Missing email and password",
			clientData: goodClient,
			args: args{userData: &UserData{
				UserId: "john.doe",
			}},
			expectedBody: "userid=john.doe",
			statusCode:   0,
			response:     "",
			testOptions:  DefaultTestOptions(),
			want:         false,
			wantErr:      true,
		},
		{
			name:       "Bad credentials",
			clientData: badClient,
			args: args{userData: &UserData{
				UserId: "john.doe",
				Email:  "john.doe@example.local",
			}},
			expectedBody: "email=john.doe%40example.local&userid=john.doe",
			statusCode:   0,
			response:     "",
			testOptions:  DefaultTestOptions(),
			want:         false,
			wantErr:      true,
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PostResponder(fmt.Sprintf("%s/ocs/v1.php/cloud/users", HOST), tt.expectedBody, tt.statusCode, tt.response, tt.testOptions)

			c := &Client{
				HostURL:    tt.clientData.HostURL,
				HTTPClient: tt.clientData.HTTPClient,
				username:   tt.clientData.username,
				password:   tt.clientData.password,
			}
			got, err := c.CreateUser(tt.args.userData)
			CheckForResponderError(t, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetUsers(t *testing.T) {
	tests := []struct {
		name         string
		clientData   clientData
		statusCode   int
		responseBody string
		testOptions  RequestTestOptions
		want         []string
		wantErr      bool
	}{
		{
			name:         "Successful request",
			clientData:   goodClient,
			statusCode:   200,
			responseBody: "<?xmlversion=\"1.0\"?><ocs><meta><status>ok</status><statuscode>100</statuscode><message>OK</message><totalitems></totalitems><itemsperpage></itemsperpage></meta><data><users><element>john.doe</element><element>jane.doe</element></users></data></ocs>",
			testOptions:  DefaultTestOptions(),
			want:         []string{"john.doe", "jane.doe"},
			wantErr:      false,
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetResponder(fmt.Sprintf("%s/ocs/v1.php/cloud/users", HOST), tt.statusCode, tt.responseBody, tt.testOptions)
			c := &Client{
				HostURL:    tt.clientData.HostURL,
				HTTPClient: tt.clientData.HTTPClient,
				username:   tt.clientData.username,
				password:   tt.clientData.password,
			}
			got, err := c.GetUsers()
			CheckForResponderError(t, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetUserDetails(t *testing.T) {
	type args struct {
		userId string
	}
	tests := []struct {
		name         string
		clientData   clientData
		args         args
		statusCode   int
		responseBody string
		testOptions  RequestTestOptions
		want         *UserDetailsResponse
		wantErr      bool
	}{
		{
			name:         "Successful request",
			clientData:   goodClient,
			args:         args{userId: "john.doe"},
			statusCode:   200,
			responseBody: `<?xmlversion="1.0"?><ocs><meta><status>ok</status><statuscode>100</statuscode><message>OK</message><totalitems></totalitems><itemsperpage></itemsperpage></meta><data><enabled>1</enabled><storageLocation>/var/www/html/data/john.doe</storageLocation><id>john.doe</id><lastLogin>1618156321000</lastLogin><backend>Database</backend><subadmin><element>employees</element></subadmin><quota><free>549184147456</free><used>16792345</used><total>549200939801</total><relative>0</relative><quota>-3</quota></quota><email>john.doe@example.local</email><displayname>John Doe</displayname><phone>+1555123</phone><address></address><website>example.local</website><twitter></twitter><groups><element>developers</element><element>employees</element></groups><language>en</language><locale></locale><backendCapabilities><setDisplayName>1</setDisplayName><setPassword>1</setPassword></backendCapabilities></data></ocs>`,
			testOptions:  DefaultTestOptions(),
			want: &UserDetailsResponse{
				XMLName: xml.Name{
					Local: "ocs",
				},
				RequestMeta: &MetaFragment{
					StatusCode:   100,
					Status:       "ok",
					Message:      "OK",
					TotalItems:   0,
					ItemsPerPage: 0,
				},
				Enabled:         true,
				StorageLocation: "/var/www/html/data/john.doe",
				Id:              "john.doe",
				LastLogin:       "1618156321000",
				Backend:         "Database",
				SubadminGroups:  []string{"employees"},
				Quota: &UserQuota{
					Free:     549184147456,
					Used:     16792345,
					Total:    549200939801,
					Relative: 0,
					Quota:    "-3",
				},
				Email:       "john.doe@example.local",
				DisplayName: "John Doe",
				Phone:       "+1555123",
				Address:     "",
				Website:     "example.local",
				Twitter:     "",
				Groups:      []string{"developers", "employees"},
				Language:    "en",
				Locale:      "",
				BackendCapabilities: &UserDetailsBackendCapabilities{
					SetDisplayName: true,
					SetPassword:    true,
				},
			},
			wantErr: false,
		},
		{
			name:         "No such user",
			clientData:   goodClient,
			args:         args{userId: "jack.nobody"},
			statusCode:   200,
			responseBody: `<?xml version="1.0"?><ocs><meta><status>failure</status><statuscode>404</statuscode><message>User does not exist</message><totalitems></totalitems><itemsperpage></itemsperpage></meta><data/></ocs>`,
			testOptions:  DefaultTestOptions(),
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "Bad credentials",
			clientData:   badClient,
			args:         args{userId: "john.doe"},
			statusCode:   401,
			responseBody: badLoginResponse,
			testOptions:  DefaultTestOptions(),
			want:         nil,
			wantErr:      true,
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetResponder(fmt.Sprintf("%s/ocs/v1.php/cloud/users/%s", HOST, tt.args.userId), tt.statusCode, tt.responseBody, tt.testOptions)
			c := &Client{
				HostURL:    tt.clientData.HostURL,
				HTTPClient: tt.clientData.HTTPClient,
				username:   tt.clientData.username,
				password:   tt.clientData.password,
			}
			got, err := c.GetUserDetails(tt.args.userId)
			CheckForResponderError(t, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				//<editor-fold desc="Check pointers manually">
				if !reflect.DeepEqual(*got.RequestMeta, *tt.want.RequestMeta) {
					t.Errorf("GetUserDetails() RequestMeta got = %v, want %v", got, tt.want)
				}
				got.RequestMeta = nil
				tt.want.RequestMeta = nil
				if !reflect.DeepEqual(*got.Quota, *tt.want.Quota) {
					t.Errorf("GetUserDetails() Quota got = %v, want %v", got, tt.want)
				}
				got.Quota = nil
				tt.want.Quota = nil
				if !reflect.DeepEqual(*got.BackendCapabilities, *tt.want.BackendCapabilities) {
					t.Errorf("GetUserDetails() BackendCapabilities got = %v, want %v", got, tt.want)
				}
				got.BackendCapabilities = nil
				tt.want.BackendCapabilities = nil
				//</editor-fold>

				if !reflect.DeepEqual(*got, *tt.want) {
					t.Errorf("GetUserDetails() got = %v, want %v", *got, *tt.want)
				}
			}
		})
	}
}
