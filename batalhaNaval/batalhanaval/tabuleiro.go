package batalhanaval

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
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

//orientacao enum para as direções
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
	var i, j int
	var direcao orientacao
	colocou := false
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	for !colocou {
		seed = rand.NewSource(time.Now().UnixNano())
		random = rand.New(seed)
		i = random.Intn(TamanhoTabuleiro)
		seed = rand.NewSource(time.Now().UnixNano())
		random = rand.New(seed)
		j = random.Intn(TamanhoTabuleiro)
		//se já tem um navio
		if t.tabuleiro[i][j] != 'N' {

			switch {
			//Ver se consegue colocar a oeste
			case i-tamanhoNavio >= 0:
				direcao = oeste
				//verifica se todo o espaço para o navio está liberado. Se estiver, já coloca o navio
				if t.espacoLiberado(i, j, tamanhoNavio, direcao) {
					colocou = true
				}
			//Ver se consegue colocar ao norte
			case j+tamanhoNavio < TamanhoTabuleiro:
				direcao = norte
				//verifica se todo o espaço para o navio está liberado. Se estiver, já coloca o navio
				if t.espacoLiberado(i, j, tamanhoNavio, direcao) {
					colocou = true
				}
			//Ver se consegue colocar a leste
			case i+tamanhoNavio < TamanhoTabuleiro:
				direcao = leste
				//verifica se todo o espaço para o navio está liberado. Se estiver, já coloca o navio
				if t.espacoLiberado(i, j, tamanhoNavio, direcao) {
					colocou = true
				}
			//Ver se consegue colocar ao sul
			case j-tamanhoNavio >= 0:
				direcao = sul
				//verifica se todo o espaço para o navio está liberado. Se estiver, já coloca o navio
				if t.espacoLiberado(i, j, tamanhoNavio, direcao) {
					colocou = true
				}
			}
		}
	}
}

//verifica se todo o espaço para o navio está liberado. Se estiver, já coloca o navio
//já foi previamente verificado se na direção e coordenada passada, iria sair do tabuleiro pelo tamanho do navio
func (t *Tabuleiro) espacoLiberado(i, j, tamanhoNavio int, direcao orientacao) bool {
	switch direcao {
	case oeste:
		for c := i; c > i-tamanhoNavio; c-- {
			if t.tabuleiro[c][j] == 'N' {
				return false
			}
		}
	case norte:
		for c := j; c < j+tamanhoNavio; c++ {
			if t.tabuleiro[i][c] == 'N' {
				return false
			}
		}
	case leste:
		for c := i; c < i+tamanhoNavio; c++ {
			if t.tabuleiro[c][j] == 'N' {
				return false
			}
		}
	case sul:
		for c := j; c > j-tamanhoNavio; c-- {
			if t.tabuleiro[i][c] == 'N' {
				return false
			}
		}
	}

	//Se chegou até aqui, é porque pode colocar
	//Vai colocar e retornar verdadeiro
	switch direcao {
	case oeste:
		for c := i; c > i-tamanhoNavio; c-- {
			t.tabuleiro[c][j] = 'N'
		}
	case norte:
		for c := j; c < j+tamanhoNavio; c++ {
			t.tabuleiro[i][c] = 'N'
		}
	case leste:
		for c := i; c < i+tamanhoNavio; c++ {
			t.tabuleiro[c][j] = 'N'
		}
	case sul:
		for c := j; c > j-tamanhoNavio; c-- {
			t.tabuleiro[i][c] = 'N'
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

//ReceberTiro função que que recebe um tiro e verifica se acertou
func (t *Tabuleiro) ReceberTiro(i, j int) bool {
	acertou := t.tabuleiro[i][j] == 'N' || t.tabuleiro[i][j] == 'V' //Assim, o jogador vai poder repetir o tiro que deu (isso fica a critério dele)
	t.RegistrarTiro(acertou, i, j)
	return acertou
}

//RegistrarTiro função que registra um tiro realizado já com o resultado se acertou ou não
func (t *Tabuleiro) RegistrarTiro(acertou bool, i, j int) {
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

//ParseTiro vai receber um tiro em string do tipo A5 e retornar um par de ints do tipo equivalente (0, 4)
func ParseTiro(tiro string) (int, int) {
	var i, j int
	iStr := byte(tiro[0])
	jStr := tiro[1:len(tiro)]
	//Parse do i
	if iStr >= 'A' && iStr <= 'J' {
		i = int(iStr - 'A')
	} else {
		i = -1
	}

	var err error
	//Parse do j
	j, err = strconv.Atoi(jStr)
	if err == nil {
		if j >= 1 && j <= 10 {
			j--
		}
	} else {
		j = -1
	}
	return i, j
}

//PosicaoDesconhecida retorna verdadeiro se a posição representar água e falso caso contrário
func (t *Tabuleiro) PosicaoDesconhecida(i, j int) bool {
	return t.tabuleiro[i][j] == '-'
}
