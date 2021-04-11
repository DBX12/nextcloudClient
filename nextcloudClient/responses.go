package nextcloudClient

import "encoding/xml"

type MetaFragment struct {
	// StatusCode Nextcloud internal status code
	StatusCode int `xml:"statuscode"`
	// Status Human readable error description
	Status       string `xml:"status"`
	Message      string `xml:"message"`
	TotalItems   int    `xml:"totalitems"`
	ItemsPerPage int    `xml:"itemsperpage"`
}

type SimpleResponse struct {
	XMLName     xml.Name      `xml:"ocs"`
	RequestMeta *MetaFragment `xml:"meta"`
}

type GetUsersResponse struct {
	XMLName     xml.Name      `xml:"ocs"`
	RequestMeta *MetaFragment `xml:"meta"`
	UserNames   []string      `xml:"data>users>element"`
}

type UserQuota struct {
	Free     int     `xml:"free"`
	Used     int     `xml:"used"`
	Total    int     `xml:"total"`
	Relative float32 `xml:"relative"`
	Quota    string  `xml:"quota"`
}

type UserDetailsBackendCapabilities struct {
	SetDisplayName bool `xml:"setDisplayName"`
	SetPassword    bool `xml:"setPassword"`
}

type UserDetailsResponse struct {
	XMLName         xml.Name      `xml:"ocs"`
	RequestMeta     *MetaFragment `xml:"meta"`
	Enabled         bool          `xml:"data>enabled"`
	StorageLocation string        `xml:"data>storageLocation"`
	Id              string        `xml:"data>id"`
	LastLogin       string        `xml:"data>last_login"`
	Backend         string        `xml:"data>backend"`
	// SubadminGroups list of groupIds this user is an admin of
	SubadminGroups []string   `xml:"subadmin>element"`
	Quota          *UserQuota `xml:"data>quota"`
	Email          string     `xml:"data>email"`
	DisplayName    string     `xml:"data>displayname"`
	Phone          string     `xml:"data>phone"`
	Address        string     `xml:"data>address"`
	Website        string     `xml:"data>website"`
	Twitter        string     `xml:"data>twitter"`
	// Groups list of groupIds this user belongs to
	Groups              []string                       `xml:"data>groups>element"`
	Language            string                         `xml:"data>language"`
	Locale              string                         `xml:"data>locale"`
	BackendCapabilities UserDetailsBackendCapabilities `xml:"data>backendCapabilities"`
}

type UserGroupsResponse struct {
	XMLName     xml.Name      `xml:"ocs"`
	RequestMeta *MetaFragment `xml:"meta"`
	GroupNames  []string      `xml:"data>groups>element"`
}

type UserSubadminGroupsResponse struct {
	XMLName     xml.Name      `xml:"ocs"`
	RequestMeta *MetaFragment `xml:"meta"`
	GroupNames  []string      `xml:"data>element"`
}

type GetGroupsResponse struct {
	XMLName     xml.Name      `xml:"ocs"`
	RequestMeta *MetaFragment `xml:"meta"`
	GroupNames  []string      `xml:"data>groups>element"`
}

type GetGroupMembersResponse struct {
	XMLName     xml.Name      `xml:"ocs"`
	RequestMeta *MetaFragment `xml:"meta"`
	UserNames   []string      `xml:"data>users>element"`
}

type GetGroupSubadminsResponse struct {
	XMLName     xml.Name      `xml:"ocs"`
	RequestMeta *MetaFragment `xml:"meta"`
	UserNames   []string      `xml:"data>element"`
}
