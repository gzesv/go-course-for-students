package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestChangeUserInfo(t *testing.T) {
	client := getTestClient()
	_, _ = client.createUser(123, "user", "somemail@mail.com")

	response, err := client.changeUserInfo(123, "namenew", "somemailnew@mail.com")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Nickname, "namenew")
	assert.Equal(t, response.Data.Email, "somemailnew@mail.com")
}

func TestDeleteUserByID(t *testing.T) {
	client := getTestClient()

	a, _ := client.createUser(123, "user", "somemail@mail.com")

	response, err := client.deleteUserByID(123)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.ID, a.Data.ID)
	assert.Equal(t, response.Data.Nickname, a.Data.Nickname)
	assert.Equal(t, response.Data.Email, a.Data.Email)

	_, err = client.deleteUserByID(3)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestGetAdsByTitle(t *testing.T) {
	client := getTestClient()

	_, _ = client.createUser(123, "user", "somemail@mail.com")

	a, _ := client.createAd(123, "title", "text")
	b, _ := client.createAd(123, "title", "text")
	_, _ = client.createAd(123, "text", "title")

	ads, err := client.getAdsByTitle("title")
	assert.NoError(t, err)
	assert.Equal(t, ads.Data[0].ID, a.Data.ID)
	assert.Equal(t, ads.Data[0].Title, a.Data.Title)
	assert.Equal(t, ads.Data[0].Text, a.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, a.Data.AuthorID)
	assert.Equal(t, ads.Data[0].Published, a.Data.Published)

	assert.Equal(t, ads.Data[1].ID, b.Data.ID)
	assert.Equal(t, ads.Data[1].Title, b.Data.Title)
	assert.Equal(t, ads.Data[1].Text, b.Data.Text)
	assert.Equal(t, ads.Data[1].AuthorID, b.Data.AuthorID)
	assert.Equal(t, ads.Data[1].Published, b.Data.Published)
}
