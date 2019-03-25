package main

import (
	"fmt"

	"github.com/mia0x75/yql"
)

func main() {
	exampleMatch()
	exampleRule()
}

func exampleMatch() {
	rawYQL := `name='deen' and age>=23 and (hobby in ('soccer', 'swim') or score>90)`
	result, _ := yql.Match(rawYQL, map[string]interface{}{
		"name":  "deen",
		"age":   23,
		"hobby": "basketball",
		"score": int64(100),
	})
	fmt.Println("rawYQL=", rawYQL, "result=", result)
	rawYQL = `score ∩ (7,1,9,5,3)`
	result, _ = yql.Match(rawYQL, map[string]interface{}{
		"score": []int64{3, 100, 200},
	})
	fmt.Println("rawYQL=", rawYQL, "result=", result)
	rawYQL = `score in (7,1,9,5,3)`
	result, _ = yql.Match(rawYQL, map[string]interface{}{
		"score": []int64{3, 5, 2},
	})
	fmt.Println("rawYQL=", rawYQL, "result=", result)
	rawYQL = `score.sum() > 10`
	result, _ = yql.Match(rawYQL, map[string]interface{}{
		"score": []int{1, 2, 3, 4, 5},
	})
	fmt.Println("rawYQL=", rawYQL, "result=", result)
	//Output:
	//true
	//true
	//false
	//true
}

func exampleRule() {
	rawYQL := `name='deen' and age<>22 and (hobby in ('soccer', 'swim') or score>90)`
	ruler, _ := yql.Rule(rawYQL)
	result, _ := ruler.Match(map[string]interface{}{
		"name":  "deen",
		"age":   23,
		"hobby": "basketball",
		"score": int64(100),
	})
	fmt.Println("rawYQL=", rawYQL, "result=", result)
	result, _ = ruler.Match(map[string]interface{}{
		"name":  "deen",
		"age":   23,
		"hobby": "basketball",
		"score": int64(90),
	})
	fmt.Println("rawYQL=", rawYQL, "result=", result)
	//Output:
	//true
	//false
}
