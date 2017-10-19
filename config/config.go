// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period   time.Duration `config:"period"`
	Target   string        `config:"target"`
	Filename string        `config:"filename"`
}

var DefaultConfig = Config{
	Period:   1 * time.Second,
	Target:   "https://httpbin.org/post",
	Filename: "payload.id",
}
