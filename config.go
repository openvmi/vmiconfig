package vmiconfig

import (
	"errors"

	"github.com/spf13/viper"
)

type ConfigurationParser struct {
	hasLoadConfig bool
	loadConfig    bool
	ViperInstance *viper.Viper
	ConfigName    string
	ConfigType    string
	ConfigPath    []string
	DefaultValue  map[string]interface{}
}

func (app *ConfigurationParser) LoadConfiguration() error {
	if app.ViperInstance == nil {
		app.ViperInstance = viper.New()
	}
	app.hasLoadConfig = true
	for k, v := range app.DefaultValue {
		app.ViperInstance.SetDefault(k, v)
	}
	app.ViperInstance.SetConfigName(app.ConfigName)
	app.ViperInstance.SetConfigType(app.ConfigType)
	for _, p := range app.ConfigPath {
		app.ViperInstance.AddConfigPath(p)
	}
	if err := app.ViperInstance.ReadInConfig(); err != nil {
		app.loadConfig = false
		return err
	}
	app.loadConfig = true
	return nil
}

func (app *ConfigurationParser) GenerateDefaultConfiguration(defaultPath string) error {
	if app.hasLoadConfig && !app.loadConfig {
		return app.ViperInstance.WriteConfigAs(defaultPath)
	}
	if !app.hasLoadConfig {
		err := "please call LoadConfiguration before generate default configuration"
		return errors.New(err)
	}
	if app.loadConfig {
		err := "configuration load success, can not generate default configuration"
		return errors.New(err)
	}
	return nil
}

func (app *ConfigurationParser) GetValue(key string) (interface{}, error) {
	if !app.hasLoadConfig {
		if err := app.LoadConfiguration(); err != nil {
			return 0, err
		}
	}
	if !app.ViperInstance.IsSet(key) {
		return 0, errors.New("don't have specified filed")
	}
	return app.ViperInstance.Get(key), nil
}
