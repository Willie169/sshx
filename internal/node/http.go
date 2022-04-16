package node

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/suutaku/sshx/internal/impl"
	"github.com/suutaku/sshx/internal/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (n *Node) ServeHTTPAndWS() {
	r := mux.NewRouter()
	s := http.StripPrefix("/", http.FileServer(http.Dir(n.ConfManager.Conf.VNCStaticPath)))
	r.PathPrefix("/").Handler(s)
	http.Handle("/", r)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		deviceId := r.URL.Query()["device"]
		logrus.Debug(deviceId)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logrus.Fatal(err)
			return
		}
		defer conn.Close()

		dal := impl.NewVNCImpl()
		param := impl.ImplParam{
			Config: *n.ConfManager.Conf,
			HostId: deviceId[0],
		}
		dal.Init(param)
		defer dal.Close()

		err = dal.Dial()
		if err != nil {
			logrus.Fatal(err)
			return
		}
		utils.Pipe(conn.UnderlyingConn(), *dal.Conn())
		logrus.Debug("end of gorutine")

	})
	logrus.Info("servce http at port ", n.ConfManager.Conf.LocalHTTPPort)
	http.ListenAndServe(fmt.Sprintf(":%d", n.ConfManager.Conf.LocalHTTPPort), nil)
}