package config

import "errors"

const FileName = "velo.toml"

var DirNames = []string{"./", "./.config", "./config"}

var ErrConfigNotFound = errors.New("config not found")
var ErrInvalidConfig = errors.New("could not parse config file")
