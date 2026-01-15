package api

import "net/http"

type Folder struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type FoldersResponse struct {
	Folders []Folder `json:"folders"`
}

func GetFolders(c *Client, spaceID string) ([]Folder, error) {
	resp, err := Do[any, FoldersResponse](c, http.MethodGet, "/space/"+spaceID+"/folder", nil)
	if err != nil {
		return nil, err
	}
	return resp.Folders, nil
}
