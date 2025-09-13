package ssh

import (
	"log"
	"ssh-storage-monitor/config"
	"strconv"

	"golang.org/x/crypto/ssh"
)

func GetSshConnection(hostConfig config.SshHostConfig) *ssh.Client {
	config := &ssh.ClientConfig{
		User: hostConfig.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(hostConfig.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", hostConfig.Host+":"+strconv.Itoa(hostConfig.Port), config)
	if err != nil {
		log.Fatalf("Verbindung fehlgeschlagen: %v", err)
	}
	return client
}

func ExecuteSshCommand(client *ssh.Client, command string) ([]byte, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return session.CombinedOutput(command)
}
