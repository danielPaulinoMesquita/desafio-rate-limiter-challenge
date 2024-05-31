### Olá galera da Full Cycle, Como vocês estão? Espero que bem.

O Desafio foi feito, darei um resumo de como foi feito e as decisões que tomei. 


### OBS: Foi solicitado o strategy, mas eu não tinha feito, eu tinha criado a instância do redis direto, 
#### mas agora coloquei o padrão strategy. COMEÇA NA LINHA 22 E TERMINA NA LINHA 53

## CONSTRUÇÃO: 
#### 1- O sistema foi contruído com base na descrição e nos requisitos.
#### 2- Variáveis de ambiente 
   * As variáveis e o limite foram definidos no environment do docker-compose,
   * Ou seja, caso queiram mudar o limite do ip basta mudar o valor da chave LIMIT, no momento o valor é 6
   * O tempo de expiração do LIMITE de requisições é 5 minutos, isso já foi definido ao salvar o limite no "REDIS", ou seja, quando o limite 
de requisições for alcançado, só depois de 5 minutos poderá ser feito novas requisições.
   * "Token" pode ser informado pelo header como vocês sugeriram na descrição, considerei que o "token" é a quantidade de requisições, 
portanto você passa o limite daquele "token" que pode ser 5,6 ou 100, o que vocês acharem melhor. Pensei em criptografar o token, 
mas como não ficou claro a descrição com relação a isso, preferi deixar da forma mais simples, e com intuito de tornar as coisas mais fáceis caso vocês queiram mudar o LIMITE que vai no token.


## EXECUÇÃO:
#### 1- Foi criado um docker-compose e um Dockerfile.
#### 2- Para executar o projeto basta rodar comando do docker-compose no terminal, uma que o próprio compose já vai executar do Dockerfile   
``docker-compose up -d``

* Caso acontece algum problema, tente um comando parecido:  
``docker-compose up --build``

#### 3- O sistema roda na porta 8080, você pode acessar o localhost, exemplo:
http://localhost:8080


### Boa revisão para vocês e Boa Sorte para mim! rsrs