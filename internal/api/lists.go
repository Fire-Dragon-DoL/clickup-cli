package api

import "net/http"

type List struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListsResponse struct {
	Lists []List `json:"lists"`
}

func GetLists(c *Client, folderID string) ([]List, error) {
	resp, err := Do[any, ListsResponse](c, http.MethodGet, "/folder/"+folderID+"/list", nil)
	if err != nil {
		return nil, err
	}
	return resp.Lists, nil
}
