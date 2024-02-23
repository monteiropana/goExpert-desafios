-Ambiente Dev:
Usado Image DockerFile + Docker-compose 
Baixe o repositorio/Entre na pasta desafios/lab: sistema-temperatura-cep
Rode o comando: docker-compose up -d 
Rode o comando docker-compose exec goapp bash
Agora so executar o go run main.go
Abra o navegador e digite a seguinte URL: http://localhost:8080/clima?cep=(CEP)
Teste: Digite um CEP Valido ex: 01153000
Digite um cep com 8 digitos mas que seja um cep invalido
Digite um Cep Inexistente

-Ambiente Prod:

passo 1: Baixe o repositorio 
passo 2: Navegue até o pasta desafios-lab/sistema-temperatura-cep
passo 3: Abra o navegador e digite a seguinte URL: https://go-expert-desafios-lab-m2d6odgxia-uc.a.run.app/clima?cep=(CEP)
Teste: Digite um CEP Valido ex: 01153000
Digite um cep com 8 digitos mas que seja um cep invalido
Digite um Cep Inexistente
