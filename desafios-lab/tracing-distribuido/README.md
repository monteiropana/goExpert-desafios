--Usado Image DockerFile + Docker-compose :

Baixe o repositorio/Entre na pasta desafios/lab: sistema-temperatura-cep
Para iniciar todos os serviços rode o comando: docker-compose up --build
Digite um CEP Valido ex: 01153000

Servico A: Passando o cep no body, a resposta é city + informacoes de clima
 - Metodo POST: http://localhost:8080/consulta
body: {
    "cep": "011530001"
    }
Servico B: Passando o cep no parametro a resposta é as informacoes do clima de dado cep
- Metodo GET: http://localhost:8081/clima?cep=011530001
Digite um CEP Valido ex: 01153000

--Resultados de Tracing:

Acesse o Zipkin em seu navegador: http://localhost:9411/zipkin/.
Na interface do Zipkin, você pode buscar por traces, visualizar detalhes dos spans, tempos de resposta e a estrutura do trace.
Os traces irão mostrar a requisição desde o Serviço A e Servico B

totaloperation - mostra o tempo total da chamada
Pode usar o filtro ServiceName= servico-a e ServiceName= servico-b para trazer os detalhes de ambos servicos.
