package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Config file not found: %s", err))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(viper.GetString("group.key")))

		fmt.Println(viper.GetStringMapString("group"))
	})

	port := viper.GetString("port")

	http.ListenAndServe(port, nil)

}
