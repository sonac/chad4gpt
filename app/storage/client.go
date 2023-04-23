package storage

import (
	"encoding/json"
	"io/ioutil"

	"github.com/rs/zerolog/log"
)

type Storage struct {
  Chats []Chat 
}

type Chat struct {
  Name string
  Language string
  ChatId int64
}

const storageFile = "data.json"

func NewStorage() *Storage {
  var chats []Chat
  d, err := ioutil.ReadFile(storageFile)
  if err != nil {
    log.Fatal().Err(err).Msg("[Error] while reading storage file")
  }
  err = json.Unmarshal(d, &chats)
  if err != nil {
    log.Fatal().Err(err).Msg("[Error] while unmarshalling storage file")
  }

  log.Printf("data parsed: %+v", chats)
  return &Storage{Chats: chats}
}

func (strg *Storage) Save() {
  fileData, _ := json.MarshalIndent(strg.Chats, "", " ")
  err := ioutil.WriteFile(storageFile, fileData, 0644)
  if err != nil {
    log.Fatal().Err(err).Msg("[Error] while saving file")
  }
}
