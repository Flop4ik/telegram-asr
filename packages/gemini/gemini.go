package gemini

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"google.golang.org/genai"
)

func RecognizeText(path string, recType string) (string, error) {

	var result string
	var promt string

	godotenv.Load()

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}
	myfile, err := client.Files.UploadFromPath(
		ctx,
		path,
		&genai.UploadFileConfig{
			MIMEType: "audio/ogg",
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("myfile=%+v\n", myfile)

	switch recType {
	case "transcribe":
		promt = "Пожалуйста, транскрибируй эту аудиозапись на русский язык как можно точнее. Не используй цензуру в словах, которые тебе кажутся непристойными"
	case "summarize":
		promt = "Пожалуйста, транскрибируй эту аудиозапись на русский язык как можно точнее. Не используй цензуру в словах, которые тебе кажутся непристойными" + "Напиши краткое содержание исходного текста (НЕ ПИШИ НИЧЕГО ПЕРЕД КРАТКИМ СОДЕРЖАНИЕМ, ТОЛЬКО САМО СОДЕРЖАНИЕ), пиши по пунктам, перед каждым пунктом <Номер пункта>., но не используй символы ** нигде кроме выделения важных слов в содержании (выделяй слова как *слово*).\n\nОтправь в формате краткое содержание ||| перевод звука в текст. ничего кроме текстов и символа <--||--> между текстами, ВООБЩЕ НИЧЕГО НЕ ПИШИ КРОМЕ ЭТОГО\n\nНЕ ПИШИ ПЕРЕД КРАТКИМ СОДРЕЖАНИЕМ НИЧЕГО, ТОЛЬКО САМ ТЕКСТ!!!!!!!"
	default:
		promt = "Пожалуйста, транскрибируй эту аудиозапись на русский язык как можно точнее. Не используй цензуру в словах, которые тебе кажутся непристойными"
	}

	parts := []*genai.Part{
		genai.NewPartFromURI(myfile.URI, myfile.MIMEType),
		genai.NewPartFromText(promt),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, "user"),
	}

	response, err := client.Models.GenerateContent(ctx, "gemini-1.5-flash", contents, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка при генерации ответа, пожалуйста подождите несколько минут и попробуйте снова")
	}
	text := response.Text()
	fmt.Printf("rРезультат:\n\n%s\n", text)

	switch recType {
	case "transcribe":
		result = fmt.Sprintf("🎤 *Расшифровка текста:*\n\n%s", text)
	case "summarize":
		splittedResponse := strings.Split(text, "|||")
		result = fmt.Sprintf("🎯 *Ключевые моменты:*\n\n%s\n\n\n🔍 *Расшифровка текста:*\n\n%s", splittedResponse[0], splittedResponse[1])
	default:
		result = fmt.Sprintf("🎤 *Расшифровка текста:*\n\n%s", text)

	}

	return result, nil
}
