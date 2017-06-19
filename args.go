package flow

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/wanliu/flow/context"
	goflow "github.com/wanliu/goflow"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type PortValue struct {
	Chan  interface{}
	Value reflect.Value
}
type PortsValues map[string]PortValue

func (in PortsValues) Set(value string) error {
	parts := strings.Split(value, ",")
	for _, part := range parts {
		var name, val string
		if exp := strings.Split(part, "="); len(exp) == 2 {
			name = exp[0]
			val = exp[1]

		} else {
			name = exp[0]
		}

		ch, n, v := GuessType(name, val)
		in[n] = PortValue{
			Chan:  ch,
			Value: reflect.ValueOf(v),
		}

		// in[name] = PortValue{Chan: make(chan string)}
	}

	return nil
}

func (in PortsValues) String() string {
	return ""
}

func (in PortsValues) SetInPorts(net *goflow.Graph) {

	for port, chVal := range in {
		net.SetInPort(port, chVal.Chan)
	}
}

func (in PortsValues) SetOutPorts(net *goflow.Graph) {
	for port, chVal := range in {
		net.SetOutPort(port, chVal.Chan)
	}
}

func (in PortsValues) Send() {
	for _, chVal := range in {
		// c := reflect.MakeChan(typ, 0)
		if chVal.Value.IsValid() {
			reflect.ValueOf(chVal.Chan).Send(chVal.Value)
		}
	}
}

func (in PortsValues) Wait() {
	for _, chVal := range in {
		v := reflect.ValueOf(chVal.Chan)
		v.Recv()
		// v.Recv()
	}
}

func (in PortsValues) Close() {
	for _, chVal := range in {
		v := reflect.ValueOf(chVal.Chan)
		v.Close()
	}
}

func Inports(s kingpin.Settings) PortsValues {
	var result = make(map[string]PortValue)
	s.SetValue(PortsValues(result))
	return result
}

func Outports(s kingpin.Settings) PortsValues {
	var result = make(map[string]PortValue)
	s.SetValue(PortsValues(result))
	return result
}

func GuessType(typ string, val string) (ch interface{}, name string, v interface{}) {
	var err error
	vals := strings.Split(typ, ":")

	l := len(vals)
	if l == 2 {
		name = vals[0]
		typ := vals[1]
		switch typ {
		case "string":
			if len(val) > 0 {
				if v, err = strconv.Unquote(val); err != nil {
					v = val
				}
			} else {
				v = nil
			}
			ch = make(chan string)

			return
		case "int":
			v, err = strconv.Atoi(val)

			if err != nil {
				panic(fmt.Sprintf("invalid int value %s", err))
			}

			ch = make(chan int)
			return
		case "bool":
			v, err = strconv.ParseBool(val)
			if err != nil {
				panic(fmt.Sprintf("invalid bool value %s", err))
			}
			ch = make(chan bool)
			return
		case "float":
			v, err = strconv.ParseFloat(val, 64)
			if err != nil {
				panic(fmt.Sprintf("invalid bool value %s", err))
			}
			ch = make(chan float64)
			return
		case "context":
			v, _ = context.NewContext()
			ch = make(chan context.Context)

			return
		default:
			ch = make(chan string)
			if len(val) > 0 {
				if v, err = strconv.Unquote(val); err != nil {
					v = val
				}
			} else {
				v = nil
			}
			return
		}
	} else if l < 2 {
		name = typ
		if v, err = strconv.ParseBool(val); err == nil {
			ch = make(chan bool)
			return
		}

		if v, err = strconv.ParseInt(val, 10, 64); err == nil {
			ch = make(chan int)
			return
		}

		if v, err = strconv.ParseFloat(val, 64); err == nil {
			ch = make(chan float64)
			return
		}
		if len(val) > 0 {
			if v, err = strconv.Unquote(val); err != nil {
				v = val
			}
		} else {
			v = nil
		}
		ch = make(chan string)
		return
	} else {
		panic("split val failed")
	}
}
