package models

// ChatRequest represents a chat message from the user
type ChatRequest struct {
	Message          string   `json:"message"`
	BotType          BotType  `json:"botType"`                    // Loại bot để xử lý yêu cầu
	ConversationID   string   `json:"conversationId,omitempty"`   // ID for conversation tracking
	PreviousMessages []string `json:"previousMessages,omitempty"` // Previous message history
	Language         string   `json:"language,omitempty"`         // Language preference (en, vi)
	UserID           string   `json:"userId,omitempty"`           // Optional user identification for personalization
}

// ChatResponse represents the response to a chat message
type ChatResponse struct {
	Message           string            `json:"message"`                     // The text response
	Data              interface{}       `json:"data"`                        // Raw data from OpenFDA
	ChartData         interface{}       `json:"chartData"`                   // Processed data for charts
	ChartType         string            `json:"chartType"`                   // Type of chart to display (line, bar, etc.)
	ChartTitle        string            `json:"chartTitle"`                  // Title for the chart
	SourceDetails     string            `json:"sourceDetails"`               // Information about the data source
	DrugInfo          *DrugInfo         `json:"drugInfo,omitempty"`          // Drug information when applicable
	SideEffectReport  *SideEffectReport `json:"sideEffectReport,omitempty"`  // Side effect reporting info when applicable
	BotType           BotType           `json:"botType"`                     // Loại bot đã xử lý yêu cầu
	ConversationID    string            `json:"conversationId,omitempty"`    // ID for conversation tracking
	IsLoading         bool              `json:"isLoading,omitempty"`         // Indicates if additional data is being fetched
	FollowUpQuestions []string          `json:"followUpQuestions,omitempty"` // Suggested follow-up questions
	MedicationImages  []MedicationImage `json:"medicationImages,omitempty"`  // Images related to medications
	PriorityLevel     PriorityLevel     `json:"priorityLevel,omitempty"`     // Priority level for alerts (high for warnings)
	Language          string            `json:"language,omitempty"`          // Response language (en, vi)
}

// MedicationImage represents an image of a medication
type MedicationImage struct {
	URL         string `json:"url"`         // Image URL
	Description string `json:"description"` // Image description
	Source      string `json:"source"`      // Image source
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

// AnalysisResult represents the parsed response from Gemini for data analysis
type AnalysisResult struct {
	Summary       string      `json:"summary"`
	ChartType     string      `json:"chartType"`
	ChartTitle    string      `json:"chartTitle"`
	ChartData     interface{} `json:"chartData"`
	SourceDetails string      `json:"sourceDetails"`
}

// DrugInfo represents information about a specific drug
type DrugInfo struct {
	DrugName          string   `json:"drugName"`
	Indications       string   `json:"indications"`                 // Công dụng (dựa trên nhãn thuốc)
	SideEffects       []string `json:"sideEffects"`                 // Tác dụng phụ phổ biến
	Warnings          string   `json:"warnings"`                    // Cảnh báo quan trọng
	Recalls           []string `json:"recalls"`                     // Thông tin thu hồi
	Contraindications string   `json:"contraindications"`           // Chống chỉ định
	SourceDetails     string   `json:"sourceDetails"`               // Nguồn dữ liệu
	SeveritySorted    bool     `json:"severitySorted,omitempty"`    // If side effects are sorted by severity
	ImageURL          string   `json:"imageUrl,omitempty"`          // URL to drug image if available
	CommonlyAskedWith []string `json:"commonlyAskedWith,omitempty"` // Related medications often asked about together
	InteractionsWith  []string `json:"interactionsWith,omitempty"`  // Medications with known interactions
}

// SideEffectReport represents information about reporting a side effect
type SideEffectReport struct {
	DrugName       string `json:"drugName"`       // Tên thuốc gây tác dụng phụ
	SideEffectName string `json:"sideEffectName"` // Tên tác dụng phụ
	IsCommon       bool   `json:"isCommon"`       // Tác dụng phụ có phổ biến không
	Frequency      string `json:"frequency"`      // Tần suất báo cáo (cao, trung bình, thấp)
	ReportingSteps string `json:"reportingSteps"` // Hướng dẫn cách báo cáo
	ReportingLink  string `json:"reportingLink"`  // Đường dẫn tới trang báo cáo FDA
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

// BotType represents the type of chatbot being used
type BotType string

// Bot type constants
const (
	HealthTrendBot   BotType = "HealthTrendBot"   // Bot phân tích xu hướng
	MediBot          BotType = "MediBot"          // Bot thông tin thuốc
	SideEffectHelper BotType = "SideEffectHelper" // Bot hỗ trợ báo cáo tác dụng phụ
)
