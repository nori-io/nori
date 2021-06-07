package helper

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
	"github.com/nori-io/nori/internal/domain/helper"
)

type GithubHelper struct {
	client *github.Client
}

type Params struct {
	Client *github.Client
}

func New(params Params) helper.GithubHelperInterface {
	return &GithubHelper{
		client: params.Client,
	}
}

// формат github.com/nori-plugins/plugin@v1.0.0
func (g GithubHelper) GetReleaseAsset(ctx context.Context, getReleaseAssetsData helper.GetReleaseAssetsData) (*helper.GithubAsset, error) {
	ReleaseByTag, _, err := g.client.Repositories.GetReleaseByTag(context.Background(), getReleaseAssetsData.Owner, getReleaseAssetsData.Repo, getReleaseAssetsData.ReleaseVersion)
	if err != nil {
		return nil, err
	}

	if ReleaseByTag == nil {
		return nil, errors.New("no releases found")
	}

	asset := getReleaseAssetsData.Repo + "_" + getReleaseAssetsData.GoVersion + "." + getReleaseAssetsData.GoosVersion + "-" + getReleaseAssetsData.GoarchVersion + ".so"

	for _, v := range ReleaseByTag.Assets {
		if *v.Name == asset {
			return &helper.GithubAsset{
				ID:   *v.ID,
				Name: *v.Name,
			}, nil
		}
	}
	return nil, errors.New("no asset found")
}

func (g GithubHelper) Parse(url string) (owner string, repo string, releaseVersion string, err error) {
	if !strings.HasPrefix(url, "github.com/") {
		return "", "", "", errors.New("unsupported host")
	}
	url = strings.TrimPrefix(url, "github.com/")

	urlSeparated := strings.Split(url, "@")

	if len(urlSeparated) == 0 {
		return "", "", "", errors.New("owner, repo and release's version not specified")
	}

	if len(urlSeparated) == 1 {
		return "", "", "", errors.New("owner and repo or release's version not specified")
	}

	urlParts := strings.Split(urlSeparated[0], "/")

	if len(urlParts) == 1 {
		return "", "", "", errors.New("repo not defined")
	}

	if len(urlParts) > 2 {
		return "", "", "", errors.New("owner and repo specified, but extra part of url exists")
	}

	owner = urlParts[0]
	repo = urlParts[1]
	releaseVersion = urlSeparated[1]

	return owner, repo, releaseVersion, nil
}

func (g GithubHelper) DownloadReleaseAsset(ctx context.Context, owner string, repo string, asset helper.GithubAsset, path string) error {
	rc, redirectUrl, err := g.client.Repositories.DownloadReleaseAsset(ctx, owner, repo, asset.ID)
	if err != nil {
		return err
	}

	if (rc == nil) && (redirectUrl == "") {
		return errors.New("get empty URL for downloading and no data")
	}

	out, err := os.Create(filepath.Join(path, asset.Name))
	if err != nil {
		return err
	}
	defer out.Close()

	if rc != nil {
		if err := saveFile(rc, out); err != nil {
			return err
		}
	}

	if redirectUrl != "" {
		if err = downloadFile(redirectUrl, out); err != nil {
			return err
		}
	}

	return nil
}

func downloadFile(redirectURL string, file *os.File) error {
	resp, err := http.Get(redirectURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := saveFile(resp.Body, file); err != nil {
		return err
	}

	return nil
}

func saveFile(rc io.Reader, file *os.File) error {
	var data []byte

	_, err := rc.Read(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(file.Name(), data, 0644)
	if err != nil {
		return err
	}

	return nil
}
