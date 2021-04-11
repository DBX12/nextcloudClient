package nextcloudClient

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	StatusSuccess  = 100
	InvalidRequest = 999
	NotAuthorized  = 997
)

const QuotaUnlimited = "none"

type UserData struct {
	UserId           string
	Password         string
	DisplayName      string
	Email            string
	GroupIds         []string
	SubadminGroupIds []string
	Quota            string
	Language         string
}

func (userData *UserData) Validate() (bool, []string) {
	var problems []string
	result := true
	if userData.UserId == "" {
		problems = append(problems, "UserId must not be empty")
		result = false
	}
	if userData.Email == "" && userData.Password == "" {
		problems = append(problems, "Either Password or Email must be set")
		result = false
	}
	return result, problems
}

func (c *Client) GetUsers() ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/cloud/users", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	response := GetUsersResponse{}

	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.RequestMeta.StatusCode != StatusSuccess {
		return nil, errors.New(fmt.Sprintf("Api returned a status code %d indicating failure", response.RequestMeta.StatusCode))
	}

	return response.UserNames, nil
}

func (c *Client) CreateUser(userData *UserData) (bool, error) {
	bodyData := url.Values{}

	result, problems := userData.Validate()
	if result == false {
		return false, errors.New(strings.Join(problems, "\n"))
	}

	bodyData.Set("userid", userData.UserId)
	if userData.Password != "" {
		bodyData.Set("password", userData.Password)
	}
	if userData.DisplayName != "" {
		bodyData.Set("displayName", userData.DisplayName)
	}
	if userData.Email != "" {
		bodyData.Set("email", userData.Email)
	}
	if len(userData.GroupIds) > 0 {
		for _, groupId := range userData.GroupIds {
			bodyData.Add("groups[]", groupId)
		}
	}
	if len(userData.SubadminGroupIds) > 0 {
		for _, subadminGroupId := range userData.SubadminGroupIds {
			bodyData.Add("subadmin[]", subadminGroupId)
		}
	}
	if userData.Quota != "" {
		bodyData.Set("quota", userData.Quota)
	}
	if userData.Language != "" {
		bodyData.Set("language", userData.Language)
	}

	return doSimpleRequest(
		c,
		http.MethodPost,
		fmt.Sprintf("%s/cloud/users", c.HostURL),
		&bodyData,
	)
}

func (c *Client) GetUserDetails(userId string) (*UserDetailsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/cloud/users/%s", c.HostURL, userId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	response := UserDetailsResponse{}

	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.RequestMeta.StatusCode != StatusSuccess {
		return nil, errors.New(fmt.Sprintf("Api returned a status code %d indicating failure", response.RequestMeta.StatusCode))
	}

	return &response, nil
}

func (c *Client) UpdateUserDetail(userId string, attribute string, value string) (bool, error) {

	bodyData := url.Values{}
	bodyData.Set("key", attribute)
	bodyData.Set("value", value)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/cloud/users/%s", c.HostURL, userId), strings.NewReader(bodyData.Encode()))
	if err != nil {
		return false, err
	}
	c.addHeadersForBody(req, len(bodyData.Encode()))
	body, err := c.doRequest(req)
	if err != nil {
		return false, err
	}

	response := SimpleResponse{}
	if err := xml.Unmarshal(body, &response); err != nil {
		return false, err
	}

	if response.RequestMeta.StatusCode != StatusSuccess {
		return false, errors.New(fmt.Sprintf("Api returned a status code %d indicating failure", response.RequestMeta.StatusCode))
	}

	return true, nil
}

func (c *Client) DisableUser(userId string) (bool, error) {
	return doSimpleRequest(
		c,
		http.MethodPut,
		fmt.Sprintf("%s/cloud/users/%s/disable", c.HostURL, userId),
		nil,
	)
}

func (c *Client) EnableUser(userId string) (bool, error) {

	return doSimpleRequest(
		c,
		http.MethodPut,
		fmt.Sprintf("%s/cloud/users/%s/enable", c.HostURL, userId),
		nil,
	)
}

func (c *Client) DeleteUser(userId string) (bool, error) {
	return doSimpleRequest(
		c,
		http.MethodDelete,
		fmt.Sprintf("%s/cloud/users/%s", c.HostURL, userId),
		nil,
	)
}

func (c *Client) GetUserGroups(userId string) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/cloud/users/%s/groups", c.HostURL, userId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	response := UserGroupsResponse{}

	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.RequestMeta.StatusCode != StatusSuccess {
		return nil, errors.New(fmt.Sprintf("Api returned a status code %d indicating failure", response.RequestMeta.StatusCode))
	}

	return response.GroupNames, nil
}

func (c *Client) AddUserToGroup(userId string, groupId string) (bool, error) {

	bodyData := url.Values{}
	bodyData.Set("groupid", groupId)

	return doSimpleRequest(
		c,
		http.MethodPost,
		fmt.Sprintf("%s/cloud/users/%s/groups", c.HostURL, userId),
		&bodyData,
	)
}

func (c *Client) RemoveUserFromGroup(userId string, groupId string) (bool, error) {

	bodyData := url.Values{}
	bodyData.Set("groupid", groupId)

	return doSimpleRequest(
		c,
		http.MethodDelete,
		fmt.Sprintf("%s/cloud/users/%s/groups", c.HostURL, userId),
		&bodyData,
	)
}

func (c *Client) PromoteToSubadmin(userId, groupId string) (bool, error) {
	bodyData := url.Values{}
	bodyData.Set("groupid", groupId)
	return doSimpleRequest(
		c,
		http.MethodPost,
		fmt.Sprintf("%s/cloud/users/%s/subadmins", c.HostURL, userId),
		&bodyData,
	)
}

func (c *Client) DemoteFromSubadmin(userId, groupId string) (bool, error) {
	bodyData := url.Values{}
	bodyData.Set("groupid", groupId)
	return doSimpleRequest(
		c,
		http.MethodDelete,
		fmt.Sprintf("%s/cloud/users/%s/subadmins", c.HostURL, userId),
		&bodyData,
	)
}

func (c *Client) GetSubadminGroups(userId string) ([]string, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/cloud/users/%s/subadmins", c.HostURL, userId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	response := UserSubadminGroupsResponse{}

	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.RequestMeta.StatusCode != StatusSuccess {
		return nil, errors.New(fmt.Sprintf("Api returned a status code %d indicating failure", response.RequestMeta.StatusCode))
	}

	return response.GroupNames, nil
}

func (c *Client) ResendWelcomeMail(userId string) (bool, error) {
	return doSimpleRequest(
		c,
		http.MethodPost,
		fmt.Sprintf("%s/cloud/users/%s/welcome", c.HostURL, userId),
		nil,
	)
}
