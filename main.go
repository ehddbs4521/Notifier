package main

import "fmt"

type Notification struct {
	ID      int
	UserId  string
	Channel string
	Title   string
	Body    string
}

func send(n Notification) {
	fmt.Printf("[%s] #%d → %s: %s - %s\n", n.Channel, n.ID, n.UserId, n.Title, n.Body)
}

func main() {
	notifications := []Notification{
		{ID: 1, UserId: "user-1", Channel: "push", Title: "새 채팅", Body: "안녕하세요!"},
		{ID: 2, UserId: "user-2", Channel: "sms", Title: "인증번호", Body: "123456"},
		{ID: 3, UserId: "user-3", Channel: "email", Title: "주간 소식", Body: "이번 주 인기글을 확인하세요"},
	}

	for _, n := range notifications {
		send(n)
	}
}
