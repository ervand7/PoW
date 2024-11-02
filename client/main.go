package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

const aliveSignalTries = 1_000_000

func solvePoW(challenge string, difficulty int) (string, error) {
	var nonce int64
	target := strings.Repeat("0", difficulty)
	for {
		data := challenge + strconv.FormatInt(nonce, 10)
		hash := sha256.Sum256([]byte(data))
		hashStr := hex.EncodeToString(hash[:])
		if strings.HasPrefix(hashStr, target) {
			return strconv.FormatInt(nonce, 10), nil
		}
		nonce++
		if nonce%aliveSignalTries == 0 {
			log.Info().Msg("still searching for nonce...")
		}
	}
}

func main() {
	serverAddress := os.Getenv("SERVER_ADDRESS")
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// get data from server
	message, err := reader.ReadString('\n')
	if err != nil {
		log.Error().Msgf("error reading challenge: %v", err)
		return
	}
	message = strings.TrimSpace(message)
	parts := strings.Split(message, ":")
	if len(parts) != 2 {
		log.Error().Msg("invalid challenge format")
		return
	}
	challenge := parts[0]
	difficulty, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Error().Msg("invalid difficulty format")
		return
	}
	log.Info().Msgf("received challenge: %s", challenge)
	log.Info().Msgf("difficulty level: %d", difficulty)

	// make PoW
	nonce, err := solvePoW(challenge, difficulty)
	if err != nil {
		log.Error().Msgf("Error solving PoW: %v", err)
		return
	}
	_, err = writer.WriteString(nonce + "\n")
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	err = writer.Flush()
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	// get final response from server
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Error().Msgf("error reading response: %v", err)
		return
	}
	log.Info().Msgf("server response: %v", response)
}
