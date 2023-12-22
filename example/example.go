package main

import (
	"context"
	"fmt"
	f "fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	api "github.com/weavel-ai/promptmodel-go/api"
	promptmodel_client "github.com/weavel-ai/promptmodel-go/client"
)

// formatting fstring for each prompts
func formatWithMap(format string, values map[string]string) string {
	for key, value := range values {
		placeholder := fmt.Sprintf("{%s}", key)
		format = strings.ReplaceAll(format, placeholder, value)
	}
	return format
}


func main() {
	// Create PromptModel Client Instance
	pmClient := promptmodel_client.NewClient() // API 키 설정

	// fetch FunctionMoelVersion with promptmodel_client
	responseChan := make(chan *api.FetchFunctionModelVersionResponseInstance)
	errorChan := make(chan error)

	// Add your code here
	var function_model_name = "summarize" // Name for your own FunctionModel in Promptmodel Dashboard
	var function_model_version = "1" // Version for your own FunctionModel in Promptmodel Dashboard

	var request = &api.FetchFunctionModelVersionRequest{
		FunctionModelName: function_model_name, 
		Version:           &function_model_version,
	}
	
	go pmClient.Api.FetchFunctionModelVersionAsync(context.Background(), request, responseChan, errorChan) 

	var functionModelVersionConfig *api.FetchFunctionModelVersionResponseInstance

	select {
	case response := <-responseChan:
		// Handle response
		functionModelVersionConfig = response
	case err := <-errorChan:
		// Handle error
		f.Println(err)
		return
	}

	// Create OpenAI Client Instance
	var openai_api_key = os.Getenv("OPENAI_API_KEY")

	oaClient := openai.NewClient(openai_api_key)

	// openAI API call
	var prompts = functionModelVersionConfig.Prompts
	var inputs = map[string]string{"text": "Large Language Models are few shot learners: ..."}

	f.Println("Prompts: ", prompts)

	// formatting string for each prompt
	for i, prompt := range prompts {
		prompts[i].Content = formatWithMap(prompt.Content, inputs)
	}

	var prompts_for_openai []openai.ChatCompletionMessage
	for _, prompt := range prompts {
		prompts_for_openai = append(prompts_for_openai, openai.ChatCompletionMessage{Role: prompt.Role, Content: prompt.Content})
	}

	resp, err := oaClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: functionModelVersionConfig.FunctionModelVersions[0].Model,
			Messages: prompts_for_openai,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}