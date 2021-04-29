package nextcloudClient

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetGroups() ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/cloud/groups", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	response := GetGroupsResponse{}

	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.RequestMeta.StatusCode != StatusSuccess {
		return nil, errors.New(fmt.Sprintf("Api returned a status code %d indicating failure", response.RequestMeta.StatusCode))
	}

	return response.GroupNames, nil
}

func (c *Client) GetGroup(groupId string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/cloud/groups?search=%s", c.HostURL, groupId), nil)
	if err != nil {
		return "", err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	response := GetGroupsResponse{}

	if err = xml.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if response.RequestMeta.StatusCode != StatusSuccess {
		return "", errors.New(fmt.Sprintf("Api returned a status code %d indicating failure", response.RequestMeta.StatusCode))
	}

	for _, name := range response.GroupNames {
		if name == groupId {
			return groupId, nil
		}
	}

	return "", errors.New(fmt.Sprintf("No group with the id %s was found", groupId))
}

func (c *Client) CreateGroup(groupId string) (bool, error) {

	bodyData := url.Values{}
	bodyData.Set("groupid", groupId)

	return doSimpleRequest(
		c,
		http.MethodPost,
		fmt.Sprintf("%s/cloud/groups", c.HostURL),
		&bodyData,
	)
}

func (c *Client) DeleteGroup(groupId string) (bool, error) {
	return doSimpleRequest(
		c,
		http.MethodDelete,
		fmt.Sprintf("%s/cloud/groups/%s", c.HostURL, groupId),
		nil,
	)
}

func (c *Client) GetGroupMembers(groupId string) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/cloud/groups/%s", c.HostURL, groupId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	response := GetGroupMembersResponse{}

	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.RequestMeta.StatusCode != StatusSuccess {
		return nil, errors.New(fmt.Sprintf("Api returned a status code %d indicating failure", response.RequestMeta.StatusCode))
	}

	return response.UserNames, nil
}

func (c *Client) GetGroupSubadmins(groupId string) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/cloud/groups/%s/subadmins", c.HostURL, groupId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	response := GetGroupSubadminsResponse{}

	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.RequestMeta.StatusCode != StatusSuccess {
		return nil, errors.New(fmt.Sprintf("Api returned a status code %d indicating failure", response.RequestMeta.StatusCode))
	}

	return response.UserNames, nil
}
