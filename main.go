package main
import (
//	"fmt"
	"os/exec"
	"flag"
	"io/ioutil"
	"log"
	"bytes"
	"net"
)


var emulateCmk = flag.Bool("cmk", true, "emulate check_mk functionality")
var scriptDir = flag.String("sdir", "/tmp/test-dir/", "directory for the scripts to run")
var tcpPort = flag.String("tport", "6556", "standard port for tcp server")



func parsescripts()(ouput []string) {
	flag.Parse()
	var output []string
	files, _:= ioutil.ReadDir(*scriptDir)
		for _, ff := range files {
			scriptname := *scriptDir+ff.Name()
			cmd := exec.Command(scriptname)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				log.Print(err)
			}
			add := out.String()
			output = append(output, add)
		}
	return output
}

func handleRequest(conn net.Conn) {
	for _, f := range parsescripts() {
		conn.Write([]byte(f))
	}
		conn.Close()
}

func checkMK() {
	flag.Parse()
	l, err := net.Listen("tcp", "localhost"+":"+*tcpPort)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
		}
		go handleRequest(conn)
	}

	
}



func main() {
	flag.Parse()
	if *emulateCmk {
		 checkMK()
	}
}
