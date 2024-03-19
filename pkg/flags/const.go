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
	ModelGeminiPro                  = "models/gemini-pro"
	ModelGeminiProVision            = "models/gemini-pro-vision"
	ModelEmbedding001               = "models/embedding-001"
	ModelGemini1dot0Pro             = "models/gemini-1.0-pro"
	ModelGemini1dot0Pro001          = "models/gemini-1.0-pro-001"
	ModelGemini1dot0ProLatest       = "models/gemini-1.0-pro-latest"
	ModelGemini1dot0ProVisionLatest = "models/gemini-1.0-pro-vision-latest"
)

const (
	HarmProbabilityUnspecified = "unspecified"
	HarmProbabilityNegligible  = "negligible"
	HarmProbabilityLow         = "low"
	HarmProbabilityMedium      = "medium"
	HarmProbabilityHigh        = "high"
)
