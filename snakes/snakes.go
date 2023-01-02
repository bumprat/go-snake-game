package snakes

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/eiannone/keyboard"
)

type matter string

const (
	food      matter = "🍲"
	wall      matter = "🧱"
	vacant    matter = "⬛"
	snake     matter = "🟢"
	wreck     matter = "😱"
	snakeHead matter = "🤢"
)

type direction byte

const (
	up    direction = iota
	down  direction = iota
	left  direction = iota
	right direction = iota
)

type cell struct {
	x      byte
	y      byte
	matter matter
}

type stageType struct {
	environment        []cell
	width              byte
	height             byte
	snake              []cell
	snakeDirection     direction
	snakeNextDirection direction
	snakeStepTime      int
}

var stage stageType

func InitDefault() (err error) {
	return Init(10, 10)
}

func Init(width byte, height byte) (err error) {
	if width < 5 || height < 5 {
		err = errors.New("width and height must be greater than 5")
		return
	}
	stage = stageType{
		[]cell{},
		width,
		height,
		[]cell{{1, 1, snakeHead}},
		right,
		right,
		50,
	}
	return
}

func render() {
	var i, j byte
	var result string
	for i = 1; i <= stage.width; i++ {
		result += string(wall)
	}
	result += string(wall+wall) + "\n"
	for j = 1; j <= stage.height; j++ {
		result += string(wall)
	nextCell:
		for i = 1; i <= stage.width; i++ {
			for _, cell := range stage.environment {
				if cell.x == i && cell.y == j {
					result += string(cell.matter)
					continue nextCell
				}
			}
			for _, cell := range stage.snake {
				if cell.x == i && cell.y == j {
					result += string(cell.matter)
					continue nextCell
				}
			}
			result += string(vacant)
		}
		result += string(wall) + "\n"
	}
	for i = 1; i <= stage.width; i++ {
		result += string(wall)
	}
	result += string(wall+wall) + "\n"
	fmt.Print("\033[H")
	fmt.Print(result)
}

var gameover = false

func Start() {
	rand.Seed(time.Now().UnixNano())
	genRandomWall()
	genRandomFood()
	genRandomFood()
	fmt.Print("\033[H\033[2J") // clear screen
	render()
	go keyboardListener()
	ticker := time.NewTicker(time.Microsecond * 200)
	stepTimeCount := 0
	for range ticker.C {
		stepTimeCount++
		if gameover {
			ticker.Stop()
			render()
			fmt.Println("game over")
			break
		}
		if stepTimeCount >= stage.snakeStepTime {
			stepTimeCount = 0
			step()
			render()
		}
	}
}

func step() {
	var moveTo cell = stage.snake[0]
	stage.snakeDirection = stage.snakeNextDirection
	switch stage.snakeDirection {
	case up:
		moveTo.y -= 1
	case down:
		moveTo.y += 1
	case left:
		moveTo.x -= 1
	default:
		moveTo.x += 1
	}

	if moveTo.x <= 0 || moveTo.x > stage.width ||
		moveTo.y <= 0 || moveTo.y > stage.height {
		gameover = true
		stage.snake[0].matter = wreck
		return
	}

	for _, c := range stage.snake {
		if c.x == moveTo.x && c.y == moveTo.y {
			gameover = true
			stage.snake[0].matter = wreck
			return
		}
	}

	for i, c := range stage.environment {
		if c.x == moveTo.x && c.y == moveTo.y {
			if c.matter == food {
				stage.snake[0].matter = snake
				stage.snake = append([]cell{moveTo}, stage.snake...)
				stage.environment = append(stage.environment[:i], stage.environment[i+1:]...)
				genRandomFood()
				render()
				return
			}
			if c.matter == wall {
				gameover = true
				stage.snake[0].matter = wreck
				return
			}
		}
	}

	for i := len(stage.snake) - 1; i >= 0; i-- {
		cell := &stage.snake[i]
		if i == 0 {
			cell.x, cell.y = moveTo.x, moveTo.y
		} else {
			forwardCell := stage.snake[i-1]
			cell.x, cell.y = forwardCell.x, forwardCell.y
			cell.matter = snake
		}
	}
}

func genRandomWall() {
	var i, j byte
	for j = 1; j <= stage.height; j++ {
		for i = 1; i <= stage.width; i++ {
			if rand.Float64() > 0.95 {
				stage.environment = append(stage.environment, cell{i, j, wall})
			}
		}
	}
}

func genRandomFood() {
loopFood:
	for {
		var i, j byte
		i = byte(rand.Intn(int(stage.width)) + 1)
		j = byte(rand.Intn(int(stage.height)) + 1)
		for _, c := range stage.environment {
			if c.x == i && c.y == j {
				continue loopFood
			}
		}
		for _, c := range stage.snake {
			if c.x == i && c.y == j {
				continue loopFood
			}
		}
		stage.environment = append(stage.environment, cell{i, j, food})
		return
	}
}

func keyboardListener() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		switch {
		case key == 0xFFED && stage.snakeDirection != down:
			stage.snakeNextDirection = up
		case key == 0xFFEC && stage.snakeDirection != up:
			stage.snakeNextDirection = down
		case key == 0xFFEB && stage.snakeDirection != right:
			stage.snakeNextDirection = left
		case key == 0xFFEA && stage.snakeDirection != left:
			stage.snakeNextDirection = right
		}
		if key == keyboard.KeyEsc {
			gameover = true
			break
		}
	}
}
