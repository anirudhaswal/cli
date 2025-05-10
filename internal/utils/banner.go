/*
Copyright © 2025 SuprSend
*/
package utils

import (
	log "github.com/sirupsen/logrus"
)

var version = "dev"

func Banner() {
	banner := `
   _____                  _____                __
  / ___/__  ______  _____/ ___/___  ____  ____/ /
  \__ \/ / / / __ \/ ___/\__ \/ _ \/ __ \/ __  / 
 ___/ / /_/ / /_/ / /   ___/ /  __/ / / / /_/ /  
/____/\__,_/ .___/_/   /____/\___/_/ /_/\__,_/   
          /_/                                    
`
	log.Info(banner)
	// Print the current version of the CLI
	log.Infof("SuprSend CLI - v.%s", version)
}
