package batalhanaval

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

//Cliente struct para definir um jogador real
type Cliente struct {
	tabuleiroAtaque *Tabuleiro
	tabuleiroDefesa *Tabuleiro
}

//NovoJogador construtor de um jogador real
func (c *Cliente) NovoJogador() {
	tabDefesa := leTabuleiroArquivo()
	c.tabuleiroAtaque = NovoTabDefesa(tabDefesa)
	c.tabuleiroAtaque = NovoTabAtaque()
}

//leTabuleiroArquivo função que lê um tabuleiro de um arquivo
func leTabuleiroArquivo() [TamanhoTabuleiro][TamanhoTabuleiro]byte {
	var tabuleiro [TamanhoTabuleiro][TamanhoTabuleiro]byte
	arquivo, err := os.Open("entrada.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
		os.Exit(-1)
	}

	leitor := bufio.NewReader(arquivo)
	i := 0
	for {
		linha, err := leitor.ReadBytes('\n')
		pos := 0
		j := 0
		for pos < len(linha) {
			if linha[pos] != ' ' {
				tabuleiro[i][j] = linha[pos]
				j++
			}
			pos++
		}

		if err == io.EOF {
			break
		}
		i++
	}

	arquivo.Close()
	return tabuleiro
}

//Atirar função que realiza um tiro
func (c *Cliente) Atirar() (int, int) {
	var tiro string
	fmt.Print("Digite seu tiro: ")
	fmt.Scanf("%s", &tiro)

	var i, j int
	i = int(tiro[0] - 'A')
	if len(tiro) > 2 {
		j = 9
	} else {
		j = int(tiro[1] - '1')
	}

	//EnviarTiro()

	return i, j
}

//Ganhou função que indica se o jogador corrente ganhou
func (c *Cliente) Ganhou() bool {
	return c.tabuleiroAtaque.AfundouTodos()
}
