package tg_test

import (
	"testing"

	"github.com/michurin/cnbot/pkg/tg"
	"github.com/stretchr/testify/assert"
)

const ordinaryMessage = `
{
  "ok": true,
  "result": [
    {
      "update_id": 50,
      "message": {
        "message_id": 1000,
        "from": {
          "id": 100,
          "is_bot": false,
          "first_name": "Alexey",
          "last_name": "Michurin",
          "username": "AlexeyMichurin",
          "language_code": "en"
        },
        "chat": {
          "id": 101,
          "first_name": "Alexey",
          "last_name": "Michurin",
          "username": "AlexeyMichurin",
          "type": "private"
        },
        "date": 1600000000,
        "text": "Text"
      }
    }
  ]
}`

const forwardFromUser = `{
  "ok": true,
  "result": [
    {
      "update_id": 50,
      "message": {
        "message_id": 1000,
        "from": {
          "id": 100,
          "is_bot": false,
          "first_name": "Alexey",
          "last_name": "Michurin",
          "username": "AlexeyMichurin",
          "language_code": "en"
        },
        "chat": {
          "id": 101,
          "first_name": "Alexey",
          "last_name": "Michurin",
          "username": "AlexeyMichurin",
          "type": "private"
        },
        "date": 1600000000,
        "forward_from": {
          "id": 500,
          "is_bot": false,
          "first_name": "user",
          "username": "username"
        },
        "forward_date": 1600000000,
        "text": "Text"
      }
    }
  ]
}`

const forwardFromBot = `{
  "ok": true,
  "result": [
    {
      "update_id": 50,
      "message": {
        "message_id": 1000,
        "from": {
          "id": 100,
          "is_bot": false,
          "first_name": "Alexey",
          "last_name": "Michurin",
          "username": "AlexeyMichurin",
          "language_code": "en"
        },
        "chat": {
          "id": 101,
          "first_name": "Alexey",
          "last_name": "Michurin",
          "username": "AlexeyMichurin",
          "type": "private"
        },
        "date": 1600000000,
        "forward_from": {
          "id": 500,
          "is_bot": true,
          "first_name": "net",
          "username": "net_bot"
        },
        "forward_date": 1600000000,
        "text": "Text"
      }
    }
  ]
}`

const forwardFromChannel = `{
  "ok": true,
  "result": [
    {
      "update_id": 50,
      "message": {
        "message_id": 970,
        "from": {
          "id": 100,
          "is_bot": false,
          "first_name": "Alexey",
          "last_name": "Michurin",
          "username": "AlexeyMichurin",
          "language_code": "ru"
        },
        "chat": {
          "id": 101,
          "first_name": "Alexey",
          "last_name": "Michurin",
          "username": "AlexeyMichurin",
          "type": "private"
        },
        "date": 1600000000,
        "forward_from_chat": {
          "id": -500,
          "title": "Title",
          "username": "username",
          "type": "channel"
        },
        "forward_from_message_id": 6129,
        "forward_date": 1600000000,
        "text": "Text"
      }
    }
  ]
}`

const contactMessage = `
{
  "ok": true,
  "result": [
    {
      "update_id": 50,
      "message": {
        "message_id": 1000,
        "from": {
          "id": 100,
          "is_bot": false,
          "first_name": "Alexey",
          "last_name": "Michurin",
          "username": "AlexeyMichurin",
          "language_code": "en"
        },
        "chat": {
          "id": 101,
          "first_name": "Alexey",
          "last_name": "Michurin",
          "username": "AlexeyMichurin",
          "type": "private"
        },
        "date": 1600000000,
        "contact": {
          "phone_number": "79999009999",
          "first_name": "Contact",
          "last_name": "last name",
          "user_id": 200
        }
      }
    }
  ]
}`

const contactMessageWithoutUserID = `
{
  "ok": true,
  "result": [
    {
      "update_id": 50,
      "message": {
        "message_id": 1000,
        "from": {
          "id": 100,
          "is_bot": false,
          "first_name": "Alexey",
          "last_name": "Michurin",
          "username": "AlexeyMichurin",
          "language_code": "en"
        },
        "chat": {
          "id": 101,
          "first_name": "Alexey",
          "last_name": "Michurin",
          "username": "AlexeyMichurin",
          "type": "private"
        },
        "date": 1600000000,
        "contact": {
          "phone_number": "79999009999",
          "first_name": "Contact",
          "last_name": "last name"
        }
      }
    }
  ]
}`

func TestDecodeGetUpdates(t *testing.T) {
	t.Parallel()
	t.Run("invalid_payload", func(t *testing.T) {
		mm, u, err := tg.DecodeGetUpdates([]byte("}"), 10, "one")
		assert.NotNil(t, err)
		assert.Equal(t, int64(10), u)
		assert.Len(t, mm, 0)
	})
	t.Run("not_ok", func(t *testing.T) {
		mm, u, err := tg.DecodeGetUpdates([]byte("{}"), 10, "one")
		assert.NotNil(t, err)
		assert.Equal(t, int64(10), u)
		assert.Len(t, mm, 0)
	})
	t.Run("no_result", func(t *testing.T) {
		mm, u, err := tg.DecodeGetUpdates([]byte(`{"ok": true}`), 10, "one")
		assert.Nil(t, err)
		assert.Equal(t, int64(10), u)
		assert.Len(t, mm, 0)
	})
	t.Run("max_update_id", func(t *testing.T) {
		mm, u, err := tg.DecodeGetUpdates([]byte(`{"ok": true, "result": [{"update_id": 50}, {"update_id": 100}, {"update_id": 40}]}`), 10, "one")
		assert.Nil(t, err)
		assert.Equal(t, int64(101), u)
		assert.Len(t, mm, 3)
	})
	for _, cs := range []struct {
		name, body, text, stype, sname string
		sid                            int64
	}{
		{"ordinary_message", ordinaryMessage, "Text", "", "", 0},
		{"forward_from_user", forwardFromUser, "Text", "user", "user", 500},
		{"forward_from_bot", forwardFromBot, "Text", "bot", "net", 500},
		{"forward_from_channel", forwardFromChannel, "Text", "channel", "Title", -500},
		{"contact_message", contactMessage, "", "contact", "Contact", 200},
		{"contact_message_without_user_id", contactMessageWithoutUserID, "", "", "", 0},
	} {
		cs := cs
		t.Run(cs.name, func(t *testing.T) {
			mm, u, err := tg.DecodeGetUpdates([]byte(cs.body), 10, "one")
			assert.Nil(t, err)
			assert.Equal(t, int64(51), u)
			assert.Len(t, mm, 1)
			assert.Equal(t, tg.Message{
				BotName:       "one",
				Text:          cs.text,
				FromID:        100,
				FromFirstName: "Alexey",
				ChatID:        101,
				SideType:      cs.stype,
				SideID:        cs.sid,
				SideName:      cs.sname,
			}, mm[0])
		})
	}
}

func TestEncodeGetUpdates(t *testing.T) {
	t.Parallel()
	t.Run("invalid_timeout", func(t *testing.T) {
		req, err := tg.EncodeGetUpdates(10, 0)
		assert.NotNil(t, err)
		assert.Nil(t, req)
	})
	t.Run("zero_offset", func(t *testing.T) {
		req, err := tg.EncodeGetUpdates(0, 10)
		assert.Nil(t, err)
		assert.Equal(t, &tg.Request{
			Method:      "getUpdates",
			ContentType: "application/json",
			Body:        []byte(`{"timeout":10,"allowed_updates":["message"]}`),
		}, req)
	})
	t.Run("none_zero_offset", func(t *testing.T) {
		req, err := tg.EncodeGetUpdates(100, 10)
		assert.Nil(t, err)
		assert.Equal(t, &tg.Request{
			Method:      "getUpdates",
			ContentType: "application/json",
			Body:        []byte(`{"offset":100,"timeout":10,"allowed_updates":["message"]}`),
		}, req)
	})
}
