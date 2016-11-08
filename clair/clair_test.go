package clair

import (
	"fmt"
	"testing"
	"net/http"
	"net/http/httptest"
	"errors"
)

var page1 = `{
	"Notification": {
		"NextPage": "page2",
		"New": {
			"Vulnerability": {
				"Name": "vnew",
				"NamespaceName": "namespace1",
				"Severity": "critical"
			},
			"LayersIntroducingVulnerability": ["layer1", "layer2", "layer3"]
		},
		"Old": {
			"Vulnerability": {
				"Name": "vnew",
				"NamespaceName": "namespace1",
				"Severity": "critical"
			},
			"LayersIntroducingVulnerability": ["layer1", "layer2", "layer3"]
		}
	}
}`

var page2 = `{
	"Notification": {
		"NextPage": "page3",
		"New": {
			"Vulnerability": {
				"Name": "vnew",
				"NamespaceName": "namespace1",
				"Severity": "critical"
			},
			"LayersIntroducingVulnerability": ["BOOM!", "layer2", "layer3"]
		},
		"Old": {
			"Vulnerability": {
				"Name": "vnew",
				"NamespaceName": "namespace1",
				"Severity": "critical"
			},
			"LayersIntroducingVulnerability": ["layer1", "layer2", "layer3"]
		}
	}
}`

var page3 = `{
	"Notification": {
		"NextPage": "",
		"New": {
			"Vulnerability": {
				"Name": "vnew",
				"NamespaceName": "namespace1",
				"Severity": "critical"
			},
			"LayersIntroducingVulnerability": ["layer1", "layer2", "layer3"]
		},
		"Old": {
			"Vulnerability": {
				"Name": "vnew",
				"NamespaceName": "namespace1",
				"Severity": "critical"
			},
			"LayersIntroducingVulnerability": ["layer1", "layer2", "layer3"]
		}
	}
}`

var pages = map[string]string {
	"": page1,
	"page2": page2,
	"page3": page3,
}

func TestProcessNotification(t *testing.T) {
	fmt.Println("testing")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryValues := r.URL.Query()
		page := ""
		if pageVals, ok := queryValues["page"]; ok && len(pageVals) > 0 {
			page = pageVals[0]
		}
		fmt.Fprintln(w, pages[page])
	}))

	var err error
	nPagesProcessed := 0

	var robustPageProcessor = func(newLayers, oldLayers []string) error {
		nPagesProcessed += 1
		return nil
	}

	err = ProcessNotification(server.URL, "notificationName", robustPageProcessor)
	if err != nil {
		t.Error("there should be no error")
	}
	if nPagesProcessed != 3 {
		t.Error("all 3 pages should be processed")
	}

	nPagesProcessed = 0

	var failingPageProcessor = func(newLayers, oldLayers []string) error {
		nPagesProcessed += 1
		if len(newLayers) > 0 && newLayers[0] == "BOOM!" {
			return errors.New("bad error")
		}
		return nil
	}
	
	err = ProcessNotification(server.URL, "notificationName", failingPageProcessor)
	if err == nil {
		t.Error("an error should be reported")
	}
	if nPagesProcessed != 3 {
		t.Errorf("all 3 pages should be processed, but only %d were", nPagesProcessed)
	}
}
