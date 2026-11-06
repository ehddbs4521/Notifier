# Notifier

Go를 배우면서 만드는 작은 **알림 발송기** 토이 프로젝트입니다.
실제 푸시·문자·이메일을 보내지는 않고, 콘솔 출력으로 발송을 흉내 냅니다.
실제 알림 시스템의 핵심 흐름(알림 생성 → 채널별 발송 → 실패 처리 → 대량 동시 발송)을 단계별로 작게 구현해 봅니다.

## 실행 방법

```bash
go run .
```

Go 1.27 이상에서 동작합니다.

## 단계별 진행

각 단계는 이전 단계 코드에 요구사항을 하나씩 더해 가는 방식입니다.

| 단계 | 내용 | 연습하는 Go 문법 | 상태 |
| --- | --- | --- | --- |
| 1 | 알림 구조체 만들고 발송 출력 | 구조체, 슬라이스, 함수, `for range`, `Printf` | 진행 중 |
| 2 | 채널별 발송기 (push / sms / email) | 인터페이스, 메서드, 맵, `v, ok` | |
| 3 | 발송 실패 처리 + 결과 집계 | `error`, `errors.New`, `fmt.Errorf`의 `%w`, `errors.Is` | |
| 4 | 알림 100개 동시 발송 | 고루틴, `sync.WaitGroup`, `sync.Mutex`, `go run -race` | |
| 5 | 느린 발송기 타임아웃 + 재시도 | `select`, `time.After`, 버퍼 채널 | |
| 6 | 수신 거부 사용자 제외 + 중복 발송 막기 | 맵, Mutex | |
| 7 (선택) | HTTP로 알림 요청 받기 | `net/http`, JSON, `context` | |

## 단계별 요구사항 요약

### 1단계: 알림 구조체 + 발송 출력
- `Notification` 구조체: `ID`, `UserID`, `Channel`, `Title`, `Body`
- `send(n Notification)` 함수로 `[채널] #ID → UserID: 제목 - 내용` 형식 출력
- 알림 3개를 슬라이스에 넣고 반복문으로 발송

### 2단계: 채널별 발송기
- `Sender` 인터페이스: `Send(n Notification) error`
- `PushSender`, `SMSSender`, `EmailSender`가 각자 다른 형식으로 발송
- `map[string]Sender`로 채널 이름에 맞는 발송기를 찾아 호출
- 없는 채널(`kakao`)은 에러 반환

### 3단계: 실패 처리 + 집계
- 미리 정해둔 에러: `ErrUnknownChannel`, `ErrEmptyTitle`, `ErrBodyTooLong`
- 푸시·이메일은 제목이 비면 실패, 문자는 내용이 20자를 넘으면 실패
- `errors.Is`로 실패 사유별 건수 집계

### 4단계: 동시 발송
- 발송 1건에 0.1초가 걸린다고 가정
- 순차 발송과 고루틴 동시 발송의 걸린 시간 비교
- Mutex 없이 결과를 집계할 때 생기는 경쟁 상태를 `-race`로 확인하고 수정

### 5단계: 타임아웃 + 재시도
- 응답이 늦은 발송기는 일정 시간이 지나면 포기
- 실패하면 정해진 횟수만큼 다시 시도

### 6단계: 수신 거부 + 중복 방지
- 수신 거부한 사용자에게는 보내지 않기
- 같은 알림 ID가 두 번 들어와도 한 번만 발송

## 참고 자료

- [A Tour of Go](https://go.dev/tour)
- [The Zen of Go](https://the-zen-of-go.netlify.app)
- [뱅크샐러드 Go 코딩 컨벤션](https://blog.banksalad.com/tech/go-best-practice-in-banksalad)
- [Go by Example: Context](https://gobyexample.com/context)