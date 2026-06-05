package errorList

import (
	"fmt"
	"html"
	"time"
)

type typeHTMLData struct {
	Identifier string
	Title      string
	Detail     string
	Status     int
	BaseURL    string
	UpdatedAt  time.Time
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

func buildTypeHTML(data typeHTMLData) string {
	identifier := html.EscapeString(data.Identifier)
	title := html.EscapeString(data.Title)
	detail := html.EscapeString(data.Detail)
	typeLink := "/v1/error?filter=type&identifier=" + identifier
	if data.BaseURL != "" {
		typeLink = data.BaseURL + typeLink
	}

	updatedAt := data.UpdatedAt.UTC().Format(time.RFC3339)

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>%s</title>
  <style>
    body { margin: 0; font-family: "Georgia", serif; background: linear-gradient(135deg, #f7f0e8, #fff); color: #2e2a27; }
    main { max-width: 860px; margin: 48px auto; padding: 32px; background: #fff; box-shadow: 0 20px 50px rgba(46, 42, 39, 0.12); border-radius: 18px; }
    h1 { margin: 0 0 12px; font-size: 32px; }
    .meta { display: grid; grid-template-columns: 1fr 1fr; gap: 8px 20px; margin-top: 20px; font-size: 14px; color: #5b524b; }
    .pill { display: inline-block; padding: 6px 12px; background: #f2e4d6; border-radius: 999px; font-size: 13px; margin-top: 8px; }
    .note { margin-top: 24px; font-size: 15px; color: #5b524b; }
    a { color: #8b5e3c; text-decoration: none; }
  </style>
</head>
<body>
  <main>
    <span class="pill">Error type</span>
    <h1>%s</h1>
    <p>%s</p>
    <div class="meta">
      <div><strong>Identifier:</strong> %s</div>
      <div><strong>Status:</strong> %d</div>
      <div><strong>Updated:</strong> %s</div>
      <div><strong>Doc:</strong> <a href="%s">%s</a></div>
    </div>
    <p class="note">If this keeps happening, please share the identifier with support so we can assist faster.</p>
  </main>
</body>
</html>`, title, title, detail, identifier, data.Status, updatedAt, typeLink, identifier)
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
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>%s</title>
  <style>
    body { margin: 0; font-family: "Georgia", serif; background: radial-gradient(circle at top, #fff6ec, #f8efe6); color: #2e2a27; }
    main { max-width: 920px; margin: 48px auto; padding: 32px; background: #fff; box-shadow: 0 20px 50px rgba(46, 42, 39, 0.12); border-radius: 18px; }
    h1 { margin: 0 0 12px; font-size: 30px; }
    .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px 24px; margin-top: 20px; font-size: 14px; color: #5b524b; }
    .badge { display: inline-block; padding: 6px 12px; background: #f0e2d3; border-radius: 999px; font-size: 13px; margin-top: 8px; }
    .note { margin-top: 24px; font-size: 15px; color: #5b524b; }
    a { color: #8b5e3c; text-decoration: none; }
  </style>
</head>
<body>
  <main>
    <span class="badge">Error instance</span>
    <h1>%s</h1>
    <p>%s</p>
    <div class="grid">
      <div><strong>Status:</strong> %d</div>
      <div><strong>Occurred:</strong> %s</div>
      <div><strong>Request:</strong> %s %s</div>
      <div><strong>Type:</strong> <a href="%s">%s</a></div>
      <div><strong>Trace ID:</strong> %s</div>
      <div><strong>User Agent:</strong> %s</div>
    </div>
    <p class="note">Share the trace ID with support so we can locate the exact failure.</p>
  </main>
</body>
</html>`, title, title, detail, data.Status, occurredAt, method, requestURL, typeLink, typeValue, traceID, userAgent)
}
