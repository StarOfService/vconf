package test

import (
  "encoding/json"
  "testing"

  "github.com/pkg/errors"

  "github.com/starofservice/vconf"
  "github.com/starofservice/vconf/test/latest"
  "github.com/starofservice/vconf/test/v1beta1"
  "github.com/starofservice/vconf/test/v1alpha2"
  "github.com/starofservice/vconf/test/v1alpha1"
)

var schemaVersions = map[string]func() vconf.ConfigInterface{
  latest.Version: latest.NewConfig,
  v1beta1.Version: v1beta1.NewConfig,
  v1alpha2.Version: v1alpha2.NewConfig,
  v1alpha1.Version: v1alpha1.NewConfig,
}

func GetCurrentVersion(data []byte) (string, error) {
  type VersionStruct struct {
    Version string `json:"version"`
  }
  version := &VersionStruct{}
  if err := json.Unmarshal(data, version); err != nil {
    return "", errors.Wrap(err, "parsing api version")
  }
  return version.Version, nil
}

func TestParseConfig(t *testing.T) {
  suite := `{"version":"v1alpha1","fieldA":"fieldA","fieldB":"fieldB"}`
  assert := `{"version":"v1","fields":{"A":"fieldA","B":"fieldB","C":"fieldC"}}`

  current, err := GetCurrentVersion([]byte(suite))
  if err != nil {
    t.Errorf(err.Error())
  }

  sh := vconf.NewSchemaHandler(latest.Version)
  for k, v := range schemaVersions {
    sh.RegVersion(k, v)
  }

  cfg, err := sh.GetLatestConfig(current, []byte(suite))
  if err != nil {
    t.Errorf(err.Error())
  }

  parsedCfg := cfg.(*latest.Config)
  result, err := json.Marshal(parsedCfg)
  if err != nil {
    t.Errorf(err.Error())
  }

  if assert != string(result) {
    t.Errorf("Test suite object %v doesn't match to the generated data %v", assert, string(result))
  }
}