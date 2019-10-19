
package wsgi

import (
  // "fmt"
  "time"
  "os"
  "os/exec"
  "net"
  "io"
  "log"
  "runtime"
  "crypto/rand"
  utils "github.com/sagara_battousai/net/netutils"
  proto "github.com/golang/protobuf/proto"
)

const (
    WSGI string = "paperpython"
  //"/src/github.com/sagara_battousai/net/wsgi/paperPython/paperpython/__init__.py"
  ACK byte = 0x06
  ID_MASK uint64 = 0xFFFFFFFFFFFFFF00
  ID_BITS uint = 9
  BITS_IN_A_BYTE uint = 8
)

type WsgiServer struct {
  *exec.Cmd
  in io.WriteCloser
  Out io.ReadCloser
  Err io.ReadCloser
  *net.TCPConn
  errorDaemon *time.Timer
  // daemonDuration time.Duration
}

func New() *WsgiServer {
  wsgi_app, is_set := os.LookupEnv("WSGI_SERVER")
  if !is_set {
    // wsgi_app = os.ExpandEnv("$GOPATH" + WSGI)
    wsgi_app = WSGI
  }

  python, is_set := os.LookupEnv("WSGI_PYTHON")
  if !is_set {
    python = "pythonw"
  }

  if wsgi_app[len(wsgi_app) - 3 : len(wsgi_app)] == ".py" {
    return NewWithArgs(python, wsgi_app)
  } else {
    return NewWithArgs(python, "-m", wsgi_app)
  }
}

func NewWithArgs(python string, arg ...string) *WsgiServer {
  cmd := exec.Command(python, arg...)
  // Dont need error checking because its all over if they dont work?
  // Prefer stdPipes over cmdline args so on wsgi can serve multiple servers
  in, _ := cmd.StdinPipe()
  out, _ := cmd.StdoutPipe()
  err, _ := cmd.StderrPipe()

  return &WsgiServer{Cmd: cmd,in: in,Out: out,Err: err}
}

func (w *WsgiServer) Handshake(l *net.TCPListener) {
  addr:= l.Addr().(*net.TCPAddr)

  ipAddr := addr.IP
  port := uint32(addr.Port)
  randomKey := getRandomKey()
  numWorkers := uint32(runtime.GOMAXPROCS(0))

  config := &Config{Ip: ipAddr.String(),
                    Port: port,
                    IsIPv6: utils.IsIPv6(ipAddr),
                    IdChecksum: &IdChecksum{IdChecksum: randomKey},
                    NumWorkers: numWorkers}

  out, err := proto.Marshal(config)
  if err != nil {
    log.Fatalln("Failed to encode ip and port:", err)
  }

  n, err := w.in.Write(out)
  if err != nil {
    log.Fatalln("Error writing to wsgi on stdin, has the server been started?", n)
  }

  w.in.Close() //Sends EOF to python so reading stops

  c, err := l.Accept()
  if err != nil {
    log.Fatalln("Couldn't connect to Wsgi Server")
  }

  conn, okay := c.(*net.TCPConn)

  if !okay {
    log.Fatalln("Connection is wrong type")
  }

  idChecksum := &IdChecksum{}
  idChecksum_bytes := make([]byte, ID_BITS)
  conn.Read(idChecksum_bytes)
  if err := proto.Unmarshal(idChecksum_bytes, idChecksum); err != nil {
    log.Fatalln("Failed to parse idChecksum:", err)
  }

  //Maybe change to while in order to get the correct connection 
  //(to stop adversary)
  if idChecksum.IdChecksum != randomKey {
    log.Fatalln("Wrong Response")
  }

  conn.Write([]byte{ACK})

  w.TCPConn = conn
}


func (w *WsgiServer) StartErrorDaemon(d time.Duration) {
  // w.daemonDuration = d
  errorBytes := make([]byte, 2048)

  var errorCheck func()

  errorCheck = func() {
    n, err := w.Err.Read(errorBytes)
    if err != nil {
      return //maybe
    }
    if n == 0 {
      w.errorDaemon = time.AfterFunc(d, errorCheck)
    } else {
      log.Fatalln("Python Wsgi Server Has Thrown Error:", string(errorBytes))
    }
  }

  w.errorDaemon = time.AfterFunc(0, errorCheck)
}

func (w *WsgiServer) StopErrorDaemon() {
  w.errorDaemon.Stop()
  w.errorDaemon = nil
}






















func getRandomKey() uint64 {
  cryptoBytes := make([]byte, ID_BITS)
  _, err := rand.Read(cryptoBytes)
  if err != nil {
    log.Fatalln("Crypto Gen Error:", err)
  }
  var randomKey uint64 = 0x00
  for i, b := range cryptoBytes {
    randomKey |= uint64(b << (uint(i) * BITS_IN_A_BYTE))
  }

  return randomKey

}






