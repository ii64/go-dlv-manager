package server

import "io"

func (s *Server) createStream(reader io.Reader, writer io.Writer, client2dap bool) (err error) {
	var n int
	var n2 int // write nb
	var buf = make([]byte, s.opt.BufferSize)
	for {
		n, err = reader.Read(buf)
		if err != nil { // EOF, anything.
			return
		}
		tmp := buf[:n] // slice header copy
		tmp = s.peekData(tmp, client2dap)
	writeBack:
		n2, err = writer.Write(tmp)
		if err != nil { // EOF, anything.
			return
		}
		if n2 < n {
			tmp = tmp[n2:]
			n -= n2
			goto writeBack
		}
		// conn = s.connMod(conn)
	}
}

func (s *Server) peekData(data []byte, client2dap bool) []byte {
	if client2dap {
		s.logger.Debug().Bytes("data", data).Msg("packet->")
	} else {
		s.logger.Debug().Bytes("data", data).Msg("packet<-")
	}
	return data
}
