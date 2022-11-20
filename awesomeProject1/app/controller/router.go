package controller

import (
	"net/http"
	"os"
)

type Router interface {
	FetchPoints(w http.ResponseWriter, r *http.Request)
	AddPoint(w http.ResponseWriter, r *http.Request)
	ChangePoint(w http.ResponseWriter, r *http.Request)
	DeletePoint(w http.ResponseWriter, r *http.Request)
}

type router struct {
	tc PointController
}

func CreateRouter(tc PointController) Router {
	return &router{tc}
}

func (ro *router) FetchPoints(w http.ResponseWriter, r *http.Request) {
	// プリフライトリクエスト用に設定している。
	w.Header().Set("Access-Control-Allow-Headers", "*")
	// CORSエラー対策。APIを叩くフロント側のURLを渡す。
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ORIGIN"))

	// 返却する値のContent-Typeを設定。
	w.Header().Set("Content-Type", "application/json")

	// controllerを呼び出す。
	ro.tc.FetchPoints(w, r)
}

func (ro *router) AddPoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ORIGIN"))
	w.Header().Set("Content-Type", "application/json")

	// preflightでAPIが二度実行されてしまうことを防ぐ。
	if r.Method == "OPTIONS" {
		return
	}

	ro.tc.AddPoint(w, r)
}

func (ro *router) DeletePoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ORIGIN"))
	w.Header().Set("Content-Type", "application/json")

	// preflightでAPIが二度実行されてしまうことを防ぐ。
	if r.Method == "OPTIONS" {
		return
	}

	ro.tc.DeletePoint(w, r)
}

func (ro *router) ChangePoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ORIGIN"))
	w.Header().Set("Content-Type", "application/json")

	// preflightでAPIが二度実行されてしまうことを防ぐ。
	if r.Method == "OPTIONS" {
		return
	}

	ro.tc.ChangePoint(w, r)
}
