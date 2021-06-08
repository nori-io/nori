package helper

import (
	"context"

	"github.com/google/go-github/github"
)

type GithubHelperInterface interface {
	GetReleaseAsset(ctx context.Context, getReleaseAssetData GetReleaseAssetData) (*GithubAsset, error)
	Parse(url string) (owner string, repo string, releaseVersion string, err error)
	DownloadReleaseAsset(ctx context.Context, owner, repo string, asset GithubAsset, path string) error
}

type GithubRelease struct {
	ID     int64
	Assets []github.ReleaseAsset
}

type GithubAsset struct {
	ID   int64
	Name string
}

type GetReleaseAssetData struct {
	Owner          string
	Repo           string
	ReleaseVersion string
	CommonVersion  string
	GoVersion      string
	GoosVersion    string
	GoarchVersion  string
}
