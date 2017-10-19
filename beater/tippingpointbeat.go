package beater

import (
	"bytes"
	"fmt"
	//      "encoding/json"
	"io"
	"io/ioutil"
	"time"
	//	"log"
	"mime/multipart"
	"net/http"
	"os"
	//      "path/filepath"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/kfortney/tippingpointbeat/config"
)

func (bt *Tippingpointbeat) postFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil
}

type Tippingpointbeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

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

	filename := bt.config.Filename
	target := bt.config.Target

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
			"target":     target,
			"filename":   filename,
		}
		bt.client.PublishEvent(event)
		logp.Info("Event sent")
		counter++
		bt.postFile(bt.config.Filename, bt.config.Target)
	}
}

func (bt *Tippingpointbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
