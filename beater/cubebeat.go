package beater

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"gitlab.com/astrid-repositories/cubebeat/config"
)

type Cubebeat struct {
	done       chan struct{}
	config     *config.Config
	cubeInputs map[string]*config.CubeInput
	client     beat.Client
}

func New(b *beat.Beat, cfg *common.Config) (bt beat.Beater, err error) {
	c := config.DefaultConfig
	if err = cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("[Cubebeat] error reading config file: %v", err)
	}

	if c.ConfigInputs == nil {
		return nil, fmt.Errorf("[Cubebeat] no config inputs defined")
	}

	bt = &Cubebeat{
		done:       make(chan struct{}),
		config:     &c,
		cubeInputs: make(map[string]*config.CubeInput),
	}

	return bt, nil
}

func (bt *Cubebeat) Run(b *beat.Beat) (err error) {
	logp.Info("[Cubebeat] hit CTRL-C to stop")

	if bt.client, err = b.Publisher.Connect(); err != nil {
		return err
	}

	if err = RunProcessCubes(bt); err != nil {
		return err
	}

	logp.Info("[Cubebeat] reload config each %s", bt.config.ConfigInputs.Reload.Period)
	ticker := time.NewTicker(bt.config.ConfigInputs.Reload.Period)
	counterReload := 1

	for {
		select {
		case <-bt.done:
			logp.Info("[Cubebeat] terminated")
			return nil
		case <-ticker.C:
		}

		logp.Info("[Cubebeat] reload num.: %d", counterReload)
		if err = RunProcessCubes(bt); err != nil {
			return err
		}
		counterReload++
	}

	return nil
}

func (bt *Cubebeat) Stop() {
	logp.Info("[Cubebeat] stopping")
	bt.client.Close()
	close(bt.done)
}

func RunProcessCubes(bt *Cubebeat) (err error) {
	var cubeInputs map[string]*config.CubeInput
	if cubeInputs, err = config.LoadCubeInputs(bt.config.ConfigInputs.Path); err != nil {
		return err
	}

	for cubeName, cubeInput := range cubeInputs {
		if _, contains := bt.cubeInputs[cubeName]; contains {
			bt.cubeInputs[cubeName].Reload = cubeInput
		} else {
			go ProcessCube(bt, cubeInput)
		}
	}

	for cubeName, cubeInput := range bt.cubeInputs {
		if _, contains := cubeInputs[cubeName]; !contains {
			cubeInput.Done = true
		}
	}

	return nil
}

func ProcessCube(bt *Cubebeat, cubeInput *config.CubeInput) error {
	if err := InitCube(cubeInput, bt); err != nil {
		logp.Critical("[Cube %s] %s", cubeInput.Name, err)
		return FinishCube(bt, cubeInput, err)
	} else {
		counterProcessCube := 1

		for {
			select {
			case <-bt.done:
				return FinishCube(bt, cubeInput, nil)
			case <-cubeInput.Ticker.C:
			}

			if cubeInput.Done {
				return FinishCube(bt, cubeInput, nil)
			}

			ReloadCube(cubeInput)

			if cubeInput.Enabled && cubeInput.Req != nil {
				logp.Info("[Cube %s] process num.: %d", cubeInput.Name, counterProcessCube)

				resp, err := http.DefaultClient.Do(cubeInput.Req)
				if err != nil {
					logp.Critical("[Cube %s] %s", cubeInput.Name, err)
				} else {
					defer resp.Body.Close()

					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						logp.Critical("[Cube %s] %s", cubeInput.Name, err)
					} else if data, err := ParseData(string(body)); err != nil {
						logp.Critical("[Cube %s] %s", cubeInput.Name, err)
					} else {
						event := beat.Event{
							Timestamp: time.Now(),
							Fields:    data,
						}
						bt.client.Publish(event)
						logp.Info("[Cube %s] event sent: %v", cubeInput.Name, event)
						counterProcessCube++
					}
				}
			}
		}
	}
	return nil
}

func InitCube(cubeInput *config.CubeInput, bt *Cubebeat) error {
	NewConnection(cubeInput)

	logp.Info("[Cube %s] Process interval: %s", cubeInput.Name, cubeInput.Period)
	cubeInput.Ticker = time.NewTicker(cubeInput.Period)

	if _, contains := bt.cubeInputs[cubeInput.Name]; contains {
		return fmt.Errorf("already present")
	} else {
		bt.cubeInputs[cubeInput.Name] = cubeInput
	}

	return nil
}

func ReloadCube(cubeInput *config.CubeInput) {
	newCubeInput := cubeInput.Reload
	if newCubeInput != nil {
		if cubeInput.Req == nil || cubeInput.PolycubeAPIURL != newCubeInput.PolycubeAPIURL {
			cubeInput.PolycubeAPIURL = newCubeInput.PolycubeAPIURL
			NewConnection(cubeInput)
		}

		if cubeInput.Period != newCubeInput.Period {
			cubeInput.Period = newCubeInput.Period
			cubeInput.Ticker.Stop()
			logp.Info("[Cube %s] process new interval: %s", cubeInput.Name, cubeInput.Period)
			cubeInput.Ticker = time.NewTicker(cubeInput.Period)
		}

		if cubeInput.Enabled != newCubeInput.Enabled {
			cubeInput.Enabled = newCubeInput.Enabled
			if cubeInput.Enabled {
				logp.Info("[Cube %s] enabled", cubeInput.Name)
			} else {
				logp.Info("[Cube %s] disabled", cubeInput.Name)
			}
		}

		cubeInput.Reload = nil
	}
}

func NewConnection(cubeInput *config.CubeInput) {
	logp.Info("[Cube %s] create HTTP Connection to Polycube API URL: %s", cubeInput.Name, cubeInput.PolycubeAPIURL)
	var err error
	if cubeInput.Req, err = http.NewRequest(http.MethodGet, cubeInput.PolycubeAPIURL, nil); err != nil {
		logp.Critical("[Cube %s] not possible to connect to Polycube API URL: %s", cubeInput.Name, cubeInput.PolycubeAPIURL, err)
	}
}

func FinishCube(bt *Cubebeat, cubeInput *config.CubeInput, err error) error {
	delete(bt.cubeInputs, cubeInput.Name)
	if err != nil {
		logp.Critical("[Cube %s] terminated: %s", cubeInput.Name, err)
	} else {
		logp.Info("[Cube %s] terminated", cubeInput.Name)
	}
	return err
}

func ParseData(input string) (output map[string]interface{}, err error) {
	if err = json.Unmarshal([]byte(input), &output); err != nil {
		return nil, err
	}
	return output, nil
}
