package main

import (
	"fmt"
	"net"
	"os"
)

const tamanhoMaxMensagem = 512

func main() {
	porta := os.Args[1]

	enderecoUDP, err := net.ResolveUDPAddr("udp", "localhost:"+porta)
	if err != nil {
		fmt.Println("Endereço incorreto:", err)
		os.Exit(-1)
	}

	socket, err := net.ListenUDP("udp", enderecoUDP)
	if err != nil {
		fmt.Println("Ocorreu um erro ao ouvir a porta:", err)
		os.Exit(-1)
	}
	fmt.Println("Servidor ouvindo na porta", porta)
	defer socket.Close() //vai garantir que irá fechar o listener assim que fechar o programa

	for {
		buffer := make([]byte, tamanhoMaxMensagem)

		tamMensagem, endereco, err := socket.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Erro ao ler mensagem do cliente:", err)
		}

		fmt.Println("Cliente : ", endereco)
		fmt.Println("Recebido do cliente:", string(buffer[:tamMensagem]))

		mensagem := []byte((buffer[:tamMensagem]))
		_, err = socket.WriteToUDP(mensagem, endereco)
		if err != nil {
			fmt.Println("Erro ao tentar enviar uma mensagem para o cliente:", err)
		}
	}
}
