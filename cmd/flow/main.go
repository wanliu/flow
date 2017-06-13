package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

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

	inports  = flow.Inports(run.Flag("in", "Inports of flow graph"))
	outports = flow.Outports(run.Flag("out", "Outports of flow graph"))

	register = kingpin.Command("register", "Register a flow components package.")
	pkgfile  = register.Arg("package", "components file of package").Required().ExistingFile()

	deregister = kingpin.Command("deregister", "Deregister a flow components package.")
	depkgfile  = deregister.Arg("package", "components file of package").Required().ExistingFile()
)

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Hysios Hu")
	kingpin.CommandLine.Help = "Wanliu Bot Dialogs Server."

	switch kingpin.Parse() {
	case "run":

		if bpk, err := flow.LoadbuiltinPackage(); err != nil {
			log.Fatalf("load builtin package failed: %s", err)
		} else {
			if err := bpk.RegisterComponents(); err != nil {
				log.Fatalf("load builtin  components failed: %s", err)
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
		if net == nil {
			log.Fatalf("load graph file failed")
		}
		// log.Printf("net %# v", pretty.Formatter(net))
		// net.in
		inports.SetInPorts(net)
		outports.SetOutPorts(net)
		goflow.RunNet(net)

		// Wait for the network setup
		<-net.Ready()

		inports.Send()
		inports.Close()

		WaitNetEnd(net, outports)
		// // if len(outports) > 0 {
		// outports.Wait()
		// // }

		// <-net.Wait()

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
	case "deregister":
		rt, err := flow.LoadConfig(*rtfile)
		if err != nil {
			log.Fatalf("load Config failed: %s", err)
		}
		pkg, err := rt.Deregister(*depkgfile)
		if err != nil {
			log.Fatalf("register pkgfile failed: %s", err)
		}

		if err := rt.SaveTo(*rtfile); err != nil {
			log.Fatalf("save to Config:%s failed: %s", *rtfile, err)
		}
		fmt.Printf("Uninstalled component's package '%s#%s' successful\n", pkg.Name, pkg.Version)
	}
}

func WaitNetEnd(net *goflow.Graph, ports flow.PortsValues) error {
	var cases = make([]reflect.SelectCase, 0)
	cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(net.Wait())})

	for _, chVal := range ports {
		v := reflect.ValueOf(chVal.Chan)
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: v})
	}

	for {
		chosen, recv, recvOK := reflect.Select(cases)
		if !recvOK {
			if chosen == 0 {
				// Net Wait signal
				return nil
			} else {
				log.Printf("recv close chosen: %d", chosen)
			}
		} else {
			log.Printf("recv: %v", recv)
		}
	}
}

// func loadConfig(rtfile string) error {
// 	flow.Register(componentName, constructor)
// 	// flow.RegisterJSON(componentName, filePath)
// }
