package latest

import (
  "encoding/json"
  "github.com/pkg/errors"

  "github.com/starofservice/vconf"
)

const Version string = "v1"

type ConfigFields struct {
  A string
  B string
  C string
}

type Config struct {
  Version string       `json:"version"`
  Fields  ConfigFields `json:"fields"`
}

func NewConfig() vconf.ConfigInterface {
  return new(Config)
}

func (c *Config) GetVersion() string {
  return c.Version
}

func (c *Config) Parse(data []byte) error {
  if err := json.Unmarshal(data, c); err != nil {
    return err
  }

  return nil
}

func (c *Config) Upgrade() (vconf.ConfigInterface, error) {
  return nil, errors.New("not implemented yet")
}