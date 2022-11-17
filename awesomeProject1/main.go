package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

//type UserResForHTTPGet struct {
//Id   string `json:"id"`
//Name string `json:"name"`
//Age  int    `json:"age"`
//}

// ① GoプログラムからMySQLへ接続
//var db *sql.DB

//func init() {
// ①-1
//	mysqlUser := os.Getenv("MYSQL_USER")
//	mysqlPwd := os.Getenv("MYSQL_PWD")
//	mysqlHost := os.Getenv("MYSQL_HOST")
//	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

//	connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
//	_db, err := sql.Open("mysql", connStr)

// ①-2 接続時のパラメータを指定
//_db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(localhost:3306)/%s", mysqlUser, mysqlUserPwd, mysqlDatabase))
//	if err != nil {
//		log.Fatalf("fail: sql.Open, %v\n", err)
//	}
// ①-3　接続できているか確認、接続情報に誤りがあればエラーになる
//	if err := _db.Ping(); err != nil {
//		log.Fatalf("fail: _db.Ping, %v\n", err)
//	}
//	_ = _db
//}

// ② /userでリクエストされたらnameパラメーターと一致する名前を持つレコードをJSON形式で返す
func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodPost:
		// ②-1
		name := r.URL.Query()
		log.Println("message:Hello")
		if name == nil {
			log.Println("fail: name is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// ②-2　GETリクエストのクエリパラメータから条件を満たすデータを取得
		rows, err := db.Query("SELECT id, name, age FROM user WHERE name = ?", name)
		if err != nil {
			log.Printf("fail: db.Query, %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// ②-3　取得した複数のレコード行分繰り返す
		users := make([]UserResForHTTPGet, 0)
		for rows.Next() {
			var u UserResForHTTPGet
			if err := rows.Scan(&u.Id, &u.Name, &u.Age); err != nil {
				log.Printf("fail: rows.Scan, %v\n", err)

				if err := rows.Close(); err != nil { // 500を返して終了するが、その前にrowsのClose処理が必要
					log.Printf("fail: rows.Close(), %v\n", err)
				}
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			users = append(users, u)
		}

		// ②-4　レスポンス用ユーザースライスをJSONへ変換
		bytes, err := json.Marshal(users)
		if err != nil {
			log.Printf("fail: json.Marshal, %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	default:
		log.Printf("fail: HTTP Method is %s\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func main() {
	// ② /userでリクエストされたらnameパラメーターと一致する名前を持つレコードをJSON形式で返す
	http.HandleFunc("/user", handler)

	// ③ Ctrl+CでHTTPサーバー停止時にDBをクローズする
	closeDBWithSysCall()

	// 8000番ポートでリクエストを待ち受ける
	log.Println("Listening...")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

// ③ Ctrl+CでHTTPサーバー停止時にDBをクローズする
//func closeDBWithSysCall() {
//	sig := make(chan os.Signal, 1)
//	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
//	go func() {
//		s := <-sig
//log.Printf("received syscall, %v", s)
//
//		if err := db.Close(); err != nil {
//			log.Fatal(err)
//		}
//		log.Printf("success: db.Close()")
//		os.Exit(0)
//	}()
//}
