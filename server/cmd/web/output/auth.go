package output

import "gotodo.rasc.ch/internal/models"

type LoginOutput struct {
	Authority models.AppUserAuthority `json:"authority"`
}
