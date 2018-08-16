package config

// Configuration to be filled by Gonfig
type Configuration struct {
	SQLitePath      string
	Port            int
	ShortCodeLength int
	Domain          string
	Public          bool
}
