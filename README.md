# name

A go lang library that creates random Dungeon and Dragons names based off a spreadsheet

### Prerequisites

1. [Download & install Go Lang](https://golang.org/doc/install)

### Installing

Get the name library from github using go's get command:
```
go get github.com/MarkBerube/name
```

## Testing

Can run the testing using the go's testing command:

```
go test
```

## How to use

1. import the name library using the `import` syntax:

```
import "github.com/MarkBerube/name"
```

2. Create a configuration file using the NewConfig function. Configuration arguments are as follows:

```
// SheetConfig is a struct that holds configuration values on how the name spreadsheet is requested.
// Configuration only supports Google's V4 sheet API currently, please see README for appropriate default values.
type SheetConfig struct {
	URL               string // URL to the spreadsheet API
	ID                string // ID of the spreadsheet
	APIKey            string // API key to call the spreadsheet API
	Name              string // Name of the sheet in the spreadsheet that has the name list
	MsgLimit          int    // Maximum length of strings returned by GetRandomNameList
	MsgHeader         string // Header attached to the top of the first message string for GetRandomNameList
	Titles            bool   // Allow titles to be used by GetRandomNameList
	SecondNameAppends bool   // Allow second name appends to be used by GetRandomNameList
}
```

3. Run the GetRandomList() function with the amount of names to generate, your configuration and the GoogleSheet API struct provided in the library. It will return an array of strings with the names. Here's an example:

```
	conf = name.NewConfig("https://test.url",
		"123foo",
		"456bar",
		"NameList",
		1,
		"testing:\n")

	result := GetRandomNameList(1, conf, name.GoogleSheetAPIClient{})
```

## Authors

* **Mark Berube**

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
