package bitbucket

import (
	"time"

	"github.com/mitchellh/mapstructure"
)

type Deployments struct {
	c *Client
}

type Deployment struct {
	Deployable struct {
		Type   string `mapstructure:"deployable" json:"deployable"`
		Name   string `mapstructure:"name" json:"name"`
		Commit struct {
			Hash string `mapstructure:"hash" json:"hash"`
		} `mapstructure:"commit" json:"commit"`
		CreatedOn *time.Time `mapstructure:"created_on" json:"created_on"`
	} `mapstructure:"deployable" json:"deployable"`
	Number  int `mapstructure:"number" json:"number"`
	Release struct {
		Name      string     `mapstructure:"name" json:"name"`
		CreatedOn *time.Time `mapstructure:"created_on" json:"created_on"`
	}
	State struct {
		Type string `mapstructure:"deployment_state_undeployed" json:"deployment_state_undeployed"`
		Name string `mapstructure:"name" json:"name"`
	}
}

func (p *Deployments) Gets(po *DeploymentsOptions) ([]*Deployment, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/deployments/"

	resp, err := p.c.executePaginated("GET", urlStr, "", nil)

	if err != nil {
		return nil, err
	}

	return decodeDeployments(resp)
}

func (p *Deployments) Get(po *DeploymentsOptions) (*Deployment, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/deployments/" + po.Uuid
	resp, err := p.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}

	return decodeDeployment(resp)
}

func decodeDeployments(in interface{}) ([]*Deployment, error) {
	var deployments []*Deployment

	d, ok := in.(map[string]interface{})["values"]

	var err error
	if ok {
		values, ok := d.([]interface{})
		if ok {
			for _, val := range values {
				deploy, err := decodeDeployment(val)
				if err == nil {
					deployments = append(deployments, deploy)
				}
			}
		}

	}

	return deployments, err
}

func decodeDeployment(in interface{}) (*Deployment, error) {
	var Deployment Deployment

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:   nil,
		Result:     &Deployment,
		DecodeHook: mapstructure.StringToTimeHookFunc("2006-01-02T15:04:05.999999Z"),
	})
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(in)

	return &Deployment, err
}
