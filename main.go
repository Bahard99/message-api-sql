package main

import (
	"encoding/hex"
	"log"
	"encoding/json"
	"database/sql"
	"time"
	"crypto/md5"

	"github.com/gofiber/fiber"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID 		string	`json:"id_user,omitempty"`
	Acode	string	`json:"auth_code,omitempty"`
	Uname	string	`json:"username,omitempty"`
}

type Args struct {
	Id_user1 string	`json:"id_user1,omitempty"`
	Acode	 string	`json:"auth_code,omitempty"`
	Msg		 string `json:"message,omitempty"`
	Id_user2 string	`json:"id_user2,omitempty"`
	Id_conv  string	`json:"id_conv,omitempty"`
}

type Message struct {
	Username  string	`json:"username,omitempty"`
	Message   string	`json:"message_value,omitempty"`
	Time	  string	`json:"timestamp,omitempty"`
}

func conn() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/tesrakamin")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main()  {
	log.Println("Starting server at localhost:8080")

	app := fiber.New()
	app.Post("/story1", story1)
	app.Post("/story2", story2)
	app.Get("/story3/:id", story3)
	app.Get("/story4/:id", story4)

	app.Listen(8080)

}

func story1(c *fiber.Ctx) {
	db, err := conn()
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	defer db.Close()

	var args Args
	json.Unmarshal([]byte(c.Body()), &args)

	if args.Msg == "" {
		c.Status(500).Send("Message cannot be empty!!")
		return
	}

	var each = User{}

	err = db.QueryRow("SELECT * FROM user WHERE id_user = ? AND auth_code = ?", args.Id_user1, args.Acode).Scan(&each.ID, &each.Acode, &each.Uname)
	switch {
	case err == sql.ErrNoRows:
		c.Status(500).Send("Credential Error!!")
		return
	case err != nil:
		c.Status(500).Send("Error Select : ", err)
		return
	}

	dt := time.Now().Format(time.RFC3339)
	str := args.Id_user1 + args.Id_user2 
	md5str := md5.Sum([]byte(str))
	idconv := hex.EncodeToString(md5str[:])

	_, err = db.Exec("INSERT INTO conv VALUES (?, ?, ?, ?)", idconv, args.Id_user1, args.Id_user2, dt)
	if err != nil {
		c.Status(500).Send("Error Insert 1 : ", err)
		return
	}

	_, err = db.Exec("INSERT INTO message (id_conv, id_user, message_value, timestamp) VALUES (?, ?, ?, ?)", idconv, args.Id_user1, args.Msg, dt)
	if err != nil {
		c.Status(500).Send("Error Insert 2 : ", err)
		return
	}

	res := "Message from " + each.Uname + " : '" + args.Msg +"'"

	c.Status(200).Send(res)
}

func story2(c *fiber.Ctx) {
	db, err := conn()
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	defer db.Close()

	var args Args
	json.Unmarshal([]byte(c.Body()), &args)

	if args.Msg == "" {
		c.Status(500).Send("Message cannot be empty!!")
		return
	}

	var each = User{}

	err = db.QueryRow("SELECT * FROM user WHERE id_user = ? AND auth_code = ?", args.Id_user1, args.Acode).Scan(&each.ID, &each.Acode, &each.Uname)
	switch {
	case err == sql.ErrNoRows:
		c.Status(500).Send("Credential Error!!")
		return
	case err != nil:
		c.Status(500).Send("Error Select : ", err)
		return
	}

	dt := time.Now().String()

	_, err = db.Exec("INSERT INTO message (id_conv, id_user, message_value, timestamp) VALUES (?, ?, ?, ?)", args.Id_conv, args.Id_user1, args.Msg, dt)
	if err != nil {
		c.Status(500).Send("Error Insert : ", err)
		return
	}

	res := "Message from " + each.Uname + " : '" + args.Msg +"'"

	c.Status(200).Send(res)
}

func story3(c *fiber.Ctx) {
	db, err := conn()
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	defer db.Close()

	idconv := c.Params("id")

	messg, err := db.Query("SELECT u.username, m.message_value, m.timestamp FROM message m, user u WHERE m.id_user = u.id_user AND m.id_conv = ?", idconv)

    defer messg.Close()

    if err != nil {
        log.Fatal(err)
    }

	var allmessage []Message

    for messg.Next() {

        var message Message
        err := messg.Scan(&message.Username, &message.Message, &message.Time)

        if err != nil {
            log.Fatal(err)
        }

        allmessage = append(allmessage, message)
    }

	json, _ := json.Marshal(allmessage)

	c.Status(200).Send(json)
}

func story4(c *fiber.Ctx) {
	db, err := conn()
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	defer db.Close()

	idusr := c.Params("id")

	idconvlist, err := db.Query("SELECT id_conv FROM conv WHERE id_user1 = ? OR id_user2 = ?", idusr, idusr)

    defer idconvlist.Close()

    if err != nil {
        log.Fatal(err)
    }

	var listconv []Args

    for idconvlist.Next() {

        var idconv Args
        err := idconvlist.Scan(&idconv.Id_conv)

        if err != nil {
            log.Fatal(err)
        }

        listconv = append(listconv, idconv)
    }

	var listMessage []Message

	for _, convId := range listconv {

		var msgg = Message{}

		err = db.QueryRow("SELECT u.username, m.message_value, m.timestamp FROM message m, user u WHERE m.id_user = u.id_user AND m.id_conv = ? ORDER BY m.timestamp DESC LIMIT 1", convId.Id_conv).Scan(&msgg.Username, &msgg.Message, &msgg.Time)
		if err != nil {
			c.Status(500).Send("Error Select : ", err)
			return
		}

		listMessage = append(listMessage, msgg)
	}

	json, _ := json.Marshal(listMessage)

	c.Status(200).Send(json)
}
