package clair

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func PushLayer(clairUrl string, json []byte) (err error) {
	_, err = http.Post(fmt.Sprintf("%v/v1/layers", clairUrl), "application/json", bytes.NewReader(json))
	// TODO check response status
	return
}

func GetLayer(clairUrl, layerId string) (json []byte, err error) {
	resp, err := http.Get(fmt.Sprintf("%v/v1/layers/%v?vulnerabilities&features", clairUrl, layerId))
	if err != nil {
		return
	}

	// TODO check http response status
	// TODO check json message for Error response

	return ioutil.ReadAll(resp.Body)
}

type notificationDetailEnvelope struct {
	Notification struct {
		NextPage string
		New      struct {
			LayersIntroducingVulnerability []string
		}
		Old struct {
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

		// TODO check http response status

		detailsBytes, err := ioutil.ReadAll(details.Body)
		if err != nil {
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

func DeleteNotification(clairUrl, notificationName string) (err error) {
	_, err = http.NewRequest("DELETE", fmt.Sprintf("%v/v1/notifications/%v", clairUrl, notificationName), nil)
	// TODO check response status
	return
}
