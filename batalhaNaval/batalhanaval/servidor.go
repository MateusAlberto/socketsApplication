package batalhanaval

import "math/rand"

//Tamanhos dos navios
const tamanhoPortaAvioes = 5
const tamanhoNaviosTanque = 4
const tamanhoContraTorpedeiros = 3
const tamanhoSubmarinos = 2

//Quantidades dos navios
const quantidadePortaAvioes = 1
const quantidadeNaviosTanque = 2
const quantidadeContraTorpedeiros = 3
const quantidadeSubmarinos = 4

//Servidor struct para definir um jogador real
type Servidor struct {
	tabuleiroAtaque *Tabuleiro
	tabuleiroDefesa *Tabuleiro
	ultimoTiro      [2]int
}

//NovoJogador construtor de um jogador real
func (s *Servidor) NovoJogador() {
	tabDefesa := gerarTabuleiroAleatorio()
	s.tabuleiroAtaque = NovoTabDefesa(tabDefesa)
	s.tabuleiroAtaque = NovoTabAtaque()
	s.ultimoTiro[0] = rand.Int() % TamanhoTabuleiro
	s.ultimoTiro[1] = rand.Int() % TamanhoTabuleiro
}

//gerarTabuleiroAleatorio
func gerarTabuleiroAleatorio() [TamanhoTabuleiro][TamanhoTabuleiro]byte {
	var tab [10][10]byte
	return tab
}

//Atirar função que realiza um tiro
func (s *Servidor) Atirar() (int, int) {
	i := s.ultimoTiro[0]
	j := s.ultimoTiro[1]

	//parei aqui

	for s.tabuleiroAtaque.tabuleiro[i][j] != '-' {
		s.ultimoTiro[0] = rand.Int() % TamanhoTabuleiro
		s.ultimoTiro[1] = rand.Int() % TamanhoTabuleiro
	}

	//EnviarTiro()

	return i, j
}

//Ganhou função que indica se o jogador corrente ganhou
func (s *Servidor) Ganhou() bool {
	return s.tabuleiroAtaque.AfundouTodos()
}
