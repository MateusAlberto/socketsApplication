package batalhanaval

import (
	"math/rand"
	"time"
)

//JogadorBot struct para definir um jogador real
type JogadorBot struct {
	TabuleiroAtaque    *Tabuleiro
	TabuleiroDefesa    *Tabuleiro
	UltimoTiroCerteiro [2]int
}

//IniciarJogador construtor de um jogador real
func (jogador *JogadorBot) IniciarJogador() {
	jogador.TabuleiroDefesa = NovoTabVazio()
	jogador.TabuleiroDefesa.GerarTabuleiroAleatorio()
	jogador.TabuleiroAtaque = NovoTabVazio()
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	jogador.UltimoTiroCerteiro[0] = random.Intn(TamanhoTabuleiro)
	seed = rand.NewSource(time.Now().UnixNano())
	random = rand.New(seed)
	jogador.UltimoTiroCerteiro[1] = random.Intn(TamanhoTabuleiro)
}

//Tiro função que realiza um tiro do jogador bot
//tem 20% de chance de realizar um tiro completamente aleatório
//ou então de realizar um tiro próximo do tiro anterior
func (jogador *JogadorBot) Tiro() (int, int) {
	atirou := false
	i := jogador.UltimoTiroCerteiro[0]
	j := jogador.UltimoTiroCerteiro[1]
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	aleatorio := random.Intn(100)
	//tem 20% de chande de gerar uma coordenada aleatória ou um tiro
	if aleatorio < 20 {
		for jogador.TabuleiroAtaque.tabuleiro[i][j] != '-' {
			seed = rand.NewSource(time.Now().UnixNano())
			random = rand.New(seed)
			i = random.Intn(TamanhoTabuleiro)
			seed = rand.NewSource(time.Now().UnixNano())
			random = rand.New(seed)
			j = random.Intn(TamanhoTabuleiro)
		}
		atirou = true
	} else {
		atirou = false
		for distancia := 1; !atirou && distancia < TamanhoTabuleiro; distancia++ {
			switch { //switch vazio equivale a vários if-else aninhados
			//Tentar à esquerda
			case i-distancia >= 0 && jogador.TabuleiroAtaque.tabuleiro[i-distancia][j] == '-':
				i -= distancia
				atirou = true
			//Tentar em cima
			case j+distancia < TamanhoTabuleiro && jogador.TabuleiroAtaque.tabuleiro[i][j+distancia] == '-':
				j += distancia
				atirou = true
			//Tentar à direita
			case i+distancia < TamanhoTabuleiro && jogador.TabuleiroAtaque.tabuleiro[i+distancia][j] == '-':
				i += distancia
				atirou = true
			//Tentar embaixo
			case j-distancia >= 0 && jogador.TabuleiroAtaque.tabuleiro[i][j-distancia] == '-':
				j -= distancia
				atirou = true
			default:
				atirou = false
			}
		}
	}

	if !atirou {
		for jogador.TabuleiroAtaque.tabuleiro[i][j] != '-' {
			seed = rand.NewSource(time.Now().UnixNano())
			random = rand.New(seed)
			i = random.Intn(TamanhoTabuleiro)
			seed = rand.NewSource(time.Now().UnixNano())
			random = rand.New(seed)
			j = random.Intn(TamanhoTabuleiro)
		}
		atirou = true
	}

	return i, j
}
