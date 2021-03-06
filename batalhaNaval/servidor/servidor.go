package main

import (
	"fmt"
	"net"
	"os"
	"socketsApplication/batalhaNaval/batalhanaval"
	"strconv"
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
	mensagem := make([]byte, tamanhoMaxMensagem)
	mensagemAEnviar := make([]byte, tamanhoMaxMensagem)
	for {
		zerarBuffer(mensagem)
		tamMensagem, err := cliente.Read(mensagem)
		if err != nil {
			servidor.descadastrar <- cliente
			cliente.Close()
			break
		}
		fmt.Println(string(mensagem[:tamMensagem]))
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
				var ganhou bool
				iCliente, jCliente := batalhanaval.ParseTiro(string(mensagem[2:tamMensagem]))
				acertou := servidor.jogos[cliente].TabuleiroDefesa.ReceberTiro(iCliente, jCliente)

				if acertou {
					ganhou = servidor.jogos[cliente].TabuleiroDefesa.AfundouTodos()
				}

				//Se o cliente ganhou, vai enviar para ele apenas o número 1
				//Se não, vai enviar o número 0, seguido de 0/1 indicando se ele acertou e depois do tiro do servidor
				if ganhou {
					mensagemAEnviar = []byte("1")
				} else {
					var acertouAEnviar string
					if acertou {
						acertouAEnviar = "1 "
					} else {
						acertouAEnviar = "0 "
					}
					iServidor, jServidor := servidor.jogos[cliente].Tiro()
					mensagemAEnviar = []byte("0 " + acertouAEnviar + strconv.Itoa(iServidor) + "," + strconv.Itoa(jServidor))
				}
				cliente.Write(mensagemAEnviar)
				if ganhou {
					servidor.encerrarJogo <- cliente
				}
			//comando para receber o resultado do tiro do cliente
			case 'T':
				acertou := mensagem[2] == '1'
				i := int(mensagem[4] - '0')
				j := int(mensagem[6] - '0')
				if acertou {
					servidor.jogos[cliente].UltimoTiroCerteiro[0] = i
					servidor.jogos[cliente].UltimoTiroCerteiro[1] = j
				}
				servidor.jogos[cliente].TabuleiroAtaque.RegistrarTiro(acertou, i, j)
			//comando para encerrar o jogo com o cliente passado como parâmetro
			case 'S':
				servidor.encerrarJogo <- cliente
				mensagemAEnviar = []byte("O jogo foi encerrado.")
				cliente.Write(mensagemAEnviar)
			}
		}
	}
}

//Pequena função para zerar o buffer
func zerarBuffer(array []byte) {
	for i := 0; i < len(array); i++ {
		array[i] = 0
	}
}
