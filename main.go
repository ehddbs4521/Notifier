package main

import "fmt"

type Notification struct {
	ID      int
	UserId  string
	Channel string
	Title   string
	Body    string
}

type Sender interface {
	Send(n Notification) error
}

type PushSender struct {
}

func (p PushSender) Send(n Notification) error {
	fmt.Printf("[push] #%d → %s: %s - %s\n", n.ID, n.UserId, n.Title, n.Body)
	return nil
}

type SMSSender struct {
}

func (s SMSSender) Send(n Notification) error {
	fmt.Printf("[sms] #%d → %s: %s\n", n.ID, n.UserId, n.Body)
	return nil
}

type EmailSender struct {
}

func (e EmailSender) Send(n Notification) error {
	fmt.Printf("[email] #%d → %s / 제목: %s / 내용: %s\n", n.ID, n.UserId, n.Title, n.Body)
	return nil
}

func dispatch(n Notification, senders map[string]Sender) error {
	sender, ok := senders[n.Channel]
	if !ok {
		return fmt.Errorf("알 수 없는 채널: %s", n.Channel)
	}
	return sender.Send(n)
}

func main() {

	senders := make(map[string]Sender)
	senders["push"] = PushSender{}
	senders["sms"] = SMSSender{}
	senders["email"] = EmailSender{}

	notifications := []Notification{
		{ID: 1, UserId: "user-1", Channel: "push", Title: "새 채팅", Body: "안녕하세요!"},
		{ID: 2, UserId: "user-2", Channel: "sms", Title: "인증번호", Body: "123456"},
		{ID: 3, UserId: "user-3", Channel: "email", Title: "주간 소식", Body: "이번 주 인기글을 확인하세요"},
		{ID: 4, UserId: "user-4", Channel: "kakao", Title: "이벤트", Body: "쿠폰이 도착했어요"},
	}

	for _, n := range notifications {
		err := dispatch(n, senders)
		if err != nil {
			fmt.Printf("발송 실패 #%d: %s\n", n.ID, err.Error())
		}
	}
}
