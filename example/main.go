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
	q := `name='deen' and age>=23 and (hobby in ('soccer', 'swim') or score>90)`
	result, _ := yql.Match(q, map[string]any{
		"name":  "deen",
		"age":   23,
		"hobby": "basketball",
		"score": int64(100),
	})
	fmt.Println("q=", q, "result=", result)
	q = `score intersect (7,1,9,5,3)`
	result, _ = yql.Match(q, map[string]any{
		"score": []int64{3, 100, 200},
	})
	fmt.Println("q=", q, "result=", result)
	q = `score in (7,1,9,5,3)`
	result, _ = yql.Match(q, map[string]any{
		"score": []int64{3, 5, 2},
	})
	fmt.Println("q=", q, "result=", result)
	q = `score.sum() > 10`
	result, _ = yql.Match(q, map[string]any{
		"score": []int{1, 2, 3, 4, 5},
	})
	fmt.Println("q=", q, "result=", result)
	//Output:
	//true
	//true
	//false
	//true
}

func exampleRule() {
	q := `name='deen' and age<>22 and (hobby in ('soccer', 'swim') or score>90)`
	ruler, _ := yql.Rule(q)
	result, _ := ruler.Match(map[string]any{
		"name":  "deen",
		"age":   23,
		"hobby": "basketball",
		"score": int64(100),
	})
	fmt.Println("q=", q, "result=", result)
	result, _ = ruler.Match(map[string]any{
		"name":  "deen",
		"age":   23,
		"hobby": "basketball",
		"score": int64(90),
	})
	fmt.Println("q=", q, "result=", result)
	//Output:
	//true
	//false
}
