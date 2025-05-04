package models

// ChatRequest represents a chat message from the user
type ChatRequest struct {
	Message          string   `json:"message"`
	ConversationID   string   `json:"conversationId,omitempty"`   // ID for conversation tracking
	PreviousMessages []string `json:"previousMessages,omitempty"` // Previous message history
	Language         string   `json:"language,omitempty"`         // Language preference (en, vi)
}

// ChatResponse represents the response to a chat message
type ChatResponse struct {
	Message           string        `json:"message"`                     // The text response
	Data              interface{}   `json:"data"`                        // Raw data from OpenFDA
	ConversationID    string        `json:"conversationId,omitempty"`    // ID for conversation tracking
	IsLoading         bool          `json:"isLoading,omitempty"`         // Indicates if additional data is being fetched
	FollowUpQuestions []string      `json:"followUpQuestions,omitempty"` // Suggested follow-up questions
	PriorityLevel     PriorityLevel `json:"priorityLevel,omitempty"`     // Priority level for alerts (high for warnings)
	Language          string        `json:"language,omitempty"`          // Response language (en, vi)
}

// PriorityLevel represents importance level of a message
type PriorityLevel string

// Priority level constants
const (
	PriorityHigh   PriorityLevel = "high"   // Critical warnings, recalls
	PriorityMedium PriorityLevel = "medium" // Important but not critical
	PriorityLow    PriorityLevel = "low"    // Informational content
)

// GeminiQueryResponse represents the parsed response from Gemini for query interpretation
type GeminiQueryResponse struct {
	QueryIntent        string      `json:"queryIntent"`
	SearchTerms        interface{} `json:"searchTerms"`
	TimeFrame          interface{} `json:"timeFrame"`
	DemographicFilters interface{} `json:"demographicFilters"`
	DataAggregation    interface{} `json:"dataAggregation"`
	IsDrugInfo         bool        `json:"isDrugInfo,omitempty"` // Flag for drug information lookup
}

// ConversationState tracks the state of a conversation with the chatbot
type ConversationState struct {
	ConversationID     string              `json:"conversationId"`
	UserID             string              `json:"userId,omitempty"`
	Messages           []MessageEntry      `json:"messages"`
	CurrentDrug        string              `json:"currentDrug,omitempty"`
	CurrentCondition   string              `json:"currentCondition,omitempty"`
	Language           string              `json:"language"`
	LastActivity       int64               `json:"lastActivity"` // Unix timestamp
	SuggestedFollowUps []string            `json:"suggestedFollowUps,omitempty"`
	CurrentContext     ConversationContext `json:"currentContext"`
}

// MessageEntry represents a single message in a conversation
type MessageEntry struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"` // Unix timestamp
	IsBot     bool   `json:"isBot"`     // Whether the message is from the bot
}

// ConversationContext represents the current context of the conversation
type ConversationContext string

// Context constants
const (
	ContextGeneral      ConversationContext = "general"
	ContextDrugInfo     ConversationContext = "drug_info"
	ContextSideEffect   ConversationContext = "side_effect"
	ContextHealthTrend  ConversationContext = "health_trend"
	ContextInteractions ConversationContext = "interactions"
)
