package handler

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pliliya111/go_sprint2_project/internal/model"
)

var (
	timeAdditionMS       = getEnvInt("TIME_ADDITION_MS", 1000)
	timeSubtractionMS    = getEnvInt("TIME_SUBTRACTION_MS", 1000)
	timeMultiplicationMS = getEnvInt("TIME_MULTIPLICATIONS_MS", 1000)
	timeDivisionMS       = getEnvInt("TIME_DIVISIONS_MS", 1000)
	validExpressionRegex = regexp.MustCompile(`^[\d\s\+\-\*\/\(\)]+$`)
)
var (
	Expressions = make(map[string]*model.Expression) // Хранение выражений
	Tasks       []*model.Task                        // Очередь задач
	TaskMap     = make(map[string]*model.Task)       // Мапа для быстрого поиска задач по ID
	Results     = make(map[string]interface{})       // Результаты задач
	Mutex       = &sync.Mutex{}                      // Мьютекс для потокобезопасности
)

func getEnvInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func parseExpression(expression string) []string {
	re := regexp.MustCompile(`\d+|\+|\-|\*|\/`)
	return re.FindAllString(expression, -1)
}
func AddExpression(c *gin.Context) {
	var request struct {
		Expression string `json:"expression"`
	}

	// Привязка JSON-запроса
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid data"})
		return
	}

	// Валидация выражения
	if !validExpressionRegex.MatchString(request.Expression) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Expression is not valid"})
		return
	}

	// Создание нового выражения
	expressionID := uuid.New().String()
	expr := &model.Expression{
		ID:         expressionID,
		Expression: request.Expression,
		Status:     "pending",
	}

	// Сохранение выражения
	Mutex.Lock()
	Expressions[expressionID] = expr
	Mutex.Unlock()

	// Разбор выражения на задачи
	tokens := parseExpression(request.Expression)

	// Обработка операций умножения и деления
	for i := 0; i < len(tokens); i++ {
		if tokens[i] == "*" || tokens[i] == "/" {
			arg1ID := tokens[i-1]
			arg2ID := tokens[i+1]
			taskID := uuid.New().String()
			task := &model.Task{
				ID:        taskID,
				Arg1:      arg1ID,
				Arg2:      arg2ID,
				Operation: tokens[i],
			}
			TaskMap[taskID] = task
			Tasks = append(Tasks, task)
			tokens[i-1] = taskID
			tokens[i] = ""
			tokens[i+1] = ""
		}
	}

	// Удаление пустых токенов
	filteredTokens := []string{}
	for _, token := range tokens {
		if token != "" {
			filteredTokens = append(filteredTokens, token)
		}
	}
	tokens = filteredTokens

	// Обработка операций сложения и вычитания
	for i := 0; i < len(tokens); i++ {
		if tokens[i] == "+" || tokens[i] == "-" {
			arg1ID := tokens[i-1]
			arg2ID := tokens[i+1]
			taskID := uuid.New().String()
			task := &model.Task{
				ID:        taskID,
				Arg1:      arg1ID,
				Arg2:      arg2ID,
				Operation: tokens[i],
			}
			TaskMap[taskID] = task
			Tasks = append(Tasks, task)
			tokens[i-1] = taskID
			tokens[i] = ""
			tokens[i+1] = ""
		}
	}

	// Установка результата выражения
	if len(Tasks) > 0 {
		expr.Result = Tasks[len(Tasks)-1].ID
	}
	expr.Status = "in_progress"

	// Логирование задач (для отладки)
	fmt.Println("All tasks:")
	for _, task := range Tasks {
		fmt.Printf("Task ID: %s, Arg1: %s, Arg2: %s, Operation: %s\n", task.ID, task.Arg1, task.Arg2, task.Operation)
	}

	// Возврат ответа
	c.JSON(http.StatusCreated, gin.H{"id": expressionID})
}

// func AddExpression(c *gin.Context) {
// 	var request struct {
// 		Expression string `json:"expression"`
// 	}
// 	if err := c.BindJSON(&request); err != nil {
// 		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid data"})
// 		return
// 	}

// 	expressionID := uuid.New().String()
// 	expr := &model.Expression{
// 		ID:         expressionID,
// 		Expression: request.Expression,
// 		Status:     "pending",
// 	}

// 	Mutex.Lock()
// 	Expressions[expressionID] = expr
// 	Mutex.Unlock()

// 	tokens := parseExpression(request.Expression)

// 	for i := 0; i < len(tokens); i++ {
// 		if tokens[i] == "*" || tokens[i] == "/" {
// 			arg1ID := tokens[i-1]
// 			arg2ID := tokens[i+1]
// 			taskID := uuid.New().String()
// 			task := &model.Task{
// 				ID:        taskID,
// 				Arg1:      arg1ID,
// 				Arg2:      arg2ID,
// 				Operation: tokens[i],
// 			}
// 			TaskMap[taskID] = task
// 			Tasks = append(Tasks, task)
// 			tokens[i-1] = taskID
// 			tokens[i] = ""
// 			tokens[i+1] = ""
// 		}
// 	}

// 	filteredTokens := []string{}
// 	for _, token := range tokens {
// 		if token != "" {
// 			filteredTokens = append(filteredTokens, token)
// 		}
// 	}
// 	tokens = filteredTokens

// 	for i := 0; i < len(tokens); i++ {
// 		if tokens[i] == "+" || tokens[i] == "-" {
// 			arg1ID := tokens[i-1]
// 			arg2ID := tokens[i+1]
// 			taskID := uuid.New().String()
// 			task := &model.Task{
// 				ID:        taskID,
// 				Arg1:      arg1ID,
// 				Arg2:      arg2ID,
// 				Operation: tokens[i],
// 			}
// 			TaskMap[taskID] = task
// 			Tasks = append(Tasks, task)

// 			tokens[i-1] = taskID
// 			tokens[i] = ""
// 			tokens[i+1] = ""
// 		}
// 	}

// 	if len(Tasks) > 0 {
// 		expr.Result = Tasks[len(Tasks)-1].ID
// 	}
// 	expr.Status = "in_progress"

// 	fmt.Println("All tasks:")
// 	for _, task := range Tasks {
// 		fmt.Printf("Task ID: %s, Arg1: %s, Arg2: %s, Operation: %s\n", task.ID, task.Arg1, task.Arg2, task.Operation)
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"id": expressionID})
// }

func GetExpressions(c *gin.Context) {
	Mutex.Lock()
	defer Mutex.Unlock()

	expressionsList := make([]gin.H, 0, len(Expressions))
	for _, expr := range Expressions {
		expressionsList = append(expressionsList, gin.H{
			"id":     expr.ID,
			"status": expr.Status,
			"result": expr.Result,
		})
	}

	c.JSON(http.StatusOK, gin.H{"expressions": expressionsList})
}

func GetExpressionByID(c *gin.Context) {
	expressionID := c.Param("id")

	Mutex.Lock()
	expr, exists := Expressions[expressionID]
	Mutex.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "expression not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"expression": gin.H{
			"id":     expr.ID,
			"status": expr.Status,
			"result": expr.Result,
		},
	})
}

func getOperationTime(operation string) int {
	switch operation {
	case "+":
		return timeAdditionMS
	case "-":
		return timeSubtractionMS
	case "*":
		return timeMultiplicationMS
	case "/":
		return timeDivisionMS
	default:
		return 0
	}
}

func GetTask(c *gin.Context) {
	Mutex.Lock()
	defer Mutex.Unlock()

	if len(Tasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no tasks available"})
		return
	}

	task := Tasks[0]
	Tasks = Tasks[1:]

	c.JSON(http.StatusOK, gin.H{
		"task": gin.H{
			"id":             task.ID,
			"arg1":           task.Arg1,
			"arg2":           task.Arg2,
			"operation":      task.Operation,
			"operation_time": getOperationTime(task.Operation),
			"result":         task.Result,
		},
	})
}
func processTask(task *model.Task, result interface{}) {
	Mutex.Lock()
	defer Mutex.Unlock()

	task.Result = result
	Results[task.ID] = result

	for _, expr := range Expressions {
		if expr.Result == task.ID {
			expr.Status = "completed"
			expr.Result = result
			break
		}
	}

	for _, t := range Tasks {
		if arg1ID, ok := t.Arg1.(string); ok && arg1ID == task.ID {
			t.Arg1 = result
		}
		if arg2ID, ok := t.Arg2.(string); ok && arg2ID == task.ID {
			t.Arg2 = result
		}
	}
}

func SubmitTaskResult(c *gin.Context) {
	var request struct {
		ID     string      `json:"id"`
		Result interface{} `json:"result"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid data"})
		return
	}

	Mutex.Lock()
	Results[request.ID] = request.Result
	fmt.Println(TaskMap)
	task, exists := TaskMap[request.ID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	Results[request.ID] = request.Result
	go processTask(task, request.Result)
	Mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "result submitted"})
}
