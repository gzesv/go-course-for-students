package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"homework10/internal/adapters/adfilters"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/userrepo"
	"homework10/internal/ads"
	"homework10/internal/app"
	"testing"
)

func TestApp_InsertAdd(t *testing.T) {
	var AuthorID int64 = 123

	type TestApp struct {
		name     string
		title    string
		text     string
		userID   int64
		expected ads.Ad
	}

	appTests := [...]TestApp{
		{"Successful addition", "title", "text", AuthorID,
			ads.Ad{AuthorID: 123, Title: "title", Text: "text", Published: false}},
		{"Can't create", "title1", "text1", AuthorID + 1,
			ads.Ad{}},
		{"Validate Error", "", "text2", AuthorID,
			ads.Ad{}},
		{"Validate Error", "title3", "", AuthorID,
			ads.Ad{}},
	}

	for _, test := range appTests {
		t.Run(test.name, func(t *testing.T) {
			a := app.NewApp(adrepo.New(), userrepo.New(), adfilters.New())
			_, err := a.CreateUser(context.Background(), "nickname",
				"somemail@mail.ru", AuthorID)
			assert.NoError(t, err)
			ad, _ := a.CreateAd(context.Background(), test.title, test.text, test.userID)
			assert.Equal(t, ad.Title, test.expected.Title)
			assert.Equal(t, ad.Text, test.expected.Text)
			assert.Equal(t, ad.AuthorID, test.expected.AuthorID)
			assert.Equal(t, ad.Published, test.expected.Published)
		})
	}
}
