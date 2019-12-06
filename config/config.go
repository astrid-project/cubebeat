package config

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

type Config struct {
	ConfigInputs *ConfigInputs `config:"config.inputs"`
}

type ConfigInputs struct {
	Path   string `config:"path"`
	Reload Reload `config:"reload"`
}

type Reload struct {
	Enabled bool          `config:"enabled"`
	Period  time.Duration `config:"period"`
}

type CubeInput struct {
	Name           string        `config:"name"`
	Enabled        bool          `config:"enabled"`
	Period         time.Duration `config:"period"`
	PolycubeAPIURL string        `config:"polycube.api-url"`
	Done           bool
	Reload         *CubeInput
	Req            *http.Request
	Ticker         *time.Ticker
}

var (
	DefaultConfig = Config{}
)

func LoadCubeInputs(path string) (cubeInputs map[string]*CubeInput, err error) {
	var cubeInputFilenames []string

	if cubeInputFilenames, err = filepath.Glob(path); err != nil {
		return nil, err
	}
	logp.Info("[Cubebeat] load cube inputs from %+v", cubeInputFilenames)

	var c *common.Config
	c, err = common.LoadFiles(cubeInputFilenames...)

	var cubeInputList []CubeInput
	if err = c.Unpack(&cubeInputList); err != nil {
		return nil, err
	}

	cubeInputs = make(map[string]*CubeInput)
	var cubeInputNames []string
	for idx, cubeInput := range cubeInputList {
		if _, contains := cubeInputs[cubeInput.Name]; contains {
			return nil, fmt.Errorf("[Cubebeat] duplicate cube name: %s", cubeInput.Name)
		}

		cubeInput.Done = false
		cubeInput.Reload = nil
		cubeInputs[cubeInput.Name] = &cubeInputList[idx]
		cubeInputNames = append(cubeInputNames, cubeInput.Name)
	}
	logp.Info("[Cubebeat] found %d cube inputs: %v", len(cubeInputs), cubeInputNames)

	return cubeInputs, nil
}
