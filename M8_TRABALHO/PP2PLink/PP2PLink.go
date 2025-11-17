package PP2PLink

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

type PP2PLink_Req_Message struct {
	To      string
	Message string
}

type PP2PLink_Ind_Message struct {
	From    string
	Message string
}

type PP2PLink struct {
	Ind   chan PP2PLink_Ind_Message
	Req   chan PP2PLink_Req_Message
	Run   bool
	dbg   bool
	Cache map[string]net.Conn // cache de conexoes
	mu    sync.Mutex          // protege Cache
}

func NewPP2PLink(_address string, _dbg bool) *PP2PLink {
	p2p := &PP2PLink{
		Req:   make(chan PP2PLink_Req_Message, 16),
		Ind:   make(chan PP2PLink_Ind_Message, 16),
		Run:   true,
		dbg:   _dbg,
		Cache: make(map[string]net.Conn)}
	p2p.outDbg(" Init PP2PLink!")
	p2p.Start(_address)
	return p2p
}

func (module *PP2PLink) outDbg(s string) {
	if module.dbg {
		fmt.Println(". . . . . . . . . . . . . . . . . [ PP2PLink msg : " + s + " ]")
	}
}

func (module *PP2PLink) Start(address string) {
	// LISTENER (recebimento)
	go func() {
		listen, err := net.Listen("tcp4", address)
		if err != nil {
			module.outDbg("erro ao criar listener: " + err.Error())
			return
		}
		module.outDbg("ok   : listener criado em " + address)
		for {
			conn, err := listen.Accept()
			if err != nil {
				// não mata o listener por erro temporário
				module.outDbg("erro accept: " + err.Error())
				time.Sleep(200 * time.Millisecond)
				continue
			}
			module.outDbg("ok   : conexao aceita com outro processo.")
			// garante configuração TCP em conexões aceitas
			if tcpConn, ok := conn.(*net.TCPConn); ok {
				_ = tcpConn.SetKeepAlive(true)
				_ = tcpConn.SetKeepAlivePeriod(30 * time.Second)
				_ = tcpConn.SetLinger(0)
			}
			// handler para cada conexão aceita
			go module.handleIncoming(conn)
		}
	}()

	// SENDER (envio a partir do canal Req)
	go func() {
		for {
			message := <-module.Req
			if !module.Run {
				return
			}
			module.Send(message)
		}
	}()
}

// handleIncoming lê mensagens de uma conexão aceita e repassa ao módulo superior.
// Fecha a conexão quando a outra ponta cai; não derruba o listener.
func (module *PP2PLink) handleIncoming(conn net.Conn) {
	defer func() {
		_ = conn.Close()
		module.outDbg("conexao handler finalizado: " + conn.RemoteAddr().String())
	}()

	for {
		// le 4 bytes do tamanho
		bufTam := make([]byte, 4)
		_, err := io.ReadFull(conn, bufTam)
		if err != nil {
			if err == io.EOF {
				module.outDbg("conexao fechada pelo remoto: " + conn.RemoteAddr().String())
			} else {
				module.outDbg("erro leitura tam: " + err.Error() + " de " + conn.RemoteAddr().String())
			}
			return
		}

		tam, err := strconv.Atoi(string(bufTam))
		if err != nil || tam < 0 {
			module.outDbg("tamanho inválido recebido: '" + string(bufTam) + "' de " + conn.RemoteAddr().String())
			return
		}

		bufMsg := make([]byte, tam)
		_, err = io.ReadFull(conn, bufMsg)
		if err != nil {
			module.outDbg("erro leitura msg: " + err.Error() + " de " + conn.RemoteAddr().String())
			return
		}

		msg := PP2PLink_Ind_Message{
			From:    conn.RemoteAddr().String(),
			Message: string(bufMsg)}
		// envia para módulo superior (pode bloquear se nao houver receptor; isso é intencional)
		module.Ind <- msg
	}
}

// getCachedConn retorna uma conexão da cache (protegida) se existir
func (module *PP2PLink) getCachedConn(addr string) (net.Conn, bool) {
	module.mu.Lock()
	defer module.mu.Unlock()
	c, ok := module.Cache[addr]
	return c, ok
}

// setCachedConn guarda uma conexão na cache (protegida)
func (module *PP2PLink) setCachedConn(addr string, c net.Conn) {
	module.mu.Lock()
	module.Cache[addr] = c
	module.mu.Unlock()
}

// delCachedConn remove e fecha a conexão cache se existir
func (module *PP2PLink) delCachedConn(addr string) {
	module.mu.Lock()
	if c, ok := module.Cache[addr]; ok {
		_ = c.Close()
		delete(module.Cache, addr)
	}
	module.mu.Unlock()
}

// dialWithOpts tenta abrir conexão TCP com opções e reconexões curtas
func (module *PP2PLink) dialWithOpts(addr string) (net.Conn, error) {
	var conn net.Conn
	var err error
	backoff := []time.Duration{0, 150 * time.Millisecond, 300 * time.Millisecond}
	for _, delay := range backoff {
		if delay > 0 {
			time.Sleep(delay)
		}
		conn, err = net.DialTimeout("tcp", addr, 2*time.Second)
		if err == nil {
			// garante opções TCP
			if tcpConn, ok := conn.(*net.TCPConn); ok {
				_ = tcpConn.SetKeepAlive(true)
				_ = tcpConn.SetKeepAlivePeriod(30 * time.Second)
				_ = tcpConn.SetLinger(0)
			}
			module.outDbg("ok   : conexao iniciada com outro processo")
			return conn, nil
		}
		// loga e tenta novamente
		module.outDbg("erro dial " + addr + " : " + err.Error())
	}
	return nil, err
}

// safeWrite faz múltiplas tentativas de escrita (uma única re-dial se necessário)
func (module *PP2PLink) safeWrite(conn net.Conn, data []byte) error {
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_, err := conn.Write(data)
	if err == nil {
		return nil
	}
	// se erro, tenta reescrever após pequena espera (mas sem fechar de imediato aqui)
	module.outDbg("erro write: " + err.Error())
	return err
}

// Send envia uma mensagem ao endereço. tenta reaproveitar conexao; em caso de erro tenta reconectar.
func (module *PP2PLink) Send(message PP2PLink_Req_Message) {
	// prepara payload: 4-digit length + message bytes
	msgBytes := []byte(message.Message)
	lengthStr := strconv.Itoa(len(msgBytes))
	for len(lengthStr) < 4 {
		lengthStr = "0" + lengthStr
	}
	header := []byte(lengthStr)

	// tenta usar conexao cacheada
	conn, ok := module.getCachedConn(message.To)
	if ok && conn != nil {
		// tentativa de escrita com a conexao existente
		if err := module.safeWrite(conn, header); err == nil {
			if err := module.safeWrite(conn, msgBytes); err == nil {
				return // sucesso usando conexao cacheada
			}
		}
		// se falhou, removemos a conexao considerada ruim e fechamos
		module.outDbg("conexao cacheada apresentou erro, removendo cache e tentando reconectar: " + message.To)
		module.delCachedConn(message.To)
	}

	// tenta dial/reconectar e escrever
	conn, err := module.dialWithOpts(message.To)
	if err != nil {
		// falha em conectar
		module.outDbg("falha ao conectar para " + message.To + " : " + err.Error())
		return
	}
	// armazena na cache antes de usar para outras goroutines
	module.setCachedConn(message.To, conn)

	// tenta enviar; se falhar, tenta uma reconexao única
	if err := module.safeWrite(conn, header); err != nil {
		module.outDbg("erro ao escrever header apos reconectar: " + err.Error())
		module.delCachedConn(message.To)
		// tenta uma reconexao final
		conn2, err2 := module.dialWithOpts(message.To)
		if err2 != nil {
			module.outDbg("reconexao final falhou: " + err2.Error())
			return
		}
		module.setCachedConn(message.To, conn2)
		if err := module.safeWrite(conn2, header); err != nil {
			module.outDbg("erro ao escrever header apos reconexao final: " + err.Error())
			return
		}
		if err := module.safeWrite(conn2, msgBytes); err != nil {
			module.outDbg("erro ao escrever msg apos reconexao final: " + err.Error())
			return
		}
		return
	}

	// envia corpo
	if err := module.safeWrite(conn, msgBytes); err != nil {
		module.outDbg("erro ao escrever msg apos conectar: " + err.Error())
		// remove cache para tentar reconectar na proxima vez
		module.delCachedConn(message.To)
		return
	}
}
