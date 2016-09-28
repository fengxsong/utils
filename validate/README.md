## Example

    var Name = validate.Use(
        validate.Len{
            Min: 6,
            Max: 32,
        }.Validate,
        validate.ValidChars("abcdefghijklmnopq", true),
        validate.ValidChoice([]string{"Golang"}, true),
        )

    fmt.Println(Name("golang"))

    //OR

    var Name01 = validate.New(
        validate.ValidChoice([]string{"Google"}, false),
    )
    Name01.Validate("google")
    fmt.Println(Name01.HasError())
    fmt.Println(Name01.Errors())

    //OR
    // `validate:"Len:6,32;Choice:TestVal,TestVal2"`
    // `validate:"Match:IP"` or Email/URL/Domain/MAC
    type T struct {
        Name string `validate:"Len:6,10;Match:feng.*"`
        Country string `validate:"Choice:UK,US"`
    }
    for _, err := range Validate(T{"fengxsong", "China"}) {
		fmt.Println(err)
	}
