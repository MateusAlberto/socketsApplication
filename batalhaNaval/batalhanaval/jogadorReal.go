package batalhanaval

import (
	"fmt"
)

//JogadorReal struct para definir um jogador real
type JogadorReal struct {
	TabuleiroAtaque *Tabuleiro
	TabuleiroDefesa *Tabuleiro
}

//IniciarJogador inicia um jogador real colocando seus tabuleiros
func (jogador *JogadorReal) IniciarJogador(tabuleiro [][]byte) {
	tabDefesa := tabuleiro
	jogador.TabuleiroDefesa = NovoTabDefesa(tabDefesa)
	jogador.TabuleiroAtaque = NovoTabVazio()
}

//ImprimirTabuleiros Imprime os dois tabuleiros do jogador
func (jogador *JogadorReal) ImprimirTabuleiros() {
	fmt.Print("\nLEGENDA:\n",
		" - : Água\n",
		" N : Parte de um Navio\n",
		" X : Tiro na Água\n",
		" V : Tiro Certeiro\n\n",
		"Seu tabuleiro de Ataque:\n\n")
	jogador.TabuleiroAtaque.Imprimir()
	fmt.Print("\nSeu tabuleiro de Defesa:\n\n")
	jogador.TabuleiroDefesa.Imprimir()
}
