package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/mitchellh/mapstructure"
)

type PullRequests struct {
	c *Client
}

/*
type Account struct {
	Links map[string]struct {
		Href string `mapstructure:"href" json:"href"`
		Name string `mapstructure:"name" json:"name"`
	} `mapstructure:"links" json:"links"`
	CreatedOn   *time.Time `mapstructure:"created_on" json:"created_on"`
	DisplayName string     `mapstructure:"display_name" json:"display_name"`
	Username    string     `mapstructure:"string" json:"string"`
	UUID        string     `mapstructure:"uuid" json:"uuid"`
}
*/

type PullRequest struct {
	Links map[string]struct {
		Href string `mapstructure:"href" json:"href"`
		Name string `mapstructure:"name" json:"name"`
	} `mapstructure:"links" json:"links"`
	Id       int    `mapstructure:"id" json:"id"`
	Title    string `mapstructure:"title" json:"title"`
	Rendered map[string]struct {
		Raw    string `mapstructure:"raw" json:"raw"`
		Markup string `mapstructure:"markup" json:"markup"`
		Html   string `mapstructure:"html" json:"html"`
	} `mapstructure:"rendered" json:"rendered"`
	State   string `mapstructure:"state" json:"state"`
	Summary struct {
		Raw    string `mapstructure:"raw" json:"raw"`
		Markup string `mapstructure:"markup" json:"markup"`
		Html   string `mapstructure:"html" json:"html"`
	} `mapstructure:"summary" json:"summary"`
	/*
		Author struct {
			Type string `mapstructure:"type" json:"type"`
			Account
		} `mapstructure:"author" json:"author"`
	*/
	CreatedOn *time.Time `mapstructure:"created_on" json:"created_on"`
	UpdatedOn *time.Time `mapstructure:"updated_on" json:"updated_on"`

	// Extra attributes added by client
	RepoSlug string `json:"repo_slug"`
}

func (p *PullRequests) Create(po *PullRequestsOptions) (interface{}, error) {
	data, err := p.buildPullRequestBody(po)
	if err != nil {
		return nil, err
	}
	urlStr := p.c.requestUrl("/repositories/%s/%s/pullrequests/", po.Owner, po.RepoSlug)
	return p.c.execute("POST", urlStr, data)
}

func (p *PullRequests) Update(po *PullRequestsOptions) (interface{}, error) {
	data, err := p.buildPullRequestBody(po)
	if err != nil {
		return nil, err
	}
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID
	return p.c.execute("PUT", urlStr, data)
}

func (p *PullRequests) Gets(po *PullRequestsOptions) ([]*PullRequest, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/"

	if po.States != nil && len(po.States) != 0 {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		for _, state := range po.States {
			query.Add("state", state)
		}
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}

	if po.Query != "" {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		query.Set("q", po.Query)
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}

	if po.Sort != "" {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		query.Set("sort", po.Sort)
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}

	resp, err := p.c.executePaginated("GET", urlStr, "", nil)

	if err != nil {
		return nil, err
	}

	return decodePullRequests(resp, po.RepoSlug)
}

func (p *PullRequests) Get(po *PullRequestsOptions) (*PullRequest, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID
	resp, err := p.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}

	return decodePullRequest(resp, po.RepoSlug)
}

func (p *PullRequests) Activities(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/activity"
	return p.c.executePaginated("GET", urlStr, "", nil)
}

func (p *PullRequests) Activity(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/activity"
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) Commits(po *PullRequestsOptions) ([]*Commit, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/commits"
	resp, err := p.c.executePaginated("GET", urlStr, "", nil)
	if err != nil {
		return nil, err
	}

	return decodeCommits(resp)
}

func (p *PullRequests) Patch(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/patch"
	return p.c.executeRaw("GET", urlStr, "")
}

func (p *PullRequests) Diff(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/diff"
	return p.c.executeRaw("GET", urlStr, "")
}

func (p *PullRequests) Merge(po *PullRequestsOptions) (interface{}, error) {
	data, err := p.buildPullRequestBody(po)
	if err != nil {
		return nil, err
	}
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/merge"
	return p.c.execute("POST", urlStr, data)
}

func (p *PullRequests) Decline(po *PullRequestsOptions) (interface{}, error) {
	data, err := p.buildPullRequestBody(po)
	if err != nil {
		return nil, err
	}
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/decline"
	return p.c.execute("POST", urlStr, data)
}

func (p *PullRequests) Approve(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/approve"
	return p.c.execute("POST", urlStr, "")
}

func (p *PullRequests) UnApprove(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/approve"
	return p.c.execute("DELETE", urlStr, "")
}

func (p *PullRequests) RequestChanges(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/request-changes"
	return p.c.execute("POST", urlStr, "")
}

func (p *PullRequests) UnRequestChanges(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/request-changes"
	return p.c.execute("DELETE", urlStr, "")
}

func (p *PullRequests) AddComment(co *PullRequestCommentOptions) (interface{}, error) {
	data, err := p.buildPullRequestCommentBody(co)
	if err != nil {
		return nil, err
	}

	urlStr := p.c.requestUrl("/repositories/%s/%s/pullrequests/%s/comments", co.Owner, co.RepoSlug, co.PullRequestID)
	return p.c.execute("POST", urlStr, data)
}

func (p *PullRequests) UpdateComment(co *PullRequestCommentOptions) (interface{}, error) {
	data, err := p.buildPullRequestCommentBody(co)
	if err != nil {
		return nil, err
	}

	urlStr := p.c.requestUrl("/repositories/%s/%s/pullrequests/%s/comments/%s", co.Owner, co.RepoSlug, co.PullRequestID, co.CommentId)
	return p.c.execute("PUT", urlStr, data)
}

func (p *PullRequests) GetComments(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/comments/"
	return p.c.executePaginated("GET", urlStr, "", nil)
}

func (p *PullRequests) GetComment(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/comments/" + po.CommentID
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) Statuses(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/statuses"
	if po.Query != "" {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		query.Set("q", po.Query)
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}

	if po.Sort != "" {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		query.Set("sort", po.Sort)
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}
	return p.c.executePaginated("GET", urlStr, "", nil)
}

func (p *PullRequests) buildPullRequestBody(po *PullRequestsOptions) (string, error) {
	body := map[string]interface{}{}
	body["source"] = map[string]interface{}{}
	body["destination"] = map[string]interface{}{}
	body["reviewers"] = []map[string]string{}
	body["title"] = ""
	body["description"] = ""
	body["message"] = ""
	body["close_source_branch"] = false

	if n := len(po.Reviewers); n > 0 {
		body["reviewers"] = make([]map[string]string, n)
		for i, uuid := range po.Reviewers {
			body["reviewers"].([]map[string]string)[i] = map[string]string{"uuid": uuid}
		}
	}

	if po.SourceBranch != "" {
		body["source"].(map[string]interface{})["branch"] = map[string]string{"name": po.SourceBranch}
	}

	if po.SourceRepository != "" {
		body["source"].(map[string]interface{})["repository"] = map[string]interface{}{"full_name": po.SourceRepository}
	}

	if po.DestinationBranch != "" {
		body["destination"].(map[string]interface{})["branch"] = map[string]interface{}{"name": po.DestinationBranch}
	}

	if po.DestinationCommit != "" {
		body["destination"].(map[string]interface{})["commit"] = map[string]interface{}{"hash": po.DestinationCommit}
	}

	if po.Title != "" {
		body["title"] = po.Title
	}

	if po.Description != "" {
		body["description"] = po.Description
	}

	if po.Message != "" {
		body["message"] = po.Message
	}

	if po.CloseSourceBranch == true || po.CloseSourceBranch == false {
		body["close_source_branch"] = po.CloseSourceBranch
	}

	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (p *PullRequests) buildPullRequestCommentBody(co *PullRequestCommentOptions) (string, error) {
	body := map[string]interface{}{}
	body["content"] = map[string]interface{}{
		"raw": co.Content,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func decodePullRequests(in interface{}, repoSlug string) ([]*PullRequest, error) {
	var pullRequests []*PullRequest

	d, ok := in.(map[string]interface{})["values"]

	var err error
	if ok {
		values, ok := d.([]interface{})
		if ok {
			for _, prVal := range values {
				pr, err := decodePullRequest(prVal, repoSlug)
				if err == nil {
					pullRequests = append(pullRequests, pr)
				}
			}
		}

	}

	return pullRequests, err
}

func decodePullRequest(in interface{}, repoSlug string) (*PullRequest, error) {
	var pullRequest PullRequest

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:   nil,
		Result:     &pullRequest,
		DecodeHook: stringToTimeHookFunc,
	})
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(in)
	pullRequest.RepoSlug = repoSlug

	return &pullRequest, err
}

func decodeCommits(in interface{}) ([]*Commit, error) {
	var commits []*Commit

	d, ok := in.(map[string]interface{})["values"]

	var err error
	if ok {
		values, ok := d.([]interface{})
		if ok {
			for _, val := range values {
				commit, err := decodeCommit(val)
				if err != nil {
					fmt.Printf("error decoding PR: %v\n", err)
					fmt.Printf("%v", val)
				} else {
					commits = append(commits, commit)
				}
			}
		}

	}

	return commits, err
}

func decodeCommit(in interface{}) (*Commit, error) {
	var commit Commit

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:   nil,
		Result:     &commit,
		DecodeHook: mapstructure.StringToTimeHookFunc("2006-01-02T15:04:05+00:00"),
	})
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(in)

	return &commit, err
}
