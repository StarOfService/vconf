package v1alpha2

import (
  "encoding/json"

  "github.com/starofservice/vconf"
  next "github.com/starofservice/vconf/test/v1beta1"
)

const Version string = "v1alpha2"

type Config struct {
  Version string `json:"version"`
  Fa      string `json:"fa"`
  Fb      string `json:"fb`
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
    Fa: c.Fa,
    Fb: c.Fb,
    Fc: "fieldC",
  }, nil
}
