//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"
)

type task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

var (
	tasks          []task
	nextID         = 1
	itemCallbacks  []js.Func
	staticHandlers []js.Func
)

func main() {
	defer releaseCallbacks()

	doc := js.Global().Get("document")
	if !doc.Truthy() {
		fmt.Println("document not found")
		return
	}

	loadTasks()
	setupForm(doc)
	renderTasks(doc)

	js.Global().Get("document").Get("body").Get("classList").Call("add", "wasm-ready")
	js.Global().Get("document").Call("getElementById", "loading").Set("textContent", "WASM todo app is ready.")

	<-make(chan struct{})
}

func setupForm(doc js.Value) {
	form := doc.Call("getElementById", "todo-form")
	input := doc.Call("getElementById", "todo-input")
	status := doc.Call("getElementById", "todo-status")

	if !form.Truthy() || !input.Truthy() {
		fmt.Println("todo form not found")
		return
	}

	submitHandler := js.FuncOf(func(this js.Value, args []js.Value) any {
		args[0].Call("preventDefault")
		text := strings.TrimSpace(input.Get("value").String())
		if text == "" {
			status.Set("textContent", "Please type something before adding.")
			return nil
		}

		addTask(text)
		input.Set("value", "")
		status.Set("textContent", "Task added.")
		renderTasks(doc)

		return nil
	})

	form.Call("addEventListener", "submit", submitHandler)
	staticHandlers = append(staticHandlers, submitHandler)
}

func addTask(text string) {
	tasks = append(tasks, task{
		ID:   nextID,
		Text: text,
	})
	nextID++
	persistTasks()
}

func toggleTask(id int) {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = !tasks[i].Done
			break
		}
	}
	persistTasks()
}

func deleteTask(id int) {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}
	persistTasks()
}

func renderTasks(doc js.Value) {
	resetItemCallbacks()

	list := doc.Call("getElementById", "todo-list")
	count := doc.Call("getElementById", "todo-count")
	empty := doc.Call("getElementById", "todo-empty")

	if !list.Truthy() {
		fmt.Println("todo list element missing")
		return
	}

	list.Set("innerHTML", "")

	if len(tasks) == 0 {
		empty.Get("classList").Call("remove", "hidden")
		count.Set("textContent", "Nothing to do yet.")
		return
	}

	empty.Get("classList").Call("add", "hidden")

	for _, t := range tasks {
		// Capture ID per-iteration for event handlers.
		taskID := t.ID
		row := doc.Call("createElement", "li")
		row.Get("classList").Call("add", "todo-item")
		if t.Done {
			row.Get("classList").Call("add", "todo-done")
		}

		label := doc.Call("createElement", "span")
		label.Set("textContent", t.Text)

		checkbox := doc.Call("createElement", "input")
		checkbox.Set("type", "checkbox")
		checkbox.Set("checked", t.Done)
		cbHandler := js.FuncOf(func(this js.Value, args []js.Value) any {
			toggleTask(taskID)
			renderTasks(doc)
			return nil
		})
		checkbox.Call("addEventListener", "change", cbHandler)

		deleteBtn := doc.Call("createElement", "button")
		deleteBtn.Get("classList").Call("add", "ghost")
		deleteBtn.Set("textContent", "Delete")
		delHandler := js.FuncOf(func(this js.Value, args []js.Value) any {
			deleteTask(taskID)
			renderTasks(doc)
			return nil
		})
		deleteBtn.Call("addEventListener", "click", delHandler)

		row.Call("appendChild", checkbox)
		row.Call("appendChild", label)
		row.Call("appendChild", deleteBtn)
		list.Call("appendChild", row)

		itemCallbacks = append(itemCallbacks, cbHandler, delHandler)
	}

	plural := "s"
	if len(tasks) == 1 {
		plural = ""
	}
	count.Set("textContent", fmt.Sprintf("%d item%s", len(tasks), plural))
}

func persistTasks() {
	storage := js.Global().Get("localStorage")
	if !storage.Truthy() {
		return
	}

	data, err := json.Marshal(tasks)
	if err != nil {
		return
	}
	storage.Call("setItem", "kinglet-wasm-todos", string(data))
}

func loadTasks() {
	storage := js.Global().Get("localStorage")
	if !storage.Truthy() {
		return
	}

	raw := storage.Call("getItem", "kinglet-wasm-todos")
	if !raw.Truthy() {
		return
	}

	var saved []task
	if err := json.Unmarshal([]byte(raw.String()), &saved); err != nil {
		return
	}

	tasks = saved
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	nextID = maxID + 1
}

func releaseCallbacks() {
	for _, fn := range staticHandlers {
		fn.Release()
	}
	resetItemCallbacks()
}

func resetItemCallbacks() {
	for _, fn := range itemCallbacks {
		fn.Release()
	}
	itemCallbacks = nil
}
