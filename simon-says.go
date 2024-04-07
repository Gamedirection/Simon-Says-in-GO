package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	playerName   string
	reader       = bufio.NewReader(os.Stdin)
	simonActions []string // This will be populated from the file
)

func loadActionsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var actions []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		action := scanner.Text()
		if action != "" { // Ignore empty lines
			actions = append(actions, action)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return actions, nil
}

func promptInput(prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func initializeSimonNames(playerName string) []string {
	// Initialize with "Simon Says" and two placeholders for random name (including "nothing") and player's name
	return []string{"Simon Says", "", "", playerName} // Adjusted for probability distribution
}

func executeSimonCommand(simonNames []string, playerName string) (bool, int) {
	rand.Seed(time.Now().UnixNano())
	// Generate a random index for simonNames to simulate the probability distribution
	index := rand.Intn(100)
	action := simonActions[rand.Intn(len(simonActions))]
	commandPrefix := ""

	if index < 50 { // 50% chance for "Simon Says"
		commandPrefix = "Simon Says"
	} else if index < 75 { // 25% chance for a random name or "nothing"
		// Exclude the last element (player's name) for random name selection
		commandPrefix = simonNames[rand.Intn(len(simonNames)-1)]
	}

	command := action
	if commandPrefix != "" {
		command = fmt.Sprintf("%s, %s", commandPrefix, action)
	}
	fmt.Println(command)

	userInput := promptInput("> ")
	if commandPrefix == "Simon Says" {
		if strings.ToLower(userInput) == strings.ToLower(action) {
			return false, 1 // Correct response
		}
		// Incorrect response to a Simon Says command
		fmt.Printf("Simon didn't say \"%s\". You got ", userInput)
		return true, 0
	}

	// If command didn't start with "Simon Says" and player responded
	if userInput != "" {
		fmt.Printf("Simon didn't ask you to \"%s\". You got ", userInput)
		return true, 0
	}

	return false, 0 // Correctly did nothing
}

func welcome() string {
	fmt.Println("Welcome to Simon Says. Let's start. Simon says, 'What is your name?'")
	playerName := promptInput("> ")
	fmt.Printf("Nice to meet you, %s. Remember, only do as Simon Says.\n", playerName)
	return playerName
}

func main() {
	var err error
	simonActions, err = loadActionsFromFile("simonsActions.txt")
	if err != nil {
		fmt.Println("Error loading actions:", err)
		return
	}

	playerName := welcome()
	simonNames := initializeSimonNames(playerName)
	score := 0
	gameStartTime := time.Now() // Capture the start time of the game

	for {
		roundStartTime := time.Now() // Capture the start time of the round
		gameOver, roundScore := executeSimonCommand(simonNames, playerName)
		score += roundScore
		if gameOver {
			roundDuration := time.Since(roundStartTime) // Calculate the round duration
			fmt.Printf("%d points! Simon Says, Type \"Yes\" if you would like to play again. You played this round for %d minutes and %d seconds.\n",
				score, roundDuration/time.Minute, roundDuration%time.Minute/time.Second)
			if strings.ToLower(promptInput("> ")) != "yes" {
				totalDuration := time.Since(gameStartTime) // Calculate total game duration
				fmt.Printf("Thanks for playing, goodbye! You played for a total of %d minutes and %d seconds.\n",
					totalDuration/time.Minute, totalDuration%time.Minute/time.Second)
				break
			}
			score = 0                  // Reset score for a new game
			gameStartTime = time.Now() // Reset the game start time for the new game
		}
	}
}
