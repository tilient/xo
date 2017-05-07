package main

import (
	//	"fmt"
	//	"github.com/nats-io/gnatsd/logger"
	//	"github.com/nats-io/gnatsd/server"
	//	"github.com/nats-io/go-nats"
	"log"
	//	"net/url"
	"os"
	"os/exec"
	//	"strconv"
	//	"time"
	"github.com/koding/multiconfig"
)

type Xo struct {
	Cmd           []string `default:"deploy,run,kill"`
	Source        string   `default:"test"`
	Target        string   `default:"~/kashbah/test"`
	Servers       []string `default:"tilient.org,dev.tilient.org"`
	NatsPort      string   `default:"44222"`
	NatsRoutePort string   `default:"22444"`
	TmuxSession   string   `default:"xo"`
}

func main() {
	log.Print("=== xo ===========================")
	m := multiconfig.NewWithPath("xo.cfg")
	xo := Xo{}
	m.Load(&xo)
	xo.updateCmd()
	log.Print("Cmd           : ", xo.Cmd)
	log.Print("Source        : ", xo.Source)
	log.Print("Target        : ", xo.Target)
	log.Print("Servers       : ", xo.Servers)
	log.Print("NatsPort      : ", xo.NatsPort)
	log.Print("NatsRoutePort : ", xo.NatsRoutePort)
	log.Print("TmuxSession   : ", xo.TmuxSession)
	for _, cmd := range xo.Cmd {
		log.Print("----------------------------------")
		switch cmd {
		case "deploy":
			xo.deployToServers()
		case "run":
			xo.runOnServers("tmux new -d -s " + xo.TmuxSession +
				" '" + xo.Target + "'")
		case "kill":
			xo.runOnServers("tmux kill-session -t " + xo.TmuxSession)
		default:
			log.Print("don't know commdand: ", cmd)
		}
	}
	log.Print("=== done =========================")
}

//------------------------------------------------------------

func (xo *Xo) updateCmd() {
	for ix, arg := range os.Args[1:] {
		if arg[0] != '-' {
			xo.Cmd = os.Args[(1 + ix):]
			return
		}
	}
}

//------------------------------------------------------------

func (xo *Xo) deployToServers() {
	for _, server := range xo.Servers {
		xo.deployToServer(server)
	}
}

func (xo *Xo) deployToServer(server string) {
	log.Print("scp ", xo.Source, " ", server+":"+xo.Target)
	err := exec.Command(
		"scp", xo.Source, server+":"+xo.Target).Run()
	if err != nil {
		log.Print("Error: ", err)
	}
}

//------------------------------------------------------------

func (xo *Xo) runOnServers(cmd string) {
	for _, server := range xo.Servers {
		xo.runOnServer(server, cmd)
	}
}

func (xo *Xo) runOnServer(server, cmd string) {
	log.Print("ssh ", "-4 ", server, " ", cmd)
	err := exec.Command("ssh", "-4", server, cmd).Run()
	if err != nil {
		log.Print("Error: ", err)
	}
}

//------------------------------------------------------------
// func clientMain(ec *nats.EncodedConn, id int) {
// 	countCh := make(chan int)
// 	ec.BindRecvQueueChan("count", "pool", countCh)
// 	ec.BindSendChan("count", countCh)
// 	defer close(countCh)
//
// 	cmdCh := make(chan string)
// 	ec.BindRecvChan("cmd", cmdCh)
// 	defer close(cmdCh)
//
// 	for {
// 		select {
// 		case cmd := <-cmdCh:
// 			if cmd == "quit" {
// 				return
// 			}
// 		case count := <-countCh:
// 			fmt.Println("count:", count)
// 			if count > 0 {
// 				countCh <- (count - 1)
// 			}
// 		case <-time.After(5 * time.Second):
// 			fmt.Println(".")
// 		}
// 	}
// }

//------------------------------------------------------------

// func run(servers []string, sync bool) {
// 	log.Print("=== Starting ===")
// 	natsServer := runNATSServer(servers)
// 	if natsServer.ReadyForConnections(10 * time.Second) {
// 		log.Print("=== Server started OK  ===")
// 	} else {
// 		log.Print("=== Server did NOT start OK ===")
// 	}
// 	nc, err := nats.Connect("nats://localhost:44222")
// 	if err != nil {
// 		fmt.Println("1>>", err)
// 	}
// 	ec, err := nats.NewEncodedConn(nc, "json")
// 	if err != nil {
// 		fmt.Println("2>>", err)
// 	}
//
// 	if sync {
// 		startCh := make(chan string)
// 		ec.BindRecvChan("start", startCh)
// 		defer close(startCh)
//
// 		resp := "xxxx"
// 		for resp == "xxxx" {
// 			ec.Request("alive", "ok", &resp, 1*time.Second)
// 		}
//
// 		<-startCh
// 	}
//
// 	id := -1
// 	if len(os.Args) > 2 {
// 		id, _ = strconv.Atoi(os.Args[2])
// 	}
// 	clientMain(ec, id)
// 	nc.Flush()
// 	ec.Close()
// 	natsServer.Shutdown()
// 	log.Print("=== Done ===")
// }

// func deployAndRun(target string, servers []string) {
// 	deployClients(target, servers)
// 	runClients(target, servers)
// 	syncClients(target, servers)
// 	log.Print("=== Waiting 5 minutes ===")
// 	time.Sleep(5 * 60 * time.Second)
// 	killClients(target, servers)
// 	log.Print("=== Done ===")
// }

// func deployClients(target string, servers []string) {
// 	log.Print("=== Deploying Executable ===")
// 	exePath, _ := os.Executable()
// 	copyToServers(exePath, target, servers)
// }

// func runClients(target string, servers []string) {
// 	log.Print("=== Launching Executable ===")
// 	tmuxCmd := "tmux new -d -s ses '" + target + " syncrun %d'"
// 	runOnServers(tmuxCmd, servers)
// }

// func syncClients(target string, servers []string) {
// 	log.Print("=== Waiting for Executables to come up  ===")
// 	nc, err := nats.Connect(natsServers(servers),
// 		nats.MaxReconnects(5), nats.ReconnectWait(30*time.Second))
// 	if err != nil {
// 		fmt.Println("1>>", err)
// 	}
// 	ec, err := nats.NewEncodedConn(nc, "json")
// 	if err != nil {
// 		fmt.Println("2>>", err)
// 	}
// 	defer ec.Close()
//
// 	cnt := 0
// 	ec.Subscribe("alive", func(subj, reply string, str string) {
// 		ec.Publish(reply, "ok")
// 		cnt += 1
// 	})
// 	for cnt < len(servers) {
// 		time.Sleep(1 * time.Second)
// 	}
// 	ec.Publish("start", "start")
// 	log.Print("=== All Executables did come up  ===")
// }

// func killClients(target string, servers []string) {
// 	log.Print("=== Killing Executables  ===")
// 	tmuxCmd := "tmux kill-session -t ses "
// 	runOnServers(tmuxCmd, servers)
// }

//------------------------------------------------------------

// func runNATSServer(servers []string) *server.Server {
// 	var opts = server.Options{}
// 	opts.Host = "::"
// 	opts.Port = 44222
// 	opts.Cluster.Host = "::"
// 	opts.Cluster.Port = 22444
// 	opts.Routes = natsRouteServerList(servers)
// 	log := logger.NewStdLogger(false, true, false, true, false)
//
// 	s := server.New(&opts)
// 	s.SetLogger(log, false, false)
// 	go s.Start()
// 	return s
// }

// func natsRouteServerList(servers []string) []*url.URL {
// 	return server.RoutesFromStr(natsRouteServers(servers))
// }

// func natsRouteServers(servers []string) string {
// 	srvrs := ""
// 	for _, server := range servers {
// 		srvrs = srvrs + ",nats://" + server + ":22444"
// 	}
// 	return srvrs[1:]
// }
//
// func natsServers(servers []string) string {
// 	srvrs := ""
// 	for _, server := range servers {
// 		srvrs = srvrs + ",nats://" + server + ":44222"
// 	}
// 	return srvrs[1:]
// }

//------------------------------------------------------------
