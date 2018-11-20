package githubsearch

import (
	"fmt"
	"github.com/chennqqi/gshark/models"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	gitClient, _, _ := GetGithubClient()
	ownerName := "neal1991"
	user := GetGithubUserInfo(gitClient, &ownerName)
	assert.Equal(t, *user.Blog, "https://madneal.com")
}

func TestInsertAllRepos(t *testing.T) {
	assets, err := models.ListInputInfo()
	if err == nil {
		for index, asset := range assets {
			fmt.Println("index:" + strconv.Itoa(index))
			fmt.Println(asset.Content)
			length := len(strings.Split(asset.Content, "/"))
			fmt.Println(strings.Split(asset.Content, "/")[length-1])
		}
	}
}
