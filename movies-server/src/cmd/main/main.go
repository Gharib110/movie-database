package main

type Config struct {
	Port int
	HostName string
}

func main() {
	run()
	return
}

func run()  {
	config := &Config{
		Port:     8080,
		HostName: "localhost",
	}

	
}