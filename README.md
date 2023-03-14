## Repositório feito em golang para replicar mensagens de tópicos de um kafka para o outro

### Configurável por meio de variáveis de ambiente. Caso preciso incluir configuração de certificado ssl, apenas incluir as respectivas envs:

### Variáveis para o kafka-replicator mandar as mensagens para o kafka replica
| Variáveis  Producer TLS        | Descrição                                                                                   |
|------------------|-----------------------------------------------------------------------------------------------|
| KAFKA_PRODUCER_TLS_ENABLED | Caso a conexão com o kafka producer precise de certifcado ssl, colocar como true |
| KAFKA_PRODUCER_SSL_CERTIFICATE_LOCATION | Caminho onde ficará o certificado público.Ex:/etc/ssl/ca.crt |
| KAFKA_PRODUCER_SSL_PRIVATE_KEY_LOCATION | Caminho onde ficará a chave privada.Ex:/etc/ssl/private.key |
| KAFKA_PRODUCER_SSL_PRIVATE_KEY_PASSWORD | A senha da chave privada.Ex:topSecretRandomPassword |

### Caso não precise de conexão ssl, apenas incluir as envs para o producer de tópico e brokers
| Variáveis  Producer        | Descrição                                                                                   |
|------------------|-----------------------------------------------------------------------------------------------|
| KAFKA_PRODUCER_BOOSTRAP_SERVERS | Brokers do kafka que receberá as mensagens replicadas pelo kafka-replicator |
| KAFKA_PRODUCER_TOPIC | Nome do tópico replica |

### Número de partições,réplicas que o tópico usará
| Variáveis  Producer Repluicação       | Descrição                                                                                   |
|------------------|-----------------------------------------------------------------------------------------------|
| KAFKA_PRODUCER_REPLICATOR_PARTITION_NUMBER | Número de replicas para o tópico.(Valor default 3) |
| KAFKA_PRODUCER_REPLICATOR_FACTOR | Fator de réplica do tópico.(Valor default 3) |
| KAFKA_PRODUCER_TIMEOUT | Valor de timeout do producer na criação do tópico.(Valor default 60s) |

### Caso essa env não tenha valor true, as mensagens não serão mandadas pro kafka replica, apenas será mostrado nos logs do container.
| Variáveis  Producer Replicação       | Descrição                                                                                   |
|------------------|-----------------------------------------------------------------------------------------------|
| KAFKA_REPLICATION_ENABLED | Ativar a replicação com valor true,caso não seja declarada ou não tenha varlo true, as mensagens não serão replicadas para o kafka replica |

### Variáveis para o que o kafka-replicator leia o tópico e puxe as mensagens do kafka original
| Variáveis  Consumer       | Descrição                                                                                   |
|------------------|-----------------------------------------------------------------------------------------------|
| KAFKA_CONSUMER_TLS_ENABLED | Caso a conexão com o kafka consumer precise de certifcado ssl, colocar como true |
| KAFKA_CONSUMER_SSL_CERTIFICATE_LOCATION | Caminho onde ficará o certificado público.Ex:/etc/ssl/ca.crt |
| KAFKA_CONSUMER_SSL_PRIVATE_KEY_LOCATION | Caminho onde ficará a chave privada.Ex:/etc/ssl/private.key |
| KAFKA_CONSUMER_SSL_PRIVATE_KEY_PASSWORD | A senha da chave privada.Ex:topSecretRandomPassword  |

### Caso não precise de conexão ssl, apenas incluir as envs para o kafka-replicator
| Variáveis  Consumer       | Descrição                                                                                   |
|------------------|-----------------------------------------------------------------------------------------------|
| KAFKA_CONSUMER_BOOSTRAP_SERVERS | Brokers do kafka que mandará as mensagens para kafka-replicator |
| KAFKA_CONSUMER_TOPIC | Nome do tópico alvo que será replicado |
| KAFKA_CONSUMER_GROUP_ID | GroupID do consumer do tópico |

