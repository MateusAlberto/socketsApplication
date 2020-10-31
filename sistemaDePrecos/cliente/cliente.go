package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

const tamanhoMaxMensagem = 512

func main() {
	enderecoEPorta := os.Args[1]

	enderecoServidor, err := net.ResolveUDPAddr("udp", enderecoEPorta)
	if err != nil {
		fmt.Println("Endereço incorreto:", err)
		os.Exit(-1)
	}

	fmt.Println("Conectando no servidor", enderecoEPorta)
	//Conectando no endereço e porta especificados (pelo padrão endereço:porta)
	socket, err := net.DialUDP("udp", nil, enderecoServidor)
	if err != nil {
		fmt.Println("Ocorreu um erro ao tentar se conectar ao servidor:", err)
		os.Exit(-1)
	}
	fmt.Println("Conexão estabelecida com", enderecoEPorta)
	fmt.Println("Endereço do servidor:", socket.RemoteAddr().String())
	fmt.Println("Endereço local do cliente:", socket.LocalAddr().String())

	defer socket.Close()

	for {
		// write a message to server
		johnLennon := bufio.NewReader(os.Stdin)
		mensagem, err := johnLennon.ReadString('\n')
		if err == io.EOF {
			return //se deu EOF na leitura padrão, é porque o programa cliente foi fechado
		}

		_, err = socket.Write([]byte(mensagem))

		if err != nil {
			fmt.Println("Erro ao tentar escrever uma mensagem para o servidor:", err)
		}

		//recebendo uma mensagem do servidor
		buffer := make([]byte, tamanhoMaxMensagem)
		tamMensagem, endereco, err := socket.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Erro ao tentar ler uma mensagem do servidor:", err)
		}

		fmt.Println("Servidor UDP:", endereco)
		fmt.Println("Recebido do servidor UDP:", string(buffer[:tamMensagem]))
	}
}
