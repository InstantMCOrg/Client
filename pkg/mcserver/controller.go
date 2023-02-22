package mcserver

import (
	"fmt"
	"github.com/instantmc/client/pkg/constants"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var isGeneratingWorld = false
var worldGeneratingStatus = 0
var WorldGenerationChan = make(chan int)
var ServerLogsChan = make(chan string)
var serverIsUpAndRunning = false
var stdin io.WriteCloser
var cmd *exec.Cmd
var serverStartupLock sync.WaitGroup

var targetRamSize = constants.MinimumRamMb

const (
	PORT = 25565
)

func StartServer(ram int) {
	log.Printf("Starting Minecraft Server with %dmb of ram...\n", ram)
	maximumRamStr := fmt.Sprintf("-Xmx%dm", ram)
	minAllocRamStr := fmt.Sprintf("-Xms%dm", constants.StartupRamAllocation)
	targetRamSize = ram
	cmd = exec.Command("java", maximumRamStr, minAllocRamStr, "-jar", "server.jar")
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal("Couldn't redirect output", err)
	}

	stdin, err = cmd.StdinPipe()

	if err != nil {
		log.Fatal("Couldn't redirect input", err)
	}

	errPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal("Couldn't connect to the error pipe:", err)
	}

	serverStartupLock.Add(1)
	err = cmd.Start()

	if err != nil {
		log.Fatal("Couldn't start Minecraft Server", err)
	}
	log.Println("Started Minecraft Server successfully. Waiting...")

	go readErrorPipe(errPipe)
	go readServerOutput(stdout)
	go crashCatcher()
}

func crashCatcher() {
	err := cmd.Wait()
	log.Println("Minecraft Server unexpectedly exited:", err, "Restarting...")
	StartServer(targetRamSize)
}

func SendCommand(command string) {
	if _, err := io.WriteString(stdin, command+"\n\r"); err != nil {
		log.Println("Couldn't write to Minecraft command pipe:", err)
	}
}

func WaitForStop() {
	cmd.Wait()
	serverIsUpAndRunning = false
	log.Println("Minecraft Server stopped")
}

func WaitUntilServerIsReady() {
	serverStartupLock.Wait()
}

func readErrorPipe(pipe io.ReadCloser) {
	for {
		tmp := make([]byte, 1024)
		_, err := pipe.Read(tmp)
		if err != nil {
			fmt.Println("Error logs pipe closed unexpectedly:", err)
			break
		}
		log.Println(string(tmp))
	}
}

func readServerOutput(pipe io.ReadCloser) {
	for {
		tmp := make([]byte, 1024)
		_, err := pipe.Read(tmp)
		if err != nil {
			fmt.Println("Minecraft logs pipe closed unexpectedly:", err)
			break
		}
		parseMinecraftLog(tmp)
	}
}

func parseMinecraftLog(output []byte) {
	minecraftLog := string(output)
	// this select is need bc we can't guarantee that there is a receiver for the channel
	select {
	case ServerLogsChan <- minecraftLog:
		// A websocket is currently listening on the logs...
		break
	default:
		// No websocket is currently listening
		break
	}

	if strings.Contains(minecraftLog, "Preparing spawn area: ") {
		spawnAreaPreparingPercent, err := strconv.Atoi(strings.Split(strings.Split(minecraftLog, "area: ")[1], "%")[0])
		log.Printf("Spawning area is currently being prepared at %d%s\n", spawnAreaPreparingPercent, "%")
		isGeneratingWorld = true
		if err != nil {
			worldGeneratingStatus = spawnAreaPreparingPercent
		}
		select {
		case WorldGenerationChan <- spawnAreaPreparingPercent:
			// A websocket is currently listening...
			break
		default:
			// No websocket is currently listening
			break
		}
	} else if strings.Contains(minecraftLog, "Done") {
		isGeneratingWorld = false
		serverIsUpAndRunning = true
		log.Println("Minecraft Server is up and running")
		select {
		case WorldGenerationChan <- 100:
			// A websocket is currently listening...
			break
		default:
			// No websocket is currently listening
			break
		}
		serverStartupLock.Done()
	}
}
