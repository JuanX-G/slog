package configUtil

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort string 
	DBport string
}

func NewAppConfig() *AppConfig {
	return &AppConfig{ServerPort: "8109", DBport: "5432"}	
}

func(a *AppConfig) ParseConfig(fName string) error {
	fBytes, err := os.ReadFile(fName)
	if err != nil {
		return err
	}
	fLines := strings.Split(string(fBytes), "\n") 
	for _, v := range fLines {
		lineSplit := strings.Split(v, "=")
		if len(lineSplit) != 2 {
			continue
		}
		key := strings.TrimSpace(lineSplit[0])
		val := strings.TrimSpace(lineSplit[1])
		switch key {
		case "server-port":
			serverPortInt, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			if serverPortInt < 1 || serverPortInt > 65535 {
				return fmt.Errorf("invalid port number %d", serverPortInt)
			}
			a.ServerPort = val
		case "db-port": 
			dbPortInt, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			if dbPortInt < 1 || dbPortInt > 65535 {
				return  fmt.Errorf("invalid port number %d", dbPortInt)
			}
			a.DBport = val
		}
	}
	return nil
}

