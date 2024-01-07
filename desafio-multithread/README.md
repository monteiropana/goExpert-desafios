# goExpert-desafio2-Multithreading
desafio 2 Multithreading

Duas requisições serão feitas simultaneamente para as seguintes APIs:
* https://viacep.com.br/ws/"CEP"/json/
* https://brasilapi.com.br/api/cep/v1/"CEP"

Os requisitos para este desafio são:

- Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.

- O resultado da request deverá ser exibido no command line com os dados do endereço, bem como qual API a enviou.

- Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.