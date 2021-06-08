package github_test

import (
	"context"
	"testing"

	gogithub "github.com/google/go-github/github"
	"github.com/nori-io/nori/internal/domain/helper"
	"github.com/nori-io/nori/internal/helper/github"
	"github.com/stretchr/testify/assert"
)

func TestReleaseAsset(t *testing.T) {
	assert := assert.New(t)

	url := "github.com/secure2work/http@0.0.85"

	client := gogithub.NewClient(nil)

	githubHelper := github.New(github.Params{Client: client})

	ctx := context.Background()

	owner, repo, releaseVersion, err := githubHelper.Parse(url)
	assert.Nil(err)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	getReleaseAssetData := helper.GetReleaseAssetData{
		Owner:          owner,
		Repo:           repo,
		ReleaseVersion: releaseVersion,
		CommonVersion:  "4.0.0",
		GoVersion:      "1.16.3",
		GoosVersion:    "linux",
		GoarchVersion:  "amd64",
	}

	asset, err := githubHelper.GetReleaseAsset(ctx, getReleaseAssetData)
	assert.NotEqual(asset, nil)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	assert.Nil(err)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	testAsset := getReleaseAssetData.Repo + releaseVersion + "_" +
		"common" + getReleaseAssetData.CommonVersion + "_" +
		"go" + getReleaseAssetData.GoVersion + "_" +
		getReleaseAssetData.GoosVersion + "_" +
		getReleaseAssetData.GoarchVersion + ".so"
	assert.Equal(testAsset, asset.Name)
	assert.Equal(int64(38264290), asset.ID)

	err = githubHelper.DownloadReleaseAsset(ctx, owner, repo, *asset, ``)
	assert.Nil(err)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}
