package GoDaemonSkeleton

import (
	"fmt"
	"os"
	"strings"
)

var (
	apps = make([]*App, 0)
)

type App struct {
	Name     string
	Handover *func()
}

func GetApp() (*App, []string) {
	appNames := func() []string {
		out := []string{}
		for _, app := range apps {
			out = append(out, app.Name)
		}

		return out
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "No App specified. Must be one of: %s\n", strings.Join(appNames(), " "))
		os.Exit(1)
	}

	var subAppArgs []string
	var subApp *App
OuterLoop:
	for i, arg := range os.Args[1:] {
		for _, subAppTemp := range apps {
			if arg == subAppTemp.Name {
				subApp = subAppTemp
				subAppArgs = os.Args[i+2:]
				os.Args = os.Args[:i+1]
				break OuterLoop
			}
		}
	}

	if subApp == nil {
		fmt.Fprintf(os.Stderr, "No App specified. Must be one of: %s\n", strings.Join(appNames(), " "))
		os.Exit(1)
	}

	return subApp, subAppArgs
}

func AppRegister(app *App) {
	if app == nil {
		panic("nil subapp supplied")
	}
	for _, dup := range apps {
		if dup.Name == app.Name {
			panic("Register called twice for subApp " + app.Name)
		}
	}
	apps = append(apps, app)
}
