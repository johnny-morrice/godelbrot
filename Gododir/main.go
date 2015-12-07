package main

import (
    do "gopkg.in/godo.v2"
)

const packageRoot = "functorama.com/demo"

var toolkit = []string{
    "base",
    "draw",
    "sequence",
    "region",
    "sharedregion",
}

var nativeArithmetic = []string{
    "nativebase",
    "nativesequence",
    "nativeregion",
    "nativesharedregion",
}

var bigArithmetic = []string{
   "bigbase",
   "bigsequence",
   "bigregion",
   "bigsharedregion",
}

var appBase = []string {
    "libgodelbrot",
}

var apps = []string{
    "configbrot",
    "renderbrot",
    "godelbrot",
    "webdelbrot",
    "autobrot",
}

var all []string

// Group all packages in one slice
func init() {
    subsystems := [][]string{
        toolkit,
        nativeArithmetic,
        bigArithmetic,
        appBase,
        apps,
    }
    for _, sub := range subsystems {
        all = append(all, sub...)
    }
}

func tasks(p *do.Project) {
    units := map[string][]string{
        "toolkit": toolkit,
        "native": nativeArithmetic,
        "big": bigArithmetic,
        "appBase": appBase,
        "apps": apps,
        "all": all,
    }
    // Install/Test for each subsystem
    for subsystem, components := range units {
        buildFeatures(p, subsystem, components)
    }
    // Default task is install all
    p.Task("default", do.S{"allInstall"}, nil)
}

func buildFeatures(p *do.Project, subsystem string, components []string) {
    var componentsInstall do.S
    var componentsTest do.S
    for _, module := range components {
        install := installTaskName(module)
        test := testTaskName(module)
        componentsInstall = append(componentsInstall, install)
        componentsTest = append(componentsTest, test)

        p.Task(install, nil, func(unit string) func(c *do.Context) {
            return func (c *do.Context) {
                goInstall(c, unit)
            }
        }(module))
        p.Task(test, nil, func(unit string) func(c *do.Context) {
             return func (c *do.Context) {
                goTest(c, unit)
            }
        }(module))
    }

    p.Task(installTaskName(subsystem), componentsInstall, nil)
    p.Task(testTaskName(subsystem), componentsTest, nil)
}

func goInstall(c *do.Context, unit string) {
    command := "go install " + fullUnitPath(unit)
    c.Run(command)
}

func goTest(c *do.Context, unit string) {
    command := "go test " + fullUnitPath(unit)
    c.Run(command)
}

func fullUnitPath(module string) string {
    return packageRoot + "/" + module
}

func installTaskName(module string) string {
    return module + "Install"
}

func testTaskName(module string) string {
    return module + "Test"
}

func main() {
    do.Godo(tasks)
}