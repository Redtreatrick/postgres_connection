package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/url"
	"os"
)

func main() {
	//Создание строки для подключения
	connStr :=
		fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
			"postgres",
			url.QueryEscape("db_user"),
			url.QueryEscape("pwd123"),
			"localhost",
			"54320",
			"db_test",
			5)
	//Контекст с отменой
	ctx, _ := context.WithCancel(context.Background())

	//Сконфигурируем пул, задав для него максимальное количество соединений
	poolConfig, _ := pgxpool.ParseConfig(connStr)
	poolConfig.MaxConns = 20
	poolConfig.MinConns = 20

	// Получаем пул соединений, используя контекст и конфиг
	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connect to database failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connection OK!")

	for i := 0; i < 10; i++ {
		go func(count int) {
			_, err = pool.Exec(ctx, ";") // Exec выполняет запрос к БД, потому-что pgxpool не умеет в Ping()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Ping failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(count, "Query OK!")
			fmt.Printf("Connections - Max: %d, Iddle: %d, Total: %d \n",
				pool.Stat().MaxConns(),
				pool.Stat().IdleConns(),
				pool.Stat().TotalConns())
		}(i)
	}
	select {}
}

// main2 is an example of simple connection which is worse that pool request solution in main
func main2() {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
		"postgres",
		url.QueryEscape("db_user"),
		url.QueryEscape("pwd123"),
		"localhost",
		"54320",
		"db_test",
		5)
	// Задаём контекст для завершения работы
	ctx, _ := context.WithCancel(context.Background())
	// Настраиваем подключение к базе
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connect to database failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connection OK!")
	// выполним простой запрос к базе данных
	err = conn.Ping(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ping failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Query OK!")
	//закрываем подключение
	conn.Close(ctx)
	// или оставляем запрос чтобы программа работала
	//select{}
}
