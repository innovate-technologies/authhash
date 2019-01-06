package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	recaptcha "github.com/dpapathanasiou/go-recaptcha"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	resty "gopkg.in/resty.v0"
)

type config struct {
	APIKey              string `required:"true"`
	Port                string `default:"8080"`
	RecaptchaPrivateKey string `required:"true" envconfig:"RECAPTCHA_PRIVATE_KEY"`
}

var conf config

func main() {
	err := envconfig.Process("authhash", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	e := echo.New()
	e.Static("/", "web")
	e.POST("/api/create", handleCreate)
	e.Logger.Fatal(e.Start(":" + conf.Port))
}

type createInfo struct {
	Recaptcha   string `json:"g-recaptcha-response" form:"g-recaptcha-response" query:"g-recaptcha-response"`
	StationName string `json:"stationname" form:"stationname" query:"stationname"`
	Genre       string `json:"genre" form:"genre" query:"genre"`
	Email       string `json:"email" form:"email" query:"email"`
	Website     string `json:"website" form:"website" query:"website"`
	Description string `json:"description" form:"description" query:"description"`
}

func handleCreate(c echo.Context) error {
	content := createInfo{}
	c.Bind(&content)

	recaptcha.Init(conf.RecaptchaPrivateKey)
	verified, err := recaptcha.Confirm(c.RealIP(), content.Recaptcha)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if !verified {
		return c.String(http.StatusBadRequest, "error verifying captcha")
	}

	hash, err := createHash(content.StationName, content.Genre, content.Email, "", "", content.Website, content.Description, "", "")

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, hash)
}

func createHash(name, genre, email, langid, countryiso, website, description, keywords, city string) (string, error) {
	r := resty.R()
	r = r.SetQueryParam("k", conf.APIKey)
	r = r.SetQueryParam("stationname", name)
	r = r.SetQueryParam("genre", genre)
	r = r.SetQueryParam("email", email)
	r = r.SetQueryParam("langid", langid)
	r = r.SetQueryParam("countryiso", countryiso)
	r = r.SetQueryParam("website", website)
	r = r.SetQueryParam("description", description)
	r = r.SetQueryParam("keywords", keywords)
	r = r.SetQueryParam("city", city)

	resp, _ := r.Get("https://yp.shoutcast.com/createauthhash")

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("HTTP error %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	yp := ypResponse{}
	xml.Unmarshal(resp.Body(), &yp)
	if yp.StatusCode != "200" {
		return "", fmt.Errorf("YP error %s: %s", yp.StatusCode, yp.StatusText)
	}

	return yp.Data.Authhash, nil

}

type ypResponse struct {
	XMLName    xml.Name `xml:"response"`
	StatusCode string   `xml:"statusCode"`
	StatusText string   `xml:"statusText"`
	Data       ypData   `xml:"data"`
}

type ypData struct {
	Authhash string `xml:"authhash"`
}
