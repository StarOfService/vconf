# VConf [![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg)](https://godoc.org/github.com/StarOfService/vconf)

This Go library is aimed to simplify working with versioned configs and metadata. 

VConf isn't bounded to a specific configuration/metadata format. Instead, it expects that your application converts a string to an object. VConf just provides a standard way to organize your code and orchestrates the process.

VConf supports two formats of versioning: Semantic and Kubernetes.

Brief description of the workflow
=================================
1) You create a package for your config/metadata and a subpackage for each version.
2) Each such version must implement `ConfigInterface` interface.
3) `ConfigInterface` interface requires `Upgrade` method which implements upgrade up to the next version.
4) After all versions are implemented and registered, you are able to provide any version of your config/metadata and will receive the latest object as an output.

Example
=======
You can find a working example at `test` subfolder. I'll describe here only the key points.

Versions
-------- 
- Every version is located in its own subpackage:
```
yourconfig/latest/
...
yourconfig/v1alpha1/
yourconfig/v1alpha2/
yourconfig/v1beta1/
yourconfig/yourconfig.go
```
- Every version Implements `ConfigInterface`:
```
type ConfigInterface interface {
  GetVersion() string
  Parse([]byte) error
  Upgrade() (ConfigInterface, error)
}
```
- `Upgrade` method creates an item of the next version and thus implements upgrade:
```
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
```

Main package
------------
- Taking into account that VConf doesn't implement any specific format, it's your responsibility to extract a version of the inbound document. `GetCurrentVersion`  function shows how it can be done for JSON format. For YAML it will be pretty much the same.
```
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
```
- Next you create a new instance of `vconf.SchemaHandler` (using `vconf.NewSchemaHandler` function)
```
sh := vconf.NewSchemaHandler(latest.Version)
```
- All available versions have to be registered at the SchemaHandler using `RegVersion` method.
It's handy to manage this list at the top of the package, next to the global constants and variables
```
var schemaVersions = map[string]func() vconf.ConfigInterface{
  latest.Version: latest.NewConfig,
  v1beta1.Version: v1beta1.NewConfig,
  v1alpha2.Version: v1alpha2.NewConfig,
  v1alpha1.Version: v1alpha1.NewConfig,
}
```
and later to register these versions after the `SchemaHandler` instance is created:
```
  for k, v := range schemaVersions {
    sh.RegVersion(k, v)
  }
```
- Now you are ready to run `GetLatestConfig` method. It expects current version and the initial document as arguments and retures an object which represents the latest version of your config.
```
  cfg, err := sh.GetLatestConfig(current, []byte(suite))
  if err != nil {
    return err
  }
```
- `GetLatestConfig` return a vconf.ConfigInterface, so a very last step is to convert the received object to the latest struct type:
```
parsedCfg := cfg.(*latest.Config)
```

If you need additional fields or methods for the config entity, you can create a new structure at the top level of your config package:
```
type LatestConfig struct {
  Data *latest.Config
  ExtraField1 string
  ...
}
...
  parsedCfg := cfg.(*latest.Config)
  &LatestConfig {
    Data: *parsedCfg,
    ExtraField1: "...",
  }
...
func (c *LatestConfig) DoSomething {
  ...
}
```
