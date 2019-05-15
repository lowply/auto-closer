package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type AutoCloser struct {
	Config Config
	Issues []Issue
}

type Config struct {
	Token        string
	Repository   string
	ListEndpoint string
	Keep         int
	Label        string
}

type Issue struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func (a *AutoCloser) setConfig() error {
	if os.Getenv("GITHUB_TOKEN") == "" {
		return errors.New("GITHUB_TOKEN is empty")
	}

	if os.Getenv("GITHUB_REPOSITORY") == "" {
		return errors.New("GITHUB_REPOSITORY is empty")
	}

	if os.Getenv("AC_LABEL") == "" {
		return errors.New("AC_LABEL is empty")
	}

	if os.Getenv("AC_KEEP") == "" {
		a.Config.Keep = 1
	} else {
		keep, err := strconv.Atoi(os.Getenv("AC_KEEP"))
		if err != nil {
			return err
		}
		a.Config.Keep = keep
	}

	a.Config.Token = os.Getenv("GITHUB_TOKEN")
	a.Config.Repository = os.Getenv("GITHUB_REPOSITORY")
	a.Config.Label = os.Getenv("AC_LABEL")
	a.Config.ListEndpoint = "https://api.github.com/repos/" + a.Config.Repository + "/issues?labels=" + a.Config.Label

	return nil
}

func (a *AutoCloser) getIssuesList() error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, a.Config.ListEndpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "token "+a.Config.Token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Error getting a list of isues: " + resp.Status)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &a.Issues)
	if err != nil {
		return err
	}

	return nil
}

func (a *AutoCloser) closeIssues() error {
	oldIssues := a.Issues[a.Config.Keep:len(a.Issues)]

	for i := range oldIssues {
		patchData := `{"state":"closed"}`

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPatch, oldIssues[i].URL, strings.NewReader(patchData))
		if err != nil {
			return err
		}

		req.Header.Add("Accept", "application/vnd.github.v3+json")
		req.Header.Add("Authorization", "token "+a.Config.Token)

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return errors.New("Error posting a comment: " + resp.Status)
		}

		fmt.Println("Closed issue:\n" + oldIssues[i].Title)

		time.Sleep(1 * time.Second)
	}

	return nil
}

func main() {
	a := AutoCloser{}

	err := a.setConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = a.getIssuesList()
	if err != nil {
		log.Fatal(err)
	}

	err = a.closeIssues()
	if err != nil {
		log.Fatal(err)
	}
}
