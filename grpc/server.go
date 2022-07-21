package main

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type Server struct {
	pipe  *PipeListener
	gpipe *grpc.Server

	tcp  net.Listener
	gtcp *grpc.Server

	proxyMux *runtime.ServeMux
}

func NewServer(addr string) (s *Server, e error) {
	tcp, e := net.Listen(`tcp`, addr)
	if e != nil {
		return
	}

	pipe := ListenPipe()
	clientConn, e := grpc.Dial(`pipe`,
		grpc.WithInsecure(),
		grpc.WithContextDialer(func(c context.Context, s string) (net.Conn, error) {
			return pipe.DialContext(c, `pipe`, s)
		}),
	)
	if e != nil {
		return
	}
	proxyMux := runtime.NewServeMux()

	s = &Server{
		pipe:     pipe,
		tcp:      tcp,
		gpipe:    newGRPCServer(proxyMux, clientConn),
		gtcp:     newGRPCServer(nil, nil),
		proxyMux: proxyMux,
	}
	return
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contextType := r.Header.Get(`Content-Type`)
	if r.ProtoMajor == 2 && strings.Contains(contextType, `application/grpc`) {
		s.gtcp.ServeHTTP(w, r) // application/grpc 路由給 grpc
	} else {
		s.proxyMux.ServeHTTP(w, r) // 非 grpc 路由給 gateway
	}
}
func (s *Server) Serve() (e error) {
	go s.gpipe.Serve(s.pipe)

	// 配置 h2c
	var httpServer http.Server
	var http2Server http2.Server
	e = http2.ConfigureServer(&httpServer, &http2Server)
	if e != nil {
		return
	}
	httpServer.Handler = h2c.NewHandler(s, &http2Server)
	// http.Serve 不支持 h2c
	// 如果直接使用 http.Serve 將使用 grpc 客戶端 無法正常訪問
	e = httpServer.Serve(s.tcp)
	return
}
func (s *Server) ServeTLS(certFile, keyFile string) (e error) {
	go s.gpipe.Serve(s.pipe)

	e = http.ServeTLS(s.tcp, s, certFile, keyFile)
	return
}
