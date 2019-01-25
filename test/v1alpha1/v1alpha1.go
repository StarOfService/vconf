package v1alpha1

import (
  "encoding/json"

  "github.com/starofservice/vconf"
  next "github.com/starofservice/vconf/test/v1alpha2"
)

const Version string = "v1alpha1"

type Config struct {
  Version string `json:"version"`
  FieldA  string `json:"fieldA"`
  FieldB  string `json:"fieldB"`
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
  return &next.Config{
    Version: next.Version,
    Fa: c.FieldA,
    Fb: c.FieldB,
  }, nil
}
