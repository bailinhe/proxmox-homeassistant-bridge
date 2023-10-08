package proxmox

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/luthermonson/go-proxmox"
	"go.uber.org/zap"
)

type Server struct {
	tokenID   string
	secret    string
	serverURL string
	logger    *zap.Logger
	insecure  bool

	client *proxmox.Client
}

type Opt func(s *Server)

func NewServer(serverURL string, opts ...Opt) *Server {
	s := &Server{
		serverURL: serverURL,
		insecure:  false,
		logger:    &zap.Logger{},
	}

	for _, opt := range opts {
		opt(s)
	}

	s.logger = s.logger.With(zap.String("component", "proxmox"))

	return s
}

func WithAPIToken(tokenID, secret string) Opt {
	return func(s *Server) {
		s.secret = secret
		s.tokenID = tokenID
	}
}

func WithLogger(l *zap.Logger) Opt {
	return func(s *Server) {
		s.logger = l
	}
}

func WithInsecure() Opt {
	return func(s *Server) {
		s.insecure = true
	}
}

func (s *Server) Start() {
	opts := []proxmox.Option{
		proxmox.WithAPIToken(s.tokenID, s.secret),
	}

	if s.insecure {
		opts = append(opts, proxmox.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}))
	}

	s.client = proxmox.NewClient(
		fmt.Sprintf("%s/api2/json", s.serverURL),
		opts...,
	)

	version, err := s.client.Version()
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	s.logger.Sugar().Debugf("proxmox version is: %s", version.Release)

	NodeStatuses, err := s.client.Nodes()
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	vms := proxmox.VirtualMachines{}

	for _, st := range NodeStatuses {
		node, err := s.client.Node(st.Node)
		if err != nil {
			s.logger.Error(err.Error())
			continue
		}

		vm, err := node.VirtualMachines()
		if err != nil {
			s.logger.Error(err.Error())
			continue
		}

		vms = append(vms, vm...)
	}

	for _, vm := range vms {
		j, _ := json.MarshalIndent(vm, "", "  ")
		s.logger.Sugar().Debug(string(j))
	}
}
