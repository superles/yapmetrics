package logger

import "log"

// ExampleInitialize Пример инициализации конфига
func ExampleInitialize() {
	err := Initialize(DebugLevel)
	if err != nil {
		log.Fatal(err)
	}
	Log.Info("log info")

	// Output:
	//
}
