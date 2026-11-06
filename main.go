package main

import (
	"errors"
	"fmt"
)

var ErrUnknownChannel = errors.New("알 수 없는 채널")
var ErrEmptyTitle = errors.New("제목이 비어 있음")
var ErrBodyTooLong = errors.New("내용이 너무 김")

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
	if n.Title == "" {
		return fmt.Errorf("push 발송 실패: %w", ErrEmptyTitle)
	}
	fmt.Printf("[push] #%d → %s: %s - %s\n", n.ID, n.UserId, n.Title, n.Body)
	return nil
}

type SMSSender struct {
}

func (s SMSSender) Send(n Notification) error {
	length := len([]rune(n.Body))
	if length > 20 {
		return fmt.Errorf("sms 발송 실패(%d자): %w", length, ErrBodyTooLong)
	}
	fmt.Printf("[sms] #%d → %s: %s\n", n.ID, n.UserId, n.Body)
	return nil
}

type EmailSender struct {
}

func (e EmailSender) Send(n Notification) error {
	if n.Title == "" {
		return fmt.Errorf("email 발송 실패: %w", ErrEmptyTitle)
	}
	fmt.Printf("[email] #%d → %s / 제목: %s / 내용: %s\n", n.ID, n.UserId, n.Title, n.Body)
	return nil
}

func dispatch(n Notification, senders map[string]Sender) error {
	sender, ok := senders[n.Channel]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownChannel, n.Channel)
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
		{ID: 5, UserId: "user-5", Channel: "sms", Title: "인증번호", Body: "인증번호는 123456입니다. 타인에게 절대 알려주지 마세요."},
		{ID: 6, UserId: "user-6", Channel: "push", Title: "", Body: "제목 없는 알림"},
	}

	successCount := 0
	ErrUnknownChannelCount := 0
	ErrEmptyTitleCount := 0
	ErrBodyTooLongCount := 0
	for _, n := range notifications {
		err := dispatch(n, senders)
		if err != nil {
			if errors.Is(err, ErrEmptyTitle) {
				ErrEmptyTitleCount++
			}
			if errors.Is(err, ErrBodyTooLong) {
				ErrBodyTooLongCount++
			}
			if errors.Is(err, ErrUnknownChannel) {
				ErrUnknownChannelCount++
			}
			fmt.Printf("발송 실패 #%d: %s\n", n.ID, err.Error())
		} else {
			successCount++
		}
	}

	fmt.Println("---")
	fmt.Printf("성공 %d건, 실패 %d건\n", successCount, ErrUnknownChannelCount+ErrEmptyTitleCount+ErrBodyTooLongCount)
	fmt.Printf("실패 사유: 알 수 없는 채널 %d건, 내용 오류 %d건\n", ErrUnknownChannelCount, ErrEmptyTitleCount+ErrBodyTooLongCount)
}
