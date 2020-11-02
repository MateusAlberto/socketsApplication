# socketsApplication
Trabalho da disciplina de Redes I do curso Ciência da Computação na PUC Minas. Este trabalho envolve a criação de dois servidores para o uso de sockets com os protocolos de transporte TCP e UDP, feito na linguagem Go

Os programas recebem argumentos para executarem. O servidor irá receber a porta onde ele estará ouvindo a comunicação. Já o cliente deve receber o endereço IPv4 e a porta em que o servidor estará ouvindo.

Para executar o servidor:

Primeiro, compile usando o seguinte comando:

```shell
go build servidor.go
```
Depois, execute da seguinte forma(supondo que o executável se chame servidor.exe):

```shell
./servidor.exe [PORTA]
//Por exemplo:
./servidor.exe 11111
```
Assim, o servidor estará ouvindo na porta 11111.

Para executar o cliente:

Primeiro, compile usando o seguinte comando:

```shell
go build cliente.go
```
Depois, execute da seguinte forma(supondo que o executável se chame cliente.exe):
```shell
./cliente.exe [ENDEREÇO_IP:PORTA]
//Por exemplo:
./servidor.exe 127.0.0.1:11111
```
Assim, o cliente tentará se conectar com o servidor com endereço IP 127.0.0.1 na porta 11111.
