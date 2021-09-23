package hub

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/openshift-pipelines/pipelines-as-code/pkg/params"
)

var (
	tektonCatalogHubName = `tekton`
	hubBaseURL           = `https://api.hub.tekton.dev/v1`
)

func getURL(ctx context.Context, cli *params.Run, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return []byte{}, err
	}
	res, err := cli.Clients.HTTP.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	statusOK := res.StatusCode >= 200 && res.StatusCode < 300
	if !statusOK {
		return nil, fmt.Errorf("Non-OK HTTP status: %d", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func getSpecificVersion(ctx context.Context, cli *params.Run, task string) (string, error) {
	split := strings.Split(task, ":")
	version := split[len(split)-1]
	taskName := split[0]
	hr := new(hubResourceVersion)
	data, err := getURL(ctx, cli,
		fmt.Sprintf("%s/resource/%s/task/%s/%s", hubBaseURL, tektonCatalogHubName, taskName, version))
	if err != nil {
		return "", fmt.Errorf("could not fetch specific task version from the hub %s:%s: %w", task, version, err)
	}
	err = json.Unmarshal(data, &hr)
	if err != nil {
		return "", err
	}
	return *hr.Data.RawURL, nil
}

func getLatestVersion(ctx context.Context, cli *params.Run, task string) (string, error) {
	hr := new(hubResource)
	data, err := getURL(ctx, cli, fmt.Sprintf("%s/resource/%s/task/%s", hubBaseURL, tektonCatalogHubName, task))
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(data, &hr)
	if err != nil {
		return "", err
	}

	return *hr.Data.LatestVersion.RawURL, nil
}

func GetTask(ctx context.Context, cli *params.Run, task string) (string, error) {
	var rawURL string
	var err error

	if strings.Contains(task, ":") {
		rawURL, err = getSpecificVersion(ctx, cli, task)
	} else {
		rawURL, err = getLatestVersion(ctx, cli, task)
	}
	if err != nil {
		return "", fmt.Errorf("could not fetch remote task %s: %w", task, err)
	}

	data, err := getURL(ctx, cli, rawURL)
	if err != nil {
		return "", fmt.Errorf("could not fetch remote task %s: %w", task, err)
	}
	return string(data), err
}
