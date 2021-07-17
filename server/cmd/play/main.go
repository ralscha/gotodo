package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type BookForm struct {
	Title         string `name:"title" validate:"max=255,email,hexcolor"`
	Author        string `json:"author" validate:"required,max=255"`
	PublishedDate string `json:"published_date" validate:"required"`
	ImageUrl      string `json:"image_url" validate:"url"`
	Description   string `json:"description"`
}

func main() {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("name")
	})

	form := BookForm{
		Title:         "",
		Author:        "",
		PublishedDate: "",
		ImageUrl:      "",
		Description:   "",
	}

	if err := validate.Struct(form); err != nil {
		fmt.Println(err)
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Field())
			fmt.Println(err.Tag())
			fmt.Println()
		}
	}
}
