package main

import (
	"fmt"
	"net"
	"os"
)

const tamanhoMaxMensagem = 512

func main() {
	porta := os.Args[1]

	listener, err := net.Listen("tcp", ":"+porta)
	if err != nil {
		fmt.Println("Ocorreu um erro ao ouvir a porta:", err)
		os.Exit(-1)
	}
	fmt.Println("Servidor ouvindo na porta", porta)
	defer listener.Close() //vai garantir que irá fechar o listener assim que fechar o programa

	servidor := Servidor{
		clientes:     make(map[net.Conn]bool),
		cadastrar:    make(chan net.Conn),
		descadastrar: make(chan net.Conn),
	}

	go servidor.iniciar()
	for {
		socket, err := listener.Accept()
		if err != nil {
			fmt.Println("Ocorreu um erro ao tentar conectar com um cliente:", err)
		}
		servidor.cadastrar <- socket
		go servidor.receber(socket)
	}
}

//Servidor struct para definir um servidor TCP
type Servidor struct {
	clientes     map[net.Conn]bool //clientes conectados no Servidor
	cadastrar    chan net.Conn     //canal para registrar um novo cliente
	descadastrar chan net.Conn     //canal para cancelar o registro de um cliente que se desconectou
}

//Funcão que irá iniciar o cadastro e o descadastro dos clientes (acontece em paralelo por uma goroutine)
func (servidor *Servidor) iniciar() {
	for {
		select {
		//se houver um cliente novo no canal de cadastro, vai adicionar isso no mapa de clientes
		case socket := <-servidor.cadastrar:
			servidor.clientes[socket] = true
			fmt.Println("Novo cliente conectado.")
		//se houver um cliente no canal de descadastro, vai retirar do mapa e fechar a conexão com o cliente
		case socket := <-servidor.descadastrar:
			_, existe := servidor.clientes[socket]
			if existe {
				delete(servidor.clientes, socket)
				fmt.Println("Um cliente foi desconectado.")
			}
		}
	}
}

//Função que acontecerá o tempo todo em paralelo e será responsável por receber as mensagens dos clientes
func (servidor *Servidor) receber(cliente net.Conn) {
	for {
		mensagem := make([]byte, tamanhoMaxMensagem)
		tamMensagem, err := cliente.Read(mensagem)
		if err != nil {
			servidor.descadastrar <- cliente
			cliente.Close()
			break
		}
		if tamMensagem > 0 {
			fmt.Println("Recebido do cliente:", string(mensagem))
			servidor.enviar(cliente, mensagem)
		}
	}
}

//Função que acontecerá o tempo todo em paralelo e será responsável por enviar mensagens aos clientes
func (servidor *Servidor) enviar(cliente net.Conn, mensagem []byte) {
	cliente.Write(mensagem)
}
