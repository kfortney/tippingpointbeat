package beater

import (
	"fmt"
	"bytes"
	"encoding/json"
	"time"
	"io/ioutil"
	"net/http"
	
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/kfortney/tippingpointbeat/config"
)

type Tippingpointbeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
	lastIndexTime time.Time
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Tippingpointbeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Tippingpointbeat) Run(b *beat.Beat) error {
	logp.Info("tippingpointbeat is running! Hit CTRL-C to stop it.")
	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
			"counter":    counter,
		}
		bt.client.PublishEvent(event)
		logp.Info("Event sent")
		counter++
	jsonData := map[string]string{"firstname": "John", "lastname": "Doe"}
      jsonValue, _ := json.Marshal(jsonData)
      request, _ := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewBuffer(jsonValue))
      request.Header.Set("Content-Type", "multipart/form-data")
      client := &http.Client{}
      response, err := client.Do(request)
      if err != nil {
          fmt.Printf("The HTTP request failed with error %s\n", err)
      } else {
          data, _ := ioutil.ReadAll(response.Body)
          fmt.Println(string(data))
      }


	}
}

func (bt *Tippingpointbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
