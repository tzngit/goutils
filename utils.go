package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tzngit/go-sh"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func CurAbsDir() string {
	file, _ := exec.LookPath(os.Args[0])
	path := filepath.Dir(file)
	return path
}

func LocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	var ip string
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return ip
}

func File2String(file string) string {
	fi, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		panic(err)
	}
	return string(fd)
}

func ExecCmd(success string, cmd string, args ...string) (error, string) {
	session := sh.NewSession()
	session.ShowCMD = true
	stdout, stderr, err := session.Command(cmd, args).Output()

	output := fmt.Sprintf("stdout:%s\nstderr:%s\n", string(stdout), string(stderr))

	cmdstr := strings.Join(args, " ")
	cmdstr = cmd + " " + cmdstr
	if err != nil {
		log.Printf("exec cmd[%s] fail! error:\n%s", cmdstr, err.Error())
		return err, output
	}
	if success != "" && !strings.Contains(output, success) {
		return errors.New("no success flag found"), output
	}
	return err, output
}

func ExecCmdInDir(successFlag []string, dir string, cmd string, args ...string) (error, string, string) {
	session := sh.NewSession()
	session.ShowCMD = true
	session.SetDir(dir)
	stdout_b, stderr_b, err := session.Command(cmd, args).Output()
    stdout := string(stdout_b)
    stderr := string(stderr_b)

	cmdstr := strings.Join(args, " ")
	cmdstr = cmd + " " + cmdstr
	if err != nil {
		log.Printf("exec cmd[%s] fail! error:\n%s", cmdstr, err.Error())
		return err, out + "\n" + stderr
	}

	if len(successFlag) > 0 {
		for _, str := range successFlag {
			if !strings.Contains(stdout, str) {
				//log.Printf("exec cmd[%s] no success flag!\n%s", cmdstr, stdout)
				return errors.New("no success flag found"), stdout
			}
		}
	}
	return err, stdout, stderr
}

func SavePid(pidFile string) {
	pid := os.Getpid()
	f, err := os.Create(pidFile)
	if err == nil {
		f.WriteString(strconv.Itoa(pid))
	}
}

func ResponseJson(w http.ResponseWriter, v interface{}) {
	result, err := json.Marshal(v)
	if err != nil {
		log.Print(err)
	} else {
		w.Write(result)
	}
}

func CorsHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
			w.Header().Set("Content-Type", "application/json")
		} else {
			fn(w, r)
		}
	}
}

func ParseJsonRequest(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&v)
	if err != nil {
		return err
	}
	return nil
}
