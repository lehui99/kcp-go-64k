// +build !linux

package kcp

import (
	"sync/atomic"
)

func (s *UDPSession) txLoop() {
	for {
		select {
		case txqueue := <-s.chTxQueue:
			nbytes := 0
			for k := range txqueue {
				if n, err := s.conn.WriteTo(txqueue[k].Buffers[0], txqueue[k].Addr); err == nil {
					nbytes += n
				} else {
					s.notifyWriteError(err)
					return
				}
			}
			atomic.AddUint64(&DefaultSnmp.OutPkts, uint64(len(txqueue)))
			atomic.AddUint64(&DefaultSnmp.OutBytes, uint64(nbytes))
		case <-s.die:
			return
		}
	}
}

func (l *Listener) txLoop() {
	for {
		select {
		case txqueue := <-l.chTxQueue:
			nbytes := 0
			for k := range txqueue {
				if n, err := l.conn.WriteTo(txqueue[k].Buffers[0], txqueue[k].Addr); err == nil {
					nbytes += n
					xmitBuf.Put(txqueue[k].Buffers[0])
				} else {
					l.Close()
					return
				}
			}
			atomic.AddUint64(&DefaultSnmp.OutPkts, uint64(len(txqueue)))
			atomic.AddUint64(&DefaultSnmp.OutBytes, uint64(nbytes))
		case <-l.die:
			return
		}
	}
}
