package main

import (
	// standard
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	// external
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func grabRaw(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}

func main() {
	fmt.Println("Loading API token")
	fh, err := os.Open("oauth-token")
	if err != nil {
		panic(err)
	}

	tdata, err := ioutil.ReadAll(fh)
	if err != nil {
		panic(err)
	}

	token := strings.Trim(string(tdata), "\n")
	fh.Close()

	fmt.Println("TOKEN: " + token)

	fmt.Println("Initializing Oauth")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	fmt.Println("Acquiing Gists")
	gists, _, err := client.Gists.ListAll(nil)
	if err != nil {
		panic(err)
	}

	for _, gist := range gists {
		fmt.Println(*gist.HTMLURL)
		for _, filename := range gist.Files {
			fmt.Println(*filename.RawURL)
			err = grabRaw(*filename.RawURL)
			if err != nil {
				panic(err)
			}
		}
		fmt.Println()
	}
}
