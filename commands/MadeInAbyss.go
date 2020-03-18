package commands

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"strings"

	"github.com/jozsefsallai/thimble-bot-telegram/config"
	"github.com/jozsefsallai/thimble-bot-telegram/utils"
	"github.com/lucsky/cuid"
	tb "gopkg.in/tucnak/telebot.v2"
)

func parseCommandParams(payload string) string {
	components := strings.Split(payload, "\n")[1:]

	source := ""
	if len(components) == 1 {
		source = components[0]
	}

	return source
}

func upload(reader io.Reader, character string) (string, error) {
	ctx := context.Background()
	firebase := utils.Firebase()

	storage, err := firebase.Storage(ctx)
	if err != nil {
		return "", err
	}

	bucket, _ := storage.Bucket("thimble-bot.appspot.com")

	name := cuid.New()
	object := bucket.Object(fmt.Sprintf("%s/%s", character, name))
	writer := object.NewWriter(ctx)

	if _, err := io.Copy(writer, reader); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	return name, nil
}

func writeToFirestore(character string, objectName string, source string) (string, error) {
	ctx := context.Background()
	firebase := utils.Firebase()

	firestore, err := firebase.Firestore(ctx)
	if err != nil {
		return "", err
	}

	collection := firestore.Collection(character)

	data := map[string]interface{}{
		"objectName": objectName,
		"source":     source,
	}

	ref, _, err := collection.Add(ctx, data)
	if err != nil {
		return "", err
	}

	return ref.ID, nil
}

func addWrapper(bot *tb.Bot, m *tb.Message, character string) {
	if !utils.HasPermission(m.Sender.ID, config.GetConfig().Permissions.CanUploadMIA) {
		bot.Send(m.Chat, "You don't have permission to this command.")
		return
	}

	if m.Photo == nil {
		bot.Send(m.Chat, "Please provide an image to upload.")
		return
	}

	bot.Send(m.Chat, "Ok, uploading!")

	photo, err := bot.GetFile(m.Photo.MediaFile())
	if err != nil {
		bot.Send(m.Chat, "Failed to process your image.")
		return
	}

	name, err := upload(photo, character)
	if err != nil {
		bot.Send(m.Chat, fmt.Sprintf("Firebase Storage Error: %s", err))
		return
	}

	source := parseCommandParams(m.Caption)

	ref, err := writeToFirestore(character, name, source)
	if err != nil {
		bot.Send(m.Chat, fmt.Sprintf("Firebase Cloud Firestore Error: %s", err))
		return
	}

	bot.Send(m.Chat, fmt.Sprintf("Uploaded! %s", ref))
}

func getWrapper(character string) (string, string, io.Reader, error) {
	ctx := context.Background()
	firebase := utils.Firebase()

	firestore, err := firebase.Firestore(ctx)
	if err != nil {
		return "", "", nil, err
	}

	collection := firestore.Collection(character)
	refs, err := collection.DocumentRefs(ctx).GetAll()
	if err != nil {
		return "", "", nil, err
	}

	var refIDs []string
	for _, ref := range refs {
		refIDs = append(refIDs, ref.ID)
	}

	target := refIDs[rand.Intn(len(refIDs))]

	document, err := collection.Doc(target).Get(ctx)
	if err != nil {
		return "", "", nil, err
	}

	data := document.Data()

	objectKey := fmt.Sprintf("%s/%s", character, data["objectName"].(string))
	source := data["source"].(string)

	storage, err := firebase.Storage(ctx)
	if err != nil {
		return "", "", nil, err
	}

	bucket, _ := storage.Bucket("thimble-bot.appspot.com")

	object := bucket.Object(objectKey)
	reader, err := object.NewReader(ctx)
	if err != nil {
		return "", "", nil, err
	}

	return target, source, reader, nil
}

// AddNanaCommand will upload a picture of Nanachi to Firebase
func AddNanaCommand(bot *tb.Bot) func(*tb.Message) {
	return func(m *tb.Message) {
		addWrapper(bot, m, "nanachi")
	}
}

// AddFapuCommand will upload a picture of Faputa to Firebase
func AddFapuCommand(bot *tb.Bot) func(*tb.Message) {
	return func(m *tb.Message) {
		addWrapper(bot, m, "faputa")
	}
}

// NanachiCommand will return a random picture of Nanachi from Firebase
func NanachiCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		ref, source, reader, err := getWrapper("nanachi")
		if err != nil {
			bot.Send(m.Chat, "Failed to fetch.")
			fmt.Println(err)
		}

		caption := fmt.Sprintf("🐰 Naaaaaaaaa~\n\n🔑 %s", ref)

		if len(source) > 0 {
			caption += fmt.Sprintf("\n\n🔗 Source: %s", source)
		}

		image := &tb.Photo{
			File:    tb.FromReader(reader),
			Caption: caption,
		}

		bot.Send(m.Chat, image)
	}
}

// FaputaCommand will return a random picture of Faputa from Firebase
func FaputaCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		ref, source, reader, err := getWrapper("faputa")
		if err != nil {
			bot.Send(m.Chat, "Failed to fetch.")
			fmt.Println(err)
		}

		caption := fmt.Sprintf("🔑 %s", ref)

		if len(source) > 0 {
			caption += fmt.Sprintf("\n\n🔗 Source: %s", source)
		}

		image := &tb.Photo{
			File:    tb.FromReader(reader),
			Caption: caption,
		}

		bot.Send(m.Chat, image)
	}
}
