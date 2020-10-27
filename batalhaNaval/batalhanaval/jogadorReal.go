package batalhanaval

import (
	"fmt"
)

//JogadorReal struct para definir um jogador real
type JogadorReal struct {
	tabuleiroAtaque *Tabuleiro
	tabuleiroDefesa *Tabuleiro
}

//Atirar função que realiza um tiro
func (jogador *JogadorReal) Atirar() (int, int) {
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
func (jogador *JogadorReal) Ganhou() bool {
	return jogador.tabuleiroAtaque.AfundouTodos()
}

//IniciarJogador inicia um jogador real colocando seus tabuleiros
func (jogador *JogadorReal) IniciarJogador(tabuleiro [][]byte) {
	tabDefesa := tabuleiro
	jogador.tabuleiroDefesa = NovoTabDefesa(tabDefesa)
	jogador.tabuleiroAtaque = NovoTabAtaque()
}

//ImprimirTabuleiros Imprime os dois tabuleiros do jogador
func (jogador *JogadorReal) ImprimirTabuleiros() {
	fmt.Print("\nLEGENDA:\n",
		" - : Água\n",
		" N : Parte de um Navio\n",
		" X : Tiro na Água\n",
		" V : Tiro Certeiro\n\n",
		"Seu tabuleiro de Ataque:\n\n")
	jogador.tabuleiroAtaque.Imprimir()
	fmt.Print("\nSeu tabuleiro de Defesa:\n\n")
	jogador.tabuleiroDefesa.Imprimir()
}
