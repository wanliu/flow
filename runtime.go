package flow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Masterminds/semver"
)

type Runtime struct {
	Name        string
	RuntimeType string
	Packages    []Package
}

type PType int

const (
	BinaryType PType = iota
	JavascriptType
)

func NewRuntime() *Runtime {
	return &Runtime{}
}

func LoadRuntime(filename string) (*Runtime, error) {
	var rt Runtime
	buf, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf, &rt)
	if err != nil {
		return nil, err
	}

	return &rt, nil
}

func (rt *Runtime) Register(pkgfile string) error {
	pkg, err := LoadPackage(pkgfile)
	if err != nil {
		return err
	}

	for i, pk := range rt.Packages {
		if pkg.Name == pk.Name {
			v1, err := semver.NewConstraint(pk.Version)
			if err != nil {
				return err
			}

			v2, err := semver.NewVersion(pkg.Version)
			if err != nil {
				return err
			}

			if !v1.Check(v2) {
				return fmt.Errorf("invalid version %s don't match %s", pk.Version, pkg.Version)
			}

			rt.Packages[i] = *pkg
			return nil
		}
	}
	rt.Packages = append(rt.Packages, *pkg)
	return nil
}

func (rt *Runtime) SaveTo(filename string) error {
	var buf []byte

	buf, err := json.Marshal(rt)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, buf, 0644)
}
