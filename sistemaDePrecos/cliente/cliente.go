package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
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
	johnLennon := bufio.NewReader(os.Stdin)

	for {
		exibirMenu()

		mensagem, err := johnLennon.ReadString('\n')
		if err == io.EOF {
			return //se deu EOF na leitura padrão, é porque o programa cliente foi fechado
		}

		mensagem = strings.ToUpper(strings.Trim(mensagem, " \r\n"))

		switch mensagem {
		case "D":
			fmt.Print("\nEntrada de dados do tipo: combustível preço latitude longitude\n",
				"combustível: 0 - diesel, 1 - álcool, 2- gasolina\n",
				"preço: valor x 1000 (sem vírgula)\n",
				"coordenadas: lagitude e longitude em graus.\n\n",
				"Digite a entrada de dados: ")
			lerEEnviar(socket, "D ")
		case "P":
			fmt.Print("\nEntrada de dados do tipo: combustível raio latitude longitude\n",
				"combustível: 0 - diesel, 1 - álcool, 2- gasolina\n",
				"raio: raio de busca em quilômetros\n",
				"coordenadas: latitude e longitude em graus.\n\n",
				"Digite a entrada de dados: ")
			lerEEnviar(socket, "P ")
		case "S":
			fmt.Println("\nSaindo do programa...")
			return //irá fechar o socket por causa do comando defer
		default:
			fmt.Println("\nComando não existe. Por favor, digite um comando correto.")
		}
	}
}

//Exibe o menu para que a pessoa possa escolher o que fazer
func exibirMenu() {
	fmt.Print("\n------ MENU SISTEMA DE PREÇOS ------\n",
		"Digite sua opção:\n",
		"d - Enviar dados ao servidor\n",
		"p - Pesquisar preço de combustível\n",
		"s - Sair do programa\n\n",
		"Digite sua opção: ")
}

//Lê o dado para enviar e chama a função enviarMensagemAoServidor para enviá-lo
func lerEEnviar(socket *net.UDPConn, tipoOperacao string) {
	johnLennon := bufio.NewReader(os.Stdin)

	entrada, err := johnLennon.ReadString('\n')
	if err == io.EOF {
		return //se deu EOF na leitura padrão, é porque o programa cliente foi fechado
	}
	entrada = strings.Trim(entrada, " \r\n")

	//Escolhi a porta do cliente como um identificador da mensagem
	porta := strings.Split(socket.LocalAddr().String(), ":")[1]

	mensagem := tipoOperacao + porta + " " + entrada
	enviarEReceberDoServidor(mensagem, socket)
}

//Envia a mensagem passada como parâmetro para o servidor, e se algo estiver errado, tenta reenviar uma vez
func enviarEReceberDoServidor(mensagem string, socket *net.UDPConn) {
	socket.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_, err := socket.Write([]byte(mensagem))
	if err != nil {
		fmt.Println("\nErro ao tentar escrever uma mensagem para o servidor. Tentando retransmitir...")
		socket.SetWriteDeadline(time.Now().Add(5 * time.Second))
		_, err = socket.Write([]byte(mensagem))
		if err != nil {
			fmt.Println("\nErro ao tentar retransmitir a mensagem para o servidor. Tente novamente mais tarde.")
			//return
		} else {
			fmt.Println("\nConseguiu retransmitir a mensagem para o servidor.")
		}

	}

	buffer := make([]byte, tamanhoMaxMensagem)
	socket.SetReadDeadline(time.Now().Add(5 * time.Second))
	tamMensagem, _, err := socket.ReadFromUDP(buffer)

	if err != nil {
		fmt.Println("\nErro ao tentar ler uma mensagem do servidor. Tentando retransmitir...")
		//Retransmissão
		socket.SetReadDeadline(time.Now().Add(5 * time.Second))
		tamMensagem, _, err = socket.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("\nNão foi possível ler a mensagem do servidor.")
		} else {
			fmt.Println("\nConseguiu ler a mensagem do servidor.")
		}
	} else {
		if tamMensagem > 0 {
			fmt.Println("\nServidor:", string(buffer[0:tamMensagem]))
		}
	}
}
