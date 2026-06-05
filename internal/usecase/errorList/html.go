package errorList

import (
	"fmt"
	"html"
	"time"
)

type typeHTMLData struct {
	Identifier   string
	Title        string
	Detail       string
	Status       int
	BaseURL      string
	UpdatedAt    time.Time
	Resolution   string
	SupportEmail string
}

type instanceHTMLData struct {
	Title      string
	Detail     string
	Type       string
	Status     int
	RequestURL string
	Method     string
	UserAgent  string
	TraceID    string
	OccurredAt time.Time
	TypeLink   string
}

func getHTTPStatusText(code int) string {
	statusTexts := map[int]string{
		400: "Bad Request",
		401: "Unauthorized",
		403: "Forbidden",
		404: "Not Found",
		500: "Internal Server Error",
		502: "Bad Gateway",
		503: "Service Unavailable",
	}
	if text, ok := statusTexts[code]; ok {
		return text
	}
	return "Error"
}

func buildTypeHTML(data typeHTMLData) string {
	identifier := html.EscapeString(data.Identifier)
	title := html.EscapeString(data.Title)
	detail := html.EscapeString(data.Detail)
	resolution := html.EscapeString(data.Resolution)
	supportEmail := html.EscapeString(data.SupportEmail)

	typeLink := "/v1/error?filter=type&identifier=" + identifier
	if data.BaseURL != "" {
		typeLink = data.BaseURL + typeLink
	}

	updatedAt := data.UpdatedAt.UTC().Format(time.RFC3339)

	// Build resolution section if provided
	resolutionSection := ""
	if resolution != "" {
		resolutionSection = fmt.Sprintf(`      <h2 class="text-xl font-semibold mt-8 mb-4">How to resolve?</h2>
      <p>%s</p>`, resolution)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en-us">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Erro: %s</title>
  <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-50 text-gray-800 font-sans">
  <div class="max-w-3xl mx-auto py-12 px-6">
    <header class="border-b pb-6 mb-8">
      <h1 class="text-3xl font-bold text-red-600">Erro %d: %s</h1>
      <p class="text-lg text-gray-600 mt-2">%s</p>
    </header>
    <section class="prose prose-slate">
      <h2 class="text-xl font-semibold mb-4">What does it mean?</h2>
      <p>%s</p>
%s
    </section>
    <div class="mt-8 pt-4 text-sm text-gray-600">
      <p><strong>Identifier:</strong> <code class="bg-gray-100 px-2 py-1 rounded">%s</code></p>
      <p><strong>Updated:</strong> %s</p>
    </div>
    <footer class="mt-12 pt-6 border-t text-sm text-gray-500 italic">Technical Documentation API - Version 1.0.0 - Support: %s</footer>
  </div>
</body>
</html>`, title, data.Status, title, detail, detail, resolutionSection, identifier, updatedAt, supportEmail)
}

func buildInstanceHTML(data instanceHTMLData) string {
	title := html.EscapeString(data.Title)
	detail := html.EscapeString(data.Detail)
	typeValue := html.EscapeString(data.Type)
	requestURL := html.EscapeString(data.RequestURL)
	method := html.EscapeString(data.Method)
	userAgent := html.EscapeString(data.UserAgent)
	traceID := html.EscapeString(data.TraceID)
	occurredAt := data.OccurredAt.UTC().Format(time.RFC3339)
	typeLink := html.EscapeString(data.TypeLink)

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en-us">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Erro: %s</title>
  <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-50 text-gray-800 font-sans">
  <div class="max-w-3xl mx-auto py-12 px-6">
    <header class="border-b pb-6 mb-8">
      <h1 class="text-3xl font-bold text-red-600">Erro %d: %s</h1>
      <p class="text-lg text-gray-600 mt-2">%s</p>
    </header>
    <section class="prose prose-slate">
      <h2 class="text-xl font-semibold mb-4">Request Details</h2>
      <ul class="list-disc pl-6 space-y-2">
        <li><strong>Occurred:</strong> %s</li>
        <li><strong>Method:</strong> %s</li>
        <li><strong>URL:</strong> %s</li>
        <li><strong>Trace ID:</strong> <code class="bg-gray-100 px-2 py-1 rounded">%s</code></li>
        <li><strong>User Agent:</strong> %s</li>
        <li><strong>Type:</strong> <a href="%s" class="text-blue-600">%s</a></li>
      </ul>
    </section>
    <div class="mt-8 pt-4 text-sm text-gray-600">
      <p><strong>Share your Trace ID with support</strong> so we can locate the exact failure.</p>
    </div>
    <footer class="mt-12 pt-6 border-t text-sm text-gray-500 italic">Technical Documentation API - Version 1.0.0 - Support: doceriadeliciasdaluoficial@gmail.com</footer>
  </div>
</body>
</html>`, title, data.Status, title, detail, occurredAt, method, requestURL, traceID, userAgent, typeLink, typeValue)
}
