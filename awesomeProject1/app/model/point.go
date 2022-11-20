package model

import (
	"database/sql"
	"github.com/oklog/ulid"
	"math/rand"
	"net/http"
	"time"
)

type PointModel interface {
	FetchPoints() ([]*Todo, error)
	AddPoint(r *http.Request) (sql.Result, error)
	ChangePoint(r *http.Request) (sql.Result, error)
	DeletePoint(r *http.Request) (sql.Result, error)
}

type pointModel struct {
}

type Todo struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

/*type Point struct {
	Id         string `json:"id"`
	sendName   string `json:"sendName"`
	sendPoint  int    `json:"sendPoint"`
	givenName  string `json:"givenName"`
	givenPoint int    `json:"givenPoint"`
}*/

func CreatePointModel() PointModel {
	return &pointModel{}
}

func (tm *pointModel) FetchPoints() ([]*Todo, error) {
	sql := `SELECT id,name,status FROM todos`

	// fetchメソッドのようなSQLを実行してデータを取ってきたい場合はQueryメソッドを使う。
	rows, err := Db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []*Todo

	// 取ってきたtodoの数だけ構造体にはめ込んで、todosに入れていく。
	for rows.Next() {
		var (
			id, name, status string
		)

		// 取ってきたデータを宣言した変数にはめ込んでいく。
		if err := rows.Scan(&id, &name, &status); err != nil {
			return nil, err
		}

		points = append(points, &Todo{
			Id:     id,
			Name:   name,
			Status: status,
		})
	}

	// 構造体の入った配列をcontrollerに返却する。json化してレスポンスに書き込むのはcontrollerの役目。
	return points, nil
}

func (tm *pointModel) AddPoint(r *http.Request) (sql.Result, error) {
	// 省略
	// 乱数のseedとして現在時刻を呼び出す。
	t := time.Now()

	// エントロピー(乱雑さ)を作成。
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)

	// ulidを作成する。
	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	req := Todo{
		// Idは今生成したものであって、クライアントからのリクエストではないのでここに入れるかは悩んだが、まあヨシとする。
		Id:     id.String(),
		Name:   r.FormValue("sendName"),
		Status: r.FormValue("status"),
	}

	sql := `INSERT INTO contributes(id, sendName) VALUES(?, ?)`

	result, err := Db.Exec(sql, req.Id, req.Name, req.Status)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (tm *pointModel) ChangePoint(r *http.Request) (sql.Result, error) {
	err := r.ParseForm()

	if err != nil {
		return nil, nil
	}

	sql := `UPDATE todos SET status = ? WHERE id = ?`

	afterStatus := "作業中"
	if r.FormValue("status") == "作業中" {
		afterStatus = "完了"
	}

	result, err := Db.Exec(sql, afterStatus, r.FormValue("id"))

	if err != nil {
		return result, err
	}

	return result, nil
}

func (tm *pointModel) DeletePoint(r *http.Request) (sql.Result, error) {
	// 省略
	err := r.ParseForm()

	if err != nil {
		return nil, nil
	}

	sql := `DELETE FROM contributes WHERE id = ?`

	// deleteメソッドのようなSQLを実行してデータベースの操作だけしたい場合はExec。
	result, err := Db.Exec(sql, r.FormValue("id"))

	if err != nil {
		return result, err
	}

	return result, nil
}
