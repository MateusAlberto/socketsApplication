package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"socketsApplication/batalhaNaval/batalhanaval"
	"strconv"
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
	defer cliente.socket.Close()

	//johnLennon é um leitor
	johnLennon := bufio.NewReader(os.Stdin)

	//Loop do Menu Principal (se sair deste loop, é porque quis se desconectar do servidor)
	for {
		exibirMenuPrincipal()
		mensagem, err := johnLennon.ReadString('\n')
		if err == io.EOF {
			return //se deu EOF na leitura padrão, é porque o programa cliente foi fechado
		}
		opcao := strings.ToUpper(strings.Trim(mensagem, " \r\n"))
		switch opcao {
		case "I":
			fmt.Println("\nIniciar Jogo")
			cliente.iniciarJogo()
		case "R":
			exibirRegras()
		case "D":
			fmt.Println("Obrigado por jogar Batalha Naval!")
			return //vai fechar o socket por causa do comando defer
		default:
			fmt.Print("Por favor, digite uma opção de menu válida.\n\n")
		}
	}

}

//Cliente struct que define um cliente para se conectar no servidor via TCP
type Cliente struct {
	socket  net.Conn
	jogador batalhanaval.JogadorReal
}

//Funçao que  vai ser responsável por receber os dados vindos do servidor
func (cliente *Cliente) receber() {
	mensagem := make([]byte, tamanhoMaxMensagem)
	tamMensagem, err := cliente.socket.Read(mensagem)
	//Se tiver algum erro, fecha a conexão
	if err != nil {
		fmt.Println("Ocorreu um erro de comunicação com o servidor:", err)
		cliente.socket.Close()
	}
	if tamMensagem > 0 {
		fmt.Println("\nServidor:", string(mensagem))
	}
}

//Funçao que vai receber a resposta do seu tiro e um tiro do servidor
func (cliente *Cliente) receberTiro() (bool, bool, int, int) {
	mensagem := make([]byte, tamanhoMaxMensagem)
	tamMensagem, err := cliente.socket.Read(mensagem)
	acertou := false
	ganhou := false
	i := -1
	j := -1
	//Se tiver algum erro, fecha a conexão
	if err != nil {
		fmt.Println("Ocorreu um erro de comunicação com o servidor:", err)
		cliente.socket.Close()
	}
	if tamMensagem > 0 {
		if mensagem[0] == '1' {
			ganhou = true
			fmt.Println("\nServidor: Tiro certeiro!\n\nParabéns! Você ganhou!")
		} else {
			if mensagem[2] == '1' {
				acertou = true
			}
			i = int(mensagem[4] - '0')
			j = int(mensagem[6] - '0')
		}
	}

	if !ganhou {
		if acertou {
			fmt.Println("\nServidor: Tiro certeiro!")
		} else {
			fmt.Println("\nServidor: Tiro na água!")
		}
	}
	return ganhou, acertou, i, j
}

//Função que vai iniciar o jogo
func (cliente *Cliente) iniciarJogo() {
	mensagemAEnviar := make([]byte, tamanhoMaxMensagem)
	mensagemAEnviar = []byte("I")
	cliente.socket.Write(mensagemAEnviar)
	cliente.receber()

	johnLennon := bufio.NewReader(os.Stdin)
	cliente.carregarTabuleiro()
	for {
		zerarBuffer(mensagemAEnviar)
		exibirMenuJogo()
		mensagem, _ := johnLennon.ReadString('\n')
		mensagem = strings.Trim(strings.ToUpper(mensagem), " \r\n")
		switch mensagem {
		case "A":
			//Enviando um tiro para o servidor
			fmt.Print("Digite seu tiro: ")
			tiro, _ := johnLennon.ReadString('\n')
			tiro = strings.Trim(strings.ToUpper(tiro), " \r\n")
			iCliente, jCliente := batalhanaval.ParseTiro(tiro)

			if iCliente >= 0 && iCliente < batalhanaval.TamanhoTabuleiro && jCliente >= 0 && jCliente < batalhanaval.TamanhoTabuleiro {
				//Se ainda não atirou na posição pedida pelo usuário, vai fazer normalmente. Se não, não vai deixar ele atirar novamente
				if cliente.jogador.TabuleiroAtaque.PosicaoDesconhecida(iCliente, jCliente) {
					mensagemAEnviar = []byte("A " + tiro)
					cliente.socket.Write(mensagemAEnviar)

					//Recebendo tiro do servidor e resultado do tiro do cliente e enviando resultado do tiro do servidor
					clienteGanhou, clienteAcertou, iServidor, jServidor := cliente.receberTiro()
					cliente.jogador.TabuleiroAtaque.RegistrarTiro(clienteAcertou, iCliente, jCliente)
					if clienteGanhou {
						return //Vai sair do jogo e voltar para o menu principal
					}
					tiroServidor := string(byte(iServidor+int('A'))) + strconv.Itoa(jServidor+1)
					fmt.Println("\nTiro do oponente:", tiroServidor)
					servidorAcertou := cliente.jogador.TabuleiroDefesa.ReceberTiro(iServidor, jServidor)

					var acertouAEnviar string
					if servidorAcertou {
						//vai ver se o servidor ganhou
						if cliente.jogador.TabuleiroDefesa.AfundouTodos() {
							fmt.Println("\nSeu oponente fez um tiro certeiro!\n\nVocê perdeu! Não foi desta vez...")
							return
						}
						fmt.Println("\nSeu oponente fez um tiro certeiro!")
						acertouAEnviar = "1 "
					} else {
						fmt.Println("\nSeu oponente fez um tiro na água!")
						acertouAEnviar = "0 "
					}

					mensagemAEnviar = []byte("T " + acertouAEnviar + strconv.Itoa(iServidor) + "," + strconv.Itoa(jServidor))
					cliente.socket.Write(mensagemAEnviar)
				} else {
					fmt.Println("\nVocê já tentou este tiro! Por favor, tente um diferente.")
				}
			} else {
				fmt.Println("\nPor favor, digite um tiro válido.")
			}

		case "P":
			cliente.jogador.ImprimirTabuleiros()
		case "R":
			exibirRegras()
		case "S":
			fmt.Print("\nSaindo do jogo...\n\n")
			mensagemAEnviar = []byte("S")
			cliente.socket.Write(mensagemAEnviar)
			cliente.receber()
			return
		default:
			fmt.Println("\nPor favor, digite um comando correto.")
		}
	}
}

func (cliente *Cliente) carregarTabuleiro() {
	johnLennon := bufio.NewReader(os.Stdin)
	fmt.Print("\nPor favor, posicione seus navios no tabuleiro em um arquivo de texto.\n",
		"Coloque o caractere '-' para representar a água e 'N' para representar a parte de um navio.\n",
		"Indique o nome do arquivo onde está o seu tabuleiro montado: ")
	nomeArquivo, _ := johnLennon.ReadString('\n')
	nomeArquivo = strings.Trim(nomeArquivo, " \r\n")
	tabuleiro := LeTabuleiroArquivo(nomeArquivo)
	cliente.jogador.IniciarJogador(tabuleiro)
}

//LeTabuleiroArquivo função que lê um tabuleiro de um arquivo
func LeTabuleiroArquivo(nomeArquivo string) [][]byte {
	tabuleiro := make([][]byte, batalhanaval.TamanhoTabuleiro)
	arquivo, err := os.Open(nomeArquivo)

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
		os.Exit(-1)
	}

	leitor := bufio.NewReader(arquivo)
	for i := 0; i < batalhanaval.TamanhoTabuleiro; i++ {
		linha, err := leitor.ReadString('\n')
		linha = strings.ReplaceAll((strings.Trim(linha, " \r\n")), " ", "")
		tabuleiro[i] = []byte(linha)

		if err == io.EOF {
			break
		}
	}

	arquivo.Close()
	return tabuleiro
}

func exibirMenuPrincipal() {
	fmt.Print("\n------ MENU PRINCIPAL ------\n",
		"Digite os seguintes comandos:\n",
		"i - Iniciar o jogo\n",
		"r - Exibir Regras\n",
		"d - Desconectar\n\n",
		"Digite sua opção: ")
}

func exibirRegras() {
	fmt.Print("\nBatalha Naval é um jogo no qual dois jogadores posicionam 10 navios em um tabuleiro 10x10 e, em seguida, revezam turnos para atirarem com o objetivo de afundar um navio do oponente.\n",
		"Os navios são os seguintes:\n",
		"  - 1 porta-aviões que ocupa 5 posições\n",
		"  - 2 navios-tanque que ocupam 4 posições\n",
		"  - 3 contratorpedeiros que ocupam 3 posições\n",
		"  - 4 submarinos que ocupam 2 posições\n\n",

		"Os navios devem ser posicionados sempre na na horizontal ou na vertical, nunca na diagonal.\n",
		"Um exemplo de tabuleiro corretamente posicionado é o seguinte:\n",
		"\t  1 2 3 4 5 6 7 8 9 10\n",
		"\tA - - - - - - - N N N\n",
		"\tB N - N N - - - - - -\n",
		"\tC N - - - - - N - - N\n",
		"\tD - - - - - - N - - N\n",
		"\tE - N N N N - - - - N\n",
		"\tF - - - - - N - - - -\n",
		"\tG - - - - - N - N N N\n",
		"\tH N N - - - N - - - -\n",
		"\tI - - - - - N - - - -\n",
		"\tJ N N N N N - - - - -\n\n",

		"Depois de posicionarem seus navios nos respectivos tabuleiros, os jogadores deverão, a cada turno, indicar a coordenada que irá lançar o seu tiro.\n",
		"O tiro deve ser indicado pela letra (A-J) da linha e número (1-10) da coluna.\n",
		"Um exemplo de tiro é D7.\n\n",

		"Ganha o primeiro jogador que derrubar todos os navios do oponente.\n\n")
}

func exibirMenuJogo() {
	fmt.Print("\n------ MENU BATALHA NAVAL ------\n",
		"Digite sua opção:\n",
		"a - Atacar\n",
		"p - Imprimir Tabuleiros\n",
		"r - Exibir Regras\n",
		"s - Sair do jogo\n\n",
		"Digite sua opção: ")
}

//Pequena função para zerar o buffer
func zerarBuffer(array []byte) {
	for i := 0; i < len(array); i++ {
		array[i] = 0
	}
}
