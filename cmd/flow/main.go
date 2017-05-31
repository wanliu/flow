package main

import (
	"io/ioutil"
	"log"

	goflow "github.com/trustmaster/goflow"
	"github.com/wanliu/flow"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	// debug           = kingpin.Flag("debug", "Enable debug mode.").Bool()
	ip      = kingpin.Flag("ip", "IP address to ping.").Default("127.0.0.1").IP()
	port    = kingpin.Flag("port", "bind serve port").Default("8081").Int()
	run     = kingpin.Command("run", "Run a brain flow file.")
	runfile = run.Arg("flowfile", "file of flow programe").Required().ExistingFile()
	rtfile  = kingpin.Flag("runtime", "flow runtime environment").Default("runtime.json").ExistingFile()

	register = kingpin.Command("register", "Register a flow components package.")
	pkgfile  = register.Arg("package", "components file of package").Required().ExistingFile()
)

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Hysios Hu")
	kingpin.CommandLine.Help = "Wanliu Bot Dialogs Server."

	switch kingpin.Parse() {
	case "run":
		buf, err := ioutil.ReadFile(*runfile)
		if err != nil {
			log.Fatalf("open run file failed: %s", err)
		}
		graph := goflow.ParseJSON(buf)
		log.Printf("%#v", graph)
	case "register":
		rt, err := flow.LoadRuntime(*rtfile)
		if err != nil {
			log.Fatalf("load runtime failed: %s", err)
		}

		if err := rt.Register(*pkgfile); err != nil {
			log.Fatalf("register pkgfile failed: %s", err)
		}

		if err := rt.SaveTo(*rtfile); err != nil {
			log.Fatalf("save to runtime:%s failed: %s", *rtfile, err)
		}
	}
}

// func loadRuntime(rtfile string) error {
// 	flow.Register(componentName, constructor)
// 	// flow.RegisterJSON(componentName, filePath)
// }
