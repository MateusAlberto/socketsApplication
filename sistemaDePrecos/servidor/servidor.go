package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"

	geo "github.com/kellydunn/golang-geo"
)

const tamanhoMaxMensagem = 512

//Nomes dos arquivos
const arquivoDiesel = "diesel.txt"
const arquivoAlcool = "alcool.txt"
const arquivoGasolina = "gasolina.txt"

func main() {
	porta := os.Args[1]

	enderecoUDP, err := net.ResolveUDPAddr("udp", "localhost:"+porta)
	if err != nil {
		fmt.Println("\nEndereço incorreto:", err)
		os.Exit(-1)
	}

	socket, err := net.ListenUDP("udp", enderecoUDP)
	if err != nil {
		fmt.Println("\nOcorreu um erro ao ouvir a porta:", err)
		os.Exit(-1)
	}
	fmt.Println("Servidor ouvindo na porta", porta)
	defer socket.Close() //vai garantir que irá fechar o listener assim que fechar o programa

	//alocando memória para as variáveis necessárias no loop
	buffer := make([]byte, tamanhoMaxMensagem)
	mensagemAEnviar := make([]byte, tamanhoMaxMensagem)
	var mensagemEnviada string
	var tamMensagem int
	var endereco *net.UDPAddr

	for {
		zerarBuffer(buffer)
		zerarBuffer(mensagemAEnviar)
		tamMensagem, endereco, err = socket.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("\nErro ao ler mensagem do cliente:", err)
		}

		mensagemEnviada = string(buffer[:tamMensagem])
		fmt.Println("\nRecebido do cliente:", mensagemEnviada)

		//buffer[0] vai conter os comandos D ou P
		switch buffer[0] {
		case 'D':
			if registrarDado(mensagemEnviada) {
				mensagemAEnviar = []byte("Dado registrado com sucesso.")
			} else {
				mensagemAEnviar = []byte("Não foi possível registrar o dado.")
			}
		case 'P':
			resultadoPesquisa := pesquisar(mensagemEnviada)
			mensagemAEnviar = []byte(resultadoPesquisa)
		default:
			fmt.Println("\nO cliente enviou um comando errado.")
			mensagemAEnviar = []byte("Comando Desconhecido.")
		}

		_, err = socket.WriteToUDP(mensagemAEnviar, endereco)
		if err != nil {
			fmt.Println("\nErro ao tentar enviar uma mensagem para o cliente:", err)
		}
	}
}

/*D 451020 1 3299 42,5 65,03
 * valores[0] = D (comando)
 * valores[1] = 45102 (identificador)
 * valores[2] = 1 (tipo do combustível)
 * valores[3] = 3299 (preço do combustível x 1000)
 * valores[4] = 42,5 (latitude em graus)
 * valores[5] = 65,03 (longitude em graus)
 */
//Função para registrar um dado
//retorna um booleano indicando se conseguiu registrar o dado com sucesso ou não
func registrarDado(dado string) bool {
	var nomeArquivo string
	valores := strings.Split(dado, " ")

	//Escolha do arquivo de acordo com o tipo de combustível
	switch valores[2] {
	case "0":
		nomeArquivo = arquivoDiesel
	case "1":
		nomeArquivo = arquivoAlcool
	case "2":
		nomeArquivo = arquivoGasolina
	}

	arquivo, err := os.OpenFile(nomeArquivo, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer arquivo.Close() //vai fechar o arquivo assim que a função retornar

	if err != nil {
		fmt.Println("Erro ao abrir o arquivo", arquivo, ":", err)
		return false
	}
	_, err = arquivo.WriteString(valores[1] + " " + valores[3] + " " + valores[4] + " " + valores[5] + "\n")
	if err != nil {
		fmt.Println("Erro ao escrever no arquivo", arquivo, ":", err)
		return false
	}
	return true
}

/*P 451020 1 10 42,5 65,03
 * valores[0] = P (comando)
 * valores[1] = 45102 (identificador)
 * valores[2] = 1 (tipo do combustível)
 * valores[3] = 10 (raio de busca em km)
 * valores[4] = 42,5 (latitude em graus)
 * valores[5] = 65,03 (longitude em graus)
 */
//Função que irá realizar uma pesquisa nos arquivos da localização e preço
//do posto com preço mais baixo dentro do raio especificado
//A forma como será implementada não possui uma grande sofisticação em relação à leitura de arquivo
//Assim, ela funciona para arquivos cujo todo o conteúdo cabe na memória RAM (o que é o caso apenas para um TP_)
func pesquisar(dado string) string {
	var nomeArquivo string
	valores := strings.Split(dado, " ")

	//Escolha do arquivo de acordo com o tipo de combustível
	switch valores[2] {
	case "0":
		nomeArquivo = arquivoDiesel
	case "1":
		nomeArquivo = arquivoAlcool
	case "2":
		nomeArquivo = arquivoGasolina
	}

	arquivo, err := ioutil.ReadFile(nomeArquivo)
	if err != nil {
		fmt.Println("Ocorreu um erro ao tentar abrir o arquivo:", err)
		return "Não foi possível abrir o arquivo."
	}

	linhas := strings.Split(string(arquivo), "\n")

	raio, err := strconv.Atoi(valores[3])
	if err != nil {
		fmt.Println(err)
		return "Por favor, digite um raio válido."
	}

	//Processamento da entrada
	latitude, err := strconv.ParseFloat(strings.ReplaceAll(valores[4], ",", "."), 64)
	if err != nil {
		fmt.Println(err)
		return "Por favor, digite um valor de latitude correto."
	}
	longitude, err := strconv.ParseFloat(strings.ReplaceAll(valores[5], ",", "."), 64)
	if err != nil {
		fmt.Println(err)
		return "Por favor, digite um valor de longitude correto."
	}
	centro := geo.NewPoint(latitude, longitude)

	resp := "\n\tCombustível: " + nomeArquivo[0:len(nomeArquivo)-4] + "\n" + menorPrecoNoRaio(linhas[:len(linhas)-1], centro, raio)
	return resp
}

/* 59575 2955 102 50
 * valores[0] = 59575 (identificador)
 * valores[1] = 2955 (preço combustível x 1000)
 * valores[2] = 102 (latitude em graus)
 * valores[3] = 50 (longitude em graus)
 */
//Função que recebe os preços para um determinado combustível e calcula o posto
//com o menor preço dentro da região delimitada pelo centro e raio passados
func menorPrecoNoRaio(postos []string, centro *geo.Point, raio int) string {
	var valores []string
	var preco int

	naRegiao := postosNaRegiao(postos, centro, raio)

	if len(naRegiao) == 0 { //não tem postos dentro da região pedida
		return "\tNenhum posto encontrado na região pedida."
	}

	valores = strings.Split(naRegiao[0], " ")
	menorPreco, _ := strconv.Atoi(valores[1])
	dadoMenorPreco := naRegiao[0]

	fmt.Println("\nDado:", naRegiao[0])
	for i := 1; i < len(naRegiao); i++ {
		fmt.Println("\nDado:", naRegiao[i])
		valores = strings.Split(naRegiao[i], " ")
		preco, _ = strconv.Atoi(valores[1])
		if preco < menorPreco {
			menorPreco = preco
			dadoMenorPreco = naRegiao[i]
		}
	}
	valores = strings.Split(dadoMenorPreco, " ")
	resp := "\tPreço: " + valores[1] + " (/1000)\n\tLatitude: " + valores[2] + "°\n\tLongitude: " + valores[3] + "°"

	return resp
}

//Função que adiciona em um array de strings todos os dados de preços
//cujas coordenadas estão dentro do círculo delimitado pelo centro e raio passados
func postosNaRegiao(postos []string, centro *geo.Point, raio int) []string {
	naRegiao := make([]string, 0)
	var valores []string
	var latitude, longitude float64
	var coordenada *geo.Point

	for _, posto := range postos {
		valores = strings.Split(posto, " ")
		latitude, _ = strconv.ParseFloat(strings.ReplaceAll(valores[2], ",", "."), 64)
		longitude, _ = strconv.ParseFloat(strings.ReplaceAll(valores[3], ",", "."), 64)
		coordenada = geo.NewPoint(latitude, longitude)

		//Se a distância em quilômetros entre os dois pontos(centro e a coordenada do posto corrente)
		//for menor que o raio, o posto está na região e será adicionado para o retorno
		fmt.Println("Distância:", centro.GreatCircleDistance(coordenada))
		if centro.GreatCircleDistance(coordenada) <= float64(raio) {
			naRegiao = append(naRegiao, posto)
		}
	}

	return naRegiao
}

//Pequena função para zerar o buffer
func zerarBuffer(array []byte) {
	for i := 0; i < len(array); i++ {
		array[i] = 0
	}
}
