package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	tokenName    = "GITHUB_ACCESS_TOKEN"
	newlineInput = "unexpected newline"
)

func main() {
	baseURL := flag.String("baseURL", "https://api.github.com/", "Base URL for GitHub API")
	flag.Parse()

	token := os.Getenv(tokenName)
	if token == "" {
		log.Fatalln("Must set environment variable:", tokenName)
	}

	ctx := context.Background()
	client := github.NewClient(
		oauth2.NewClient(
			ctx,
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			),
		),
	)

	parsedURL, err := url.Parse(*baseURL)
	if err != nil {
		log.Fatalln(err)
	}
	client.BaseURL = parsedURL

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Type the repo name without the username prefix to delete.\n\n")

	listOpts := &github.RepositoryListOptions{}
	for {
		repos, resp, err := client.Repositories.List(ctx, *user.Login, listOpts)
		if err != nil {
			log.Fatalln(err)
		}

		for _, repo := range repos {
			if !*repo.Fork {
				continue
			}

			var input string
			fmt.Printf("%s: ", *repo.HTMLURL)
			_, err := fmt.Scanln(&input)
			if err != nil && err.Error() != newlineInput {
				log.Fatalln(err)
			}

			if input == *repo.Name {
				_, err = client.Repositories.Delete(ctx, *repo.Owner.Login, *repo.Name)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Printf("%s has been deleted.\n", *repo.HTMLURL)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		listOpts.Page = resp.NextPage
	}
}
