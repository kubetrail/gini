package flags

const (
	ApiKey               = "api-key"
	Model                = "model"
	AutoSave             = "auto-save"
	Render               = "render"
	AllowHarmProbability = "allow-harm-probability"
	TopK                 = "top-k"
	TopP                 = "top-p"
	Temperature          = "temperature"
	CandidateCount       = "candidate-count"
	MaxOutputTokens      = "max-output-tokens"
	File                 = "file"
	Format               = "format"
)

const (
	RenderFormatHtml     = "html"
	RenderFormatMarkdown = "markdown"
	RenderFormatPretty   = "pretty"
)

const (
	FormatPng  = "image/png"
	FormatJpeg = "image/jpeg"
	FormatHeic = "image/heic"
	FormatHeif = "image/heif"
	FormatWebp = "image/webp"
	FormatPdf  = "application/pdf"
)

const (
	MaxBlobBufferSizeBytes = 4194304
)

const (
	ApiKeyEnv = "GOOGLE_API_KEY"
)

const (
	M00 = "models/chat-bison-001"
	M01 = "models/text-bison-001"
	M02 = "models/embedding-gecko-001"
	M03 = "models/gemini-1.0-pro-latest"
	M04 = "models/gemini-1.0-pro"
	M05 = "models/gemini-pro"
	M06 = "models/gemini-1.0-pro-001"
	M07 = "models/gemini-1.0-pro-vision-latest"
	M08 = "models/gemini-pro-vision"
	M09 = "models/gemini-1.5-pro-latest"
	M10 = "models/gemini-1.5-pro-001"
	M11 = "models/gemini-1.5-pro-002"
	M12 = "models/gemini-1.5-pro"
	M13 = "models/gemini-1.5-pro-exp-0801"
	M14 = "models/gemini-1.5-pro-exp-0827"
	M15 = "models/gemini-1.5-flash-latest"
	M16 = "models/gemini-1.5-flash-001"
	M17 = "models/gemini-1.5-flash-001-tuning"
	M18 = "models/gemini-1.5-flash"
	M19 = "models/gemini-1.5-flash-exp-0827"
	M20 = "models/gemini-1.5-flash-002"
	M21 = "models/gemini-1.5-flash-8b"
	M22 = "models/gemini-1.5-flash-8b-001"
	M23 = "models/gemini-1.5-flash-8b-latest"
	M24 = "models/gemini-1.5-flash-8b-exp-0827"
	M25 = "models/gemini-1.5-flash-8b-exp-0924"
	M26 = "models/gemini-exp-1114"
	M27 = "models/embedding-001"
	M28 = "models/text-embedding-004"
	M29 = "models/aqa"
)

var Models = []string{
	M00,
	M01,
	M02,
	M03,
	M04,
	M05,
	M06,
	M07,
	M08,
	M09,
	M10,
	M11,
	M12,
	M13,
	M14,
	M15,
	M16,
	M17,
	M18,
	M19,
	M20,
	M21,
	M22,
	M23,
	M24,
	M25,
	M26,
	M27,
	M28,
	M29,
}

const (
	HarmProbabilityUnspecified = "unspecified"
	HarmProbabilityNegligible  = "negligible"
	HarmProbabilityLow         = "low"
	HarmProbabilityMedium      = "medium"
	HarmProbabilityHigh        = "high"
)
