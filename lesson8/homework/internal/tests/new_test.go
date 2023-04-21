package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	client := getTestClient()

	userResp, err := client.createUser(123, "user", "somemail1@mail.com")
	assert.NoError(t, err)
	assert.Equal(t, userResp.Data.ID, int64(123))
	assert.Equal(t, userResp.Data.Nickname, "user")
	assert.Equal(t, userResp.Data.Email, "somemail1@mail.com")
}

func TestChangeUserInfo(t *testing.T) {
	client := getTestClient()

	_, _ = client.createUser(123, "user", "somemail1@mail.com")
	response, err := client.changeUserInfo(123, "123", "somemail2@mail.ru")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Nickname, "123")
	assert.Equal(t, response.Data.Email, "somemail2@mail.ru")
}

func TestGetAdsByTitle(t *testing.T) {
	client := getTestClient()

	_, _ = client.createUser(123, "user", "somemail1@mail.com")

	_, _ = client.createAd(123, "hello1", "world1")
	_, _ = client.createAd(123, "hello1", "world2")
	_, _ = client.createAd(123, "hello1", "world3")

	ads, err := client.getAdsByTitle("hello1")
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 3)
}

/*
func TestFilterByAuthor(t *testing.T) {
	client := getTestClient()

	_, _ = client.createUser(123, "user1", "somemail1@mail.com")
	_, _ = client.createUser(124, "user2", "somemail2@mail.com")

	ad1, _ := client.createAd(123, "hello1", "world1")
	ad2, _ := client.createAd(123, "hello2", "world2")
	ad3, _ := client.createAd(124, "hello3", "world3")

	_, _ = client.changeAdStatus(123, ad1.Data.ID, true)
	_, _ = client.changeAdStatus(123, ad2.Data.ID, true)
	_, _ = client.changeAdStatus(124, ad3.Data.ID, true)

	ads, err := client.listAdsAuthor(123)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
}
*/
