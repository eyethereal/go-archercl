/*
	Logging

	For logging we are using the "github.com/op/go-logging" logging library which provides
	some modularity around the built-in framework and adds more flexibility including
	multi-destination and colored outputs.

	To set it up in your package, do something like the following at a single place inside
	your package:

		//////////////////////////////////////////////////////////////////////
		// Package level logging
		// This only needs to be in one file per-package

		var log = config.Logger("mymodule")

		//////////////////////////////////////////////////////////////////////

	To make sure that colored output is sent to the console there is a convenience
	method that can be called from init() if you wish. This should only be one in
	test code though because this should generally be controlled via the config file.


		func init() {
			config.ColoredLoggingToConsole()
		}


	Using the setup above, you will have a log variable available that has Printf style methods on
	it as per:

		log.Criticalf(fmt, ...)
		log.Errorf(fmt, ...)
		log.Warningf(fmt, ...)
		log.Noticef(fmt, ...)
		log.Infof(fmt, ...)
		log.Debugf(fmt, ...)

	NOTE: There was a breaking change in early 2016 which requires the trailing f on the level names
	if you want to use formatting strings. This matches the fmt.Printf(...) syntax. There are also
	non-f versions of the level methods which will immediately log the first value as a message. This
	requires going through our whole codebase and some things might have gotten missed. Yay for progress.

	Logging Configuration


	From the archer.acl file you can configure some global logging defaults, a list of
	backends, and individual logging level on a per module basis. Automatically at the
	end of config.LoadACLConfig the logging configuration for the whole app will be set. It
	can also be set later, but it isn't necessary.

	Here is an example configuration which logs to two destinations, both a file and stdout,
	using a different log format for stdout which includes color

		logging {
			backends {
				log_file {
					type: file
					filename="out.log"
				}

				console {
					type:stdout
					format="%{color}%{time:15:04:05.000} %{level:4.4s} %{module:8.8s} %{color:bold}%{shortfunc:10.10s}%{color:reset} ▶ %{message}"
				}
			}
		}

	Here is a full example including all backend types and all options for each of them.
	You wouldn't ever do this in practice of course. Notice that all the logging
	configuration happens inside the "logging" object


		logging {

			// If debug is set to true, then after all configuration has been done
			// a series of messages, one at each log level, will be written so
			// you can check to make sure logging is working the way you want it to.
			debug: true

			// The global logging level. Only message at this level and above will
			// be output. This can also be controlled on a per-module basis.
			level: info

			// The global format that is used by all backends unless another format
			// is specified in the definition of the individual backend. The value
			// shown here is the default value if one is not given. See go-logging
			// for more info about available commands.
			format: "%{time:15:04:05.000} %{shortfunc:10.10s} %{level:4.4s} %{module:8.8s} ▶ %{message}"

			// The modules object is used to define per-module options, which currently
			// is just the log level for that module.
			modules {
				"arpc": {
					level: warning
				}

				// The shorter ACL syntax can be used like this
				web level : debug
			}


			// Any number of backends can be specified, they just need unique
			// names. This is done instead of using an array so they can be more
			// easily addressed in a cascade of config files if necessary
			backends {

				// Just hold things in memory
				memory {
					type: memory
					size: 1000	// Number of lines to hold onto
					forTesting: false // See the go-logging docs
					format: "%{time:15:04:05} %{message}"
				}

				// Same as memory, but using a channel to store the data
				channel_memory {
					type: channelMemory
					size: 1000
					format: "%{time:15:04:05} %{message}"
				}

				// Write to a file, optionally adding color in addition to
				// whatever is specified in log format
				log_file {
					type: file
					filename: "out.log"
					color: false
					format: "%{time:15:04:05} %{message}"
				}

				// A special case of the file backend directed to stdout
				console {
					type: stdout
					color: true
					format: "%{time:15:04:05} %{message}"
				}

				// Like stdout, but for stderr
				standard_err {
					type: stderr
					color: true
					format: "%{time:15:04:05} %{message}"
				}

				// Logs to syslog
				standard_err {
					type: syslog
					prefix: "archer" // A prefix for all message. See builtin log package
					facility: user	 // Syslog facility.
									 // See the SyslogFacilities map in this file
					format: "%{time:15:04:05} %{message}"
				}
			}
		}

*/
package config

import (
	"fmt"
	"github.com/op/go-logging"
	"log"
	"log/syslog"
	"os"
)

const DEFAULT_FORMAT_STRING = "%{time:15:04:05.000} %{shortfunc:10.10s} %{level:4.4s} %{module:8.8s} ▶ %{message}"

//const DEFAULT_FORMAT_STRING = "%{color}%{time:15:04:05.000} %{shortfunc:10.10s} %{level:4.4s}%{color:reset} %{module:8.8s} ▶ %{message}"

//const DEFAULT_FORMAT_STRING = "%{time:15:04:05.000} %{shortfunc:.10s} %{level:.4s}%{color:reset} %{module:.8s} ▶ %{message}"

var loggers = make(map[string]*logging.Logger)

var loggingACL *AclNode

// Get the logger associated with a given module name.
func Logger(name string) (logger *logging.Logger) {

	// fmt.Printf("Logger(%s, )\n", name)

	logger = logging.MustGetLogger(name)

	if loggers[name] != nil && loggers[name] != logger {
		loggers[name].Notice("This logger has been replaced by a new one in the central configuration manager.")
	}

	loggers[name] = logger

	if loggingACL != nil {
		configureLogger(name, logger)
	}

	return logger
}

// Retrieves the underlying backend configured with the given name. If a custom
// format was specified, this backend is further wrapped, but this function will
// always return the basic type which can then be coercised into whatever you
// need if you want to further configure this. One use here is to retrieve the
// memory backend instance so you can get at the messages it contains.
func GetBackend(name string) logging.Backend {
	holder := backends[name]

	if holder == nil {
		return nil
	}

	return holder.Backend
}

type BackendHolder struct {
	// The name as given in the config file
	Name string

	Node *AclNode

	// Backend is the underlying backend, so Memory, File, etc.
	Backend logging.Backend

	// Formatted is either a reference to the Backend directly or it might
	// be a StringFormatter wrapper which applies custom formatting to
	// the message. It's held separately because this is what is actually
	// given to loggers, but if you want like history from a memory logger
	// you need the unwrapped version.
	Formatted logging.Backend
}

var backends = make(map[string]*BackendHolder)

var globalLevel logging.Level

func configureLogger(name string, logger *logging.Logger) {

	// fmt.Printf("Configuring logger '%s'\n", name)
	moduleCfg := loggingACL.Child("modules", name)

	// The level is set on a per-logger basis
	moduleLevel := globalLevel

	mls := moduleCfg.ChildAsString("level")
	if len(mls) > 0 {
		ml, err := logging.LogLevel(mls)
		if err == nil {
			moduleLevel = ml
		} else {
			logDelayed(logging.ERROR, "Did not understand log level for module "+name)
		}
	}

	logging.SetLevel(moduleLevel, name)
}

// SetColoredConsole is a convenience method for simple test apps that would like a reasonable
// colored log output without the need to setup other configuration stuff. Note that loading a
// configuration AFTER you have called this would cause this configuration to be overwritten by
// whatever logging is setup in that config. Thus, call this last if you don't want to specify
// logging in your config file otherwise - although you probably want to just put logging
// into your config file if you have one. This is more for when you don't have one.
func ColoredLoggingToConsole() {

	travisMarker := os.ExpandEnv("$TRAVIS")
	if len(travisMarker) != 0 {
		logDelayed(logging.INFO, "Not setting a colored config because $TRAVIS is "+travisMarker)
		return
	}

	cfg := StringToACL(`
		logging backends color_console {
			type: stdout
			color: true
		}
	`)

	SetLoggingConfig(cfg)

}

func SetLoggingConfig(acl *AclNode) {
	loggingACL = acl.Child("logging")

	var err error

	// A couple of easy global config values
	glString := loggingACL.ChildAsString("level")
	//fmt.Printf("glString=%v\n", glString)
	globalLevel, err = logging.LogLevel(glString)
	if err != nil {
		//fmt.Printf("err=%v\n", err)
		globalLevel = logging.INFO
	}
	//fmt.Printf("globalLevel = %v\n", globalLevel)
	logDelayed(logging.INFO, "Setting global logging level to "+globalLevel.String())

	fmtStr := loggingACL.ChildAsString("format")
	if len(fmtStr) > 0 {
		logging.SetFormatter(logging.MustStringFormatter(fmtStr))
	} else {
		logging.SetFormatter(logging.MustStringFormatter(DEFAULT_FORMAT_STRING))
	}

	// Remove any old backends. When we set a new backend to each logger it
	// replaces any previous ones that existed there
	backends = make(map[string]*BackendHolder)
	all := make([]logging.Backend, 0)
	beACL := loggingACL.Child("backends")
	if beACL == nil {
		be := logging.NewLogBackend(os.Stdout, "", 0)
		be.Color = true
		all = append(all, be)
	} else {

		for name, beNode := range beACL.Children {
			holder := &BackendHolder{
				Name: name,
				Node: beNode,
			}
			switch beNode.ChildAsString("type") {
			case "memory":
				makeMemoryBackend(holder)

			case "channelMemory":
				makeChannelMemoryBackend(holder)

			case "stdout":
				makeStdoutBackend(holder)

			case "stderr":
				makeStderrBackend(holder)

			case "file":
				makeFileBackend(holder)

			case "syslog":
				makeSyslogBackend(holder)

			case "loggly":
				makeLogglyBackend(holder)

			}

			if holder.Backend == nil {
				fmt.Printf("Ignoring backend named '%s'\n", holder.Name)
				continue
			}

			setupFormatter(holder)
			backends[holder.Name] = holder
			all = append(all, holder.Formatted)
		}
	}

	if len(all) > 0 {
		logging.SetBackend(all...)
	} else {
		logDelayed(logging.ERROR, "No backends were configured. Default logging configuration!")
	}

	// Now that all the backends are there, add them to the proper things
	for name, logger := range loggers {
		configureLogger(name, logger)
	}

	// And then some debugging of the config if necessary
	lcd := loggingACL.ChildAsBool("debug")
	if lcd {
		lgr := logging.MustGetLogger("log-debug")

		lgr.Critical("A critical message")
		lgr.Error("Merely an error")
		lgr.Warning("You have been warned")
		lgr.Notice("This notice has been delivered")
		lgr.Info("Informational. That is all")
		lgr.Debug("WTF hasn't this been debugged yet?")
	}
}

func setupFormatter(holder *BackendHolder) {

	fmtStr := holder.Node.ChildAsString("format")
	if len(fmtStr) == 0 {
		// fmt.Printf("Using default formatter for %v\n", holder)
		holder.Formatted = holder.Backend
		return
	}

	// fmt.Printf("Using custom formatter '%v' for %v\n", fmtStr, holder)
	formatter := logging.MustStringFormatter(fmtStr)
	holder.Formatted = logging.NewBackendFormatter(holder.Backend, formatter)
}

func makeMemoryBackend(holder *BackendHolder) {

	level, err := logging.LogLevel(holder.Node.ChildAsString("level"))
	if err != nil {
		level = globalLevel
	}

	size := holder.Node.ChildAsInt("size")
	if size == 0 {
		size = 3000
	}

	if holder.Node.ChildAsBool("forTesting") {
		holder.Backend = logging.InitForTesting(level)
	} else {
		holder.Backend = logging.NewMemoryBackend(size)
	}
}

func makeChannelMemoryBackend(holder *BackendHolder) {

	size := holder.Node.ChildAsInt("size")
	if size == 0 {
		size = 3000
	}

	holder.Backend = logging.NewChannelMemoryBackend(size)
}

func makeStdoutBackend(holder *BackendHolder) {

	be := logging.NewLogBackend(os.Stdout, "", 0)
	be.Color = holder.Node.ChildAsBool("color")

	holder.Backend = be
}

func makeStderrBackend(holder *BackendHolder) {

	be := logging.NewLogBackend(os.Stderr, "", 0)
	be.Color = holder.Node.ChildAsBool("color")

	holder.Backend = be
}

func makeFileBackend(holder *BackendHolder) {

	fName := holder.Node.ChildAsString("filename")
	if len(fName) == 0 {
		return
	}

	// TODO: More exciting things about filename such as sequence numbers, etc.

	file, err := os.OpenFile(fName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0666))
	if err != nil {
		log.Panicf("Unable to open file '%s' : %s", fName, err)
		return
	}

	be := logging.NewLogBackend(file, "", 0)
	be.Color = holder.Node.ChildAsBool("color")
	holder.Backend = be
}

var SyslogFacilities = map[string]syslog.Priority{
	"kern":     syslog.LOG_KERN,
	"user":     syslog.LOG_USER,
	"mail":     syslog.LOG_MAIL,
	"daemon":   syslog.LOG_DAEMON,
	"auth":     syslog.LOG_AUTH,
	"syslog":   syslog.LOG_SYSLOG,
	"lpr":      syslog.LOG_LPR,
	"news":     syslog.LOG_NEWS,
	"uucp":     syslog.LOG_UUCP,
	"cron":     syslog.LOG_CRON,
	"authpriv": syslog.LOG_AUTHPRIV,
	"ftp":      syslog.LOG_FTP,

	"local0": syslog.LOG_LOCAL0,
	"local1": syslog.LOG_LOCAL1,
	"local2": syslog.LOG_LOCAL2,
	"local3": syslog.LOG_LOCAL3,
	"local4": syslog.LOG_LOCAL4,
	"local5": syslog.LOG_LOCAL5,
	"local6": syslog.LOG_LOCAL6,
	"local7": syslog.LOG_LOCAL7,
}

func makeSyslogBackend(holder *BackendHolder) {
	var err error

	prefix := holder.Node.ChildAsString("prefix")
	fac := holder.Node.ChildAsString("facility")

	var facility syslog.Priority
	if len(fac) > 0 {
		facility = SyslogFacilities[fac]
	}

	if facility > 0 {
		holder.Backend, err = logging.NewSyslogBackendPriority(prefix, facility)
	} else {
		holder.Backend, err = logging.NewSyslogBackend(prefix)
	}

	if err != nil {
		log.Panicf("Unable to create syslog backend: %s", err)
	}
}

func makeLogglyBackend(holder *BackendHolder) {

	token := holder.Node.ChildAsString("token")
	if len(token) == 0 {
		log.Panicf("No token for loggly backend")
	}

	tags := holder.Node.ChildAsStringList("tags")

	client := NewLogglyClient(token, tags...)
	holder.Backend = client
}