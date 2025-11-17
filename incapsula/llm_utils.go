package incapsula

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	brtypes "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func newBedrockClient(ctx context.Context, region, profile string) (*bedrockruntime.Client, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}
	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return bedrockruntime.NewFromConfig(cfg), nil
}

func newMCPSession(ctx context.Context, cmd string, args ...string) (*mcp.ClientSession, error) {
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "go-bedrock-agent",
		Version: "0.1.0",
	}, nil)

	transport := &mcp.CommandTransport{
		Command: exec.Command(cmd, args...),
	}

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func callMCPTool(ctx context.Context, session *mcp.ClientSession, toolName string, args map[string]any) (string, error) {
	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      toolName,
		Arguments: args,
	})
	if err != nil {
		return "", fmt.Errorf("CallTool(%s) failed: %w", toolName, err)
	}
	if res.IsError {
		for _, c := range res.Content {
			if text, ok := c.(*mcp.TextContent); ok {
				fmt.Printf("MCP tool error: %s\n", text.Text)
			}
		}
		return "", fmt.Errorf("tool %s returned error: %v", toolName, res.Content)
	}

	var sb strings.Builder
	for _, c := range res.Content {
		if text, ok := c.(*mcp.TextContent); ok {
			sb.WriteString(text.Text)
			sb.WriteString("\n")
		}
	}
	return strings.TrimSpace(sb.String()), nil
}

type Agent struct {
	br      *bedrockruntime.Client
	modelID string
	mcpSess *mcp.ClientSession
}

type ToolCall struct {
	Tool      *string        `json:"tool"`
	Arguments map[string]any `json:"arguments"`
}

func (a *Agent) Answer(ctx context.Context, prompt string) (string, error) {

	messages := []brtypes.Message{
		{
			Role: brtypes.ConversationRoleUser,
			Content: []brtypes.ContentBlock{
				&brtypes.ContentBlockMemberText{Value: prompt},
			},
		},
	}

	out, err := a.br.Converse(ctx, &bedrockruntime.ConverseInput{
		ModelId:  aws.String(a.modelID),
		Messages: messages,
		InferenceConfig: &brtypes.InferenceConfiguration{
			MaxTokens:   aws.Int32(512),
			Temperature: aws.Float32(0.4),
			TopP:        aws.Float32(0.9),
		},
	})
	if err != nil {
		return "", fmt.Errorf("Converse failed: %w", err)
	}

	msg, ok := out.Output.(*brtypes.ConverseOutputMemberMessage)
	if !ok || len(msg.Value.Content) == 0 {
		return "", fmt.Errorf("unexpected Converse output")
	}

	textBlock, ok := msg.Value.Content[0].(*brtypes.ContentBlockMemberText)
	if !ok {
		return "", fmt.Errorf("first content block is not text")
	}
	return textBlock.Value, nil
}

func getToolToExecute(ctx context.Context, mcpSess *mcp.ClientSession, question string, agent Agent) (*ToolCall, error) {
	mcpTools, err := getMcpTools(ctx, mcpSess)
	if err != nil {
		log.Fatalf("failed to get MCP tools: %v", err)
	}

	toolsDesc, err := marshalToolsForPrompt(mcpTools)
	if err != nil {
		log.Fatalf("failed to marshal tools: %v", err)
	}

	// ----- First Bedrock call: which tool to use? -----
	prompt := buildToolSelectionPrompt(question, toolsDesc)

	selectionRaw, err := agent.Answer(ctx, prompt)

	if err != nil {
		log.Fatalf("Bedrock (selection) error: %v", err)
		return nil, errors.New("failed to get tool selection from Bedrock")
	}

	var tc ToolCall
	if err := json.Unmarshal([]byte(selectionRaw), &tc); err != nil {
		// Model didn't follow the JSON contract – just dump its text.
		fmt.Println("Bedrock replied (non-JSON):")
		fmt.Println(selectionRaw)
		return nil, nil
	}

	if tc.Tool == nil || *tc.Tool == "" {
		// Model decided no tool needed – just treat selectionRaw as the final answer.
		fmt.Println("Bedrock answer (no tool needed):")
		fmt.Println(selectionRaw)
		return nil, nil
	}

	return &tc, nil

}

func executeMCPTool(ctx context.Context, mcpSess *mcp.ClientSession, toolName string, argumensts map[string]any) (*mcp.CallToolResult, error) {
	return mcpSess.CallTool(ctx, &mcp.CallToolParams{
		Name:      toolName,
		Arguments: argumensts,
	})
}

func queryAgent(prompt string) (string, error) {
	ctx := context.Background()
	brClient, err := newBedrockClient(ctx, "us-west-2", "dev")
	if err != nil {
		log.Fatalf("failed to create Bedrock client: %v", err)
	}
	agent := &Agent{
		br:      brClient,
		modelID: "anthropic.claude-3-sonnet-20240229-v1:0",
	}

	finalAnswer, err := agent.Answer(ctx, prompt)
	if err != nil {
		log.Fatalf("Bedrock (answer) error: %v", err)
		return "", errors.New("llm fail to answer")
	}

	fmt.Println("Agent answered for prompt %s: %s", prompt, finalAnswer)
	return finalAnswer, nil
}

func answerWithTools(question string, api_id string, api_key string) (string, error) {
	ctx := context.Background()

	brClient, err := newBedrockClient(ctx, "us-west-2", "dev")
	if err != nil {
		log.Fatalf("failed to create Bedrock client: %v", err)
	}

	mcpArgs := strings.Fields("mcp-remote https://api.stage.impervaservices.com/cwaf-external-mcp/mcp/ --header X-Api-Id:" + api_id + " --header X-Api-Key:" + api_key)

	mcpSess, err := newMCPSession(ctx, "npx", mcpArgs...)

	if err != nil {
		log.Fatalf("failed to connect to MCP server: %v", err)
	}
	defer mcpSess.Close()

	agent := &Agent{
		br:      brClient,
		modelID: "anthropic.claude-3-sonnet-20240229-v1:0",
		mcpSess: mcpSess,
	}

	question = strings.TrimSpace(question)
	if question == "" {
		return "", errors.New("question cannot be empty")
	}
	tc, err := getToolToExecute(ctx, mcpSess, question, *agent)
	if err != nil {
		log.Fatalf("failed to get tool to execute: %v", err)
		return "", errors.New("failed to get tool to execute")
	}

	toolName := *tc.Tool
	fmt.Printf("Bedrock chose tool: %s\n", toolName)

	callToolRes, err := executeMCPTool(ctx, mcpSess, toolName, tc.Arguments)
	if err != nil {
		log.Fatalf("MCP CallTool %q failed: %v", toolName, err)
		return "", errors.New("tool call failed")
	}

	toolResultJSON, err := json.MarshalIndent(callToolRes, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal tool result: %v", err)
		return "", errors.New("tool call marshal failed")
	}

	answerPrompt := buildAnswerPrompt(question, toolName, tc.Arguments, string(toolResultJSON))

	finalAnswer, err := agent.Answer(ctx, answerPrompt)
	if err != nil {
		log.Fatalf("Bedrock (answer) error: %v", err)
		return "", errors.New("llm fail to answer")
	}

	fmt.Println("Agent answered for question %s: %s", question, finalAnswer)
	return finalAnswer, nil
}

func buildAnswerPrompt(userQuestion, toolName string, args map[string]any, toolResultJSON string) string {
	argsJSON, _ := json.MarshalIndent(args, "", "  ")
	return fmt.Sprintf(`
You are an assistant that explains results to the user.

USER QUESTION:
%s

You previously decided to call MCP tool "%s" with arguments:

%s

The MCP server returned this JSON result:

%s

Using ONLY the information above (plus general knowledge), write a clear and helpful answer to the user.
Do not show the raw JSON unless it is explicitly useful.
`, userQuestion, toolName, string(argsJSON), toolResultJSON)
}

func buildToolSelectionPrompt(userQuestion, toolsJSON string) string {
	return fmt.Sprintf(`
You are a tool-choosing agent. You have access to a Model Context Protocol (MCP) server that exposes these tools (JSON array):

%s

Each tool has:
- "name": the identifier used to call it
- "description": what it does
- "inputSchema": JSON schema for its arguments

USER QUESTION:
%s

Your job:

1. Decide whether you need to call exactly ONE of these tools.
2. If you do, respond ONLY with a single JSON object of the form:

{
  "tool": "<tool_name>",
  "arguments": { ... }
}

Where:
- "tool" is one of the tool.name values from the tools list.
- "arguments" is a JSON object matching that tool's inputSchema.

3. If no tool is required, respond with:

{
  "tool": null,
  "arguments": {}
}

IMPORTANT:
- Do NOT add extra fields.
- Do NOT wrap the JSON in backticks or text.
- Return only the raw JSON object.
`, toolsJSON, userQuestion)
}

func marshalToolsForPrompt(tools []*mcp.Tool) (string, error) {
	// We keep the schema but cut out Meta/annotations to reduce noise.
	type simpleTool struct {
		Name        string      `json:"name"`
		Title       string      `json:"title,omitempty"`
		Description string      `json:"description,omitempty"`
		InputSchema interface{} `json:"inputSchema,omitempty"`
	}

	compact := make([]simpleTool, 0, len(tools))
	for _, t := range tools {
		compact = append(compact, simpleTool{
			Name:        t.Name,
			Title:       t.Title,
			Description: t.Description,
			InputSchema: t.InputSchema,
		})
	}

	b, err := json.MarshalIndent(compact, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func getMcpTools(ctx context.Context, sess *mcp.ClientSession) ([]*mcp.Tool, error) {
	toolsIter := sess.Tools(ctx, &mcp.ListToolsParams{})
	var tools []*mcp.Tool
	for tool, err := range toolsIter {
		if err != nil {
			return nil, fmt.Errorf("error listing tools: %w", err)
		}
		tools = append(tools, tool)
	}

	return tools, nil
}
