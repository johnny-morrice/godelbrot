// +build !production

package main

import (
    "fmt"
    do "gopkg.in/godo.v2"
)

const packageRoot = "github.com/johnny-morrice/godelbrot"
const internalRoot = packageRoot + "/internal"

type pkg struct {
    name string
    root string
}

func publicPkg(name string) pkg {
    p := pkg{}
    p.name = name
    p.root = packageRoot
    return p
}

func internalPkg(name string) pkg {
    p := pkg{}
    p.name = name
    p.root = internalRoot
    return p
}

var lib = []pkg{
    internalPkg("base"),
    internalPkg("draw"),
    internalPkg("sequence"),
    internalPkg("region"),
}

var nativeArithmetic = []pkg{
    internalPkg("nativebase"),
    internalPkg("nativesequence"),
    internalPkg("nativeregion"),
}

var bigArithmetic = []pkg{
   internalPkg("bigbase"),
   internalPkg("bigsequence"),
   internalPkg("bigregion"),
}

var appBase = []pkg {
    publicPkg("libgodelbrot"),
}

var apps = []pkg{
    publicPkg("configbrot"),
    publicPkg("renderbrot"),
    publicPkg("colorbrot"),
    publicPkg(""), // Top level binary
}

var all []pkg

// Group all packages in one slice
func init() {
    subsystems := [][]pkg{
        lib,
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
    units := map[string][]pkg{
        "lib": lib,
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

func buildFeatures(p *do.Project, subsystem string, components []pkg) {
    var componentsInstall do.S
    var componentsTest do.S
    for _, module := range components {
        install := installTaskName(module.name)
        test := testTaskName(module.name)
        componentsInstall = append(componentsInstall, install)
        componentsTest = append(componentsTest, test)

        p.Task(install, nil, func(unit pkg) func(c *do.Context) {
            return func (c *do.Context) {
                goInstall(c, unit)
            }
        }(module))
        p.Task(test, nil, func(unit pkg) func(c *do.Context) {
             return func (c *do.Context) {
                goTest(c, unit)
            }
        }(module))
    }

    p.Task(installTaskName(subsystem), componentsInstall, nil)
    p.Task(testTaskName(subsystem), componentsTest, nil)
}

func goInstall(c *do.Context, unit pkg) {
    command := fmt.Sprintf("go install %v/%v", unit.root, unit.name)
    c.Run(command)
}

func goTest(c *do.Context, unit pkg) {
    command := fmt.Sprintf("go test %v/%v", unit.root, unit.name)
    c.Run(command)
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