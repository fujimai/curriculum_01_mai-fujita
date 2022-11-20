package main

import (
	"flag"
	"fmt"
	"github.com/curriculum_01_mai-fujita/awesomeProject1/controller"
	"github.com/curriculum_01_mai-fujita/awesomeProject1/model"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

// todoModelのインスタンスを作成。
var tm = model.CreatePointModel()

// todoControllerのインスタンスを作成。todoModelを注入。
var tc = controller.CreatePointController(tm)

// routerのインスタンスを作成。todoControllerを作成。
var ro = controller.CreateRouter(tc)

func migrate() {
	sql := `INSERT INTO todos(id, name, status) VALUES('00000000000000000000000000','買い物', '作業中'),('00000000000000000000000001','洗濯', '作業中'),('00000000000000000000000002','皿洗い', '完了');`

	_, err := model.Db.Exec(sql)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Migration is success!")
}

func main() {
	f := flag.String("option", "", "migrate database or not")
	flag.Parse()
	if *f == "migrate" {
		migrate()
	}
	// 省略
	http.HandleFunc("/fetch-points", ro.FetchPoints)
	http.HandleFunc("/add-point", ro.AddPoint)
	http.HandleFunc("/delete-point", ro.DeletePoint)
	http.HandleFunc("/change-point", ro.ChangePoint)
	http.ListenAndServe(":8080", nil)
}

/*
var db *sql.DB

type UserResForHTTPGet struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Points struct {
	Id         string `json:"id"`
	SendName   string `json:"send_name"`
	SendPoint  int    `json:"send_point"`
	GivenName  string `json:"given_name"`
	GivenPoint int    `json:"given_point"`
}

// ① GoプログラムからMySQLへ接続

func init() {
	// ①-1
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPwd := os.Getenv("MYSQL_PWD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
	_db, err := sql.Open("mysql", connStr)

	// ①-2 接続時のパラメータを指定
	//_db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(localhost:3306)/%s", mysqlUser, mysqlPwd, mysqlDatabase))
	if err != nil {
		log.Fatalf("fail: sql.Open, %v\n", err)
	}
	// ①-3　接続できているか確認、接続情報に誤りがあればエラーになる
	if err := _db.Ping(); err != nil {
		log.Fatalf("fail: _db.Ping, %v\n", err)
	}
	_ = _db
}

// ② /userでリクエストされたらnameパラメーターと一致する名前を持つレコードをJSON形式で返す
func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		// ②-1
		name := r.URL.Query()
		/*				log.Println("message:Hello")
								if name == nil {
									log.Println("fail: name is empty")
									w.WriteHeader(http.StatusBadRequest)
									return
						         }


		// ②-2　GETリクエストのクエリパラメータから条件を満たすデータを取得
		rows, err := db.Query("SELECT user.id, user.name, user.age FROM user WHERE user.name = ?", name)
		if err != nil {
			log.Printf("fail: db.Query, %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()
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

	case http.MethodPost:
		tx, err := db.Begin()

		if err != nil {
			tx.Rollback() //何らかの処理が失敗
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		var u string

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			tx.Rollback()
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(u)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)

		return

	default:
		log.Printf("fail: HTTP Method is %s\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// /pointでリクエストされた時
func pointHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodPost:
		var userId string
		// http.Postのレスポンス処理
		if err := json.NewDecoder(r.Body).Decode(&userId); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rows, err := db.Query("SELECT contribute.send_name, contribute.send_point FROM contribute WHERE contribute.Id =?", userId)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		points := make([]Points, 0)
		for rows.Next() {
			var p Points
			if err := rows.Scan(&p.Id, &p.SendPoint, &p.SendName, &p.GivenName, &p.GivenPoint); err != nil {
				log.Printf("points, %v\n", err)

				if err := rows.Close(); err != nil {
					log.Printf("pointsss, %v\n", err)
				}
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			points = append(points, p)
		}
		bytes, err := json.Marshal(points)
		if err != nil {
			fmt.Println("point1", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)

	case http.MethodDelete:
		var pointId string
		if err := json.NewDecoder(r.Body).Decode(&pointId); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tx, err := db.Begin()
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		queDelete := fmt.Sprintf("DELETE FROM points WHERE id = '%v'", pointId)
		stmt, err := tx.Prepare(queDelete)
		if err != nil {
			stmt.Close()
			tx.Rollback()
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if res, err := stmt.Exec(); err != nil {
			tx.Rollback()
			fmt.Println(res.RowsAffected())
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		/*case http.MethodPut:
		var edit Edit
		if err := json.NewDecoder(r.Body).Decode()

	}
}

func userSignup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodPost:
		var user string
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		queStr := fmt.Sprintf("SELECT * FROM users WHERE name ='%v'", user)
		rows, err := db.Query(queStr)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		cnt := 0
		for rows.Next() {
			cnt += 1
		}
		if cnt >= 1 {
			w.WriteHeader(http.StatusConflict)
			return
		}
		tx, err := db.Begin()
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ins, err := tx.Prepare("INSERT INTO users(id,name)VALUES (?,?)")

		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t := time.Now()
		entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
		id := ulid.MustNew(ulid.Timestamp(t), entropy).String()
		res, err := ins.Exec(id, user)

		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			fmt.Println(res.LastInsertId())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	// ② /userでリクエストされたらnameパラメーターと一致する名前を持つレコードをJSON形式で返す
	http.HandleFunc("/user", handler)
	http.HandleFunc("/point", pointHandler)
	http.HandleFunc("/signup", userSignup)
	http.ListenAndServe(":8080", nil)
	// ③ Ctrl+CでHTTPサーバー停止時にDBをクローズする
	closeDBWithSysCall()

	// 8000番ポートでリクエストを待ち受ける
	log.Println("Listening...")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

//③ Ctrl+CでHTTPサーバー停止時にDBをクローズする
func closeDBWithSysCall() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		log.Printf("received syscall, %v", s)

		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
		log.Printf("success: db.Close()")
		os.Exit(0)
	}()
}
*/
