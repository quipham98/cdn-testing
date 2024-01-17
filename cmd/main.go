package main

import "github.com/quipham98/cdn-testing/internal/manager"

func main() {
	manager.InitCommands()
	manager.Execute()
}
