package entities

// Config is struct
type Config struct {
	LogLevel   string
	Hostfile   string
	TargetPort string
	BaseURI    string
	Hosts      []string
	FileInput  string
	FileOutput string
}
