package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/paribu/acervus-cli/src/settings"
)

type GenerateBoilerplateRequest struct {
	AbiFile      string `json:"abi_file"`
	SettingsFile string `json:"settings_file"`
	GraphqlFile  string `json:"graphql_file"`
}

type GenerateBoilerplateResponse struct {
	Files []struct {
		Path     string `json:"path"`
		Contents string `json:"contents"`
	}
}

func (a *ProjectManagerAPI) GenerateBoilerplate(settingsFilePath string) (*GenerateBoilerplateResponse, error) {
	yamlFile, err := settings.NewProjectFromFile(settingsFilePath)
	if err != nil {
		return nil, err
	}

	abiFile, err := os.ReadFile(yamlFile.Sources[0].Source.Abi)
	if err != nil {
		return nil, err
	}

	graphQLFile, err := os.ReadFile(yamlFile.Schema)
	if err != nil {
		return nil, err
	}

	yamlStr, err := yamlFile.ToString()
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(GenerateBoilerplateRequest{
		SettingsFile: yamlStr,
		AbiFile:      string(abiFile),
		GraphqlFile:  string(graphQLFile),
	})
	if err != nil {
		return nil, err
	}

	resp, err := a.makeAuthenticatedAPIRequest(
		http.MethodPost,
		endpoints.generate.boilerplate,
		RequestData{Body: body},
	)
	if err != nil {
		return nil, err
	}

	var generateBoilerplateResp GenerateBoilerplateResponse
	err = json.Unmarshal(resp, &generateBoilerplateResp)
	if err != nil {
		return nil, err
	}

	return &generateBoilerplateResp, nil
}
