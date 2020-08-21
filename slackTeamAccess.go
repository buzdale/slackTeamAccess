package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

func main() {

	userflag := flag.String("u", "all", "Which user to view (Default:All)")
	pagesflag := flag.Int("p", 1, "How mmany pages of data, each page is 1000 entries (Max 100)")
	flag.Parse()
	token := viperEnvVariable("SLACK_TOKENACCESS")
	getAccessLogs(*userflag, *pagesflag, token)

}
func getAccessLogs(username string, pages int, token string) {
	for i := 1; i <= pages; i++ {
		prams := slack.AccessLogParameters{
			Count: 1000, // 1000 per page limit
			Page:  i,    // page limit is 100
		}
		api := slack.New(token)
		// If you set debugging, it will log all requests to the console
		// Useful when encountering issues
		// slack.New(token, slack.OptionDebug(true))
		logins, _, err := api.GetAccessLogs(prams)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		for _, user := range logins {
			if username == "all" || username == user.Username {
				fmt.Println(user.Username, " First Login: ", epochToLocal(user.DateFirst), "Last Login: ", epochToLocal(user.DateLast))
			}
		}
	}
}

func epochToLocal(epochtime int) time.Time {
	localtime := time.Unix(int64(epochtime), 0)
	return localtime
}
func viperEnvVariable(key string) string {

	// SetConfigFile explicitly defines the path, name and extension of the config file.
	// Viper will use this and not check any of the config paths.
	// .env - It will search for the .env file in the current directory
	viper.SetConfigFile(".env")

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	// viper.Get() returns an empty interface{}
	// to get the underlying type of the key,
	// we have to do the type assertion, we know the underlying value is string
	// if we type assert to other type it will throw an error
	value, ok := viper.Get(key).(string)

	// If the type is a string then ok will be true
	// ok will make sure the program not break
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}
