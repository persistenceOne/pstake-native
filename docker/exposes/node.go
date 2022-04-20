package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const (
	homeDirKey = "HOME_DIR"
	genesisFileKey = "GENESIS_FILE"
	validatorFileKey = "VALIDATOR_PRIV_FILE"
	portKey = "PORT"
)

func executeCommand(name string, arg ...string) string {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatalf("Error when cmd stdout pipe creation, %s\n", err)
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		log.Fatalf("Execute command: '%s' error, %s\n", cmd.String(), err)
	}

	cmdOut, _ := ioutil.ReadAll(stdout)

	return string(cmdOut)
}

func getNodeIDHandler(w http.ResponseWriter, r *http.Request) {
	homeDir := os.Getenv(homeDirKey)
	io.WriteString(w, executeCommand("pstaked", "tendermint", "show-node-id", "--home", homeDir))
}

func getPubKeyHandler(w http.ResponseWriter, r *http.Request) {
	validatorFile := os.Getenv(validatorFileKey)
	io.WriteString(w, executeCommand("bash", "-c", fmt.Sprintf("cat %s | jq '.pub_key'", validatorFile)))
}

func getGenesisHandler(w http.ResponseWriter, r *http.Request) {
	genesisFile := os.Getenv(genesisFileKey)
	jsonFile, err := os.Open(genesisFile)
	if err != nil {
		log.Fatalf("Error opening genesis file at %s", genesisFile)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(byteValue)
}

func main() {
	fmt.Println("Server started ...")
	http.HandleFunc("/node_id", getNodeIDHandler)
	http.HandleFunc("/pub_key", getPubKeyHandler)
	http.HandleFunc("/genesis", getGenesisHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv(portKey)), nil)
	if err != nil {
		log.Fatalf("Fail to start server, %s\n", err)
	}
}