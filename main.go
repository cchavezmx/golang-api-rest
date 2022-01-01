// para ejecutar el servidor → CompileDaemon -command ./go-apirest-server

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:"ID"`
	Name    string `json:"Name"`
	Contnet string `json:"Content"`
}

type allTasks []task

var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task 1",
		Contnet: "Content 1",
	},
}

func main() {
	// definimos el router
	router := mux.NewRouter().StrictSlash(true)

	// definimos las rutas
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/createTask", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTaskByID).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	// iniciamos el servidor
	log.Fatal(http.ListenAndServe(":3000", router))
	println("Server started on port 3000")
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	// se importa el modulo
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Inserte datos correctos")
	}

	json.Unmarshal(reqBody, &newTask)

	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	// tipo de dato que se va a enviar
	w.Header().Set("Content-Type", "application/json")
	// estatus del servidor
	w.WriteHeader(http.StatusCreated)
	// enviamos el json
	json.NewEncoder(w).Encode(newTask)

}

func getTaskByID(w http.ResponseWriter, r *http.Request) {

	// es coimo hacer un req.params
	vars := mux.Vars(r)
	// se obtiene el id de params
	key := vars["id"]
	// convertimos el id a int
	taskId, error := strconv.Atoi(key)

	if error != nil {
		fmt.Fprintf(w, "Inserte datos correctos")
		return
	}

	for _, task := range tasks {
		if task.ID == taskId {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}

}

func deleteTask(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["id"]

	taskId, error := strconv.Atoi(key)

	if error != nil {
		fmt.Fprintf(w, "Inserte datos correctos")
		return
	}

	for index, task := range tasks {

		if task.ID == taskId {
			// para eliminar un elemento de un arreglo
			// buscamos el elemento que queremos eliminar y lo eliminamos
			tasks = append(tasks[:index], tasks[index+1:]...)
			fmt.Fprintf(w, "Task with ID %v has been deleted", taskId)
		}
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(tasks)

}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	var updateTask task

	taskId, error := strconv.Atoi(key)

	if error != nil {
		fmt.Fprintf(w, "Inserte datos correctos")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Inserte datos correctos")
		return
	}
	// guardamos el json en una variable
	json.Unmarshal(reqBody, &updateTask)

	for index, task := range tasks {
		if task.ID == taskId {
			tasks = append(tasks[:index], tasks[index+1:]...)
			updateTask.ID = taskId
			tasks = append(tasks, updateTask)

			fmt.Fprintf(w, "Task with ID %v has been updated", taskId)
		}
	}

}
