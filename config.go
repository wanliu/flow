package flow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Masterminds/semver"
)

type Config struct {
	Name       string
	ConfigType string
	Packages   []Package
}

type PType int

const (
	BinaryType PType = iota
	JavascriptType
)

func NewConfig() *Config {
	return &Config{}
}

func LoadConfig(filename string) (*Config, error) {
	var rt Config
	buf, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf, &rt)
	if err != nil {
		return nil, err
	}

	rt.LoadPackages()
	return &rt, nil
}

func (rt *Config) Register(pkgfile string) (*Package, error) {
	pkg, err := LoadPackage(pkgfile)
	if err != nil {
		return nil, err
	}

	for i, pk := range rt.Packages {
		if pkg.Name == pk.Name {
			v1, err := semver.NewConstraint(pk.Version)
			if err != nil {
				return nil, err
			}

			v2, err := semver.NewVersion(pkg.Version)
			if err != nil {
				return nil, err
			}

			if !v1.Check(v2) {
				return nil, fmt.Errorf("invalid version %s don't match %s", pk.Version, pkg.Version)
			}

			rt.Packages[i] = *pkg
			return pkg, nil
		}
	}
	rt.Packages = append(rt.Packages, *pkg)
	return pkg, nil
}

func (rt *Config) SaveTo(filename string) error {
	var buf []byte

	buf, err := json.MarshalIndent(rt, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, buf, 0644)
}

func (rt *Config) LoadPackages() error {
	for i, pkg := range rt.Packages {
		pk, err := LoadPackage(pkg.Path)
		if err != nil {
			return err
		}

		rt.Packages[i] = *pk
	}

	return nil
}

func (rt *Config) LoadComponents() error {
	for _, pkg := range rt.Packages {
		if err := pkg.RegisterComponents(); err != nil {
			return err
		}
	}

	return nil
}
