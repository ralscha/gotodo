package output

type HealthcheckOutput struct {
	Status string `json:"status"`
}

type AppVersionOutput struct {
	BuildTime string `json:"buildTime"`
	Version   string `json:"version"`
}
