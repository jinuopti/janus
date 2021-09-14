package main

import (
	"flag"
	"fmt"
	conf "github.com/jinuopti/janus/configure"
	. "github.com/jinuopti/janus/log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	app "github.com/jinuopti/janus/application" // insert your application name
)

const JanusVersion = "0.2.0"

// global variables
var janus *Janus

// Janus application main structure
type Janus struct {
	args    *Arguments
	config 	*conf.Values
	chans 	*Channels
}

// Arguments process arguments
type Arguments struct {
	printConfig *string
	singleShot  *bool
	iniFile     *string
	test 	    *string

	/* insert your application arguments */
}

// Channels go channels
type Channels struct {
	doneChan chan bool			// true: exit application
	sigChan  chan os.Signal		// signal channel (SIGINT...)
}

func NewChannels() *Channels {
	return &Channels{}
}

func NewJanus() *Janus {
	if janus == nil {
		janus = &Janus{}
		janus.config = conf.NewValues()
		janus.chans = NewChannels()
	}
	return janus
}

func NewArguments() *Arguments {
	return &Arguments{}
}

func parseArguments(args *Arguments) bool {
	/* common arguments */
	args.printConfig = flag.String("pc", "", "Print config values [all|core|timer|log|net|...]")
	args.singleShot  = flag.Bool("ss", false, "Run once and exit the application")
	args.iniFile     = flag.String("ini", "janus.ini", "Set configuration file")
	args.test        = flag.String("test", "", "Input a string argument for the test function")

	/* insert your application arguments */

	flag.Parse()

	return true
}

func runSingleShot() {
	Logi("SingleShot Function Started")

	if len(*janus.args.test) > 0 {
		runTestCode(*janus.args.test)
	}

	// insert your singleshot application logic

	Logi("SingleShot Function Finished")
}

func runInfinite() {

	if janus.config.Kredigo.Enabled {
		go app.Run(janus.config)
	}

	// start your application

	ticker := time.NewTicker(time.Second * time.Duration(janus.config.Time.IdleTimeout))
	for {
		select {
		case sig := <-janus.chans.sigChan:
			Logd("\nSignal received: %s", sig)
			janus.chans.doneChan <- true
		case <-ticker.C:
			Logd("IDLE Timeout %d sec", janus.config.Time.IdleTimeout)
		default:
			// Logd("Log rotate test...")
			time.Sleep(20 * time.Millisecond)
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // max number of cpu on this PC

	janus = NewJanus()

	// parse application arguments
	janus.args = NewArguments()
	if parseArguments(janus.args) == false {
		return
	}

	// Exit the application after executing this function
	defer func() {
		if *janus.args.singleShot == false {
			fmt.Println("### Janus Application Finished ###")
		}
		Logi("### Janus Application Finished ###")
		Close()
	}()

	// load configuration
	_, err := janus.config.GetValueALL(*janus.args.iniFile)
	if err != nil {
		fmt.Printf("Error, open ini file, [%v]\n", err)
		return
	}

	// init log system
	InitLogger(janus.config.Log)

	// print application start message
	Logi("### Janus Application Started, Version[%s], CPU Core Num: %d ###", JanusVersion, runtime.GOMAXPROCS(0))

	// print config values
	if len(*janus.args.printConfig) > 0 {
		values := janus.config.PrintValues(*janus.args.printConfig)
		fmt.Printf(values)
		return
	}

	// single-shot logic
	if *janus.args.singleShot == true {
		runSingleShot()
		return
	}

	// make default channels
	janus.chans.sigChan = make(chan os.Signal, 1)
	janus.chans.doneChan = make(chan bool, 1)
	signal.Notify(janus.chans.sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Main Infinite Loop
	go runInfinite()
	<-janus.chans.doneChan
}

func runTestCode(str string) {
	Logd("Start TEST Logic, String: [%s]", str)
}