package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	difficultyStr := os.Getenv("DIFFICULTY")
	difficulty, err = strconv.Atoi(difficultyStr)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	quotes = strings.Split(os.Getenv("QUOTES"), ";")
	if len(quotes) == 0 {
		log.Fatal().Msg("no quotes provided")
	}
}

var (
	quotes     []string
	difficulty int
)

func main() {
	var port = os.Getenv("PORT")

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msgf("server is listening on port: %s", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error().Msgf("error accepting connection: %v", err)
			continue
		}

		go processConnection(conn)
	}
}

func processConnection(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// send challenge to client
	challenge := generateChallenge()
	message := fmt.Sprintf("%s:%d\n", challenge, difficulty)
	_, err := writer.WriteString(message)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	err = writer.Flush()
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	// receive client resp
	nonce, err := reader.ReadString('\n')
	if err != nil {
		log.Error().Msgf("error reading nonce: %v", err)
		return
	}
	nonce = strings.TrimSpace(nonce)

	// verify and response
	if verifyPoW(challenge, nonce, difficulty) {
		quote := quotes[rand.Intn(len(quotes))]
		_, err = writer.WriteString(quote + "\n")
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}
	} else {
		_, err = writer.WriteString("invalid Proof of Work\n")
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}
	}
	err = writer.Flush()
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
}

func generateChallenge() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%x", rand.Int63())
}

func verifyPoW(challenge, nonce string, difficulty int) bool {
	data := challenge + nonce
	hash := sha256.Sum256([]byte(data))
	hashStr := hex.EncodeToString(hash[:])
	return strings.HasPrefix(hashStr, strings.Repeat("0", difficulty))
}
