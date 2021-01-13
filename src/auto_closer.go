package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type autoCloser struct {
	token      string
	repository string
	endpoint   string
	keep       int
	label      string
	issues     []issue
}

func newAutoCloser() (*autoCloser, error) {
	a := &autoCloser{}

	if os.Getenv("AC_KEEP") == "" {
		a.keep = 1
	} else {
		keep, err := strconv.Atoi(os.Getenv("AC_KEEP"))
		if err != nil {
			return nil, err
		}
		if keep < 1 {
			return nil, errors.New("AC_KEEP should be larger than 0")
		}
		if keep > 99 {
			return nil, errors.New("AC_KEEP should be less than 100")
		}
		a.keep = keep
	}

	a.token = os.Getenv("GITHUB_TOKEN")
	a.repository = os.Getenv("GITHUB_REPOSITORY")
	a.label = os.Getenv("AC_LABEL")
	a.endpoint = "https://api.github.com/repos/" + a.repository + "/issues?labels=" + a.label
	return a, nil
}

func (a *autoCloser) getIssuesList() error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, a.endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "token "+a.token)

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

	err = json.Unmarshal(b, &a.issues)
	if err != nil {
		return err
	}

	return nil
}

func (a *autoCloser) closeIssues() error {
	if len(a.issues) == 0 {
		fmt.Println("No issues found with the label: " + a.label)
	}

	if len(a.issues) < a.keep {
		fmt.Printf("AC_KEEP is %v, but there are only %v open issues.", a.keep, len(a.issues))
		return nil
	}

	oldIssues := a.issues[a.keep:len(a.issues)]

	for i := range oldIssues {
		patchData := `{"state":"closed"}`

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPatch, oldIssues[i].URL, strings.NewReader(patchData))
		if err != nil {
			return err
		}

		req.Header.Add("Accept", "application/vnd.github.v3+json")
		req.Header.Add("Authorization", "token "+a.token)

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
