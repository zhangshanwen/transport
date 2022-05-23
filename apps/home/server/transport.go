package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/zhangshanwen/transport/apps/home/conf"
)

func (t *Transponder) currentLimit() (err error) {
	/*
		默认限制1秒钟最大连接数为100个
	*/
	now := time.Now()
	if t.lastConnectTime != nil && now.Sub(*t.lastConnectTime).Seconds() <= float64(t.connectTime) {
		if t.connect > t.maxConnect {
			return errors.New("超过最大连接数")
		}
		atomic.AddInt32(&t.connect, 1)
	} else {
		t.lastConnectTime = &now
		t.connect = 1
	}
	return
}

func (t *Transponder) Listen() (err error) {
	go func() {
		if err = http.ListenAndServe(fmt.Sprintf("%s:%s", conf.C.Host, conf.C.Port), t); err != nil {
			logrus.Error(err)
			return
		}
	}()
	return
}

func (t *Transponder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = t.currentLimit(); err != nil {
		t.NotFound(w, "请稍后再试......")
		logrus.Error(err)
		return
	}
	for _, module := range t.Modules {
		if strings.HasPrefix(r.URL.Path, module.Prefix) {
			t.mu.Lock()
			if int(module.Index) > len(module.Slaves)-1 {
				atomic.AddInt32(&module.Index, -module.Index)
			}
			proxy := httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: module.Scheme,
				Host:   module.Slaves[module.Index].Addr,
			})
			atomic.AddInt32(&module.Index, 1)
			logrus.Infof("转发模块%s,路径%s,第%v个副本.......", module.Name, module.Prefix, module.Index)
			proxy.ServeHTTP(w, r)
			t.mu.Unlock()
		}
	}
}

func (t *Transponder) NotFound(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	w.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprintf(w, msg)
}
