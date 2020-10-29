package batalhanaval

import (
	"math/rand"
)

//JogadorBot struct para definir um jogador real
type JogadorBot struct {
	tabuleiroAtaque *Tabuleiro
	tabuleiroDefesa *Tabuleiro
	ultimoTiro      [2]int
}

//IniciarJogador construtor de um jogador real
func (jogador *JogadorBot) IniciarJogador() {
	jogador.tabuleiroDefesa = NovoTabVazio()
	jogador.tabuleiroDefesa.GerarTabuleiroAleatorio()
	jogador.tabuleiroAtaque = NovoTabVazio()
	jogador.ultimoTiro[0] = rand.Int() % TamanhoTabuleiro
	jogador.ultimoTiro[1] = rand.Int() % TamanhoTabuleiro
}

//Atirar função que realiza um tiro
func (jogador *JogadorBot) Atirar() (int, int) {
	i := jogador.ultimoTiro[0]
	j := jogador.ultimoTiro[1]

	//parei aqui

	for jogador.tabuleiroAtaque.tabuleiro[i][j] != '-' {
		jogador.ultimoTiro[0] = rand.Int() % TamanhoTabuleiro
		jogador.ultimoTiro[1] = rand.Int() % TamanhoTabuleiro
	}

	//EnviarTiro()

	return i, j
}

//Ganhou função que indica se o jogador corrente ganhou
func (jogador *JogadorBot) Ganhou() bool {
	return jogador.tabuleiroAtaque.AfundouTodos()
}
