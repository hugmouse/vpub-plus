// Code generated by go generate; DO NOT EDIT.

package handler

var TplCommonMap = map[string]string{
	"board_form": `{{ define "board_form" }}
    <div class="field">
        <label for="forumId">Forum</label>
        <select name="forumId" id="forumId">
            {{ range .Forums }}
                <option value="{{ .Id }}" {{ if eq .Id $.ForumId }}selected{{ end }}>{{ .Name }}</option>
            {{ end }}
        </select>
    </div>
    <div class="field">
        <label for="name">Name</label>
        <input type="text" name="name" id="name" value="{{ .Name }}" autocomplete="off" maxlength="120" required
               autofocus/>
    </div>
    <div class="field">
        <label for="description">Description</label>
        <textarea class="editor" name="description" id="description" required>{{ .Description }}</textarea>
    </div>
    <div class="field">
        <label for="position">Position</label>
        <input type="number" name="position" id="position" value="{{ .Position }}" autocomplete="off" required/>
    </div>
    <div class="field">
        <label for="locked">Locked</label>
        <select name="locked" id="locked">
            <option value="false">false</option>
            <option value="true" {{ if .IsLocked }}selected{{ end }}>true</option>
        </select>
    </div>
{{ end }}`,
	"forum_form": `{{ define "forum_form" }}
    <div class="field">
        <label for="name">Name</label>
        <input type="text" name="name" id="name" value="{{ .Name }}" autocomplete="off" maxlength="120" required
               autofocus/>
    </div>
    <div class="field">
        <label for="position">Position</label>
        <input type="number" name="position" id="position" value="{{ .Position }}" autocomplete="off" required/>
    </div>
    <div>
        <input type="checkbox" id="locked" name="locked" {{ if .IsLocked }}checked{{ end }}>
        <label for="locked">Locked</label>
    </div>
{{ end }}`,
	"forum_nav": `{{ define "forum_nav" }}
    <nav class="breadcrumb">
        <a href="/">Forums</a>
        {{ if .Forum.Name }}
            {{ if .Board.Name }}
                › <a href="/forums/{{ .Forum.Id }}">{{ .Forum.Name }}</a>
                {{ if .Topic }}
                    › <a href="/boards/{{ .Board.Id }}">{{ .Board.Name }}</a> › {{ .Topic }}
                {{ else }}
                    › {{ .Board.Name }}
                {{ end }}
            {{ else }}
                › {{ .Forum.Name }}
            {{ end }}
        {{ end }}
    </nav>
{{ end }}`,
	"layout": `{{ define "layout" }}
    <!DOCTYPE html>
    <html lang="{{ .settings.Lang }}">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        {{ if .board.Description }}
            {{ if not .navigation.Topic }}
                <meta property="description" content="{{ .board.Description }}">
                <meta property="og:description" content="{{ .board.Description }}">
            {{ end }}
        {{ end }}
        <link rel="stylesheet" href="/style.css">
        {{ if .navigation.Forum.Name }}
            {{ if .navigation.Board.Name }}
                {{ if .navigation.Topic }}
                    <title>{{ .navigation.Topic }} | {{ .settings.Name }}</title>
                {{ else }}
                    <title>{{ .navigation.Board.Name }} | {{ .settings.Name }}</title>
                {{ end }}
            {{ else }}
                <title>{{ .navigation.Forum.Name }} | {{ .settings.Name }}</title>
            {{ end }}
        {{ else }}
            <title>{{ .settings.Name }}</title>
        {{ end }}
        {{ template "head" . }}
    </head>
    <body>
    <header>
        <nav>
            <a href="/">home</a> <a href="/posts">posts</a> <a href="/feed.atom">atom</a>
            {{ if logged }}
                <a href="/users/{{ .logged.Id }}">{{ .logged.Name }}</a> <a href="/account">account</a> <a
                        href="/logout">logout</a>
            {{ else }}
                <a href="/login">login</a> <a href="/register">register</a>
            {{ end }}
        </nav>
    </header>
    {{ if .errors }}
        <div class="errors flash">
            <ul>
                {{ range .errors }}
                    <li>{{ . }}</li>
                {{ end }}
            </ul>
        </div>
    {{ end }}
    {{ if .info }}
        <div class="info flash">
            <ul>
                {{ range .info }}
                    <li>{{ . }}</li>
                {{ end }}
            </ul>
        </div>
    {{ end }}
    <main>
        {{ template "content" . }}
    </main>
    {{ if .settings.Footer }}
        <footer>
            {{ html .settings.Footer }}
        </footer>
    {{ end }}
    </body>
    </html>
{{ end }}
{{ define "head" }}{{ end }}
{{ define "breadcrumb" }}{{ end }}`,
	"pagination": `{{ define "pagination" }}
    {{ if or (ne 1 .Page) .HasMore }}
        <p>
            {{ if ne 1 .Page }}
                <a href="?page={{ dec .Page }}">Previous page</a>
            {{ end }}
            {{ if .HasMore }}
                <a href="?page={{ inc .Page }}">Next page</a>
            {{ end }}
        </p>
    {{ end }}
{{ end }}`,
	"post_form": `{{ define "post_form" }}
    <input type="hidden" name="topicId" value="{{ .TopicId }}">
    <div class="field">
        <label for="subject">Subject</label>
        <input type="text" name="subject" id="subject" value="{{ .Subject }}" autocomplete="off" maxlength="115"
               required autofocus/>
    </div>
    <div class="field">
        <label for="content">Content</label>
        <textarea class="editor" name="content" id="content" required>{{ .Content }}</textarea>
    </div>
{{ end }}`,
	"topic_form": `{{ define "topic_form" }}
    <details>
        <summary>Admin options</summary>
        <div class="field">
            <label for="newBoardId">Board</label>
            <select name="newBoardId" id="newBoardId">
                {{ range .Boards }}
                    <option value="{{ .Id }}" {{ if eq .Id $.BoardId }}selected{{ end }}>{{ .Name }}</option>
                {{ end }}
            </select>
        </div>
        <div>
            <input type="checkbox" id="sticky" name="sticky" {{ if .IsSticky }}checked{{ end }}>
            <label for="sticky">Sticky</label>
        </div>
        <div>
            <input type="checkbox" id="locked" name="locked" {{ if .IsLocked }}checked{{ end }}>
            <label for="locked">Locked</label>
        </div>
    </details>
{{ end }}`,
}
