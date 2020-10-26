package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const tamanhoMaxMensagem = 512

func main() {
	enderecoEPorta := os.Args[1]

	fmt.Println("Conectando no servidor", enderecoEPorta)
	//Conectando no endereço e porta especificados (pelo padrão endereço:porta)
	socket, err := net.Dial("tcp", enderecoEPorta)
	if err != nil {
		fmt.Println("Ocorreu um erro ao tentar se conectar ao servidor:", err)
		os.Exit(-1)
	}
	fmt.Println("Conexão realizada com sucesso")

	cliente := &Cliente{socket: socket}
	go cliente.receber() //vai rodar em paralelo o tempo todo

	for {
		//johnLennon é um leitor
		johnLennon := bufio.NewReader(os.Stdin)
		mensagem, err := johnLennon.ReadString('\n')
		if err != io.EOF {
			cliente.socket.Write([]byte(strings.TrimRight(mensagem, "\n")))
		}
	}
}

//Cliente sctruct que define um cliente para se conectar no servidor via TCP
type Cliente struct {
	socket net.Conn
}

//Funçao que rodará em paralelo e vai ser responsável por receber os dados vindos do servidor
func (cliente *Cliente) receber() {
	for {
		mensagem := make([]byte, tamanhoMaxMensagem)
		tamMensagem, err := cliente.socket.Read(mensagem)
		//Se tiver algum erro, fecha a conexão
		if err != nil {
			fmt.Println("Ocorreu um erro de comunicação com o servidor:", err)
			cliente.socket.Close()
			break
		}
		if tamMensagem > 0 {
			fmt.Println("Servidor:", string(mensagem))
		}
	}
}
