package gemini

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"google.golang.org/genai"
)

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
		genai.NewPartFromText("Пожалуйста, транскрибируй эту аудиозапись на русский язык как можно точнее."),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, "user"),
	}

	response, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash-lite", contents, nil)
	if err != nil {
		log.Fatal(err)
	}
	text := response.Text()
	fmt.Printf("rРезультат:\n\n%s\n", text)

	return text, nil
}
