package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Login
func Login(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	//Read From Request Body
	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed")
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	if email == "" || password == "" {
		sendErrorResponse(w, "Masih terdapat input kosong")
		return
	}

	hash := sha256.Sum256([]byte(password))
	passwordHash := hex.EncodeToString(hash[:])

	query := "SELECT userid, Name, UserType FROM USERS WHERE Email ='" + email + "' && Password='" + passwordHash + "'"
	var user model.User
	err1 := db.QueryRow(query).Scan(&user.ID, &user.Name, &user.UserType)

	if err1 != nil {
		if err1 == sql.ErrNoRows {
			sendErrorResponse(w, "Email atau Password salah")
			return
		}
		log.Println(err1)
		sendErrorResponse(w, "Terjadi Error")
		return
	} else {
		generateToken(w, user.ID, user.Name, user.UserType)
		sendSuccessResponse(w, "Login Success", nil)
		//sendMailLogin(user2)
	}

}

func Logout(w http.ResponseWriter, r *http.Request) {
	//User id ambil pakai cookie
	userid := getUserIdFromCookie(r)
	if userid == -1 {
		sendErrorResponse(w, "Tidak ada aktivitas login sebelumnya")
	} else {
		resetUserToken(w)
		sendSuccessResponse(w, "Logout Success", nil)
	}

}

// fungsi untuk register user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	//Read From Request Body
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Failed")
		return
	}
	name := r.Form.Get("name")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	address := r.Form.Get("address")
	telephoneNo := r.Form.Get("telephone")

	if name == "" || email == "" || password == "" || address == "" || telephoneNo == "" {
		sendErrorResponse(w, "Masih terdapat input kosong")
		return
	}

	hash := sha256.Sum256([]byte(password))
	passwordHash := hex.EncodeToString(hash[:])

	user := model.User{
		Name:        name,
		Email:       email,
		Password:    passwordHash,
		Address:     address,
		TelephoneNo: telephoneNo,
	}

	query := "SELECT userid FROM USERS WHERE Email ='" + email + "'"
	err1 := db.QueryRow(query).Scan(&user.ID)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			res1, errQuery := db.Exec("INSERT INTO users(name, email, password, address, telpNo)values(?,?,?,?,?)",
				user.Name,
				user.Email,
				user.Password,
				user.Address,
				user.TelephoneNo,
			)

			if errQuery != nil {
				log.Println(errQuery)
				sendErrorResponse(w, "Register Gagal")
			} else {
				id, _ := res1.LastInsertId()
				_, errQuery2 := db.Exec("INSERT INTO CART (userid) VALUES (?)", id)
				if errQuery2 != nil {
					log.Println(errQuery)
					sendErrorResponse(w, "Register Gagal")
				} else {
					sendSuccessResponse(w, "Register Berhasil", nil)
					sendMailRegis(user)
				}
			}
		}
	} else {
		sendErrorResponse(w, "Email sudah terdaftar")
	}
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	//Read From Request Body
	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed")
		return
	}
	//Ambil password lama dan baru dari form
	oldpassword := r.Form.Get("old_password")
	newpassword := r.Form.Get("new_password")

	//User id ambil pakai cookie
	userid := getUserIdFromCookie(r)

	//Password lama user dari database untuk dibandingkan
	var password = GetUserPassword(userid)

	//Hash password lama yang diinput untuk verifikasi
	hash := sha256.Sum256([]byte(oldpassword))
	passwordHash := hex.EncodeToString(hash[:])

	if password == passwordHash {
		//hash password baru
		hash2 := sha256.Sum256([]byte(newpassword))
		passwordHash2 := hex.EncodeToString(hash2[:])
		query := "UPDATE users SET password = '" + passwordHash2 + "' WHERE userid = " + strconv.Itoa(userid)
		_, errQuery := db.Exec(query)

		if errQuery != nil {
			sendErrorResponse(w, "Failed to change password!")
		} else {
			sendSuccessResponse(w, "Password successfully changed!", nil)
		}
	} else {
		sendErrorResponse(w, "Password does not match!")
	}
}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT userid, name, email, address, telpNo FROM users"
	name := r.URL.Query().Get("name")
	userid := r.URL.Query().Get("userid")

	if userid == "" {
		userid = strconv.Itoa(getUserIdFromCookie(r))
	}
	if name != "" {
		query += " WHERE name='" + name + "'"
	}
	if userid != "" {
		query += " WHERE userid=" + userid
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		sendErrorResponse(w, "Error")
	} else {
		var user model.User
		var users []model.User
		for rows.Next() {
			if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.TelephoneNo); err != nil {
				sendErrorResponse(w, "Error while scanning rows")
				return
			} else {
				users = append(users, user)
			}
		}
		if err == nil {
			sendSuccessResponse(w, "Success", users)
		} else {
			sendErrorResponse(w, "Error")
		}
	}

}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	currentID := getUserIdFromCookie(r)

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Error while parsing form")
		return
	}
	name := r.Form.Get("name")
	email := r.Form.Get("email")
	address := r.Form.Get("address")
	telpNo := r.Form.Get("telpNo")

	_, errQuery := db.Exec("UPDATE users SET name=?, email=?, address=?, telpNo=? WHERE userid=?", name, email, address, telpNo, currentID)

	if errQuery == nil {
		rows, err := db.Query("SELECT userid, name, email, address, telpNo FROM users WHERE userid = ?", currentID)
		if err != nil {
			sendErrorResponse(w, "Error while fetching updated data")
			return
		}
		defer rows.Close()

		var user model.User
		for rows.Next() {
			if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.TelephoneNo); err != nil {
				sendErrorResponse(w, "Error while scanning rows")
				return
			}
		}
		sendSuccessResponse(w, "Profile updated", user)
	} else {
		sendErrorResponse(w, "Failed to update profile")
	}

}

func RegisterSeller(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	currentID := getUserIdFromCookie(r)

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Error while parsing form")
		return
	}
	//Kasih peringatan kalau penjual akses ini karena mereka adalah penjual
	if getUserTypeFromCookie(r) == 2 {
		sendErrorResponse(w, "User is already a seller")
		return
	}
	inputpassword := r.Form.Get("password")
	password := GetUserPassword(currentID)

	hash := sha256.Sum256([]byte(inputpassword))
	passwordHash := hex.EncodeToString(hash[:])

	if passwordHash != password {
		sendErrorResponse(w, "Password does not match")
	} else {
		_, errQuery := db.Exec("UPDATE users SET usertype = ? WHERE userid = ?", 2, currentID)

		if errQuery != nil {
			sendErrorResponse(w, "Failed to register the user as a seller")
		} else {
			sendSuccessResponse(w, "Successfully registered the user as a seller", nil)
		}
	}
}

func GetUserPassword(id int) string {
	db := connect()
	defer db.Close()

	var password string
	query := "SELECT password FROM users WHERE userid = " + strconv.Itoa(id)
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return ""
	}

	for rows.Next() {
		if err := rows.Scan(&password); err != nil {
			log.Println(err)
			return ""
		}
	}
	return password
}

// context untuk redis
var ctx = context.Background()

func setCurUserToRedis(user model.User) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if err := rdb.HSet(ctx, "curUser", user).Err(); err != nil {
		panic(err)
	}
	rdb.Expire(ctx, "curUser", 30*time.Minute)
}
func getCurUserToRedis() model.User {
	var user model.User
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := rdb.HMGet(ctx, "curUser", "email", "type").Scan(&user)
	if err != nil {
		panic(err)
	}
	return user
}
