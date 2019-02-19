package flex

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

// EmailModule ...
type EmailModule struct {
	appMetadata kinveyAppMetadata
	client      *http.Client
	baseRoute   string
}

// EmailResponse ...
type emailResponse struct {
	mailServerResponse string
}

func newEmailModule(appMetadata kinveyAppMetadata) EmailModule {
	return EmailModule{
		appMetadata: appMetadata,
		client:      &http.Client{},
		baseRoute:   "rpc",
	}
}

// Email ...
type Email struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Subject  string `json:"subject,omitempty"`
	TextBody string `json:"body,omitempty"`
	ReplyTo  string `json:"replyTo,omitempty"`
	HTMLBody string `json:"html,omitempty"`
	CC       string `json:"cc,omitempty"`
	BCC      string `json:"bcc,omitempty"`
}

// Send ...
func (m EmailModule) Send(email Email) (string, error) {
	if email.From == "" || email.To == "" || email.Subject == "" || email.TextBody == "" {
		return "", errors.New("To send an email, you must specify the 'to', 'from', 'subject' and 'body' parameters")
	}

	requestOptions, err := m.buildEmailRequest(email)
	if err != nil {
		return "", err
	}

	return m.makeRequest(requestOptions)
}

func (m EmailModule) makeRequest(req *http.Request) (string, error) {
	resp, err := m.client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 26214400))
	if err != nil {
		return "", err
	}

	emailResponse := emailResponse{}
	err = json.Unmarshal(body, &emailResponse)
	if err != nil {
		return "", err
	}

	return emailResponse.mailServerResponse, nil
}

func (m EmailModule) buildEmailRequest(email Email) (*http.Request, error) {
	url := m.appMetadata.BaaSURL + "/" + m.baseRoute + "/" + m.appMetadata.ID + "/send-email"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-kinvey-api-version", "3")

	req.SetBasicAuth(m.appMetadata.ID, m.appMetadata.MasterSecret)

	json, err := json.Marshal(email)
	if err != nil {
		return nil, err
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(json))

	return req, nil
}
