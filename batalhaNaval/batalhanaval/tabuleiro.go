package batalhanaval

import (
	"fmt"
	"math/rand"
)

//Tamanhos dos navios
const tamPortaAvioes = 5
const tamNaviosTanque = 4
const tamContraTorpedeiros = 3
const tamSubmarinos = 2

//Quantidades dos navios
const quantPortaAvioes = 1
const quantNaviosTanque = 2
const quantContraTorpedeiros = 3
const quantSubmarinos = 4

//enum para as direções
type orientacao int

const (
	oeste orientacao = iota
	norte
	leste
	sul
)

//TamanhoTabuleiro tamanho do tabuleiro
const TamanhoTabuleiro = 10

//Tabuleiro struct para definir as operações de um tabuleiro
type Tabuleiro struct {
	tabuleiro [][]byte
}

//NovoTabDefesa construtor de um tabuleiro de defesa
func NovoTabDefesa(tabuleiro [][]byte) *Tabuleiro {
	tab := &Tabuleiro{}
	tab.tabuleiro = tabuleiro
	return tab
}

//NovoTabVazio construtor de um tabuleiro de ataque
func NovoTabVazio() *Tabuleiro {
	tab := &Tabuleiro{}
	tab.tabuleiro = make([][]byte, TamanhoTabuleiro)
	for i := 0; i < TamanhoTabuleiro; i++ {
		tab.tabuleiro[i] = make([]byte, TamanhoTabuleiro)
		for j := 0; j < TamanhoTabuleiro; j++ {
			tab.tabuleiro[i][j] = '-'
		}
	}
	return tab
}

//GerarTabuleiroAleatorio gera um tabuleiro aleatório para um jogador bot
func (t *Tabuleiro) GerarTabuleiroAleatorio() {
	quantNavios := []int{quantPortaAvioes, quantNaviosTanque, quantContraTorpedeiros, quantSubmarinos}
	tamNavios := []int{tamPortaAvioes, tamNaviosTanque, tamContraTorpedeiros, tamSubmarinos}
	for i, quant := range quantNavios {
		for j := 0; j < quant; j++ {
			t.colocarNavio(tamNavios[i])
		}
	}
}

//Função para colocar um navio numa posição aleatória do tabuleiro
func (t *Tabuleiro) colocarNavio(tamanhoNavio int) {
	var x, y int
	var direcao orientacao
	colocou := false

	for !colocou {
		x = rand.Int() % TamanhoTabuleiro
		y = rand.Int() % TamanhoTabuleiro
		//se já tem um navio
		if t.tabuleiro[x][y] != 'N' {

			switch {
			//Ver se consegue colocar a oeste
			case x-tamanhoNavio >= 0:
				direcao = oeste
				//verifica se todo o espaço para o navio está liberado. Se estiver, já coloca o navio
				if t.espacoLiberado(x, y, tamanhoNavio, direcao) {
					colocou = true
				}
			//Ver se consegue colocar ao norte
			case y+tamanhoNavio < TamanhoTabuleiro:
				direcao = norte
				//verifica se todo o espaço para o navio está liberado. Se estiver, já coloca o navio
				if t.espacoLiberado(x, y, tamanhoNavio, direcao) {
					colocou = true
				}
			//Ver se consegue colocar a leste
			case x+tamanhoNavio < TamanhoTabuleiro:
				direcao = leste
				//verifica se todo o espaço para o navio está liberado. Se estiver, já coloca o navio
				if t.espacoLiberado(x, y, tamanhoNavio, direcao) {
					colocou = true
				}
			//Ver se consegue colocar ao sul
			case y-tamanhoNavio >= 0:
				direcao = sul
				//verifica se todo o espaço para o navio está liberado. Se estiver, já coloca o navio
				if t.espacoLiberado(x, y, tamanhoNavio, direcao) {
					colocou = true
				}
			}
		}
	}
}

//verifica se todo o espaço para o navio está liberado. Se estiver, já coloca o navio
//já foi previamente verificado se na direção e coordenada passada, iria sair do tabuleiro pelo tamanho do navio
func (t *Tabuleiro) espacoLiberado(x, y, tamanhoNavio int, direcao orientacao) bool {
	switch direcao {
	case oeste:
		for i := x; i > x-tamanhoNavio; i-- {
			if t.tabuleiro[i][y] == 'N' {
				return false
			}
		}
	case norte:
		for j := y; j < y+tamanhoNavio; j++ {
			if t.tabuleiro[x][j] == 'N' {
				return false
			}
		}
	case leste:
		for i := x; i < x+tamanhoNavio; i++ {
			if t.tabuleiro[i][y] == 'N' {
				return false
			}
		}
	case sul:
		for j := y; j > y-tamanhoNavio; j-- {
			if t.tabuleiro[x][j] == 'N' {
				return false
			}
		}
	}

	//Se chegou até aqui, é porque pode colocar
	//Vai colocar e retornar verdadeiro
	switch direcao {
	case oeste:
		for i := x; i > x-tamanhoNavio; i-- {
			t.tabuleiro[i][y] = 'N'
		}
	case norte:
		for j := y; j < y+tamanhoNavio; j++ {
			t.tabuleiro[x][j] = 'N'
		}
	case leste:
		for i := x; i < x+tamanhoNavio; i++ {
			t.tabuleiro[i][y] = 'N'
		}
	case sul:
		for j := y; j > y-tamanhoNavio; j-- {
			t.tabuleiro[x][j] = 'N'
		}
	}
	return true
}

//Imprimir função que imprime o tabuleiro corrente
func (t *Tabuleiro) Imprimir() {
	//Primeira linha -> índice
	fmt.Print("  ")
	for i := 1; i <= TamanhoTabuleiro; i++ {
		fmt.Print(" ", i)
	}
	fmt.Println()
	//Linhas seguintes
	for i := 0; i < TamanhoTabuleiro; i++ {
		fmt.Printf(" %c", ('A' + i))
		for j := 0; j < TamanhoTabuleiro; j++ {
			fmt.Printf(" %c", t.tabuleiro[i][j])
		}
		fmt.Println()
	}
}

//ReceberTiro função que do tabuleiro de uleiro
//que recebe um tiro e verifica se acertou e se
//afundou um navio
func (t *Tabuleiro) ReceberTiro(i, j int) bool {
	acertou := t.tabuleiro[i][j] == 'N'
	t.RegistrarTiro(i, j, acertou)
	return acertou
}

//RegistrarTiro função que registra um tiro realizado já com o resultado se acertou ou não
func (t *Tabuleiro) RegistrarTiro(i, j int, acertou bool) {
	if acertou {
		t.tabuleiro[i][j] = 'V'
	} else {
		t.tabuleiro[i][j] = 'X'
	}
}

//AfundouTodos função que verifica se afundou todos os navios
//Obs.: esta função só é implementada deste jeito porque
//o tamanho do tabuleiro é fixo e pequeno
func (t *Tabuleiro) AfundouTodos() bool {
	for i := 0; i < TamanhoTabuleiro; i++ {
		for j := 0; j < TamanhoTabuleiro; j++ {
			if t.tabuleiro[i][j] == 'N' {
				return false
			}
		}
	}
	return true
}
