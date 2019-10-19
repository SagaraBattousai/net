
package wsgi

import (
  "strings"
  "fmt"
  "net/http"
  "log"
  // utils "github.com/sagara_battousai/net/netutils"
  proto "github.com/golang/protobuf/proto"
)



func prepareHeadders(h http.Header) map[string] string {
  header := make(map[string] string)
  for k, v := range h {
    header[k] = strings.Join(v, ",")
  }
  return header
}

func (wsgiServer *WsgiServer) Handler() func (w http.ResponseWriter, r *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    host := strings.SplitN(r.Host, ":", 2)
    httpRequestHeaders := prepareHeadders(r.Header)
    wsgi := &Wsgi{
                Version:       &Wsgi_Version{Major: 1, Minor: 0},
                UrlScheme:     r.URL.Scheme,
                InputStream:   "socket",
                ErrorStream:   "stderr",
                Multithreaded: false,
                Multiprocess:  true,
                RunOnce:       false,
    }
    environ := &Environ{
                   RequestMethod:  r.Method,
                   ScriptName:     "", //Actually empty as app is root
                   PathInfo:       r.URL.Path,
                   QueryString:    r.URL.RawQuery,
                   ContentType:    "",//?
                   ContentLength:  string(r.ContentLength),
                   ServerName:     host[0],
                   ServerPort:     host[1],
                   ServerProtocol: r.Proto,
                   HttpRequestHeaders: httpRequestHeaders,
                   Wsgi: wsgi,
                   // ServerHeaders: 12,
    }
    out, _ := proto.Marshal(environ)
    wsgiServer.Write(out)

    response := &Response{}
    b := make([]byte, 4096)
    n, _ := wsgiServer.Read(b)
    if err := proto.Unmarshal(b[:n], response); err != nil {
      log.Fatalln("Failed to parse response:", err)
    }

    // //Response length
    fmt.Fprintf(w, "%s", response.Body)
  }
}

// func han(wsgi *wsgi.WsgiServer) func(http.ResponseWriter, *http.Request) {
//   return func(w http.ResponseWriter, r *http.Request) {
//     wsgi.Write([]byte("yo"))
//     b := make([]byte, 2048)
//     wsgi.Out.Read(b)
//     fmt.Fprintf(w, "%q", b)
//   }
// }



