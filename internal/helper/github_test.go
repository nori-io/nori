package helper_test

import (
	"context"
	"testing"

	"github.com/google/go-github/github"
	h "github.com/nori-io/nori/internal/domain/helper"
	"github.com/nori-io/nori/internal/helper"

	"github.com/stretchr/testify/assert"
)

func TestReleaseAsset(t *testing.T) {
	assert := assert.New(t)

	url := "github.com/secure2work/http@0.0.81"

	client := github.NewClient(nil)

	githubHelper := helper.New(helper.Params{Client: client})

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

	getReleaseAssetData := h.GetReleaseAssetData{
		Owner:          owner,
		Repo:           repo,
		ReleaseVersion: releaseVersion,
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

	testAsset := getReleaseAssetData.Repo + "_" + getReleaseAssetData.GoVersion + "." + getReleaseAssetData.GoosVersion + "-" + getReleaseAssetData.GoarchVersion + ".so"
	assert.Equal(testAsset, asset.Name)
	assert.Equal(int64(37988274), asset.ID)

	err = githubHelper.DownloadReleaseAsset(ctx, owner, repo, *asset, ``)
	assert.Nil(err)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}
