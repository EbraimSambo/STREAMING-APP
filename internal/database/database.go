package database

import (
	"context"
	"log"
	"os"
	"stream/ent"

	_ "github.com/lib/pq"
)

var Client *ent.Client

func Connect() {
	// Captura a URL do banco do ambiente
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("❌ DATABASE_URL não definida nas variáveis de ambiente")
	}

	var err error
	Client, err = ent.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	if err := Client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("Erro ao criar schema: %v", err)
	}
}
