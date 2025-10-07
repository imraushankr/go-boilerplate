package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Todo represents a single todo item
type Todo struct {
	ID        int
	Task      string
	Completed bool
}

// TodoService defines the interface for todo operations
type TodoService interface {
	Add(task string) *Todo
	List() []Todo
	Complete(id int) error
	Delete(id int) error
	Get(id int) (*Todo, error)
}

// InMemoryTodoService implements TodoService with in-memory storage
type InMemoryTodoService struct {
	todos  []Todo
	nextID int
}

func NewInMemoryTodoService() *InMemoryTodoService {
	return &InMemoryTodoService{
		todos:  make([]Todo, 0),
		nextID: 1,
	}
}

func (s *InMemoryTodoService) Add(task string) *Todo {
	todo := Todo{
		ID:        s.nextID,
		Task:      task,
		Completed: false,
	}
	s.todos = append(s.todos, todo)
	s.nextID++
	return &todo
}

func (s *InMemoryTodoService) List() []Todo {
	return s.todos
}

func (s *InMemoryTodoService) Complete(id int) error {
	for i := range s.todos {
		if s.todos[i].ID == id {
			s.todos[i].Completed = true
			return nil
		}
	}
	return fmt.Errorf("todo with ID %d not found", id)
}

func (s *InMemoryTodoService) Delete(id int) error {
	for i, todo := range s.todos {
		if todo.ID == id {
			s.todos = append(s.todos[:i], s.todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("todo with ID %d not found", id)
}

func (s *InMemoryTodoService) Get(id int) (*Todo, error) {
	for _, todo := range s.todos {
		if todo.ID == id {
			return &todo, nil
		}
	}
	return nil, fmt.Errorf("todo with ID %d not found", id)
}

// TodoCLI handles the command-line interface
type TodoCLI struct {
	service TodoService
	scanner *bufio.Scanner
}

func NewTodoCLI(service TodoService) *TodoCLI {
	return &TodoCLI{
		service: service,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (cli *TodoCLI) Run() {
	fmt.Println("=== TODO Application ===")
	cli.printHelp()

	for {
		fmt.Print("\nEnter command: ")
		if !cli.scanner.Scan() {
			break
		}

		input := strings.TrimSpace(cli.scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]

		switch command {
		case "add":
			if len(parts) < 2 {
				fmt.Println("Usage: add <task>")
				continue
			}
			task := strings.Join(parts[1:], " ")
			todo := cli.service.Add(task)
			fmt.Printf("Added todo #%d: %s\n", todo.ID, todo.Task)

		case "list":
			cli.listTodos()

		case "complete":
			if len(parts) < 2 {
				fmt.Println("Usage: complete <id>")
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid ID")
				continue
			}
			if err := cli.service.Complete(id); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Completed todo #%d\n", id)
			}

		case "delete":
			if len(parts) < 2 {
				fmt.Println("Usage: delete <id>")
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid ID")
				continue
			}
			if err := cli.service.Delete(id); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Deleted todo #%d\n", id)
			}

		case "help":
			cli.printHelp()

		case "exit", "quit":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func (cli *TodoCLI) listTodos() {
	todos := cli.service.List()
	if len(todos) == 0 {
		fmt.Println("No todos found.")
		return
	}

	fmt.Println("\nYour Todos:")
	for _, todo := range todos {
		status := " "
		if todo.Completed {
			status = "✓"
		}
		fmt.Printf("[%s] #%d: %s\n", status, todo.ID, todo.Task)
	}
}

func (cli *TodoCLI) printHelp() {
	fmt.Println(`
Available commands:
  add <task>      - Add a new todo
  list            - List all todos
  complete <id>   - Mark todo as completed
  delete <id>     - Delete a todo
  help            - Show this help message
  exit/quit       - Exit the application`)
}

func main() {
	service := NewInMemoryTodoService()
	cli := NewTodoCLI(service)
	cli.Run()
}