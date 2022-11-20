package controller

import (
	"encoding/json"
	"fmt"
	"github.com/curriculum_01_mai-fujita/awesomeProject1/model"
	"net/http"
)

type PointController interface {
	FetchPoints(w http.ResponseWriter, r *http.Request)
	AddPoint(w http.ResponseWriter, r *http.Request)
	ChangePoint(w http.ResponseWriter, r *http.Request)
	DeletePoint(w http.ResponseWriter, r *http.Request)
}

type pointController struct {
	tm model.PointModel
}

// model同様、インターフェースが戻り値の型になっているところが肝。
func CreatePointController(tm model.PointModel) PointController {
	return &pointController{tm}
}

func (tc *pointController) FetchPoints(w http.ResponseWriter, r *http.Request) {
	// modelのFetchTodosを実行。SQLを実行してtodosを取得。
	todos, err := tc.tm.FetchPoints()

	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	// json形式に変換。
	json, err := json.Marshal(todos)

	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	// レスポンスにtodosを入れる。jsonはそのままだとbyte型の配列なのでstring型へ変換。
	fmt.Fprint(w, string(json))
}

func (tc *pointController) AddPoint(w http.ResponseWriter, r *http.Request) {
	result, err := tc.tm.AddPoint(r)

	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	json, err := json.Marshal(result)

	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	fmt.Fprint(w, string(json))
}

func (tc *pointController) ChangePoint(w http.ResponseWriter, r *http.Request) {
	result, err := tc.tm.ChangePoint(r)

	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	json, err := json.Marshal(result)

	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	fmt.Fprint(w, string(json))
}

func (tc *pointController) DeletePoint(w http.ResponseWriter, r *http.Request) {
	result, err := tc.tm.DeletePoint(r)

	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	json, err := json.Marshal(result)

	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	fmt.Fprint(w, string(json))
}
