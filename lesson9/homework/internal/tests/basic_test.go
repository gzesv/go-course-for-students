package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAd(t *testing.T) {
	client := getTestClient()
	_, _ = client.createUser(123, "nickname", "example@mail.com")

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(123))
	assert.False(t, response.Data.Published)
}

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()
	_, _ = client.createUser(123, "user", "somemail@mail.com")
	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(123, response.Data.ID, true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Published)

	response, err = client.changeAdStatus(123, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()
	_, _ = client.createUser(123, "user", "somemail@mail.com")
	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(123, response.Data.ID, "привет", "мир")

	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
}

func TestListAds(t *testing.T) {
	client := getTestClient()
	_, _ = client.createUser(123, "user", "somemail@mail.com")
	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(123, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(123, "title", "text")
	assert.NoError(t, err)

	ads, er := client.listAds()
	assert.NoError(t, er)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}

func TestDeleteAd(t *testing.T) {
	client := getTestClient()

	_, _ = client.createUser(123, "user", "somemail@mail.com")
	a, _ := client.createAd(123, "title", "text")

	response, _ := client.deleteAd(a.Data.AuthorID, a.Data.ID)
	_, err := client.changeAdStatus(123, response.Data.ID, true)
	assert.ErrorIs(t, err, ErrBadRequest)
}
