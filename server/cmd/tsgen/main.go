package main

import (
	"os"
	"strings"

	"github.com/aarondl/null/v8"
	"github.com/gobuffalo/validate"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
	"gotodo.rasc.ch/cmd/web/input"
	"gotodo.rasc.ch/cmd/web/output"
	"gotodo.rasc.ch/internal/models"
)

const outputFile = "../client/src/app/api/types.ts"

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

	err := converter.ConvertToFile(outputFile)
	if err != nil {
		panic(err)
	}

	contents, err := os.ReadFile(outputFile)
	if err != nil {
		panic(err)
	}

	updated := strings.ReplaceAll(string(contents), "{[key: string]: string[]}", "Record<string, string[]>")
	err = os.WriteFile(outputFile, []byte(updated), 0o644)
	if err != nil {
		panic(err)
	}
}
