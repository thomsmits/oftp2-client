package cmd

type Options struct {
	Server   string
	Port     int
	OdetteId string
	Verbose  bool
}

var activeOptions = &Options{
	Server:   "localhost",
	Port:     3305,
	OdetteId: "LOCAL",
}
