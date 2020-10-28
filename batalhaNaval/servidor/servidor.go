package main

import (
	"fmt"
	"net"
	"os"
	"socketsApplication/batalhaNaval/batalhanaval"
	"strings"
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
		jogos:        make(map[net.Conn]*batalhanaval.JogadorBot),
		iniciarJogo:  make(chan net.Conn),
		encerrarJogo: make(chan net.Conn),
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
	clientes     map[net.Conn]bool                     //clientes conectados no Servidor
	cadastrar    chan net.Conn                         //canal para registrar um novo cliente
	descadastrar chan net.Conn                         //canal para cancelar o registro de um cliente que se desconectou
	jogos        map[net.Conn]*batalhanaval.JogadorBot //Jogos ativos
	iniciarJogo  chan net.Conn                         //canal para iniciar um novo jogo de um cliente
	encerrarJogo chan net.Conn                         //canal para encerrar um jogo com um cliente
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
		//se houver um novo cliente no canal de iniciarJogo, vai adicionar no mapa de jogos e iniciar um novo jogo com ele
		case socket := <-servidor.iniciarJogo:
			servidor.jogos[socket] = &batalhanaval.JogadorBot{}
			servidor.jogos[socket].IniciarJogador()
			fmt.Println("Novo jogo iniciado.")
		//se houver um novo cliente no canal de encerrarJogo, vai retirar do mapa de jogos para fechar o jogo com ele
		case socket := <-servidor.encerrarJogo:
			_, existe := servidor.jogos[socket]
			if existe {
				delete(servidor.jogos, socket)
				fmt.Println("Um jogo foi encerrado.")
			}
		}

	}
}

//Função que acontecerá o tempo todo em paralelo e será responsável por receber as mensagens dos clientes
func (servidor *Servidor) receber(cliente net.Conn) {
	for {
		mensagem := make([]byte, tamanhoMaxMensagem)
		mensagemAEnviar := make([]byte, tamanhoMaxMensagem)
		tamMensagem, err := cliente.Read(mensagem)
		if err != nil {
			servidor.descadastrar <- cliente
			cliente.Close()
			break
		}
		if tamMensagem > 0 {
			comando := mensagem[0]
			switch comando {
			//comando para iniciar um jogo com o cliente passado como parâmetro
			case 'I':
				servidor.iniciarJogo <- cliente
				mensagemAEnviar = []byte("Que vença o melhor.")
				cliente.Write(mensagemAEnviar)
			//comando para receber um tiro do cliente passado como parâmetro e em seguida atirar
			case 'A':
				//servidor.jogos[cliente].ReceberTiro()
				mensagemAEnviar = []byte("Comando para atirar\nTiro na coordenada: " + strings.Trim(string(mensagem[2:tamMensagem]), " \r\n"))
				fmt.Println(string(mensagemAEnviar))
				cliente.Write(mensagemAEnviar)
			//comando para encerrar o jogo com o cliente passado como parâmetro
			case 'S':
				servidor.encerrarJogo <- cliente
				mensagemAEnviar = []byte("O jogo foi encerrado.")
				cliente.Write(mensagemAEnviar)
			}
		}
	}
}
