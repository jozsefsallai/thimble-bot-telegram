package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/jozsefsallai/thimble-bot-telegram/config"
	"github.com/jozsefsallai/thimble-bot-telegram/utils"
	"github.com/lucsky/cuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func setSource(
	reader io.Reader,
	ctx *context.Context,
	snapshot *firestore.DocumentSnapshot,
	character string,
	id string,
	source string,
) string {
	var err error

	if len(source) == 0 {
		source, err = utils.ImageSourceLookup(reader)
		if err != nil {
			log.Println(err)
		}
	}

	data := snapshot.Data()
	data["source"] = source
	snapshot.Ref.Set(*ctx, data)
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

	if len(source) == 0 {
		photo2, _ := bot.GetFile(m.Photo.MediaFile())
		source, err = utils.ImageSourceLookup(photo2)
		if err != nil {
			bot.Send(m.Chat, fmt.Sprintf("Gophersauce Lookup Error: %s", err))
		}
	}

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

	rand.Seed(time.Now().UnixNano())
	target := refIDs[rand.Intn(len(refIDs))]

	document := collection.Doc(target)
	snapshot, err := document.Get(ctx)
	if err != nil {
		return "", "", nil, err
	}

	data := snapshot.Data()

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

	if len(source) == 0 {
		reader2, _ := object.NewReader(ctx)
		source = setSource(reader2, &ctx, snapshot, character, target, "")
	}

	return target, source, reader, nil
}

func deleteWrapper(sender int, character string, id string) error {
	if !utils.HasPermission(sender, config.GetConfig().Permissions.CanUploadMIA) {
		return errors.New("you don't have permission to use this command")
	}

	ctx := context.Background()
	firebase := utils.Firebase()

	firestore, err := firebase.Firestore(ctx)
	if err != nil {
		return err
	}

	collection := firestore.Collection(character)
	doc := collection.Doc(id)
	snapshot, err := doc.Get(ctx)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return errors.New("the requested document could not be found")
		}

		return err
	}

	snapshotData := snapshot.Data()
	image := snapshotData["objectName"]

	_, err = doc.Delete(ctx)
	if err != nil {
		return err
	}

	storage, err := firebase.Storage(ctx)
	if err != nil {
		return err
	}

	bucket, _ := storage.Bucket("thimble-bot.appspot.com")
	object := bucket.Object(fmt.Sprintf("%s/%s", character, image.(string)))
	err = object.Delete(ctx)
	if err != nil {
		return err
	}

	return nil
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

// DeleteNanaCommand will delete a Nanachi
func DeleteNanaCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		payload := m.Payload
		if len(payload) == 0 {
			bot.Send(m.Chat, "Please provide a Nanachi to delete.")
			return
		}

		err := deleteWrapper(m.Sender.ID, "nanachi", payload)
		if err != nil {
			bot.Send(m.Chat, fmt.Sprintf("Deletion Error: %s", err))
			return
		}

		bot.Send(m.Chat, fmt.Sprintf("Deleted %s", payload))
	}
}

// DeleteFapuCommand will delete a Faputa
func DeleteFapuCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		payload := m.Payload
		if len(payload) == 0 {
			bot.Send(m.Chat, "Please provide a Faputa to delete.")
			return
		}

		err := deleteWrapper(m.Sender.ID, "faputa", payload)
		if err != nil {
			bot.Send(m.Chat, fmt.Sprintf("Deletion Error: %s", err))
			return
		}

		bot.Send(m.Chat, fmt.Sprintf("Deleted %s", payload))
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

// SetSourceCommand will set the source of an image of Nanachi/Faputa
func SetSourceCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		if !utils.HasPermission(m.Sender.ID, config.GetConfig().Permissions.CanUploadMIA) {
			bot.Send(m.Chat, "You don't have permission to this command.")
			return
		}

		genericError := "Please provide the document ID, the character, and the source."

		payload := strings.Replace(m.Text, "/setsource ", "", 1)
		if len(payload) == 0 {
			bot.Send(m.Chat, genericError)
			return
		}

		components := strings.Split(payload, "\n")
		if len(components) != 3 {
			bot.Send(m.Chat, genericError)
			return
		}

		id := components[0]
		character := components[1]
		source := components[2]

		ctx := context.Background()
		firebase := utils.Firebase()

		firestore, err := firebase.Firestore(ctx)
		if err != nil {
			bot.Send(m.Chat, fmt.Sprintf("Error: %s", err))
			return
		}

		collection := firestore.Collection(character)
		doc := collection.Doc(id)
		snapshot, err := doc.Get(ctx)

		setSource(nil, &ctx, snapshot, character, id, source)

		bot.Send(m.Chat, "Done!")
	}
}
