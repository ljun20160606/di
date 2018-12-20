package di

import (
	"io"
)

func TomlLoadFile(path string) error {
	return ConfigLoadFile(path, TOML)
}

func TomlLoadReader(reader io.Reader) error {
	return ConfigLoadReader(reader, TOML)
}

func TomlLoad(content string) error {
	return ConfigLoad(content, TOML)
}
