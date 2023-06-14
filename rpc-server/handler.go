package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

type Message struct {
	Chat     string `json:"chat"`
	Sender   string `json:"sender"`
	Text     string `json:"text"`
	SendTime int64  `json:"sendTime"`
}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	resp := rpc.NewSendResponse()
	resp.Code, resp.Msg = saveMessage(req.Message.Chat, req.Message.Sender, req.Message.Text, req.Message.SendTime)
	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	resp := rpc.NewPullResponse()
	resp.Code, resp.Msg = pullMessage(req, resp)
	return resp, nil
}

func saveMessage(id string, sender string, msg string, sendtime int64) (int32, string) {
	db, err := sql.Open("mysql", "newuser:Doggo@tcp(docker.for.mac.localhost:3306)/db")
	if err != nil {
		panic(err)
	}
	insert, err := db.Prepare("INSERT INTO chatroom VALUES (?, ?, ?, ?)")
	defer insert.Close()
	_, err = insert.Query(id, sender, msg, sendtime)
	if err != nil {
		return 500, err.Error() + fmt.Sprintf("INSERT INTO chatroom VALUES ('%s', '%s', '%s', %d)", id, sender, msg, sendtime)
	}
	return 0, "Success"
}

func pullMessage(req *rpc.PullRequest, resp *rpc.PullResponse) (int32, string) {
	id := req.Chat
	cursor := req.Cursor
	lim := req.Limit

	db, err := sql.Open("mysql", "newuser:Doggo@tcp(docker.for.mac.localhost:3306)/db")
	if err != nil {
		panic(err)
	}
	query := ""
	if *(req.Reverse) == true {
		query = fmt.Sprintf("SELECT * FROM Chatroom WHERE id = '%s' AND time_stamp >= %d ORDER BY time_stamp ASC LIMIT %d;", id, cursor, lim+1)
	} else {
		if cursor <= 0 {
			cursor = 1<<63 - 1
		}
		query = fmt.Sprintf("SELECT * FROM Chatroom WHERE id = '%s' AND time_stamp <= %d ORDER BY time_stamp DESC LIMIT %d;", id, cursor, lim+1)
	}
	rows, err := db.Query(query)
	if err != nil {
		return 500, err.Error() + query
	}
	var iterations int32 = 0
	var b bool = false
	resp.SetHasMore(&b)
	var (
		chat     string
		sender   string
		msg      string
		sendTime int64
	)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&chat, &sender, &msg, &sendTime)
		if err != nil {
			panic(err)
		}
		iterations += 1

		if iterations == lim+1 {
			resp.SetNextCursor(&sendTime)
			b = true
			resp.SetHasMore(&b)
			break
		}

		message := rpc.NewMessage()
		message.SetChat(id)
		message.SetSender(sender)
		message.SetText(msg)
		message.SetSendTime(sendTime)

		resp.Messages = append(resp.Messages, message)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return 0, "Success"
}
