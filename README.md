![Golang](https://user-images.githubusercontent.com/6461792/111548403-6e760700-8759-11eb-8a20-ea1f49d660a4.png)
![Rabbitmq](https://user-images.githubusercontent.com/6461792/152242192-8a9f2ec1-a735-4c70-b629-49885c63c4d0.png) 

# Golang + Rabbitmq

Um simples guide com alguns exemplos de como enviar e consumir mensagens em fila

### Subir os serviços docker
```shell
make docker
```

Nesse guide tem vários exemplos e vou mostrar como usar. Acessar [rabbitmq](http://localhost:15672/)

-- --
#### Enviando e consumindo simples mensagens
```shell
make hello-sender
```

```shell
make hello-receive
```


#### Enviando e consumindo com exchange fanout
```shell
make fanout-sender
```

```shell
make fanout-receive
```


#### Enviando e consumindo com exchange direct
```shell
make direct-sender arg=warning
```

```shell
make direct-receive arg=warning
```