package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Pega configurações do .env ou usa padrão
	rabbitMQURI := os.Getenv("RABBITMQ_URI")
	apiURL := os.Getenv("API_URL")
	queueName := "weather_queue"

	if rabbitMQURI == "" {
		rabbitMQURI = "amqp://user:password@rabbitmq:5672"
	}
	if apiURL == "" {
		apiURL = "http://backend:3000/api/weather/logs"
	}

	log.Println("🐰 Worker Go Iniciando...")
	
	// Loop de tentativa de conexão (caso o RabbitMQ demore a subir)
	var conn *amqp.Connection
	var err error
	for i := 0; i < 15; i++ {
		conn, err = amqp.Dial(rabbitMQURI)
		if err == nil {
			break
		}
		log.Printf("⏳ Tentando conectar ao RabbitMQ... (%d/15)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("❌ Falha fatal ao conectar no RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Falha ao abrir canal: %v", err)
	}
	defer ch.Close()

	// Declara a fila para garantir que ela existe
	q, err := ch.QueueDeclare(
		queueName, // nome
		true,      // durável
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("❌ Falha ao declarar fila: %v", err)
	}

	// Define como consumir as mensagens
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack (FALSE = só confirma se der certo)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("❌ Falha ao registrar consumidor: %v", err)
	}

	log.Println("✅ Worker Go conectado e aguardando mensagens!")

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("📥 Recebido da Fila: %s", d.Body)

			// Enviar para a API NestJS
			resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(d.Body))
			
			if err != nil {
				log.Printf("❌ Erro ao conectar na API: %v", err)
				// Nack: Devolve para a fila para tentar de novo
				d.Nack(false, true) 
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				log.Println("💾 Sucesso: Dados salvos no Backend.")
				d.Ack(false) // Ack: Remove da fila, deu tudo certo
			} else {
				log.Printf("⚠️ Erro API (Status %d)", resp.StatusCode)
				d.Nack(false, true) // Devolve para fila
			}
		}
	}()

	<-forever
}