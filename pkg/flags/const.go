package flags

const (
	ApiKey               = "api-key"
	Model                = "model"
	AutoSave             = "auto-save"
	AllowHarmProbability = "allow-harm-probability"
	TopK                 = "top-k"
	TopP                 = "top-p"
	Temperature          = "temperature"
	CandidateCount       = "candidate-count"
	MaxOutputTokens      = "max-output-tokens"
)

const (
	ApiKeyEnv = "GOOGLE_API_KEY"
)

const (
	ModelGeminiPro       = "models/gemini-pro"
	ModelGeminiProVision = "models/gemini-pro-vision"
	ModelEmbedding001    = "models/embedding-001"
)

const (
	HarmProbabilityUnspecified = "unspecified"
	HarmProbabilityNegligible  = "negligible"
	HarmProbabilityLow         = "low"
	HarmProbabilityMedium      = "medium"
	HarmProbabilityHigh        = "high"
)
