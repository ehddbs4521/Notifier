package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrUnknownChannel = errors.New("알 수 없는 채널")
var ErrEmptyTitle = errors.New("제목이 비어 있음")
var ErrBodyTooLong = errors.New("내용이 너무 김")

type Result struct {
	Success        int
	UnknownChannel int
	EmptyTitle     int
	BodyTooLong    int
}

func (r *Result) add(err error) {
	if err == nil {
		r.Success++
	}
	if errors.Is(err, ErrUnknownChannel) {
		r.UnknownChannel++
	}
	if errors.Is(err, ErrEmptyTitle) {
		r.EmptyTitle++
	}
	if errors.Is(err, ErrBodyTooLong) {
		r.BodyTooLong++
	}
}

func makeNotifications(count int) []Notification {
	channels := []string{"push", "sms", "email", "kakao"}
	notifications := make([]Notification, 0, count)
	for i := 1; i <= count; i++ {
		title := fmt.Sprintf("알림 %d", i)
		if i%10 == 0 {
			title = ""
		}
		body := fmt.Sprintf("내용 %d", i)
		if i%7 == 0 {
			body = "아주 긴 내용입니다. 이 문자는 이십 자를 훌쩍 넘어갑니다."
		}
		notifications = append(notifications, Notification{
			ID: i, UserId: fmt.Sprintf("user-%d", i), Channel: channels[i%4], Title: title, Body: body,
		})
	}
	return notifications
}

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
	time.Sleep(100 * time.Millisecond)
	if n.Title == "" {
		return fmt.Errorf("push 발송 실패: %w", ErrEmptyTitle)
	}
	return nil
}

type SMSSender struct {
}

func (s SMSSender) Send(n Notification) error {
	time.Sleep(100 * time.Millisecond)
	length := len([]rune(n.Body))
	if length > 20 {
		return fmt.Errorf("sms 발송 실패(%d자): %w", length, ErrBodyTooLong)
	}
	return nil
}

type EmailSender struct {
}

func (e EmailSender) Send(n Notification) error {
	time.Sleep(100 * time.Millisecond)
	if n.Title == "" {
		return fmt.Errorf("email 발송 실패: %w", ErrEmptyTitle)
	}
	return nil
}

func dispatch(n Notification, senders map[string]Sender) error {
	sender, ok := senders[n.Channel]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownChannel, n.Channel)
	}
	return sender.Send(n)
}

func sendAllSequential(ns []Notification, senders map[string]Sender) Result {
	result := Result{}
	for _, n := range ns {
		err := dispatch(n, senders)
		result.add(err)
	}
	return result
}

func sendAllConcurrent(ns []Notification, senders map[string]Sender) Result {
	result := Result{}
	var wg sync.WaitGroup

	for _, n := range ns {
		wg.Add(1)
		go func(n Notification) {
			defer wg.Done()
			err := dispatch(n, senders)
			result.add(err)
		}(n)
	}

	wg.Wait()
	return result
}

func main() {

	senders := make(map[string]Sender)
	senders["push"] = PushSender{}
	senders["sms"] = SMSSender{}
	senders["email"] = EmailSender{}

	ns := makeNotifications(100)

	start := time.Now()
	r := sendAllSequential(ns, senders)
	elapsed := time.Since(start)

	failCount := r.UnknownChannel + r.EmptyTitle + r.BodyTooLong
	fmt.Printf("[순차] 성공 %d건, 실패 %d건 (알 수 없는 채널 %d, 제목 없음 %d, 내용 김 %d) / 걸린 시간 %.2fs\n", r.Success, failCount, r.UnknownChannel, r.EmptyTitle, r.BodyTooLong, elapsed.Seconds())

	start = time.Now()
	r = sendAllConcurrent(ns, senders)
	elapsed = time.Since(start)

	failCount = r.UnknownChannel + r.EmptyTitle + r.BodyTooLong
	fmt.Printf("[동시] 성공 %d건, 실패 %d건 (알 수 없는 채널 %d, 제목 없음 %d, 내용 김 %d) / 걸린 시간 %.2fs\n", r.Success, failCount, r.UnknownChannel, r.EmptyTitle, r.BodyTooLong, elapsed.Seconds())
}
