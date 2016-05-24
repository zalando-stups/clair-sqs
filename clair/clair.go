package clair

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type errorResponse struct {
	Error struct {
		Message string
	}
}

func tryMessageError(reader io.Reader) error {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	var response errorResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return err
	}

	if response.Error.Message != "" {
		return errors.New(response.Error.Message)
	}

	return nil
}

func PushLayer(clairUrl string, json []byte) error {
	resp, err := http.Post(fmt.Sprintf("%v/v1/layers", clairUrl), "application/json", bytes.NewReader(json))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = tryMessageError(resp.Body); err != nil {
		return err
	}

	return nil
}

func GetLayer(clairUrl, layerId string) (json []byte, err error) {
	resp, err := http.Get(fmt.Sprintf("%v/v1/layers/%v?vulnerabilities&features", clairUrl, layerId))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	json, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if err = tryMessageError(bytes.NewReader(json)); err != nil {
		return
	}

	return json, nil
}

type notificationDetailEnvelope struct {
	Notification struct {
		NextPage string
		New      struct {
			Vulnerability struct {
				Name          string
				NamespaceName string
				Severity      string
			}
			LayersIntroducingVulnerability []string
		}
		Old struct {
			Vulnerability struct {
				Name          string
				NamespaceName string
				Severity      string
			}
			LayersIntroducingVulnerability []string
		}
	}
}

func ProcessNotification(clairUrl, notificationName string, pageProcessor func(newLayers []string, oldLayers []string) error) error {
	notificationUrl := fmt.Sprintf("%v/v1/notifications/%v", clairUrl, notificationName)

	var page = ""
	for {
		var pageUrl string
		if page == "" {
			pageUrl = fmt.Sprintf("%v?limit=%v", notificationUrl, 10)
		} else {
			pageUrl = fmt.Sprintf("%v?limit=%v&page=%v", notificationUrl, 10, page)
		}

		details, err := http.Get(pageUrl)
		if err != nil {
			return err
		}
		defer details.Body.Close()

		detailsBytes, err := ioutil.ReadAll(details.Body)
		if err != nil {
			return err
		}

		if err = tryMessageError(bytes.NewReader(detailsBytes)); err != nil {
			return err
		}

		var detailsJson notificationDetailEnvelope
		if err = json.Unmarshal(detailsBytes, &detailsJson); err != nil {
			return err
		}

		// call the notification page processor
		if err = pageProcessor(detailsJson.Notification.New.LayersIntroducingVulnerability, detailsJson.Notification.Old.LayersIntroducingVulnerability); err != nil {
			return err
		}

		// another page to process?
		if detailsJson.Notification.NextPage != "" {
			page = detailsJson.Notification.NextPage
		} else {
			break
		}
	}

	return nil
}

func DeleteNotification(clairUrl, notificationName string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%v/v1/notifications/%v", clairUrl, notificationName), nil)
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
