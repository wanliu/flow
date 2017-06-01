package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/kr/pretty"
	"github.com/wanliu/flow"
	goflow "github.com/wanliu/goflow"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	// debug           = kingpin.Flag("debug", "Enable debug mode.").Bool()
	ip     = kingpin.Flag("ip", "IP address to ping.").Default("127.0.0.1").IP()
	port   = kingpin.Flag("port", "bind serve port").Default("8081").Int()
	rtfile = kingpin.Flag("Config", "flow Config environment").Default("config.json").ExistingFile()

	run     = kingpin.Command("run", "Run a brain flow file.")
	runfile = run.Arg("flowfile", "File of flow program").Required().ExistingFile()

	inports  = Inports(run.Flag("in", "Inports of flow graph"))
	outports = Outports(run.Flag("out", "Outports of flow graph"))

	register = kingpin.Command("register", "Register a flow components package.")
	pkgfile  = register.Arg("package", "components file of package").Required().ExistingFile()
)

type PortsValues map[string]interface{}

func (in PortsValues) Set(value string) error {
	parts := strings.Split(value, ",")
	for _, part := range parts {
		if exp := strings.Split(part, "="); len(exp) == 2 {
			name := exp[0]

			typ := exp[1]
			switch typ {
			case "string":
				in[name] = make(chan string)
				// case "int":
				// 	in[name] = reflect.TypeOf(0)
				// case "bool":
				// 	in[name] = reflect.TypeOf(true)
				// case "float":
				// 	in[name] = reflect.TypeOf(0.0)
				// default:
				// 	in[name] = reflect.TypeOf("")
			}
		} else {
			name := exp[0]
			in[name] = make(chan string)
		}
	}

	return nil
}

func (in PortsValues) String() string {
	return ""
}

func (in PortsValues) SetInPorts(net *goflow.Graph) {

	for port, ch := range in {
		// c := reflect.MakeChan(typ, 0)
		net.SetInPort(port, ch)
	}
}

func (in PortsValues) SetOutPorts(net *goflow.Graph) {

	for port, ch := range in {
		// c := reflect.MakeChan(typ, 0)
		net.SetOutPort(port, ch)
	}
}

func (in PortsValues) Wait() {
	for _, ch := range in {
		v := reflect.ValueOf(ch)
		v.Recv()
	}
}

func (in PortsValues) Close() {
	for _, ch := range in {
		v := reflect.ValueOf(ch)
		v.Close()
	}
}

func Inports(s kingpin.Settings) PortsValues {
	var result = make(map[string]interface{})
	s.SetValue(PortsValues(result))
	if len(result) == 0 {
		result["Start"] = make(chan string)
	}
	return result
}

func Outports(s kingpin.Settings) PortsValues {
	var result = make(map[string]interface{})
	s.SetValue(PortsValues(result))
	if len(result) == 0 {
		result["Out"] = make(chan string)
	}
	return result
}

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Hysios Hu")
	kingpin.CommandLine.Help = "Wanliu Bot Dialogs Server."

	switch kingpin.Parse() {
	case "run":

		if bpk, err := flow.LoadBuilitnPackage(); err != nil {
			log.Fatalf("load builitn package failed: %s", err)
		} else {
			if err := bpk.RegisterComponents(); err != nil {
				log.Fatalf("load builitn  components failed: %s", err)
			}
		}

		rt, err := flow.LoadConfig(*rtfile)

		if err != nil {
			log.Fatalf("load Config failed: %s", err)
		}

		if err := rt.LoadComponents(); err != nil {
			log.Fatalf("load components failed: %s", err)
		}

		buf, err := ioutil.ReadFile(*runfile)
		if err != nil {
			log.Fatalf("open run file failed: %s", err)
		}
		net := goflow.ParseJSON(buf)
		log.Printf("net %# v", pretty.Formatter(net))
		inports.SetInPorts(net)
		outports.SetOutPorts(net)
		goflow.RunNet(net)

		// Wait for the network setup
		<-net.Ready()

		outports.Wait()
		inports.Close()
		<-net.Wait()

	case "register":
		rt, err := flow.LoadConfig(*rtfile)
		if err != nil {
			log.Fatalf("load Config failed: %s", err)
		}
		pkg, err := rt.Register(*pkgfile)
		if err != nil {
			log.Fatalf("register pkgfile failed: %s", err)
		}

		if err := rt.SaveTo(*rtfile); err != nil {
			log.Fatalf("save to Config:%s failed: %s", *rtfile, err)
		}
		fmt.Printf("Installed component's package '%s#%s' successful\n", pkg.Name, pkg.Version)
	}
}

// func loadConfig(rtfile string) error {
// 	flow.Register(componentName, constructor)
// 	// flow.RegisterJSON(componentName, filePath)
// }
