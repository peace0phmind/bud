# bud: Simplify and Graceful Golang Code

Simplify and Graceful Golang Code. Features include: auto-wiring, stream api, ... more TODO

Feature List

- [x] wire: Runtime Dynamic Auto-Wiring Using Dependency Injection
- [x] stream: stream api like java to deal with slice
- [ ] enum: TODO
- [ ] trace-log: TODO
- [ ] config: TODO
- [ ] plugin: TODO
- [ ] util: TODO
- [ ] remote: TODO:
- [ ] sys info: TODO: github.com/shirou/gopsutil

## wire

Runtime Dynamic Auto-Wiring Using Dependency Injection

Wire is a tool for automating the wiring of components at runtime using dependency injection.
The dependencies between components can be defined by the tag marked on the field of the struct, or by the parameters of
the Init method.

### easy start

repository struct looks like

```go
package repo

import "github.com/peace0phmind/bud/factory"

func init() {
	factory.Singleton[MyRepo]() // registry MyRepo as a singleton service
}

// MyRepo to control model with db
type MyRepo struct {
	//...
}

func (mr *MyRepo) SaveData() {
	//...
}
```

service struct looks like

```go

package service

import "github.com/peace0phmind/bud/factory"
import "repo"

func init() {
	factory.Singleton[MyService]() // registry MyRepo as a singleton service
}

type MyService struct {
	myRepo *repo.MyRepo `wire:"auto"`
}

func (ms *MyService) DoSomething() {
	// ...
	ms.myRepo.SaveDate()
	// ...
}
```

main.go

```go
package main

import "github.com/peace0phmind/bud/factory"
import "service"

func main() {
	serv := factory.Get[servcie.MyService]()

	serv.DoSomething()
}

```