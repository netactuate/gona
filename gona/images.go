package gona

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// Image
type Image struct {
	ID             int          `json:"id"`
	Name           string       `json:"os"`
	Description    *string      `json:"description"`
	Type           string       `json:"type"`
	Subtype        string       `json:"subtype"`
	Bits           string       `json:"bits"`
	Tech           string       `json:"tech"`
	Size           string       `json:"size"`
	Category       string       `json:"category"`
	Enabled        *int         `json:"os_enabled"`
	ScriptBash     int          `json:"script_bash"`
	ScriptCloudinit int         `json:"script_cloudinit"`
	Created        string       `json:"created"`
	Updated        string       `json:"updated"`
	ActiveBuild    *ImageBuild  `json:"active_build"`
}

// ImageBuild (job for tracking)
type ImageBuild struct {
	ID          int    `json:"id"`
	Status      int    `json:"status"`
	Command     string `json:"command"`
	TSInsert    string `json:"ts_insert"`
	MbID        int    `json:"mb_id"`
	MbPkgID     int    `json:"mb_pkgid"`
	Params      string `json:"params"`
	BuildPacket string `json:"build_packet"`
	Response    string `json:"response"`
	Created     string `json:"created"`
	LastUpdated string `json:"last_updated"`
}

// CreateImageRequest
type CreateImageRequest struct {
	MbPkgID          int    `url:"mbpkgid"`
	ImageName        string `url:"image_name"`
	ImageDescription string `url:"image_description,omitempty"`
	KeepSSHUserdirs  bool   `url:"keep_ssh_userdirs,omitempty"`
}

// CreateImageResponse
type CreateImageResponse struct {
	QueueID int `json:"queue_id"`
}

// DeleteImageResponse
type DeleteImageResponse struct {
	QueueID int `json:"queue_id"`
}

// ImageQueueStatus
type ImageQueueStatus struct {
	Status    string `json:"status"`
	Percent   int    `json:"percent"`
	Response  string `json:"response"`
	ImageID   int    `json:"image_id"`
	ImageName string `json:"image_name"`
	ImageHelp string `json:"image_help"`
	Location  string `json:"location"`
	MbPkgID   int    `json:"mb_pkgid"`
	FQDN      string `json:"fqdn"`
	OS        string `json:"os"`
}

func (c *Client) GetMyImages() ([]Image, error) {
	var images []Image
	if err := c.get("cloud/images/my", &images); err != nil {
		return nil, err
	}
	return images, nil
}

func (c *Client) GetImage(id int) (Image, error) {
	var image Image
	if err := c.get("cloud/images/"+strconv.Itoa(id), &image); err != nil {
		return Image{}, err
	}
	return image, nil
}

func (c *Client) CreateImage(req *CreateImageRequest) (CreateImageResponse, error) {
	values := url.Values{}
	values.Add("mbpkgid", strconv.Itoa(req.MbPkgID))
	values.Add("image_name", req.ImageName)
	if req.ImageDescription != "" {
		values.Add("image_description", req.ImageDescription)
	}
	if req.KeepSSHUserdirs {
		values.Add("keep_ssh_userdirs", "1")
	}

	var resp CreateImageResponse
	if err := c.post("cloud/images/create", []byte(values.Encode()), &resp); err != nil {
		return CreateImageResponse{}, err
	}
	return resp, nil
}

func (c *Client) EditImage(id int, name, description string) error {
	body := map[string]string{
		"os":          name,
		"description": description,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return c.patch("cloud/images/"+strconv.Itoa(id)+"/edit", bodyBytes, nil)
}

func (c *Client) DeleteImage(id int) (DeleteImageResponse, error) {
	var resp DeleteImageResponse
	if err := c.delete("cloud/images/"+strconv.Itoa(id)+"/delete", nil, &resp); err != nil {
		return DeleteImageResponse{}, err
	}
	return resp, nil
}

func (c *Client) GetImageQueueStatus(queueID int) (ImageQueueStatus, error) {
	var status ImageQueueStatus
	if err := c.get("cloud/images/queue_status/"+strconv.Itoa(queueID), &status); err != nil {
		return ImageQueueStatus{}, err
	}
	return status, nil
}
