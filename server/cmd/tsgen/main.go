package main

import (
	"github.com/gobuffalo/validate"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
	"github.com/volatiletech/null/v8"
	"gotodo.rasc.ch/cmd/web/input"
	"gotodo.rasc.ch/cmd/web/output"
	"gotodo.rasc.ch/internal/models"
)

func main() {
	converter := typescriptify.New()
	converter.Add(models.Todo{})
	converter.Add(input.LoginInput{})
	converter.Add(input.PasswordResetInput{})
	converter.Add(input.PasswordResetRequestInput{})
	converter.Add(input.EmailChangeInput{})
	converter.Add(input.PasswordInput{})
	converter.Add(input.TokenInput{})
	converter.Add(input.PasswordChangeInput{})
	converter.Add(input.SignUpInput{})
	converter.Add(output.LoginOutput{})
	converter.Add(validate.Errors{})
	converter.Add(output.AppVersionOutput{})
	converter.CreateInterface = true
	converter.BackupDir = ""
	converter.ManageType(null.String{}, typescriptify.TypeOptions{TSType: "string"})

	err := converter.ConvertToFile("../client/src/app/api/types.ts")
	if err != nil {
		panic(err)
	}

}
