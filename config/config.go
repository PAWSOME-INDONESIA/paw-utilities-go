package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func New(path string, object interface{}) error {
	// - check file does exists

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.Wrapf(err, "config file %s does not exists!", path)
	}

	dir := getDirectory(path)
	file, err := getFile(path)

	if err != nil {
		return err
	}

	v := viper.New()
	v.SetConfigName(file)
	v.AddConfigPath(dir)
	v.SetConfigType("properties")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return errors.Wrapf(err, "failed to read %s file", path)
	}

	if err := v.Unmarshal(&object); err != nil {
		return errors.Wrap(err, "failed to unmarshal config to object")
	}

	return nil
}

func getDirectory(path string) string {
	splits := strings.Split(path, "/")
	return strings.Join(splits[:len(splits)-1], "/")
}

func getFile(path string) (string, error) {
	splits := strings.Split(path, "/")
	last := splits[len(splits)-1]

	files := strings.Split(last, ".")

	if len(files) != 2 {
		return "", errors.New(fmt.Sprintf("invalid config file %v", files))
	}

	return files[0], nil
}
