package v1beta1

import (
  "encoding/json"

  "github.com/starofservice/vconf"
  next "github.com/starofservice/vconf/test/latest"
)

const Version string = "v1beta1"

type Config struct {
  Version string `json:"version"`
  Fa      string `json:"fa"`
  Fb      string `json:"fb`
  Fc      string `json:"fc`
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
    Fields: next.ConfigFields {
      A: c.Fa,
      B: c.Fb,
      C: c.Fc,
    },
  }, nil
}
