package gemini

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"google.golang.org/genai"
)

func Summary(text string) (string, error) {
	godotenv.Load()

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_KEY"),
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(fmt.Sprintf("Напиши краткое содержание исходного текста (НЕ ПИШИ НИЧЕГО ПЕРЕД КРАТКИМ СОДЕРЖАНИЕМ, ТОЛЬКО САМО СОДЕРЖАНИЕ), пиши по пунктам, перед каждым пунктом <Номер пункта>., но не используй символы ** нигде кроме выделения важных слов в содержании (выделяй слова как *слово*).\n\nИсходный текст:\n%s", text)),
		nil,
	)

	if err != nil {
		return "", fmt.Errorf("ошибка при генерации ответа, пожалуйста подождите несколько минут и попробуйте снова: %w", err)
	}

	summary := result.Text()

	fmt.Printf("Результат:\n\n%s\n", summary)

	return summary, nil

}

func RecognizeText(path string) (string, error) {

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

	parts := []*genai.Part{
		genai.NewPartFromURI(myfile.URI, myfile.MIMEType),
		genai.NewPartFromText("Пожалуйста, транскрибируй эту аудиозапись на русский язык как можно точнее. Не используй цензуру в словах, которые тебе кажутся непристойными"),
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

	summary, err := Summary(text)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании краткого содержания: %w", err)
	}

	result := fmt.Sprintf("🎯 *Ключевые моменты:*\n\n%s\n\n...\n\n🔍 *Расшифровка текста:*\n\n%s", summary, text)

	return result, nil
}
