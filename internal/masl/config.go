package masl

import (
	"log"
	"os"
	"os/user"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
)

// Accounts represents the accounts section of the masl config file
type Accounts []struct {
	ID                     string `toml:"ID"`
	Name                   string `toml:"Name"`
	EnvironmentIndependent bool   `toml:"EnvironmentIndependent"`
}

// Config represents the masl config file
type Config struct {
	BaseURL      string `toml:"BaseURL"`
	ClientID     string `toml:"ClientID"`
	ClientSecret string `toml:"ClientSecret"`
	AppID        string `toml:"AppID"`
	Subdomain    string `toml:"Subdomain"`
	Username     string `toml:"Username"`
	Duration     int    `toml:"Duration"`
	Profile      string `toml:"Profile"`
	DefaultRole  string `toml:"DefaultRole"`
	Debug        bool   `toml:"Debug"`
	Environments []struct {
		Name     string   `toml:"Name"`
		Accounts []string `toml:"Accounts"`
	} `toml:"Environments"`
	Accounts Accounts `toml:"Accounts"`
}

// GetConfig reads the masl.toml configuration file for initialization.
func GetConfig(logger *logrus.Logger) Config {

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// Read masl.toml config file for initialization
	conf := Config{Profile: "masl", Debug: false, Duration: 3600} // Set default values
	if _, err := toml.DecodeFile(usr.HomeDir+string(os.PathSeparator)+"masl.toml", &conf); err != nil {
		log.Fatal(err.Error())
	}

	logger.WithFields(logrus.Fields{
		"baseURL":       conf.BaseURL,
		"clientID":      conf.ClientID,
		"clientSecret":  conf.ClientSecret,
		"appID":         conf.AppID,
		"subdomain":     conf.Subdomain,
		"username":      conf.Username,
		"profile":       conf.Profile,
		"duration":      conf.Duration,
		"debug":         conf.Debug,
		"#environments": len(conf.Environments),
		"#accounts":     len(conf.Accounts),
	}).Info("Config settings")

	return conf
}

// SearchAccounts search an account name for a given acount id
func SearchAccounts(accountInfo Accounts, accountID string) (string, bool) {

	for _, account := range accountInfo {
		if account.ID == accountID {
			return account.Name, account.EnvironmentIndependent
		}
	}
	return "", false
}

// GetEnvironmentDetails search an environment's detail for a given environment name
func GetEnvironmentDetails(conf Config, environment *string) []string {
	for _, env := range conf.Environments {
		if env.Name == *environment {
			return env.Accounts
		}
	}
	return nil
}
