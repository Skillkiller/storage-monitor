package ssh

import (
	"fmt"
	"log"
	"os"
	"ssh-storage-monitor/config"
	"strconv"

	"golang.org/x/crypto/ssh"
)

func GetSshConnection(hostConfig config.SshHostConfig) *ssh.Client {
	config := &ssh.ClientConfig{
		User: hostConfig.User,
		Auth: []ssh.AuthMethod{
			getAuthMethod(hostConfig),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", hostConfig.Host+":"+strconv.Itoa(hostConfig.Port), config)
	if err != nil {
		log.Fatalf("Verbindung fehlgeschlagen: %v", err)
	}
	return client
}

func getAuthMethod(hostConfig config.SshHostConfig) ssh.AuthMethod {
	if hostConfig.Password != "" {
		return ssh.Password(hostConfig.Password)
	}

	if hostConfig.PrivateKey != "" && hostConfig.PrivateKeyFile == "" {
		// Has private key
		auth, err := publicKeyAuthFromBytesWithPassphrase([]byte(hostConfig.PrivateKey), []byte(hostConfig.PrivateKeyPassphrase))
		if err != nil {
			log.Fatal(err)
		}
		return auth
	}

	if hostConfig.PrivateKey == "" && hostConfig.PrivateKeyFile != "" {
		// Has a private key file
		key, err := os.ReadFile(hostConfig.PrivateKeyFile)
		if err != nil {
			log.Fatal(err)
		}
		auth, err := publicKeyAuthFromBytesWithPassphrase(key, []byte(hostConfig.PrivateKeyPassphrase))
		if err != nil {
			log.Fatal(err)
		}
		return auth
	}

	log.Fatalf("No auth method set for host")
	return nil
}

func publicKeyAuthFromBytesWithPassphrase(key []byte, passphrase []byte) (ssh.AuthMethod, error) {
	var signer ssh.Signer
	var err error
	if passphrase == nil {
		signer, err = ssh.ParsePrivateKey(key)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, passphrase)
	}
	if err != nil {
		return nil, fmt.Errorf("konnte privaten Schlüssel mit Passphrase nicht parsen: %w", err)
	}
	return ssh.PublicKeys(signer), nil
}

func ExecuteSshCommand(client *ssh.Client, command string) ([]byte, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return session.CombinedOutput(command)
}
