package main

import (
	"context"
	"fmt"
	"gitlab-terraform-manager/pkg/graphql"
	"io/ioutil"
	"net/http"
	"strings"
)

type GitlabImpl struct {
	ServerAddress string
	AccessToken   string
	ProjectID     string
	ProjectPath   string
}

type GitlabStatesListImpl struct {
	Data struct {
		Project struct {
			TerraformStates struct {
				Edges []struct {
					Node struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"terraformStates"`
		} `json:"project"`
	} `json:"data"`
}

func gitlabInit(s, i, p, t string) (result *GitlabImpl) {
	result = &GitlabImpl{
		ServerAddress: s,
		ProjectID:     i,
		ProjectPath:   p,
		AccessToken:   t,
	}
	return
}

func (gitlab *GitlabImpl) GitlabGetStatesList() (result []string) {
	client := graphql.NewClient(fmt.Sprintf("%s/api/graphql", gitlab.ServerAddress))
	req := graphql.NewRequest(fmt.Sprintf(`
	{
	project(fullPath: "%s") {
		terraformStates(first: null, after: null, last: null, before: null) {
		  edges {
			node {
			  id
			  name
			}
		  }
		}
	  }
	}
    `, gitlab.ProjectPath))

	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Private-Token", gitlab.AccessToken)
	ctx := context.Background()

	var respData GitlabStatesListImpl
	if err := client.Run(ctx, req, &respData); err != nil {
		panic(err)
	}

	for _, node := range respData.Data.Project.TerraformStates.Edges {
		result = append(result, node.Node.Name)
	}

	return
}

func (gitlab *GitlabImpl) GetState(stateName string) string {
	client := &http.Client{}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/api/v4/projects/%s/terraform/state/%s",
			gitlab.ServerAddress,
			gitlab.ProjectID,
			stateName,
		),
		nil,
	)

	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Private-Token", gitlab.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if (resp.StatusCode != http.StatusOK) && (resp.StatusCode != http.StatusNoContent) {
		panic(resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(bodyBytes)
}

func (gitlab *GitlabImpl) RemoveState(stateName string) {
	client := &http.Client{}

	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf(
			"%s/api/v4/projects/%s/terraform/state/%s",
			gitlab.ServerAddress,
			gitlab.ProjectID,
			stateName,
		),
		nil,
	)

	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Private-Token", gitlab.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(resp.StatusCode)
	}
}

func (gitlab *GitlabImpl) RestoreState(stateName, stateData string) {
	client := &http.Client{}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"%s/api/v4/projects/%s/terraform/state/%s",
			gitlab.ServerAddress,
			gitlab.ProjectID,
			stateName,
		),
		strings.NewReader(stateData),
	)

	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Private-Token", gitlab.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	println("Cannot restore", stateName)
	// }
}
