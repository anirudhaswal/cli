package utils

import (
	log "github.com/sirupsen/logrus"
	"github.com/suprsend/cli/mgmnt"
	suprsend "github.com/suprsend/suprsend-go"
)

// SDKInstance is a singleton instance of the SuprSend SDK
var SDKInstance *mgmnt.SS_MgmntClient

// InitSDK initializes the SuprSend SDK
func InitSDK(serviceToken string, debug bool) {
	InitSDKWithUrls(serviceToken, "", "", debug)
}

// InitSDKWithUrls initializes the SuprSend SDK with custom URLs
func InitSDKWithUrls(serviceToken string, baseUrl string, mgmntUrl string, debug bool) {
	if SDKInstance != nil {
		log.Error("SDK already initialized")
		return
	}

	SDKInstance = mgmnt.NewClientWithUrls(serviceToken, baseUrl, mgmntUrl, debug)
}

// GetSuprSendMgmntClient returns the singleton instance of the SuprSend Mgmnt SDK
func GetSuprSendMgmntClient() *mgmnt.SS_MgmntClient {
	return SDKInstance
}

func GetSuprSendWorkspaceClient(workspace string) (*suprsend.Client, error) {
	return SDKInstance.GetWorkspaceClient(workspace)
}
