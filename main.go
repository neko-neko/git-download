package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

// GitHub api response
type assetsResponse struct {
	Assets []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"assets"`
}

var (
	repository = flag.String("repo", "", "target repository")
	saveDir    = flag.String("dir", "", "save dir path. default are set current dir")
	version    = flag.String("version", "", "specified version. default are set latest")
	include    = flag.String("include", "", "download specified file name pattern")
)

func main() {
	flag.Parse()

	if *repository == "" {
		flag.Usage()
		os.Exit(1)
	}

	wg := &sync.WaitGroup{}
	assets := readAssets().Assets
	for i := range assets {
		asset := assets[i]
		wg.Add(1)

		go func() {
			defer wg.Done()

			if *include != "" {
				if strings.Contains(asset.Name, "darwin") {
					downloadAsset(asset.Id, asset.Name)
				}
			} else {
				downloadAsset(asset.Id, asset.Name)
			}
		}()
	}
	wg.Wait()
}

func readAssets() assetsResponse {
	var tag = "latest"
	if *version != "" {
		tag = fmt.Sprintf("tags/%s", *version)
	}

	body := request(
		fmt.Sprintf("%s/repos/%s/releases/%s", generateGitHubApiUrl(), *repository, tag),
		map[string]string{"Accept": "application/json"},
	)

	var assetsRes assetsResponse
	err := json.Unmarshal(body, &assetsRes)
	if err != nil {
		abort(err)
	}

	return assetsRes
}

func downloadAsset(assetId int, assetName string) {
	fmt.Printf("Downloading %s...\n", assetName)

	body := request(
		fmt.Sprintf("%s/repos/%s/releases/assets/%d", generateGitHubApiUrl(), *repository, assetId),
		map[string]string{"Accept": "application/octet-stream"},
	)

	var dir string
	if *saveDir != "" {
		dir = fmt.Sprintf("%s/%s", *saveDir, assetName)
	} else {
		dir = fmt.Sprintf("%s", assetName)
	}

	file, err := os.OpenFile(dir, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		abort(err)
	}

	defer func() {
		file.Close()
		fmt.Printf("%s stored to %s\n", assetName, dir)
	}()
	file.Write(body)
}

func request(url string, headers map[string]string) []byte {
	req, _ := http.NewRequest(
		"GET",
		url,
		nil,
	)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		abort(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		abort(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		abort(errors.New(string(body)))
	}

	return body
}

func generateGitHubApiUrl() string {
	token := os.Getenv("GITHUB_TOKEN")
	apiUrl := os.Getenv("GITHUB_API")
	if apiUrl == "" {
		apiUrl = "api.github.com"
	}

	if token == "" {
		return fmt.Sprintf("https://%s", apiUrl)
	} else {
		return fmt.Sprintf("https://%s@%s", token, apiUrl)
	}
}

func abort(err error) {
	fmt.Fprintf(os.Stderr, "%s:%s", os.Args[0], err)
	os.Exit(1)
}
