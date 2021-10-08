package merge

import (
	"os/exec"
	"time"
)

var GethRPCEngine, _ = NewRPCNode("http://127.0.0.1:8545", func() {
	cmd := exec.Command("/home/matematik/go/src/github.com/ethereum/go-ethereum/build/bin/geth", "--dev", "--catalyst", "--http", "--http.api=\"eth,engine\"", "--override.totalterminaldifficulty=0")
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)
	go func() {
		time.Sleep(15 * time.Second)
		cmd.Process.Kill()
	}()
})
