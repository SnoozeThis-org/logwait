package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

var (
	FlagSet = pflag.CommandLine

	configFile = FlagSet.StringP("config", "c", "", "Optional JSON config file to read flags from")
)

func Parse() {
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		FlagSet.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nSettings can also be supplied as ENVIRONMENT_VARIABLES or through a JSON config file with --config")
	}

	FlagSet.VisitAll(func(f *pflag.Flag) {
		envKey := strings.ReplaceAll(strings.ToUpper(f.Name), "-", "_")
		if v := os.Getenv(envKey); v != "" {
			f.Value.Set(v)
		}
	})
	FlagSet.Parse(os.Args[1:])

	if *configFile != "" {
		if err := applyConfigFile(); err != nil {
			log.Fatalf("Error while parsing config file: %v", err)
		}
	}
}

func applyConfigFile() error {
	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return err
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	for k, v := range m {
		f := FlagSet.Lookup(k)
		if f == nil {
			return fmt.Errorf("setting %q doesn't exist", k)
		}
		if f.Changed {
			continue
		}
		switch v := v.(type) {
		case string:
			f.Value.Set(v)
		case bool, float64:
			f.Value.Set(fmt.Sprint(v))
		default:
			return fmt.Errorf("unexpected type %T for config setting %q", v, k)
		}
	}
	return nil
}
