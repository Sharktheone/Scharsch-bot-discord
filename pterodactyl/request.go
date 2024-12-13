package pterodactyl

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"io"
	"net/http"
	url2 "net/url"
)

func request(path string, method string, payload io.Reader) (*http.Response, error) {
	var (
		url, err = url2.JoinPath(conf.Config.Pterodactyl.PanelURL, path)
	)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(method, url, payload)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", conf.Config.Pterodactyl.APIKey))
	return http.DefaultClient.Do(req)
}
