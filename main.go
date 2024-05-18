package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	redisClient       *redis.Client
	ipRequestCountMap = make(map[string]int)
	mutex             sync.Mutex
)

func main() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	//log.Printf("REDIS_HOST: %s", redisHost)
	//log.Printf("REDIS_PORT: %s", redisPort)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})

	router := mux.NewRouter()

	router.Use(rateLimiterMiddleware)

	router.HandleFunc("/", handleRequest).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Fazendo Requisição!"))
}

func rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Obtém o endereço IP do cliente
		ip := strings.Split(r.RemoteAddr, ":")[0]

		token := r.Header.Get("API_KEY")
		if token == "" {
			token = "0"
		}

		// Verifica se há um limite configurado para o endereço IP
		ipLimit, err := getRateLimitFromRedis(ip)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		println("IP LIMIT: ", ipLimit)

		// Verifica se o ipLimit já está salvo
		if ipLimit == 0 {

			// Zerando o map(cache das requisições anteriormente feitas)
			restartRequestCount(ip)

			// Recupera o limite do arquivo .env
			limitFromEnv := os.Getenv("LIMIT")

			var limitNumber int

			// Verifica se o limit por ip é maior que o do token
			if limitFromEnv > token {
				limitNumber, err = strconv.Atoi(limitFromEnv)
				if err != nil {
					fmt.Println("Error FIRST:", err)
					return
				}
			} else {
				limitNumber, err = strconv.Atoi(token)
				if err != nil {
					fmt.Println("Error SECOND:", err)
					return
				}
			}

			saveRateLimit(ip, limitNumber)
			incrementRequestCount(ip)
			return
		}

		// Incrementa o contador de requisições para o IP
		incrementRequestCount(ip)

		if exceededLimit(ip, ipLimit) {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		println(ipRequestCountMap[ip])

		next.ServeHTTP(w, r)
	})
}

func exceededLimit(ip string, limit int) bool {
	mutex.Lock()
	defer mutex.Unlock()

	// Verifica se o IP excedeu o limite
	count := ipRequestCountMap[ip]
	return count > limit
}

func incrementRequestCount(ip string) {
	mutex.Lock()
	defer mutex.Unlock()

	// Incrementa o contador de requisições para o IP
	ipRequestCountMap[ip]++
}

func restartRequestCount(ip string) {
	mutex.Lock()
	defer mutex.Unlock()

	// Reseta o contador
	ipRequestCountMap[ip] = 0
}

// Esta função vai ser responsável por salvar as informações do limiter no REdis
func saveRateLimit(ipKey string, limit int) {
	err := redisClient.Set(ipKey, limit, time.Minute*5).Err()
	if err != nil {
		log.Println("Erro ao atualizar o contador no Redis:", err)
	}
}

func getRateLimitFromRedis(key string) (int, error) {
	limit, err := redisClient.Get(key).Int()
	if err != nil && err != redis.Nil {
		return 0, err
	}
	return limit, nil
}
