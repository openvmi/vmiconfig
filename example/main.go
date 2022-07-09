package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/openvmi/vmiconfig"
	"github.com/spf13/viper"
)

func CHECK_ERROR_THEN_EXIT(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(2)
}

type AppConfiguration struct {
	parser *vmiconfig.ConfigurationParser
	Port   int
	Host   string
}

func (app *AppConfiguration) Init() error {
	app.parser = &vmiconfig.ConfigurationParser{
		ViperInstance: nil,
		ConfigName:    "appConfig",
		ConfigType:    "json",
		ConfigPath:    []string{".", "/myapp/path"},
		DefaultValue: map[string]interface{}{
			"port": 23,
			"host": "127.0.0.1",
		},
	}
	if err := app.parser.LoadConfiguration(); err != nil {
		return err
	}
	_p, err := app.getPort()
	if err != nil {
		return err
	}
	app.Port = _p

	_h, err := app.getHost()
	if err != nil {
		return err
	}
	app.Host = _h
	return nil
}

func (app *AppConfiguration) getPort() (int, error) {
	p, err := app.parser.GetValue("port")
	if err != nil {
		return 0, err
	}
	fmt.Println(reflect.TypeOf(p))
	//解码过程中，number会被解码成float64
	pInt, ok := p.(float64)
	if !ok {
		fmt.Println("not int")
		return 0, errors.New("type error")
	}
	return int(pInt), nil
}

func (app *AppConfiguration) getHost() (string, error) {
	h, err := app.parser.GetValue("host")
	if err != nil {
		return "", err
	}
	hStr, ok := h.(string)
	if !ok {
		return "", errors.New("type error")
	}
	return hStr, nil
}

func main() {
	config := AppConfiguration{}
	if err := config.Init(); err != nil {
		fmt.Println(err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println(config.parser.GenerateDefaultConfiguration("./template.json"))
		}
	} else {
		fmt.Println("Read configuration success")
		//port, err := config.GetPort()
		CHECK_ERROR_THEN_EXIT(err)
		//fmt.Println(port)
		fmt.Println(config.Port)
		fmt.Println(config.Host)
	}
}
