package main

import (
	"fmt"
	"log"
	"ssh-storage-monitor/config"
	"ssh-storage-monitor/ssh"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", cfg.Hosts)

	for _, host := range cfg.Hosts {
		for _, wd := range host.WatchedDirs {
			fmt.Printf("%+v\n", wd)
		}
	}

	for _, host := range cfg.Hosts {
		fmt.Printf("%+v\n", host)
		connection := ssh.GetSshConnection(host)

		output, _ := ssh.ExecuteSshCommand(connection, "df -h")
		fmt.Println(string(output))
		connection.Close()
	}
}
